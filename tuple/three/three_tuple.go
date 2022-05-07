// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package three provides a 3-tuple type.
package three

import (
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/result"
)

// Tuple contains three values.
type Tuple[T, T2, T3 any] struct {
	V  T
	V2 T2
	V3 T3
}

// Values returns the contained values.
// This aids in assigning to variables or function arguments.
func (t Tuple[T, T2, T3]) Values() (v T, v2 T2, v3 T3) {
	return t.V, t.V2, t.V3
}

// TupleOf creates a Tuple of (v, v2, v3).
func TupleOf[T, T2, T3 any](v T, v2 T2, v3 T3) Tuple[T, T2, T3] {
	return Tuple[T, T2, T3]{V: v, V2: v2, V3: v3}
}

// TupleValueOf creates an optional.Value of (v, v2, v3) if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func TupleValueOf[T, T2, T3 any](v T, v2 T2, v3 T3, ok bool) optional.Value[Tuple[T, T2, T3]] {
	return optional.Of(TupleOf(v, v2, v3), ok)
}

// TupleResultOf creates a result.Result of either (v, v2, v3) or err.
// This aids interoperability with function return values.
func TupleResultOf[T, T2, T3 any](v T, v2 T2, v3 T3, err error) result.Result[Tuple[T, T2, T3]] {
	return result.Of(TupleOf(v, v2, v3), err)
}

// TupleMap returns a Tuple with each value replaced by the result of each function.
//
// funcs.Ident can be used to leave the value unchanged.
func TupleMap[T, T2, T3, Tp, T2p, T3p any](t Tuple[T, T2, T3], f func(T) Tp, f2 func(T2) T2p, f3 func(T3) T3p) Tuple[Tp, T2p, T3p] {
	return TupleOf(f(t.V), f2(t.V2), f3(t.V3))
}
