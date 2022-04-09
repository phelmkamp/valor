package three_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/three"
	"github.com/phelmkamp/valor/value"
)

func get() (string, int, float32, bool) {
	return "a", 1, 1.0, true
}

func Example() {
	val := value.Of(three.TupleOf(get()))
	fmt.Println(val)
	// Output: {{a 1 1} true}
}
