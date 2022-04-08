package value_test

import (
	"fmt"
	"github.com/phelmkamp/valor/value"
	"strconv"
	"testing"
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
	switch val {
	case val.OfOk():
		foo = val.MustOk()
		fmt.Println("Ok")
	case value.OfNotOk[int]():
		fmt.Println("Not Ok")
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

func TestValue_Take(t *testing.T) {
	val := value.OfOk(42)
	val2 := val.Take()
	if want2 := value.OfOk(42); val2 != want2 {
		t.Errorf("Take() = %v, want %v", val2, want2)
	}
	if want := value.OfNotOk[int](); val != want {
		t.Errorf("val after Take() = %v, want %v", val, want)
	}
}
