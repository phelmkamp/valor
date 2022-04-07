package result

import (
	"errors"
	"fmt"
	"github.com/phelmkamp/valor/value"
)

// Result contains either a value or an error.
type Result[T any] struct {
	val T
	err error
}

// Of creates a Result of either val or err.
// This aids interoperability with function return values
func Of[T any](val T, err error) Result[T] {
	if err != nil {
		return OfError[T](err)
	}
	return OfValue(val)
}

// OfValue creates a Result of val.
func OfValue[T any](val T) Result[T] {
	return Result[T]{val: val}
}

// OfError creates a Result of err.
func OfError[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// IsError returns whether r contains an error.
func (r Result[T]) IsError() bool {
	return r.err != nil
}

// String returns the underlying value or error
// formatted as a string.
func (r Result[T]) String() string {
	if r.IsError() {
		return fmt.Sprint(r.err)
	}
	return fmt.Sprint(r.val)
}

// Value returns an value.Value containing the
// contained value, or an empty value.Value if no value is present.
func (r Result[T]) Value() value.Value[T] {
	if r.IsError() {
		return value.Value[T]{}
	}
	return value.OfOk(r.val)
}

// Error returns the contained error, or nil if no error is present.
func (r Result[T]) Error() error {
	return r.err
}

// Errorf returns a Result where the contained error has
// been formatted with format.
// Does nothing if r does not contain an error.
func (r Result[T]) Errorf(format string) Result[T] {
	if r.IsError() {
		r.err = fmt.Errorf(format, r.err)
	}
	return r
}

// ErrorAs calls errors.As with the underlying error.
// Does nothing if r does not contain an error.
func (r Result[T]) ErrorAs(target any) bool {
	if !r.IsError() {
		return false
	}
	return errors.As(r.err, target)
}

// ErrorIs calls errors.Is with the underlying error.
// Does nothing if r does not contain an error.
func (r Result[T]) ErrorIs(target error) bool {
	if !r.IsError() {
		return false
	}
	return errors.Is(r.err, target)
}

// ErrorUnwrap calls errors.Unwrap with the underlying error.
// Does nothing if r does not contain an error.
func (r Result[T]) ErrorUnwrap() Result[T] {
	if !r.IsError() {
		return r
	}
	r.err = errors.Unwrap(r.err)
	return r
}
