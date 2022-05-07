// Package funcs provides common functions that
// can be used as arguments to map, zip, etc.
package funcs

// Ident is an identity function that returns v.
// This can be used as a map function that doesn't change the type.
func Ident[T any](v T) T {
	return v
}
