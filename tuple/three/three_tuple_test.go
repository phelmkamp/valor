// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package three_test

import (
	"fmt"
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
