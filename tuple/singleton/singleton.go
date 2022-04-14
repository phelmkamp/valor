// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package singleton provides a singleton set.
package singleton

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/two"
	"github.com/phelmkamp/valor/tuple/unit"
)

// Set contains at most one element.
type Set[E any] map[unit.Type]E

// String returns the Set formatted as a string.
func (s Set[E]) String() string {
	v, ok := s[unit.Unit]
	if !ok {
		return "{}"
	}
	return fmt.Sprintf("{%v}", v)
}

// SetOf creates a Set of v.
func SetOf[E any](v E) Set[E] {
	return Set[E]{unit.Unit: v}
}

func SetZip[T, T2 any](s Set[T], s2 Set[T2]) two.Tuple[T, T2] {
	return two.TupleOf(s[unit.Unit], s2[unit.Unit])
}

func SetUnzip[T, T2 any](t two.Tuple[T, T2]) (Set[T], Set[T2]) {
	return SetOf(t.V), SetOf(t.V2)
}
