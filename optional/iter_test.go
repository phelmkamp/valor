package optional_test

import (
	"testing"

	"github.com/phelmkamp/valor/optional"
)

func all[T any](s []T) func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func TestFirst(t *testing.T) {
	if got := optional.First(all([]string{})); got.IsOk() {
		t.Errorf("First() of empty slice = %v, want %v", got, optional.OfNotOk[string]())
	}
	var res string
	if got := optional.First(all([]string{"a", "b", "c"})); !got.Ok(&res) || res != "a" {
		t.Errorf("First() = %v, want %v", got, optional.OfOk("a"))
	}
}

func TestLast(t *testing.T) {
	if got := optional.Last(all([]string{})); got.IsOk() {
		t.Errorf("Last() of empty slice = %v, want %v", got, optional.OfNotOk[string]())
	}
	var res string
	if got := optional.Last(all([]string{"a", "b", "c"})); !got.Ok(&res) || res != "c" {
		t.Errorf("Last() = %v, want %v", got, optional.OfOk("c"))
	}
}

func TestMin(t *testing.T) {
	if got := optional.Min(all([]int{})); got.IsOk() {
		t.Errorf("Min() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res int
	if got := optional.Min(all([]int{2})); !got.Ok(&res) || res != 2 {
		t.Errorf("Min() of 1 element = %v, want %v", got, optional.OfOk(2))
	}
	if got := optional.Min(all([]int{3, 1, 2})); !got.Ok(&res) || res != 1 {
		t.Errorf("Min() = %v, want %v", got, optional.OfOk(1))
	}
}

func TestMinFunc(t *testing.T) {
	shorter := func(s1, s2 string) bool {
		return len(s1) < len(s2)
	}
	if got := optional.MinFunc(all([]string{}), shorter); got.IsOk() {
		t.Errorf("MinFunc() of empty slice = %v, want %v", got, optional.OfNotOk[string]())
	}
	var res string
	if got := optional.MinFunc(all([]string{"ab"}), shorter); !got.Ok(&res) || res != "ab" {
		t.Errorf("MinFunc() of 1 element = %v, want %v", got, optional.OfOk("ab"))
	}
	if got := optional.MinFunc(all([]string{"ab", "abc", "a"}), shorter); !got.Ok(&res) || res != "a" {
		t.Errorf("MinFunc() = %v, want %v", got, optional.OfOk("a"))
	}
}

func TestMax(t *testing.T) {
	if got := optional.Max(all([]int{})); got.IsOk() {
		t.Errorf("Max() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res int
	if got := optional.Max(all([]int{2})); !got.Ok(&res) || res != 2 {
		t.Errorf("Max() of 1 element = %v, want %v", got, optional.OfOk(2))
	}
	if got := optional.Max(all([]int{1, 3, 2})); !got.Ok(&res) || res != 3 {
		t.Errorf("Max() = %v, want %v", got, optional.OfOk(3))
	}
}

func TestMaxFunc(t *testing.T) {
	longer := func(s1, s2 string) bool {
		return len(s1) > len(s2)
	}
	if got := optional.MaxFunc(all([]string{}), longer); got.IsOk() {
		t.Errorf("MaxFunc() of empty slice = %v, want %v", got, optional.OfNotOk[string]())
	}
	var res string
	if got := optional.MaxFunc(all([]string{"ab"}), longer); !got.Ok(&res) || res != "ab" {
		t.Errorf("MaxFunc() of 1 element = %v, want %v", got, optional.OfOk("ab"))
	}
	if got := optional.MaxFunc(all([]string{"ab", "abc", "a"}), longer); !got.Ok(&res) || res != "abc" {
		t.Errorf("MaxFunc() = %v, want %v", got, optional.OfOk("abc"))
	}
}

func TestReduce(t *testing.T) {
	joiner := func(s1, s2 string) string {
		return s1 + s2
	}
	if got := optional.Reduce(all([]string{}), joiner); got.IsOk() {
		t.Errorf("Reduce() of empty slice = %v, want %v", got, optional.OfNotOk[string]())
	}
	var res string
	if got := optional.Reduce(all([]string{"a"}), joiner); !got.Ok(&res) || res != "a" {
		t.Errorf("Reduce() of 1 element = %v, want %v", got, optional.OfOk("a"))
	}
	if got := optional.Reduce(all([]string{"a", "b", "c"}), joiner); !got.Ok(&res) || res != "abc" {
		t.Errorf("Reduce() = %v, want %v", got, optional.OfOk("abc"))
	}
}

func TestSum(t *testing.T) {
	if got := optional.Sum(all([]int{})); got.IsOk() {
		t.Errorf("Sum() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res int
	if got := optional.Sum(all([]int{2})); !got.Ok(&res) || res != 2 {
		t.Errorf("Sum() of 1 element = %v, want %v", got, optional.OfOk(2))
	}
	if got := optional.Sum(all([]int{1, 2, 3})); !got.Ok(&res) || res != 6 {
		t.Errorf("Sum() = %v, want %v", got, optional.OfOk(6))
	}
}

func TestAvg(t *testing.T) {
	if got := optional.Avg(all([]int{})); got.IsOk() {
		t.Errorf("Avg() of empty slice = %v, want %v", got, optional.OfNotOk[int]())
	}
	var res float64
	if got := optional.Avg(all([]int{2})); !got.Ok(&res) || res != 2 {
		t.Errorf("Avg() of 1 element = %v, want %v", got, optional.OfOk(2))
	}
	if got := optional.Avg(all([]int{1, 2, 3, 4})); !got.Ok(&res) || res != 2.5 {
		t.Errorf("Avg() = %v, want %v", got, optional.OfOk(2.5))
	}
}
