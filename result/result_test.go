// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package result_test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/tuple/two"
	"github.com/phelmkamp/valor/tuple/unit"
)

// type checks
var (
	_ = result.OfError[*struct{}](errFail)
	_ = result.OfError[map[string]struct{}](errFail)
	_ = result.OfError[[]struct{}](errFail)
)

func Example() {
	var w strings.Builder
	// traditional
	if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
		fmt.Println(res.Error())
		return
	}

	// try to get value, printing wrapped error if not ok
	// note: only relevant values are in-scope after handling
	var n int
	if res := result.Of(w.Write([]byte("foo"))); !res.Value().Ok(&n) {
		fmt.Println(res.Errorf("Write() failed: %w").Error())
		return
	}
	fmt.Println(n)

	// same as above with multiple values
	var s string
	var b bool
	if res := two.TupleResultOf(multi(false)); !res.Value().Do(
		func(t two.Tuple[string, bool]) { s, b = t.Values() },
	).IsOk() {
		fmt.Println(res.Errorf("multi() failed: %w").Error())
		return
	}
	fmt.Println(s, b)

	// errors.Is
	if res := result.Of(w.Write([]byte("foo"))); res.ErrorIs(io.ErrShortWrite) {
		fmt.Println(res.Error())
		return
	}

	// errors.As
	if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
		var err *fs.PathError
		if res.ErrorAs(&err) {
			fmt.Println("path=" + err.Path)
		}
		fmt.Println(res.Error())
		return
	}

	// errors.Unwrap
	if res := result.Of(mid(true)); res.IsError() {
		fmt.Println(res.ErrorUnwrap().Error())
	}

	// switch
	switch res := result.Of(w.Write([]byte("foo"))); res {
	case res.OfOk():
		n = res.Value().MustOk()
		fmt.Println("Ok")
	case res.OfError():
		fmt.Println(res.Error())
		return
	}
	// Output: 3
	// foo true
	// fail
	// Ok
}

func mid(fail bool) (string, error) {
	var s string
	if res := result.Of(leaf(fail)); !optional.Map(res.Value(), strconv.Itoa).Ok(&s) {
		return "", res.Errorf("leaf() failed: %w").Error()
	}
	return "", nil
}

var errFail = errors.New("fail")

func leaf(fail bool) (int, error) {
	if fail {
		return 0, errFail
	}
	return 1, nil
}

func multi(fail bool) (string, bool, error) {
	if fail {
		return "", false, errFail
	}
	return "foo", true, nil
}

func TestOf(t *testing.T) {
	if got := result.Of(leaf(true)); got != result.OfError[int](errFail) {
		t.Errorf("Of() = %v, want %v", got, result.OfError[int](errFail))
	}
	if got := result.Of(leaf(false)); got != result.OfOk(1) {
		t.Errorf("Of() = %v, want %v", got, result.OfOk(1))
	}
}

func TestOfError(t *testing.T) {
	// Special case: OfError(nil) has no error AND a not-ok Value
	got := result.OfError[unit.Type](nil)
	if got.IsError() {
		t.Errorf("IsError() after OfError(nil) = %v, want %v", got.IsError(), false)
	}
	if got.Value().IsOk() {
		t.Errorf("Value().IsOk() after OfError(nil) = %v, want %v ", got.Value().IsOk(), false)
	}
}

func TestOfValue(t *testing.T) {
	if got := result.OfValue(optional.OfNotOk[string](), errFail); got != result.OfError[string](errFail) {
		t.Errorf("OfValue() = %v, want %v", got, result.OfError[string](errFail))
	}
	if got := result.OfValue(optional.OfOk("foo"), errFail); got != result.OfOk("foo") {
		t.Errorf("OfValue() = %v, want %v", got, result.OfOk("foo"))
	}
}

func TestResult_Error(t *testing.T) {
	if got := result.OfError[float64](errFail).Error(); got != errFail {
		t.Errorf("Error() error = %v, want %v", got, errFail)
	}
	if got := result.OfOk(1.0).Error(); got != nil {
		t.Errorf("Error() error = %v, want %v", got, nil)
	}
}

func TestResult_ErrorAs(t *testing.T) {
	var gotTarget *fs.PathError
	if got := result.OfOk(1).ErrorAs(&gotTarget); got {
		t.Errorf("ErrorAs() = %v, want %v", got, false)
	}
	if gotTarget != nil {
		t.Errorf("gotTarget after ErrorAs() = %v, wantTarget %v", gotTarget, nil)
	}

	if got := result.OfError[int](errFail).ErrorAs(&gotTarget); got {
		t.Errorf("ErrorAs() = %v, want %v", got, false)
	}
	if gotTarget != nil {
		t.Errorf("gotTarget after ErrorAs() = %v, wantTarget %v", gotTarget, nil)
	}

	if got := result.OfError[int](&fs.PathError{Path: "/foo/bar"}).ErrorAs(&gotTarget); !got {
		t.Errorf("ErrorAs() = %v, want %v", got, true)
	}
	wantTarget := &fs.PathError{Path: "/foo/bar"}
	if *gotTarget != *wantTarget {
		t.Errorf("gotTarget after ErrorAs() = %v, wantTarget %v", gotTarget, wantTarget)
	}
}

func TestResult_ErrorIs(t *testing.T) {
	target := &fs.PathError{Path: "/foo/bar"}
	if got := result.OfOk(1).ErrorIs(target); got {
		t.Errorf("ErrorIs() = %v, want %v", got, false)
	}
	if got := result.OfError[int](errFail).ErrorIs(target); got {
		t.Errorf("ErrorIs() = %v, want %v", got, false)
	}
	if got := result.OfError[int](target).ErrorIs(target); !got {
		t.Errorf("ErrorIs() = %v, want %v", got, true)
	}
}

func TestResult_ErrorUnwrap(t *testing.T) {
	if got := result.Of(leaf(false)).ErrorUnwrap(); got != result.OfOk(1) {
		t.Errorf("ErrorUnwrap() = %v, want %v", got, result.OfOk(1))
	}
	if got := result.Of(leaf(true)).ErrorUnwrap(); got != result.OfError[int](errFail) {
		t.Errorf("ErrorUnwrap() = %v, want %v", got, result.OfError[int](errFail))
	}
	if got := result.Of(mid(true)).ErrorUnwrap(); got != result.OfError[string](errFail) {
		t.Errorf("ErrorUnwrap() = %v, want %v", got, result.OfError[string](errFail))
	}
}

func TestResult_Errorf(t *testing.T) {
	if got := result.Of(leaf(false)).Errorf("leaf() failed: %w"); got != result.OfOk(1) {
		t.Errorf("Errorf() = %v, want %v", got, result.OfOk(1))
	}
	if got := result.Of(leaf(true)).Errorf("leaf() failed: %w"); !reflect.DeepEqual(got, result.OfError[int](fmt.Errorf("leaf() failed: %w", errFail))) {
		t.Errorf("Errorf() = %v, want %v", got, result.OfError[int](fmt.Errorf("leaf() failed: %w", errFail)))
	}
}

func TestResult_IsError(t *testing.T) {
	if got := result.OfOk("foo").IsError(); got {
		t.Errorf("IsError() = %v, want %v", got, false)
	}
	if got := result.OfError[string](errFail).IsError(); !got {
		t.Errorf("IsError() = %v, want %v", got, true)
	}
}

func TestResult_String(t *testing.T) {
	if got := result.OfOk(1.5).String(); got != "{1.5 <nil>}" {
		t.Errorf("String() = %v, want %v", got, "1.5")
	}
	if got := result.OfError[float64](errFail).String(); got != "{0 fail}" {
		t.Errorf("String() = %v, want %v", got, "{0 fail}")
	}
}

func TestResult_Value(t *testing.T) {
	if got := result.OfOk(1).Value(); got != optional.OfOk(1) {
		t.Errorf("Value() = %v, want %v", got, optional.OfOk(1))
	}
	if got := result.OfError[int](errFail).Value(); got != optional.OfNotOk[int]() {
		t.Errorf("Value() = %v, want %v", got, optional.OfNotOk[int]())
	}
}

func TestTranspose(t *testing.T) {
	if got := result.Transpose(result.OfError[optional.Value[string]](errFail)); got != optional.OfOk(result.OfError[string](errFail)) {
		t.Errorf("Transpose() = %v, want %v", got, optional.OfOk(result.OfError[string](errFail)))
	}
	if got := result.Transpose(result.OfOk(optional.OfNotOk[string]())); got != optional.OfNotOk[result.Result[string]]() {
		t.Errorf("Transpose() = %v, want %v", got, optional.OfNotOk[result.Result[string]]())
	}
	if got := result.Transpose(result.OfOk(optional.OfOk("foo"))); got != optional.OfOk(result.OfOk("foo")) {
		t.Errorf("Transpose() = %v, want %v", got, optional.OfOk(result.OfOk("foo")))
	}
}

func TestTransposeValue(t *testing.T) {
	if got := result.TransposeValue(optional.OfNotOk[result.Result[int]]()); got != result.OfOk(optional.OfNotOk[int]()) {
		t.Errorf("TransposeValue() = %v, want %v", got, result.OfOk(optional.OfNotOk[int]()))
	}
	if got := result.TransposeValue(optional.OfOk(result.OfError[int](errFail))); got != result.OfError[optional.Value[int]](errFail) {
		t.Errorf("TransposeValue() = %v, want %v", got, result.OfError[optional.Value[int]](errFail))
	}
	if got := result.TransposeValue(optional.OfOk(result.OfOk(1))); got != result.OfOk(optional.OfOk(1)) {
		t.Errorf("TransposeValue() = %v, want %v", got, result.OfOk(optional.OfOk(1)))
	}
}

func TestResult_OfOk(t *testing.T) {
	if got := result.OfError[float64](errFail).OfOk(); got != result.OfOk(0.0) {
		t.Errorf("OfOk() = %v, want %v", got, result.OfOk(0.0))
	}
	if got := result.OfOk(1.0).OfOk(); got != result.OfOk(1.0) {
		t.Errorf("OfOk() = %v, want %v", got, result.OfOk(1.0))
	}
}

func TestResult_OfError(t *testing.T) {
	if got := result.OfError[float64](errFail).OfError(); got != result.OfError[float64](errFail) {
		t.Errorf("OfError() = %v, want %v", got, result.OfError[float64](errFail))
	}
	if got := result.OfOk(1.0).OfError(); got != result.OfError[float64](nil) {
		t.Errorf("OfError() = %v, want %v", got, result.OfError[float64](nil))
	}
}

func TestResult_Unpack(t *testing.T) {
	if v, err := result.OfError[string](errFail).Unpack(); v != "" || err != errFail {
		t.Errorf("Unpack() = %v %v, want %v %v", v, err, "", errFail)
	}
	if v, err := result.OfOk("foo").Unpack(); v != "foo" || err != nil {
		t.Errorf("Unpack() = %v %v, want %v %v", v, err, "foo", nil)
	}
}
