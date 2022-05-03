package optional_test

import (
	"encoding/json"
	"fmt"
	"log"

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

// Example_json demonstrates that a Value can be marshaled to and unmarshaled from JSON.
func Example_json() {
	type Obj struct {
		Name string              `json:"name"`
		Val  optional.Value[int] `json:"val"`
	}
	b, err := json.Marshal(Obj{Name: "foo", Val: optional.OfOk(0)})
	if err != nil {
		log.Fatalf("json.Marshal() failed: %v", err)
	}
	fmt.Println(string(b))
	var obj Obj
	err = json.Unmarshal(b, &obj)
	if err != nil {
		log.Fatalf("json.Unmarshal() failed: %v", err)
	}
	fmt.Println(obj)
	fmt.Println()

	b, err = json.Marshal(Obj{Name: "foo", Val: optional.OfNotOk[int]()})
	if err != nil {
		log.Fatalf("json.Marshal() failed: %v", err)
	}
	fmt.Println(string(b))
	obj = Obj{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		log.Fatalf("json.Unmarshal() failed: %v", err)
	}
	fmt.Println(obj)
	// Output:
	// {"name":"foo","val":0}
	// {foo {0 true}}
	//
	// {"name":"foo","val":null}
	// {foo {0 false}}
}
