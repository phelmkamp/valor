// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package two_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/two"
	"testing"
)

func get() (string, int, bool) {
	return "a", 1, true
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
