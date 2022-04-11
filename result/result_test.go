// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package result_test

import (
	"errors"
	"fmt"
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/value"
	"io"
	"io/fs"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Example() {
	var w strings.Builder
	// traditional
	if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
		fmt.Println(res.Error())
		return
	}

	// try to get value, printing wrapped error if not ok
	var n int
	if res := result.Of(w.Write([]byte("foo"))); !res.Value().Ok(&n) {
		fmt.Println(res.Errorf("Write() failed: %w").Error())
		return
	}
	fmt.Println(n)

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
	// fail
	// Ok
}

func mid(fail bool) (string, error) {
	var s string
	if res := result.Of(leaf(fail)); !value.Map(res.Value(), strconv.Itoa).Ok(&s) {
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

func TestOf(t *testing.T) {
	if got := result.Of(leaf(true)); got != result.OfError[int](errFail) {
		t.Errorf("Of() = %v, want %v", got, result.OfError[int](errFail))
	}
	if got := result.Of(leaf(false)); got != result.OfOk(1) {
		t.Errorf("Of() = %v, want %v", got, result.OfOk(1))
	}
}

func TestOfValue(t *testing.T) {
	if got := result.OfValue(value.OfNotOk[string](), errFail); got != result.OfError[string](errFail) {
		t.Errorf("OfValue() = %v, want %v", got, result.OfError[string](errFail))
	}
	if got := result.OfValue(value.OfOk("foo"), errFail); got != result.OfOk("foo") {
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
	if got := result.OfOk(1).Value(); got != value.OfOk(1) {
		t.Errorf("Value() = %v, want %v", got, value.OfOk(1))
	}
	if got := result.OfError[int](errFail).Value(); got != value.OfNotOk[int]() {
		t.Errorf("Value() = %v, want %v", got, value.OfNotOk[int]())
	}
}

func TestTranspose(t *testing.T) {
	if got := result.Transpose(result.OfError[value.Value[string]](errFail)); got != value.OfOk(result.OfError[string](errFail)) {
		t.Errorf("Transpose() = %v, want %v", got, value.OfOk(result.OfError[string](errFail)))
	}
	if got := result.Transpose(result.OfOk(value.OfNotOk[string]())); got != value.OfNotOk[result.Result[string]]() {
		t.Errorf("Transpose() = %v, want %v", got, value.OfNotOk[result.Result[string]]())
	}
	if got := result.Transpose(result.OfOk(value.OfOk("foo"))); got != value.OfOk(result.OfOk("foo")) {
		t.Errorf("Transpose() = %v, want %v", got, value.OfOk(result.OfOk("foo")))
	}
}

func TestTransposeValue(t *testing.T) {
	if got := result.TransposeValue(value.OfNotOk[result.Result[int]]()); got != result.OfOk(value.OfNotOk[int]()) {
		t.Errorf("TransposeValue() = %v, want %v", got, result.OfOk(value.OfNotOk[int]()))
	}
	if got := result.TransposeValue(value.OfOk(result.OfError[int](errFail))); got != result.OfError[value.Value[int]](errFail) {
		t.Errorf("TransposeValue() = %v, want %v", got, result.OfError[value.Value[int]](errFail))
	}
	if got := result.TransposeValue(value.OfOk(result.OfOk(1))); got != result.OfOk(value.OfOk(1)) {
		t.Errorf("TransposeValue() = %v, want %v", got, result.OfOk(value.OfOk(1)))
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
