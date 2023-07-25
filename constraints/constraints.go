// Package constraints defines a set of useful constraints to be used with type parameters.
package constraints

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
