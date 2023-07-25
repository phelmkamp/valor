// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package two_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/phelmkamp/valor/funcs"
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/tuple/two"
)

func get() (string, int, bool) {
	return "a", 1, true
}

func kvs[K comparable, V any](m map[K]V) func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func Example() {
	val := two.TupleValueOf(get())
	fmt.Println(val)
	// Output: {{a 1} true}
}

func TestTuple_Values(t *testing.T) {
	tup := two.TupleOf(1, "two")
	v, v2 := tup.Values()
	if v != 1 || v2 != "two" {
		t.Errorf("Values() = %v %v, want %v %v", v, v2, 1, "two")
	}
}

func Test_TupleMap(t *testing.T) {
	tup := two.TupleOf(1, "two")
	got := two.TupleMap(tup, strconv.Itoa, funcs.Ident[string])
	if got != two.TupleOf("1", "two") {
		t.Errorf("TupleMap() = %v, want %v", got, two.TupleOf("1", "two"))
	}
}

func Test_TupleIter(t *testing.T) {
	var m map[int]string
	if got := optional.First(two.TupleIter(kvs(m))); got.IsOk() {
		t.Errorf("First() of empty map = %v, want %v", got, optional.OfNotOk[two.Tuple[int, string]]())
	}
	m = map[int]string{
		0: "a",
	}
	want := two.TupleValueOf(0, "a", true)
	if got := optional.First(two.TupleIter(kvs(m))); !optional.Equal(got, want) {
		t.Errorf("First() = %v, want %v", got, want)
	}
}
