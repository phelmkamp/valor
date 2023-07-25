// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package two provides a 2-tuple type.
package two

import (
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/tuple/four"
)

// Tuple contains two values.
type Tuple[T, T2 any] struct {
	V  T
	V2 T2
}

// Values returns the contained values.
// This aids in assigning to variables or function arguments.
func (t Tuple[T, T2]) Values() (v T, v2 T2) {
	return t.V, t.V2
}

// TupleOf creates a Tuple of (v, v2).
func TupleOf[T, T2 any](v T, v2 T2) Tuple[T, T2] {
	return Tuple[T, T2]{V: v, V2: v2}
}

// TupleValueOf creates an optional.Value of (v, v2) if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func TupleValueOf[T, T2 any](v T, v2 T2, ok bool) optional.Value[Tuple[T, T2]] {
	return optional.Of(TupleOf(v, v2), ok)
}

// TupleResultOf creates a result.Result of either (v, v2) or err.
// This aids interoperability with function return values.
func TupleResultOf[T, T2 any](v T, v2 T2, err error) result.Result[Tuple[T, T2]] {
	return result.Of(TupleOf(v, v2), err)
}

// TupleMap returns a Tuple with each value replaced by the result of each function.
//
// funcs.Ident can be used to leave the value unchanged.
func TupleMap[T, T2, Tp, T2p any](t Tuple[T, T2], f func(T) Tp, f2 func(T2) T2p) (tp Tuple[Tp, T2p]) {
	return TupleOf(f(t.V), f2(t.V2))
}

// TupleZip combines the values of t and t2 into a four.Tuple.
func TupleZip[T, T2, T3, T4 any](t Tuple[T, T3], t2 Tuple[T2, T4]) four.Tuple[T, T2, T3, T4] {
	return four.TupleOf(t.V, t2.V, t.V2, t2.V2)
}

// TupleUnzip separates the values of t into two Tuples.
func TupleUnzip[T, T2, T3, T4 any](t four.Tuple[T, T2, T3, T4]) (Tuple[T, T3], Tuple[T2, T4]) {
	return TupleOf(t.V, t.V3), TupleOf(t.V2, t.V4)
}

// TupleIter converts an iterator of two values into an iterator of Tuples.
func TupleIter[T, T2 any](iter func(yield func(T, T2) bool)) optional.Iter[Tuple[T, T2]] {
	return func(yield func(t Tuple[T, T2]) bool) {
		iter(func(v T, v2 T2) bool {
			return yield(TupleOf(v, v2))
		})
	}
}
