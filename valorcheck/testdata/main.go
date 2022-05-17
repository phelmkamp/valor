// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"github.com/phelmkamp/valor/enum"
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/result"
)

var (
	Fruit = enum.OfString("apple", "banana", "orange")
)

func main() {
	var w strings.Builder
	m := map[string]int{"foo": 1}
	val := optional.OfIndex(m, "foo")
	res := result.Of(w.Write(nil))
	f := Fruit.ValueOf("apple")

	optional.OfIndex(m, "foo").MustOk()      // UNGUARDED
	i := val.MustOk()                        // UNGUARDED
	result.Of(w.Write(nil)).Value().MustOk() // UNGUARDED
	n := res.Value().MustOk()                // UNGUARDED
	Fruit.ValueOf("apple").MustOk()          // UNGUARDED
	s := f.ValueOf("apple").MustOk()         // UNGUARDED
	val.Ok(&i)                               // UNCHECKED
	res.Value().Ok(&n)                       // UNCHECKED
	f.Ok(&s)                                 // UNCHECKED

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

	switch val := optional.Of(0, true); val {
	case val.OfOk():
		i = val.MustOk()
	case optional.OfNotOk[int]():
	}
}
