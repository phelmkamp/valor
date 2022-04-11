// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package four_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/four"
)

func get() (string, int, float32, []int, bool) {
	return "a", 1, 1.0, []int{1}, true
}

func Example() {
	val := four.TupleValueOf(get())
	fmt.Println(val)
	// Output: {{a 1 1 [1]} true}
}
