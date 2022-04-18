package enum_test

import (
	"fmt"
	"github.com/phelmkamp/valor/enum"
)

var suits = enum.Of("clubs", "diamonds", "hearts", "spades")

func Example() {
	fmt.Println(suits.Values())
	mySuit := suits.Select("foo")
	fmt.Println(mySuit)
	mySuit = suits.Select("hearts")
	fmt.Println(mySuit)
	// Output:
	// [clubs diamonds hearts spades]
	// {{0 } false}
	// {{2 hearts} true}
}
