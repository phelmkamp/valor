// Package two provides a 2-tuple type.
package two

// Tuple contains two values.
type Tuple[T, T2 any] struct {
	V  T
	V2 T2
}

// TupleOf creates a Tuple of v and v2.
// The optional ok argument aids interoperability with
// return values that follow the "comma ok" idiom.
func TupleOf[T, T2 any](v T, v2 T2, ok ...bool) (Tuple[T, T2], bool) {
	if len(ok) == 0 {
		ok = []bool{true}
	}
	return Tuple[T, T2]{V: v, V2: v2}, ok[0]
}
