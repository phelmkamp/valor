// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package two provides a 2-tuple type.
package two

import (
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/tuple/four"
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

func TupleZip[T, T2, T3, T4 any](t Tuple[T, T3], t2 Tuple[T2, T4]) four.Tuple[T, T2, T3, T4] {
	return four.TupleOf(t.V, t2.V, t.V2, t2.V2)
}

func TupleUnzip[T, T2, T3, T4 any](t four.Tuple[T, T2, T3, T4]) (Tuple[T, T3], Tuple[T2, T4]) {
	return TupleOf(t.V, t.V3), TupleOf(t.V2, t.V4)
}
