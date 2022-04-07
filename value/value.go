package value

// Value either contains a value (ok) or nothing (not ok).
type Value[T any] struct {
	v  T
	ok bool
}

// Of creates a Value of v if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func Of[T any](v T, ok bool) Value[T] {
	if ok {
		return OfOk(v)
	}
	return Value[T]{}
}

// OfOk creates an ok Value of v.
func OfOk[T any](v T) Value[T] {
	return Value[T]{v: v, ok: true}
}

// OfNotOk creates a Value that is not ok.
// This aids in comparisons, enabling the use of Value in switch statements.
func OfNotOk[T any]() Value[T] {
	return Value[T]{}
}

// IsOk returns whether v contains a value.
func (val Value[T]) IsOk() bool {
	return val.ok
}

// Ok sets dst to the underlying value if ok.
// Returns true if ok, false if not ok.
func (val Value[T]) Ok(dst *T) bool {
	if !val.IsOk() {
		return false
	}
	*dst = val.v
	return true
}

// Or returns the underlying value if ok, or def if not ok.
func (val Value[T]) Or(def T) T {
	if val.IsOk() {
		return val.v
	}
	return def
}

// OrZero returns the underlying value if ok, or the zero value if not ok.
func (val Value[T]) OrZero() T {
	return val.v
}

// OrElse returns the underlying value if ok, or the result of f if not ok.
func (val Value[T]) OrElse(f func() T) T {
	if val.IsOk() {
		return val.v
	}
	return f()
}

// OfOk creates an ok Value of the underlying value.
// This aids in comparisons, enabling the use of val in switch statements.
func (val Value[T]) OfOk() Value[T] {
	return OfOk(val.v)
}

// Inspect calls f with the underlying value if ok.
// Does nothing if not ok.
func (val Value[T]) Inspect(f func(T)) Value[T] {
	if val.IsOk() {
		f(val.v)
	}
	return val
}

// Map returns the result of f on the underlying value of v.
func Map[T, T2 any](val Value[T], f func(T) T2) Value[T2] {
	if !val.IsOk() {
		return Value[T2]{}
	}
	return OfOk(f(val.v))
}
