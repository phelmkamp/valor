package optional_test

import (
	"fmt"

	"github.com/phelmkamp/valor/optional"
)

type Map[K comparable, V any] struct {
	m map[K]V
}

func (m Map[K, V]) Index(k K) optional.Value[V] {
	return optional.OfIndex(m.m, k)
}

// Example_map demonstrates that a custom map type can be
// implemented that forces consumers to handle a missing pair.
func Example_map() {
	m := Map[string, int]{m: map[string]int{"foo": 0}}
	fmt.Println(m.Index("foo"))
	fmt.Println(m.Index("bar"))
	// Output:
	// {0 true}
	// {0 false}
}

type Chan[T any] struct {
	ch chan T
}

func (ch Chan[T]) Receive() optional.Value[T] {
	return optional.OfReceive(ch.ch)
}

func (ch Chan[T]) Close() {
	close(ch.ch)
}

// Example_channel demonstrates that a custom channel type can be
// implemented that forces consumers to handle a closed channel.
func Example_channel() {
	ch := Chan[int]{ch: make(chan int, 1)}
	ch.ch <- 0
	fmt.Println(ch.Receive())
	ch.Close()
	fmt.Println(ch.Receive())
	// Output:
	// {0 true}
	// {0 false}
}
