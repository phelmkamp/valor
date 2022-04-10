// Package three provides a 3-tuple type.
package three

import (
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/value"
)

// Tuple contains three values.
type Tuple[T, T2, T3 any] struct {
	V  T
	V2 T2
	V3 T3
}

// TupleOf creates a Tuple of (v, v2, v3).
func TupleOf[T, T2, T3 any](v T, v2 T2, v3 T3) Tuple[T, T2, T3] {
	return Tuple[T, T2, T3]{V: v, V2: v2, V3: v3}
}

// TupleValueOf creates a value.Value of (v, v2, v3) if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func TupleValueOf[T, T2, T3 any](v T, v2 T2, v3 T3, ok bool) value.Value[Tuple[T, T2, T3]] {
	return value.Of(TupleOf(v, v2, v3), ok)
}

// TupleResultOf creates a result.Result of either (v, v2, v3) or err.
// This aids interoperability with function return values.
func TupleResultOf[T, T2, T3 any](v T, v2 T2, v3 T3, err error) result.Result[Tuple[T, T2, T3]] {
	return result.Of(TupleOf(v, v2, v3), err)
}
