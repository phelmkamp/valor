package value_test

import (
	"fmt"
	"github.com/phelmkamp/valor/value"
	"strconv"
)

func Example() {
	m := map[string]int{"foo": 42}
	foo, ok := m["foo"]
	val := value.Of(foo, ok)
	fmt.Println(val.IsOk()) // true

	var foo2 int
	fmt.Println(val.Ok(&foo2), foo2) // true 42

	val2 := value.Map(val, strconv.Itoa)
	fmt.Println(val2) // {42 true}

	bar, ok := m["bar"]
	val3 := value.Of(bar, ok)
	fmt.Println(val3.Or(-1))                          // -1
	fmt.Println(val3.OrZero())                        // 0
	fmt.Println(val3.OrElse(func() int { return 1 })) // 1

	// switch
	var foo3 int
	switch val {
	case val.OfOk():
		val.Ok(&foo3)
		fmt.Println("Ok")
	case value.OfNotOk[int]():
		fmt.Println("None")
		return
	}
	// Output: true
	// true 42
	// {42 true}
	// -1
	// 0
	// 1
	// Ok
}
