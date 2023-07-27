// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package value

import (
	"fmt"

	"github.com/phelmkamp/valor/optional"
)

// Value either contains a value (ok) or nothing (not ok).
//
// Deprecated: use optional.Value instead.
type Value[T any] struct {
	optional.Value[T]
}

func (val Value[T]) String() string {
	return fmt.Sprint(val.Value)
}

// Of creates a Value of v if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func Of[T any](v T, ok bool) Value[T] {
	if ok {
		return OfOk(v)
	}
	return OfNotOk[T]()
}

// OfOk creates an ok Value of v.
func OfOk[T any](v T) Value[T] {
	return Value[T]{optional.OfOk(v)}
}

// OfNotOk creates a Value that is not ok.
// This aids in comparisons, enabling the use of Value in switch statements.
func OfNotOk[T any]() Value[T] {
	return Value[T]{}
}

// OfOk creates an ok Value of the underlying value.
// This aids in comparisons, enabling the use of val in switch statements.
func (val Value[T]) OfOk() Value[T] {
	return Value[T]{val.Value.OfOk()}
}

// Do calls f with the underlying value if ok.
// Does nothing if not ok.
func (val Value[T]) Do(f func(T)) Value[T] {
	return Value[T]{val.Value.Do(f)}
}

// Filter returns val if f returns true for the underlying value.
// Otherwise returns a not-ok Value.
func (val Value[T]) Filter(f func(T) bool) Value[T] {
	return Value[T]{val.Value.Filter(f)}
}

// Map returns a Value of the result of f on the underlying value.
// Returns a not-ok Value if val is not ok.
func Map[T, T2 any](val Value[T], f func(T) T2) Value[T2] {
	return Value[T2]{optional.Map(val.Value, f)}
}

// FlatMap returns the result of f on the underlying value.
// Returns a not-ok Value if val is not ok.
func FlatMap[T, T2 any](val Value[T], f func(T) Value[T2]) Value[T2] {
	optF := func(v T) optional.Value[T2] { return f(v).Value }
	return Value[T2]{optional.FlatMap(val.Value, optF)}
}

// Contains returns whether the underlying value equals v.
// Returns false if val is not ok.
func Contains[T comparable](val Value[T], v T) bool {
	return optional.Contains(val.Value, v)
}

// ZipWith calls f with the underlying values of val and val2 and returns a Value of the result.
// Returns a not-ok Value if either val or val2 is not ok.
func ZipWith[T, T2, T3 any](val Value[T], val2 Value[T2], f func(T, T2) T3) Value[T3] {
	return Value[T3]{optional.ZipWith(val.Value, val2.Value, f)}
}

// UnzipWith calls f with the underlying value of val and returns Values of the result.
// Does nothing and returns not-ok Values if val is not ok.
func UnzipWith[T, T2, T3 any](val Value[T], f func(T) (T2, T3)) (val2 Value[T2], val3 Value[T3]) {
	optVal2, optVal3 := optional.UnzipWith(val.Value, f)
	return Value[T2]{optVal2}, Value[T3]{optVal3}
}

// Flatten returns the underlying Value of val.
// Returns a not-ok Value if val is not ok.
func Flatten[T any](val Value[Value[T]]) Value[T] {
	if !val.IsOk() {
		return OfNotOk[T]()
	}
	return val.MustOk()
}
