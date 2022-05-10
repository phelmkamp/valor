package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "valorcheck",
	Doc:  "Checks that access to an optional value is guarded against the case where the value is not present.",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		v := visitor{
			pass:    pass,
			guarded: make(map[*types.Var]struct{}),
			checked: make(map[*ast.SelectorExpr]struct{}),
		}
		ast.Walk(&v, f)

	}
	return nil, nil
}

type visitor struct {
	pass    *analysis.Pass
	guarded map[*types.Var]struct{}
	checked map[*ast.SelectorExpr]struct{}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	// peek ahead to mark checked
	if assign, ok := node.(*ast.AssignStmt); ok {
		for _, rh := range assign.Rhs {
			call, ok := rh.(*ast.CallExpr)
			if !ok {
				continue
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			v.checked[sel] = struct{}{}
		}
		return v
	}
	if ifStmt, ok := node.(*ast.IfStmt); ok {
		call, ok := ifStmt.Cond.(*ast.CallExpr)
		if !ok {
			// might be negated
			unary, ok := ifStmt.Cond.(*ast.UnaryExpr)
			if !ok {
				return v
			}
			call, ok = unary.X.(*ast.CallExpr)
			if !ok {
				return v
			}
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return v
		}
		v.checked[sel] = struct{}{}
		return v
	}

	sel, ok := node.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	var xVar *types.Var
	if xIdent, ok := sel.X.(*ast.Ident); ok {
		xVar, ok = v.pass.TypesInfo.ObjectOf(xIdent).(*types.Var)
		if !ok {
			return v
		}
	} else {
		// might be call returning a variable
		call, ok := sel.X.(*ast.CallExpr)
		if !ok {
			return v
		}
		funSel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return v
		}
		fun, ok := v.pass.TypesInfo.ObjectOf(funSel.Sel).(*types.Func)
		if !ok {
			return v
		}
		sig, ok := fun.Type().(*types.Signature)
		if !ok {
			return v
		}
		if sig.Results().Len() != 1 {
			return v
		}
		xVar = sig.Results().At(0)
	}

	xType, ok := xVar.Type().(*types.Named)
	if !ok {
		return v
	}

	if xType.Obj().Pkg().Path() == "github.com/phelmkamp/valor/optional" {
		v.checkOpt(sel, xVar)
	}

	return v
}

func (v *visitor) checkOpt(sel *ast.SelectorExpr, optVar *types.Var) {
	switch sel.Sel.Name {
	case "IsOk":
		v.guarded[optVar] = struct{}{}
	case "MustOk":
		if _, ok := v.guarded[optVar]; !ok {
			// TODO: consider calling Report with suggested fix
			v.pass.ReportRangef(sel, "call to MustOk not guarded by IsOk might panic")
		}
	case "Ok":
		if _, ok := v.checked[sel]; !ok {
			v.pass.ReportRangef(sel, "result of Ok is not checked")
		}
	}
}
