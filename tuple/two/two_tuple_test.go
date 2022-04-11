// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package two_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/two"
)

func get() (string, int, bool) {
	return "a", 1, true
}

func Example() {
	val := two.TupleValueOf(get())
	fmt.Println(val)
	// Output: {{a 1} true}
}
