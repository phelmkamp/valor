package enum

import (
	"fmt"
	"github.com/phelmkamp/valor/optional"
)

// Enum is an enumerated type.
//
// It embeds an optional.Value of the currently selected member.
type Enum[T comparable] struct {
	optional.Value[Member[T]]
	keys map[T]int
}

// Member is a value and its associated index.
type Member[T comparable] struct {
	Index int
	Value T
}

// Of creates an Enum of the allowed values.
func Of[T comparable](allowed ...T) Enum[T] {
	e := Enum[T]{keys: make(map[T]int)}
	for i, v := range allowed {
		e.keys[v] = i
	}
	return e
}

// Select returns a copy of e with v as the currently selected member.
// The optional.Value will be ok if v is allowed, not-ok otherwise.
func (e Enum[T]) Select(v T) Enum[T] {
	i, ok := e.keys[v]
	e.Value = optional.Of(Member[T]{Index: i, Value: v}, ok)
	return e
}

// String returns e formatted as a string.
func (e Enum[T]) String() string {
	return fmt.Sprint(e.Value)
}

// Values returns the allowed values.
func (e Enum[T]) Values() []T {
	s := make([]T, len(e.keys))
	for v, i := range e.keys {
		s[i] = v
	}
	return s
}
