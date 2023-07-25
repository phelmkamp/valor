package optional

import "github.com/phelmkamp/valor/constraints"

// Iter is an iterator over a collection of elements.
type Iter[T any] func(yield func(T) bool)

// NumIter is an iterator over a collection of numbers.
type NumIter[T constraints.Integer | constraints.Float] func(yield func(T) bool)

// First returns a Value of the first element.
// Returns a not-ok Value if there are no elements.
func (iter Iter[T]) First() Value[T] {
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
func (iter Iter[T]) Last() Value[T] {
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
func (iter Iter[T]) Min(lt func(T, T) bool) Value[T] {
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

// Max returns a Value of the maximum element.
// Returns a not-ok Value if there are no elements.
func (iter Iter[T]) Max(gt func(T, T) bool) Value[T] {
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
func (iter Iter[T]) Reduce(accum func(res T, v T) T) Value[T] {
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
func (iter NumIter[T]) Sum() Value[T] {
	return Iter[T](iter).Reduce(func(v, v2 T) T {
		return v + v2
	})
}

// Avg returns a Value of the average of all elements.
// Returns a not-ok Value if there are no elements.
func (iter NumIter[T]) Avg() Value[float64] {
	n := float64(1)
	sum := Iter[T](iter).Reduce(func(v, v2 T) T {
		n++
		return v + v2
	})
	return Map(sum, func(v T) float64 { return float64(v) / n })
}
