// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package three_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/phelmkamp/valor/funcs"
	"github.com/phelmkamp/valor/tuple/three"
)

func get() (string, int, float32, bool) {
	return "a", 1, 1.0, true
}

func Example() {
	val := three.TupleValueOf(get())
	fmt.Println(val)
	// Output: {{a 1 1} true}
}

func TestTuple_Values(t *testing.T) {
	tup := three.TupleOf(1, "two", 3.0)
	v, v2, v3 := tup.Values()
	if v != 1 || v2 != "two" || v3 != 3.0 {
		t.Errorf("Values() = %v %v %v, want %v %v %v", v, v2, v3, 1, "two", 3.0)
	}
}

func Test_TupleMap(t *testing.T) {
	tup := three.TupleOf(1, "two", 3.0)
	got := three.TupleMap(tup, strconv.Itoa, funcs.Ident[string], funcs.Ident[float64])
	if got != three.TupleOf("1", "two", 3.0) {
		t.Errorf("TupleMap() = %v, want %v", got, three.TupleOf("1", "two", 3.0))
	}
}
