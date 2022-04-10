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
