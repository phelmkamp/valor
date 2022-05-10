// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/phelmkamp/valor/optional"
)

func main() {
	m := map[string]int{"foo": 1}
	val := optional.OfIndex(m, "foo")

	optional.OfIndex(m, "foo").MustOk() // UNGUARDED
	i := val.MustOk()                   // UNGUARDED
	val.Ok(&i)                          // UNCHECKED

	if val.IsOk() {
		i = val.MustOk()
	}

	if ok := val.Ok(&i); !ok {
		return
	}

	ok := val.Ok(&i)
	if !ok {
		return
	}

	if !val.Ok(&i) {
		return
	}

	if n := 0; val.Ok(&n) {
		n++
	}
}
