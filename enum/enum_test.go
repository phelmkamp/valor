package enum_test

import (
	"fmt"
	"github.com/phelmkamp/valor/enum"
	"github.com/phelmkamp/valor/tuple/two"
	"reflect"
	"testing"
	"time"
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
	Suit = enum.OfString(Clubs, Diamonds, Hearts, Spades)

	Fruit = enum.Of(
		two.TupleOf("apple", Apple),
		two.TupleOf("banana", Banana),
		two.TupleOf("orange", Orange),
	)

	Event = enum.OfText(
		time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	)
)

func Example() {
	fmt.Println(Suit.Values())
	fmt.Println(Suit.ValueOf("Foo"))
	fmt.Println(Suit.ValueOf(Hearts))
	fmt.Println(Fruit.ValueOf(Orange))
	fmt.Println(Event.ValueOf(time.UnixMilli(0).UTC()))
	// Output:
	// [clubs diamonds hearts spades]
	// { false}
	// {hearts true}
	// {orange true}
	// {1970-01-01T00:00:00Z true}
}

// Example_marshal demonstrates that an enum.Enum can be marshaled to and unmarshaled from text (and therefore JSON).
func Example_marshal() {
	var text []byte
	text, _ = Fruit.ValueOf(Apple).MarshalText()
	fav := Fruit
	_ = fav.UnmarshalText(text)
	fmt.Println(fav)
	// Output: {apple true}
}

func Test_marshal(t *testing.T) {
	fav := Fruit
	text, err := Fruit.ValueOf(-1).MarshalText()
	if err != nil || len(text) != 0 {
		t.Errorf("MarshalText() = %s %v, want %s %v", text, err, "", nil)
	}
	if err = fav.UnmarshalText(text); err != nil {
		t.Errorf("UnmarshalText() = %v, want %v", err, nil)
	}
	if !reflect.DeepEqual(fav, Fruit) {
		t.Errorf("fav after UnmarshalText() = %v, want %v", fav, Fruit)
	}

	text, err = Fruit.ValueOf(Apple).MarshalText()
	if err != nil || string(text) != "apple" {
		t.Errorf("MarshalText() = %s %v, want %s %v", text, err, "apple", nil)
	}
	if err = fav.UnmarshalText(text); err != nil {
		t.Errorf("UnmarshalText() = %v, want %v", err, nil)
	}
	if !reflect.DeepEqual(fav, Fruit.ValueOf(Apple)) {
		t.Errorf("fav after UnmarshalText() = %v, want %v", fav, Fruit.ValueOf(Apple))
	}
}
