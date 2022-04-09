// Package four provides a 4-tuple type.
package four

// Tuple contains four values.
type Tuple[T, T2, T3, T4 any] struct {
	V  T
	V2 T2
	V3 T3
	V4 T4
}

// TupleOf creates a Tuple of v, v2, v3, and v4.
// The optional ok argument aids interoperability with
// return values that follow the "comma ok" idiom.
func TupleOf[T, T2, T3, T4 any](v T, v2 T2, v3 T3, v4 T4, ok ...bool) (Tuple[T, T2, T3, T4], bool) {
	if len(ok) == 0 {
		ok = []bool{true}
	}
	return Tuple[T, T2, T3, T4]{V: v, V2: v2, V3: v3, V4: v4}, ok[0]
}
