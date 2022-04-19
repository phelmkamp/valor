package enum

import (
	"encoding"
	"fmt"
	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/tuple/two"
)

// Enum is an enumerated type.
//
// It wraps an optional.Value that is only ok if it's a member of the allowed values.
type Enum[T comparable] struct {
	optional.Value[T]
	members map[T]metadata // carry allowed values for validation
}

type metadata struct {
	i    int
	name string
}

// ComparableText is a constraint that permits
// comparable types that can be marshaled to text.
type ComparableText interface {
	comparable
	encoding.TextMarshaler
}

// Of creates an Enum of the given name-value pairs.
func Of[T comparable](pairs ...two.Tuple[string, T]) Enum[T] {
	e := Enum[T]{members: make(map[T]metadata)}
	for i, m := range pairs {
		e.members[m.V2] = metadata{i: i, name: m.V}
	}
	return e
}

// OfString creates an Enum of the given string values.
func OfString(vals ...string) Enum[string] {
	e := Enum[string]{members: make(map[string]metadata)}
	for i, v := range vals {
		e.members[v] = metadata{i: i, name: v}
	}
	return e
}

// OfText creates an Enum of the given text values.
// Panics if MarshalText returns an error for one of the values.
func OfText[T ComparableText](vals ...T) Enum[T] {
	e := Enum[T]{members: make(map[T]metadata)}
	for i, v := range vals {
		text, err := v.MarshalText()
		if err != nil {
			panic(err)
		}
		e.members[v] = metadata{i: i, name: string(text)}
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

// Values returns the allowed values.
func (e Enum[T]) Values() []T {
	s := make([]T, len(e.members))
	for v, m := range e.members {
		s[m.i] = v
	}
	return s
}

// MarshalText returns the name of the current member.
// Returns nil if e is not ok.
func (e Enum[T]) MarshalText() (text []byte, err error) {
	var v T
	if !e.Ok(&v) {
		return nil, nil
	}
	return []byte(e.members[v].name), nil
}

// UnmarshalText sets e to wrap the member with the given name.
// Sets e to not-ok if text is not the name of a valid member.
func (e *Enum[T]) UnmarshalText(text []byte) error {
	s := string(text)
	for v, m := range e.members {
		if m.name == s {
			e.Value = optional.OfOk(v)
			return nil
		}
	}
	e.Value = optional.OfNotOk[T]()
	return nil
}
