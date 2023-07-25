package optional

import "github.com/phelmkamp/valor/constraints"

// Iter is an iterator over a collection of elements.
type Iter[T any] func(yield func(T) bool)

// First returns a Value of the first element.
// Returns a not-ok Value if there are no elements.
func First[T any](iter Iter[T]) Value[T] {
	// for v := range iter {
	// 	return OfOk(v)
	// }
	// return OfNotOk[T]()

	var first T
	var ok bool
	iter(func(v T) bool {
		first = v
		ok = true
		return false
	})
	return Of(first, ok)
}

// Last returns a Value of the last element.
// Returns a not-ok Value if there are no elements.
func Last[T any](iter Iter[T]) Value[T] {
	var last T
	var ok bool
	iter(func(v T) bool {
		last = v
		ok = true
		return true
	})
	return Of(last, ok)
}

// Min returns a Value of the minimum element.
// Returns a not-ok Value if there are no elements.
func Min[T constraints.Ordered](iter Iter[T]) Value[T] {
	var min T
	var ok bool
	iter(func(v T) bool {
		if v < min {
			min = v
		}
		ok = true
		return true
	})
	return Of(min, ok)
}

// MinFunc returns a Value of the minimum element.
// Returns a not-ok Value if there are no elements.
func MinFunc[T any](iter Iter[T], lt func(T, T) bool) Value[T] {
	var min T
	var ok bool
	iter(func(v T) bool {
		if lt(v, min) {
			min = v
		}
		ok = true
		return true
	})
	return Of(min, ok)
}

// Max returns a Value of the minimum element.
// Returns a not-ok Value if there are no elements.
func Max[T constraints.Ordered](iter Iter[T]) Value[T] {
	var max T
	var ok bool
	iter(func(v T) bool {
		if v > max {
			max = v
		}
		ok = true
		return true
	})
	return Of(max, ok)
}

// MaxFunc returns a Value of the maximum element.
// Returns a not-ok Value if there are no elements.
func MaxFunc[T any](iter Iter[T], gt func(T, T) bool) Value[T] {
	var max T
	var ok bool
	iter(func(v T) bool {
		if gt(v, max) {
			max = v
		}
		ok = true
		return true
	})
	return Of(max, ok)
}

// Reduce returns a Value of the result of accumulating all elements.
// Returns a not-ok Value if there are no elements.
func Reduce[T any](iter Iter[T], accum func(res T, v T) T) Value[T] {
	var res T
	var ok bool
	iter(func(v T) bool {
		if !ok {
			res = v
		} else {
			res = accum(res, v)
		}
		ok = true
		return true
	})
	return Of(res, ok)
}

// Sum returns a Value of the sum of all elements.
// Returns a not-ok Value if there are no elements.
func Sum[T constraints.Integer | constraints.Float](iter Iter[T]) Value[T] {
	return Reduce(iter, func(v, v2 T) T {
		return v + v2
	})
}

// Avg returns a Value of the average of all elements.
// Returns a not-ok Value if there are no elements.
func Avg[T constraints.Integer | constraints.Float](iter Iter[T]) Value[float64] {
	n := float64(1)
	sum := Reduce(iter, func(v, v2 T) T {
		n++
		return v + v2
	})
	return Map(sum, func(v T) float64 { return float64(v) / n })
}
