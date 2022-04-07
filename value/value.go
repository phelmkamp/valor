package value

// Value either contains a value (ok) or nothing (not ok).
type Value[T any] struct {
	val T
	ok  bool
}

// Of creates a Value of val if ok is true.
// This aids interoperability with return values
// that follow the "comma ok" idiom.
func Of[T any](val T, ok bool) Value[T] {
	if ok {
		return OfOk(val)
	}
	return Value[T]{}
}

// OfOk creates a Value of val.
func OfOk[T any](val T) Value[T] {
	return Value[T]{val: val, ok: true}
}

//func OfNone[T any]() Value[T] {
//	return &Value[T]{ok: false}
//}

// IsOk returns whether v contains a value.
func (v Value[T]) IsOk() bool {
	return v.ok
}

//func (o Value[T]) Must() T {
//	if !v.IsOk() {
//		panic(`called Value.Value() on a not-ok value`)
//	}
//	return v.val
//}

// Ok sets dst to the underlying value if ok.
// Returns true if ok, false if not ok.
func (v Value[T]) Ok(dst *T) bool {
	if !v.IsOk() {
		return false
	}
	*dst = v.val
	return true
}

// Or returns the underlying value if ok, or def if not ok.
func (v Value[T]) Or(def T) T {
	if v.IsOk() {
		return v.val
	}
	return def
}

// OrZero returns the underlying value if ok, or the zero value if not ok.
func (v Value[T]) OrZero() T {
	return v.val
}

// OrElse returns the underlying value if ok, or the result of f if not ok.
func (v Value[T]) OrElse(f func() T) T {
	if v.IsOk() {
		return v.val
	}
	return f()
}

// Map returns the result of f on the underlying value of v.
func Map[T, T2 any](v Value[T], f func(T) T2) Value[T2] {
	if !v.IsOk() {
		return Value[T2]{}
	}
	return OfOk(f(v.val))
}
