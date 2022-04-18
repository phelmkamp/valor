package enum

import (
	"fmt"
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/tuple/two"
)

// Enum is an enumerated type.
//
// It wraps an optional.Value that is only ok if it's a member of the allowed values.
type Enum[T comparable] struct {
	optional.Value[T]
	members map[T]member // carry allowed values for validation
}

type member struct {
	i    int // used to maintain order
	name string
}

// Of creates an Enum of the given values.
func Of[T comparable](vals ...T) Enum[T] {
	e := Enum[T]{members: make(map[T]member)}
	for i, v := range vals {
		e.members[v] = member{i: i, name: fmt.Sprint(v)}
	}
	return e
}

// OfNamed creates an Enum of the given named values.
func OfNamed[T comparable](vals ...two.Tuple[string, T]) Enum[T] {
	e := Enum[T]{members: make(map[T]member)}
	for i, m := range vals {
		e.members[m.V2] = member{i: i, name: m.V}
	}
	return e
}

// ValueOf returns an Enum that wraps v if v is a member of the allowed values.
// Returns a not-ok Enum otherwise.
func (e Enum[T]) ValueOf(v T) Enum[T] {
	_, ok := e.members[v]
	e.Value = optional.Of(v, ok)
	return e
}

// String returns e formatted as a string.
func (e Enum[T]) String() string {
	val := optional.Map(e.Value, func(v T) string {
		return e.members[v].name
	})
	return fmt.Sprint(val)
}

// Members returns the allowed values.
func (e Enum[T]) Members() []T {
	s := make([]T, len(e.members))
	for v, m := range e.members {
		s[m.i] = v
	}
	return s
}
