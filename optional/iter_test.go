package optional_test

import (
	"testing"

	"github.com/phelmkamp/valor/constraints"
	"github.com/phelmkamp/valor/optional"
)

func nums[T constraints.Float | constraints.Integer](s []T) func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func TestSum(t *testing.T) {
	if got := optional.Sum(nums([]int{})); got.IsOk() {
		t.Errorf("Sum() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res int
	if got := optional.Sum(nums([]int{1, 2, 3})); !got.Ok(&res) || res != 6 {
		t.Errorf("Sum() = %v, want %v", got, optional.OfOk(6))
	}
}

func TestAvg(t *testing.T) {
	if got := optional.Avg(nums([]int{})); got.IsOk() {
		t.Errorf("Avg() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res float64
	if got := optional.Avg(nums([]int{1, 2, 3, 4})); !got.Ok(&res) || res != 2.5 {
		t.Errorf("Avg() = %v, want %v", got, optional.OfOk(2.5))
	}
}
