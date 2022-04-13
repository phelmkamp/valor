package singleton

import (
	"fmt"
	"github.com/phelmkamp/valor/tuple/unit"
)

// Set contains at most one element.
type Set[E any] map[unit.Type]E

// String returns the Set formatted as a string.
func (s Set[E]) String() string {
	v, ok := s[unit.Unit]
	if !ok {
		return "{}"
	}
	return fmt.Sprintf("{%v}", v)
}

// SetOf creates a Set of v.
func SetOf[E any](v E) Set[E] {
	return Set[E]{unit.Unit: v}
}
