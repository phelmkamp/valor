package optional_test

import (
	"testing"

	"github.com/phelmkamp/valor/optional"
)

func sliceIter[T any](s []T) optional.Iter[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func TestSum(t *testing.T) {
	if got := optional.Sum(sliceIter([]int{})); got.IsOk() {
		t.Errorf("Sum() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res int
	if got := optional.Sum(sliceIter([]int{1, 2, 3})); !got.Ok(&res) || res != 6 {
		t.Errorf("Sum() = %v, want %v", got, optional.OfOk(6))
	}
}

func TestAvg(t *testing.T) {
	if got := optional.Avg(sliceIter([]int{})); got.IsOk() {
		t.Errorf("Avg() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res float64
	if got := optional.Avg(sliceIter([]int{1, 2, 3, 4})); !got.Ok(&res) || res != 2.5 {
		t.Errorf("Avg() = %v, want %v", got, optional.OfOk(2.5))
	}
}
