package two_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/two"
	"github.com/phelmkamp/valor/value"
)

func get() (string, int, bool) {
	return "a", 1, true
}

func Example() {
	val := value.Of(two.TupleOf(get()))
	fmt.Println(val)
	// Output: {{a 1} true}
}
