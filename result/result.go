package result

import (
	"errors"
	"fmt"
	"github.com/phelmkamp/valor/value"
)

// Result contains either a value or an error.
type Result[T any] struct {
	v   T
	err error
}

// Of creates a Result of either v or err.
// This aids interoperability with function return values
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
	return Result[T]{err: err}
}

// OfValue creates a Result of val if ok or err otherwise.
// This aids in converting a value.Value to a Result.
func OfValue[T any](val value.Value[T], err error) Result[T] {
	res := Result[T]{}
	if !val.Ok(&res.v) {
		res.err = err
	}
	return res
}

// IsError returns whether r contains an error.
func (res Result[T]) IsError() bool {
	return res.err != nil
}

// String returns the underlying value or error
// formatted as a string.
func (res Result[T]) String() string {
	if res.IsError() {
		return fmt.Sprint(res.err)
	}
	return fmt.Sprint(res.v)
}

// Value returns a value.Value containing either the
// underlying value or nothing.
func (res Result[T]) Value() value.Value[T] {
	if res.IsError() {
		return value.Value[T]{}
	}
	return value.OfOk(res.v)
}

// Error returns the contained error, or nil if no error is present.
func (res Result[T]) Error() error {
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
// Does nothing if res does not contain an error.
func (res Result[T]) ErrorUnwrap() Result[T] {
	if !res.IsError() {
		return res
	}
	res.err = errors.Unwrap(res.err)
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
