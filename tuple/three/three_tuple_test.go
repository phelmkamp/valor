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
