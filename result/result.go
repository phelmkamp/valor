// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package result

import (
	"errors"
	"fmt"
	"github.com/phelmkamp/valor/optional"
)

var (
	// errNil differentiates OfError(nil) from OfOk.
	errNil = errors.New("errNil")
)

// Result contains either a value or an error.
type Result[T any] struct {
	v   T
	err error
}

// Empty is an empty value.
// It's an alias for struct{} and is equivalent to struct{} in all ways.
// This is useful when a Result can only contain an error or nothing:
//	result.OfError[result.Empty](err)
// Deprecated: use unit.Type instead.
type Empty = struct{}

// Of creates a Result of either v or err.
// This aids interoperability with function return values.
func Of[T any](v T, err error) Result[T] {
	if err != nil {
		return OfError[T](err)
	}
	return OfOk(v)
}

// OfOk creates a Result of v.
func OfOk[T any](v T) Result[T] {
	return Result[T]{v: v}
}

// OfError creates a Result of err.
func OfError[T any](err error) Result[T] {
	if err == nil {
		err = errNil
	}
	return Result[T]{err: err}
}

// OfValue creates a Result of the underlying value of val if ok or err otherwise.
// This aids in converting a optional.Value to a Result.
func OfValue[T any](val optional.Value[T], err error) Result[T] {
	var res Result[T]
	if !val.Ok(&res.v) {
		return OfError[T](err)
	}
	return res
}

// IsError returns whether r contains an error.
func (res Result[T]) IsError() bool {
	// Call Error() to handle errNil properly.
	return res.Error() != nil
}

// String returns res formatted as a string.
func (res Result[T]) String() string {
	return fmt.Sprintf("{%v %v}", res.v, res.err)
}

// Value returns a optional.Value containing either the
// underlying value or nothing.
func (res Result[T]) Value() optional.Value[T] {
	// Special case: errNil should be "not ok" so don't call IsError().
	if res.err != nil {
		return optional.OfNotOk[T]()
	}
	return optional.OfOk(res.v)
}

// Error returns the contained error, or nil if no error is present.
func (res Result[T]) Error() error {
	if res.err == errNil {
		return nil
	}
	return res.err
}

// Errorf returns a Result where the contained error has
// been formatted with format.
// Does nothing if res does not contain an error.
func (res Result[T]) Errorf(format string) Result[T] {
	if res.IsError() {
		res.err = fmt.Errorf(format, res.err)
	}
	return res
}

// ErrorAs calls errors.As with the underlying error.
// Does nothing if res does not contain an error.
func (res Result[T]) ErrorAs(target any) bool {
	if !res.IsError() {
		return false
	}
	return errors.As(res.err, target)
}

// ErrorIs calls errors.Is with the underlying error.
// Does nothing if res does not contain an error.
func (res Result[T]) ErrorIs(target error) bool {
	if !res.IsError() {
		return false
	}
	return errors.Is(res.err, target)
}

// ErrorUnwrap calls errors.Unwrap with the underlying error.
// Does nothing if res does not contain an error that wraps another error.
func (res Result[T]) ErrorUnwrap() Result[T] {
	if u := errors.Unwrap(res.err); u != nil {
		res.err = u
	}
	return res
}

// OfOk creates a Result of the underlying value, dropping any error.
// This aids in comparisons, enabling the use of res in switch statements.
func (res Result[T]) OfOk() Result[T] {
	return OfOk(res.v)
}

// OfError creates a Result of the underlying error, dropping any value.
// This aids in comparisons, enabling the use of res in switch statements.
func (res Result[T]) OfError() Result[T] {
	return OfError[T](res.err)
}

// Transpose converts res to a optional.Value of Result.
// Returns a not-ok optional.Value if the underlying optional.Value is not ok.
// Otherwise, returns an ok optional.Value of a Result that contains the underlying value or error.
func Transpose[T any](res Result[optional.Value[T]]) optional.Value[Result[T]] {
	if res.IsError() {
		return optional.OfOk(OfError[T](res.err))
	}
	if res.v.IsOk() {
		return optional.OfOk(OfOk[T](res.v.MustOk()))
	}
	return optional.OfNotOk[Result[T]]()
}

// TransposeValue converts val to a Result of optional.Value.
// Returns an error Result if the underlying value is an error.
// Otherwise, returns a Result of a optional.Value that contains the underlying value if ok or nothing if not ok.
func TransposeValue[T any](val optional.Value[Result[T]]) Result[optional.Value[T]] {
	var v Result[T]
	if !val.Ok(&v) {
		return OfOk(optional.OfNotOk[T]())
	}
	if v.IsError() {
		return OfError[optional.Value[T]](v.err)
	}
	return OfOk(optional.OfOk(v.v))
}
