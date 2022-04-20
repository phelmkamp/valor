// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package four_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/four"
	"testing"
)

func get() (string, int, float32, []int, bool) {
	return "a", 1, 1.0, []int{1}, true
}

func Example() {
	val := four.TupleValueOf(get())
	fmt.Println(val)
	// Output: {{a 1 1 [1]} true}
}

func TestTuple_Values(t *testing.T) {
	tup := four.TupleOf(1, "two", 3.0, []int{4})
	v, v2, v3, v4 := tup.Values()
	if v != 1 || v2 != "two" || v3 != 3.0 || v4[0] != 4 {
		t.Errorf("Values() = %v %v %v %v, want %v %v %v %v", v, v2, v3, v4, 1, "two", 3.0, []int{4})
	}
}
