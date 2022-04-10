// Package two provides a 2-tuple type.
package two

import (
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/value"
)

// Tuple contains two values.
type Tuple[T, T2 any] struct {
	V  T
	V2 T2
}

// TupleOf creates a Tuple of (v, v2).
func TupleOf[T, T2 any](v T, v2 T2) Tuple[T, T2] {
	return Tuple[T, T2]{V: v, V2: v2}
}

// TupleValueOf creates a value.Value of (v, v2) if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func TupleValueOf[T, T2 any](v T, v2 T2, ok bool) value.Value[Tuple[T, T2]] {
	return value.Of(TupleOf(v, v2), ok)
}

// TupleResultOf creates a result.Result of either (v, v2) or err.
// This aids interoperability with function return values.
func TupleResultOf[T, T2 any](v T, v2 T2, err error) result.Result[Tuple[T, T2]] {
	return result.Of(TupleOf(v, v2), err)
}
