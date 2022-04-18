package enum_test

import (
	"fmt"
	"github.com/phelmkamp/valor/enum"
	"github.com/phelmkamp/valor/tuple/two"
)

const (
	Clubs    = "clubs"
	Diamonds = "diamonds"
	Hearts   = "hearts"
	Spades   = "spades"

	Apple = iota
	Banana
	Orange
)

var (
	Suit = enum.Of(Clubs, Diamonds, Hearts, Spades)

	Fruit = enum.OfNamed(
		two.TupleOf("apple", Apple),
		two.TupleOf("banana", Banana),
		two.TupleOf("orange", Orange),
	)
)

func Example() {
	fmt.Println(Suit.Members())
	fmt.Println(Suit.ValueOf("Foo"))
	fmt.Println(Suit.ValueOf(Hearts))

	fmt.Println(Fruit.ValueOf(Banana))
	// Output:
	// [clubs diamonds hearts spades]
	// { false}
	// {hearts true}
	// {banana true}
}
