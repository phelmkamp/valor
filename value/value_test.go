// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package value_test

import (
	"fmt"
	"github.com/phelmkamp/valor/value"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func Example() {
	m := map[string]int{"foo": 42}
	foo, ok := m["foo"]
	val := value.Of(foo, ok)
	fmt.Println(val.IsOk()) // true

	var foo2 int
	fmt.Println(val.Ok(&foo2), foo2) // true 42

	val2 := value.Map(val, strconv.Itoa)
	fmt.Println(val2) // {42 true}

	bar, ok := m["bar"]
	val3 := value.Of(bar, ok)
	fmt.Println(val3.Or(-1))                          // -1
	fmt.Println(val3.OrZero())                        // 0
	fmt.Println(val3.OrElse(func() int { return 1 })) // 1

	// switch
	switch val {
	case val.OfOk():
		foo = val.MustOk()
		fmt.Println("Ok")
	case value.OfNotOk[int]():
		fmt.Println("Not Ok")
		return
	}
	// Output: true
	// true 42
	// {42 true}
	// -1
	// 0
	// 1
	// Ok
}

func TestValue_Ok(t *testing.T) {
	i := -1
	if got := value.OfNotOk[int]().Ok(&i); got != false {
		t.Errorf("Ok() = %v, want %v", got, false)
	}
	if i != -1 {
		t.Errorf("i after Ok() = %v, want %v", i, -1)
	}
	if got := value.OfOk(1).Ok(&i); got != true {
		t.Errorf("Ok() = %v, want %v", got, true)
	}
	if i != 1 {
		t.Errorf("i after Ok() = %v, want %v", i, 1)
	}
}

func TestContains(t *testing.T) {
	if got := value.Contains(value.OfNotOk[string](), ""); got != false {
		t.Errorf("Contains() = %v, want %v", got, false)
	}
	if got := value.Contains(value.OfOk(""), "foo"); got != false {
		t.Errorf("Contains() = %v, want %v", got, false)
	}
	if got := value.Contains(value.OfOk(""), ""); got != true {
		t.Errorf("Contains() = %v, want %v", got, true)
	}
}

func TestFlatMap(t *testing.T) {
	toVal := func(f float64) value.Value[float64] {
		return value.OfOk(f * 2.0)
	}
	if got := value.FlatMap(value.OfNotOk[float64](), toVal); got != value.OfNotOk[float64]() {
		t.Errorf("FlatMap() = %v, want %v", got, value.OfNotOk[float64]())
	}
	if got := value.FlatMap(value.OfOk(2.0), toVal); got != value.OfOk(4.0) {
		t.Errorf("FlatMap() = %v, want %v", got, value.OfOk(4.0))
	}
}

func TestFlatten(t *testing.T) {
	if got := value.Flatten(value.OfNotOk[value.Value[int]]()); got != value.OfNotOk[int]() {
		t.Errorf("Flatten() = %v, want %v", got, value.OfNotOk[bool]())
	}
	if got := value.Flatten(value.OfOk(value.OfOk(1))); got != value.OfOk(1) {
		t.Errorf("Flatten() = %v, want %v", got, value.OfOk(1))
	}
}

func TestMap(t *testing.T) {
	if got := value.Map(value.OfNotOk[int](), strconv.Itoa); got != value.OfNotOk[string]() {
		t.Errorf("Map() = %v, want %v", got, value.OfNotOk[string]())
	}
	if got := value.Map(value.OfOk(2), strconv.Itoa); got != value.OfOk("2") {
		t.Errorf("Map() = %v, want %v", got, value.OfOk("2"))
	}
}

func TestOf(t *testing.T) {
	m := sync.Map{}
	if got := value.Of(m.Load("foo")); got != value.OfNotOk[any]() {
		t.Errorf("Of() = %v, want %v", got, value.OfNotOk[any]())
	}
	m.Store("foo", "bar")
	if got := value.Of(m.Load("foo")); got != value.OfOk[any]("bar") {
		t.Errorf("Of() = %v, want %v", got, value.OfOk[any]("bar"))
	}
}

func TestOfNotOk(t *testing.T) {
	want := value.Value[uint8]{}
	if got := value.OfNotOk[uint8](); got != want {
		t.Errorf("OfNotOk() = %v, want %v", got, want)
	}
}

func TestOfOk(t *testing.T) {
	if got := value.OfOk("foo"); got.MustOk() != "foo" {
		t.Errorf("OfOk().v = %v, want %v", got.MustOk(), "foo")
	}
}

func splitFirstColon(s string) (string, int) {
	i := strings.Index(s, ":")
	if i < 0 {
		return s, i
	}
	return s[:i], i
}

func ExampleUnzipWith() {
	str := "foo:bar:baz"
	for {
		strVal, idxVal := value.UnzipWith(value.OfOk(str), splitFirstColon)
		fmt.Print(strVal.MustOk())
		idxVal = value.FlatMap(idxVal, func(i int) value.Value[int] { return value.Of(i, i >= 0) })
		if !idxVal.IsOk() {
			break
		}
		i := idxVal.MustOk()
		str = str[i+1:]
		fmt.Println(" " + str)
	}
	// Output:
	// foo bar:baz
	// bar baz
	// baz
}

func TestUnzipWith(t *testing.T) {
	gotVal2, gotVal3 := value.UnzipWith(value.OfNotOk[string](), splitFirstColon)
	if gotVal2 != value.OfNotOk[string]() {
		t.Errorf("UnzipWith() gotVal2 = %v, want %v", gotVal2, value.OfNotOk[string]())
	}
	if gotVal3 != value.OfNotOk[int]() {
		t.Errorf("UnzipWith() gotVal3 = %v, want %v", gotVal3, value.OfNotOk[string]())
	}
	gotVal2, gotVal3 = value.UnzipWith(value.OfOk("foo:bar:baz"), splitFirstColon)
	if gotVal2 != value.OfOk("foo") {
		t.Errorf("UnzipWith() gotVal2 = %v, want %v", gotVal2, value.OfOk("foo"))
	}
	if gotVal3 != value.OfOk(3) {
		t.Errorf("UnzipWith() gotVal3 = %v, want %v", gotVal3, value.OfOk(3))
	}
}

func TestValue_Do(t *testing.T) {
	var n int
	count := func(string) {
		n++
	}
	if got := value.OfNotOk[string]().Do(count); got != value.OfNotOk[string]() {
		t.Errorf("Do() = %v, want %v", got, value.OfNotOk[string]())
	}
	if n != 0 {
		t.Errorf(" n after Do() = %v, want %v", n, 0)
	}
	if got := value.OfOk("foo").Do(count); got != value.OfOk("foo") {
		t.Errorf("Do() = %v, want %v", got, value.OfOk("foo"))
	}
	if n != 1 {
		t.Errorf(" n after Do() = %v, want %v", n, 1)
	}
}

func TestValue_Filter(t *testing.T) {
	isEven := func(i int) bool {
		return i%2 == 0
	}
	if got := value.OfNotOk[int]().Filter(isEven); got != value.OfNotOk[int]() {
		t.Errorf("Filter() = %v, want %v", got, value.OfNotOk[int]())
	}
	if got := value.OfOk(3).Filter(isEven); got != value.OfNotOk[int]() {
		t.Errorf("Filter() = %v, want %v", got, value.OfNotOk[int]())
	}
	if got := value.OfOk(4).Filter(isEven); got != value.OfOk(4) {
		t.Errorf("Filter() = %v, want %v", got, value.OfOk(4))
	}
}

func TestValue_IsOk(t *testing.T) {
	if got := value.OfNotOk[string]().IsOk(); got {
		t.Errorf("IsOk() = %v, want %v", got, false)
	}
	if got := value.OfOk(1).IsOk(); !got {
		t.Errorf("IsOk() = %v, want %v", got, true)
	}
}

func recoverMustOk(val value.Value[float64]) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			msg, _ := r.(string)
			if panicked = msg == "value.MustOk(): not ok"; !panicked {
				// propagate unexpected panic
				panic(r)
			}
		}
	}()
	val.MustOk()
	return
}

func TestValue_MustOk(t *testing.T) {
	if panicked := recoverMustOk(value.OfNotOk[float64]()); !panicked {
		t.Errorf("MustOk() panicked = %v, want %v", panicked, true)
	}
	if panicked := recoverMustOk(value.OfOk(1.0)); panicked {
		t.Errorf("MustOk() panicked = %v, want %v", panicked, false)
	}
}

func TestValue_Or(t *testing.T) {
	if got := value.OfNotOk[int]().Or(-1); got != -1 {
		t.Errorf("Or() = %v, want %v", got, -1)
	}
	if got := value.OfOk(1).Or(-1); got != 1 {
		t.Errorf("Or() = %v, want %v", got, 1)
	}
}

func TestValue_OrElse(t *testing.T) {
	rand := rand.NewSource(42)
	if got := value.OfNotOk[int64]().OrElse(rand.Int63); got != 3440579354231278675 {
		t.Errorf("OrElse() = %v, want %v", got, 3440579354231278675)
	}
	if got := value.OfOk(int64(1)).OrElse(rand.Int63); got != 1 {
		t.Errorf("OrElse() = %v, want %v", got, 1)
	}
}

func TestValue_OrZero(t *testing.T) {
	if got := value.OfNotOk[string]().OrZero(); got != "" {
		t.Errorf("OrZero() = %v, want %v", got, "")
	}
	if got := value.OfOk("foo").OrZero(); got != "foo" {
		t.Errorf("OrZero() = %v, want %v", got, "foo")
	}
}

func TestZipWith(t *testing.T) {
	if got := value.ZipWith(value.OfNotOk[[]string](), value.OfOk(","), strings.Join); got != value.OfNotOk[string]() {
		t.Errorf("ZipWith() = %v, want %v", got, value.OfNotOk[string]())
	}
	if got := value.ZipWith(value.OfOk([]string{"foo", "bar"}), value.OfNotOk[string](), strings.Join); got != value.OfNotOk[string]() {
		t.Errorf("ZipWith() = %v, want %v", got, value.OfNotOk[string]())
	}
	if got := value.ZipWith(value.OfOk([]string{"foo", "bar"}), value.OfOk(","), strings.Join); got != value.OfOk("foo,bar") {
		t.Errorf("ZipWith() = %v, want %v", got, value.OfOk("foo,bar"))
	}
}
