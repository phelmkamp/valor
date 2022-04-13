package singleton_test

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/singleton"
)

func Example() {
	fmt.Println(singleton.SetOf(42))
	// Output: {42}
}
