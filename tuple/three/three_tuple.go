// Package three provides a 3-tuple type.
package three

// Tuple contains three values.
type Tuple[T, T2, T3 any] struct {
	V  T
	V2 T2
	V3 T3
}

// TupleOf creates a Tuple of v, v2, and v3.
// The optional ok argument aids interoperability with
// return values that follow the "comma ok" idiom.
func TupleOf[T, T2, T3 any](v T, v2 T2, v3 T3, ok ...bool) (Tuple[T, T2, T3], bool) {
	if len(ok) == 0 {
		ok = []bool{true}
	}
	return Tuple[T, T2, T3]{V: v, V2: v2, V3: v3}, ok[0]
}
