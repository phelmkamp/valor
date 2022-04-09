package four_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/four"
	"github.com/phelmkamp/valor/value"
)

func get() (string, int, float32, []int, bool) {
	return "a", 1, 1.0, []int{1}, true
}

func Example() {
	val := value.Of(four.TupleOf(get()))
	fmt.Println(val)
	// Output: {{a 1 1 [1]} true}
}
