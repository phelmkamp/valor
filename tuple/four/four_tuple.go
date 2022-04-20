// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package four provides a 4-tuple type.
package four

import (
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/result"
)

// Tuple contains four values.
type Tuple[T, T2, T3, T4 any] struct {
	V  T
	V2 T2
	V3 T3
	V4 T4
}

// Values returns the contained values.
// This aids in assigning to variables or function arguments.
func (t Tuple[T, T2, T3, T4]) Values() (v T, v2 T2, v3 T3, v4 T4) {
	return t.V, t.V2, t.V3, t.V4
}

// TupleOf creates a Tuple of (v, v2, v3, v4).
func TupleOf[T, T2, T3, T4 any](v T, v2 T2, v3 T3, v4 T4) Tuple[T, T2, T3, T4] {
	return Tuple[T, T2, T3, T4]{V: v, V2: v2, V3: v3, V4: v4}
}

// TupleValueOf creates an optional.Value of (v, v2, v3, v4) if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func TupleValueOf[T, T2, T3, T4 any](v T, v2 T2, v3 T3, v4 T4, ok bool) optional.Value[Tuple[T, T2, T3, T4]] {
	return optional.Of(TupleOf(v, v2, v3, v4), ok)
}

// TupleResultOf creates a result.Result of either (v, v2, v3, v4) or err.
// This aids interoperability with function return values.
func TupleResultOf[T, T2, T3, T4 any](v T, v2 T2, v3 T3, v4 T4, err error) result.Result[Tuple[T, T2, T3, T4]] {
	return result.Of(TupleOf(v, v2, v3, v4), err)
}
