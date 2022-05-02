// Copyright 2022 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package optional_test

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/phelmkamp/valor/optional"
	"github.com/phelmkamp/valor/tuple/four"
	"github.com/phelmkamp/valor/tuple/singleton"
	"github.com/phelmkamp/valor/tuple/two"
)

// type checks
var (
	_ = optional.OfNotOk[*struct{}]()
	_ = optional.OfNotOk[map[string]struct{}]()
	_ = optional.OfNotOk[[]struct{}]()
)

func Example() {
	m := map[string]int{"foo": 42}
	val := optional.OfIndex(m, "foo")
	fmt.Println(val.IsOk()) // true

	var foo int
	fmt.Println(val.Ok(&foo), foo) // true 42

	valStr := optional.Map(val, strconv.Itoa)
	fmt.Println(valStr) // {42 true}

	val = optional.OfIndex(m, "bar")
	fmt.Println(val.Or(-1))                          // -1
	fmt.Println(val.OrZero())                        // 0
	fmt.Println(val.OrElse(func() int { return 1 })) // 1

	// switch
	switch val {
	case val.OfOk():
		foo = val.MustOk()
		fmt.Println(foo)
	case optional.OfNotOk[int]():
		fmt.Println("Not Ok")
		return
	}
	// Output: true
	// true 42
	// {42 true}
	// -1
	// 0
	// 1
	// Not Ok
}

func TestValue_Ok(t *testing.T) {
	i := -1
	if got := optional.OfNotOk[int]().Ok(&i); got != false {
		t.Errorf("Ok() = %v, want %v", got, false)
	}
	if i != -1 {
		t.Errorf("i after Ok() = %v, want %v", i, -1)
	}
	if got := optional.OfOk(1).Ok(&i); got != true {
		t.Errorf("Ok() = %v, want %v", got, true)
	}
	if i != 1 {
		t.Errorf("i after Ok() = %v, want %v", i, 1)
	}
}

func TestContains(t *testing.T) {
	if got := optional.Contains(optional.OfNotOk[string](), ""); got != false {
		t.Errorf("Contains() = %v, want %v", got, false)
	}
	if got := optional.Contains(optional.OfOk(""), "foo"); got != false {
		t.Errorf("Contains() = %v, want %v", got, false)
	}
	if got := optional.Contains(optional.OfOk(""), ""); got != true {
		t.Errorf("Contains() = %v, want %v", got, true)
	}
}

func TestFlatMap(t *testing.T) {
	toVal := func(f float64) optional.Value[float64] {
		return optional.OfOk(f * 2.0)
	}
	if got := optional.FlatMap(optional.OfNotOk[float64](), toVal); got != optional.OfNotOk[float64]() {
		t.Errorf("FlatMap() = %v, want %v", got, optional.OfNotOk[float64]())
	}
	if got := optional.FlatMap(optional.OfOk(2.0), toVal); got != optional.OfOk(4.0) {
		t.Errorf("FlatMap() = %v, want %v", got, optional.OfOk(4.0))
	}
}

func TestFlatten(t *testing.T) {
	if got := optional.Flatten(optional.OfNotOk[optional.Value[int]]()); got != optional.OfNotOk[int]() {
		t.Errorf("Flatten() = %v, want %v", got, optional.OfNotOk[bool]())
	}
	if got := optional.Flatten(optional.OfOk(optional.OfOk(1))); got != optional.OfOk(1) {
		t.Errorf("Flatten() = %v, want %v", got, optional.OfOk(1))
	}
}

func TestMap(t *testing.T) {
	if got := optional.Map(optional.OfNotOk[int](), strconv.Itoa); got != optional.OfNotOk[string]() {
		t.Errorf("Map() = %v, want %v", got, optional.OfNotOk[string]())
	}
	if got := optional.Map(optional.OfOk(2), strconv.Itoa); got != optional.OfOk("2") {
		t.Errorf("Map() = %v, want %v", got, optional.OfOk("2"))
	}
}

func TestOf(t *testing.T) {
	m := sync.Map{}
	if got := optional.Of(m.Load("foo")); got != optional.OfNotOk[any]() {
		t.Errorf("Of() = %v, want %v", got, optional.OfNotOk[any]())
	}
	m.Store("foo", "bar")
	if got := optional.Of(m.Load("foo")); got != optional.OfOk[any]("bar") {
		t.Errorf("Of() = %v, want %v", got, optional.OfOk[any]("bar"))
	}
}

func TestOfNotOk(t *testing.T) {
	want := optional.Value[uint8]{}
	if got := optional.OfNotOk[uint8](); got != want {
		t.Errorf("OfNotOk() = %v, want %v", got, want)
	}
}

func TestOfOk(t *testing.T) {
	if got := optional.OfOk("foo"); got.MustOk() != "foo" {
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
		strVal, idxVal := optional.UnzipWith(optional.OfOk(str), splitFirstColon)
		fmt.Print(strVal.MustOk())
		idxVal = optional.FlatMap(idxVal, func(i int) optional.Value[int] { return optional.Of(i, i >= 0) })
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

func ExampleUnzipWith_tuple() {
	v1234 := optional.OfOk(four.TupleOf(1, "two", 3.0, []int{4}))
	v13, v24 := optional.UnzipWith(v1234, two.TupleUnzip[int, string, float64, []int])
	fmt.Println(v13.MustOk(), v24.MustOk())
	v1, v3 := optional.UnzipWith(v13, singleton.SetUnzip[int, float64])
	v2, v4 := optional.UnzipWith(v24, singleton.SetUnzip[string, []int])
	fmt.Println(v1.MustOk(), v3.MustOk(), v2.MustOk(), v4.MustOk())
	// Output:
	// {1 3} {two [4]}
	// {1} {3} {two} {[4]}
}

func TestUnzipWith(t *testing.T) {
	gotVal2, gotVal3 := optional.UnzipWith(optional.OfNotOk[string](), splitFirstColon)
	if gotVal2 != optional.OfNotOk[string]() {
		t.Errorf("UnzipWith() gotVal2 = %v, want %v", gotVal2, optional.OfNotOk[string]())
	}
	if gotVal3 != optional.OfNotOk[int]() {
		t.Errorf("UnzipWith() gotVal3 = %v, want %v", gotVal3, optional.OfNotOk[string]())
	}
	gotVal2, gotVal3 = optional.UnzipWith(optional.OfOk("foo:bar:baz"), splitFirstColon)
	if gotVal2 != optional.OfOk("foo") {
		t.Errorf("UnzipWith() gotVal2 = %v, want %v", gotVal2, optional.OfOk("foo"))
	}
	if gotVal3 != optional.OfOk(3) {
		t.Errorf("UnzipWith() gotVal3 = %v, want %v", gotVal3, optional.OfOk(3))
	}
}

func TestValue_Do(t *testing.T) {
	var n int
	count := func(string) {
		n++
	}
	if got := optional.OfNotOk[string]().Do(count); got != optional.OfNotOk[string]() {
		t.Errorf("Do() = %v, want %v", got, optional.OfNotOk[string]())
	}
	if n != 0 {
		t.Errorf(" n after Do() = %v, want %v", n, 0)
	}
	if got := optional.OfOk("foo").Do(count); got != optional.OfOk("foo") {
		t.Errorf("Do() = %v, want %v", got, optional.OfOk("foo"))
	}
	if n != 1 {
		t.Errorf(" n after Do() = %v, want %v", n, 1)
	}
}

func TestValue_Filter(t *testing.T) {
	isEven := func(i int) bool {
		return i%2 == 0
	}
	if got := optional.OfNotOk[int]().Filter(isEven); got != optional.OfNotOk[int]() {
		t.Errorf("Filter() = %v, want %v", got, optional.OfNotOk[int]())
	}
	if got := optional.OfOk(3).Filter(isEven); got != optional.OfNotOk[int]() {
		t.Errorf("Filter() = %v, want %v", got, optional.OfNotOk[int]())
	}
	if got := optional.OfOk(4).Filter(isEven); got != optional.OfOk(4) {
		t.Errorf("Filter() = %v, want %v", got, optional.OfOk(4))
	}
}

func TestValue_IsOk(t *testing.T) {
	if got := optional.OfNotOk[string]().IsOk(); got {
		t.Errorf("IsOk() = %v, want %v", got, false)
	}
	if got := optional.OfOk(1).IsOk(); !got {
		t.Errorf("IsOk() = %v, want %v", got, true)
	}
}

func recoverMustOk(val optional.Value[float64]) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			msg, _ := r.(string)
			if panicked = msg == "Value.MustOk(): not ok"; !panicked {
				// propagate unexpected panic
				panic(r)
			}
		}
	}()
	val.MustOk()
	return
}

func TestValue_MustOk(t *testing.T) {
	if panicked := recoverMustOk(optional.OfNotOk[float64]()); !panicked {
		t.Errorf("MustOk() panicked = %v, want %v", panicked, true)
	}
	if panicked := recoverMustOk(optional.OfOk(1.0)); panicked {
		t.Errorf("MustOk() panicked = %v, want %v", panicked, false)
	}
}

func TestValue_Or(t *testing.T) {
	if got := optional.OfNotOk[int]().Or(-1); got != -1 {
		t.Errorf("Or() = %v, want %v", got, -1)
	}
	if got := optional.OfOk(1).Or(-1); got != 1 {
		t.Errorf("Or() = %v, want %v", got, 1)
	}
}

// ExampleValue_OrElse demonstrates that OrElse can be used to work with a cache.
func ExampleValue_OrElse() {
	cache := make(map[string]string)
	load := func(k string) string {
		// expensive call to load value goes here
		v := "bar"
		cache[k] = v
		return v
	}

	v := optional.OfIndex(cache, "foo").OrElse(func() string { return load("foo") })
	fmt.Println(v)
	// Output: bar
}

func TestValue_OrElse(t *testing.T) {
	rand := rand.NewSource(42)
	if got := optional.OfNotOk[int64]().OrElse(rand.Int63); got != 3440579354231278675 {
		t.Errorf("OrElse() = %v, want %v", got, 3440579354231278675)
	}
	if got := optional.OfOk(int64(1)).OrElse(rand.Int63); got != 1 {
		t.Errorf("OrElse() = %v, want %v", got, 1)
	}
}

func TestValue_OrZero(t *testing.T) {
	if got := optional.OfNotOk[string]().OrZero(); got != "" {
		t.Errorf("OrZero() = %v, want %v", got, "")
	}
	if got := optional.OfOk("foo").OrZero(); got != "foo" {
		t.Errorf("OrZero() = %v, want %v", got, "foo")
	}
}

func ExampleZipWith_tuple() {
	v1 := optional.OfOk(singleton.SetOf(1))
	v2 := optional.OfOk(singleton.SetOf("two"))
	v3 := optional.OfOk(singleton.SetOf(3.0))
	v4 := optional.OfOk(singleton.SetOf([]int{4}))
	v12 := optional.ZipWith(v1, v2, singleton.SetZip[int, string])
	v34 := optional.ZipWith(v3, v4, singleton.SetZip[float64, []int])
	fmt.Println(v12.MustOk(), v34.MustOk())
	v1324 := optional.ZipWith(v12, v34, two.TupleZip[int, float64, string, []int])
	fmt.Println(v1324.MustOk())
	// Output:
	// {1 two} {3 [4]}
	// {1 3 two [4]}
}

func TestZipWith(t *testing.T) {
	if got := optional.ZipWith(optional.OfNotOk[[]string](), optional.OfOk(","), strings.Join); got != optional.OfNotOk[string]() {
		t.Errorf("ZipWith() = %v, want %v", got, optional.OfNotOk[string]())
	}
	if got := optional.ZipWith(optional.OfOk([]string{"foo", "bar"}), optional.OfNotOk[string](), strings.Join); got != optional.OfNotOk[string]() {
		t.Errorf("ZipWith() = %v, want %v", got, optional.OfNotOk[string]())
	}
	if got := optional.ZipWith(optional.OfOk([]string{"foo", "bar"}), optional.OfOk(","), strings.Join); got != optional.OfOk("foo,bar") {
		t.Errorf("ZipWith() = %v, want %v", got, optional.OfOk("foo,bar"))
	}
}

// ExampleOfPointer demonstrates how to dereference a pointer without the risk of a panic due to nil.
func ExampleOfPointer() {
	var i *int
	j := optional.Map(
		optional.OfPointer(i),
		func(p *int) int { return *p },
	).OrZero()
	fmt.Println(j)

	i = new(int)
	*i++
	j = optional.Map(
		optional.OfPointer(i),
		func(p *int) int { return *p },
	).OrZero()
	fmt.Println(j)
	// Output:
	// 0
	// 1
}

func TestOfPointer(t *testing.T) {
	i := 1
	type args struct {
		p *int
	}
	tests := []struct {
		name string
		args args
		want optional.Value[*int]
	}{
		{
			name: "nil",
			args: args{p: nil},
			want: optional.OfNotOk[*int](),
		},
		{
			name: "1",
			args: args{p: &i},
			want: optional.OfOk(&i),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optional.OfPointer(tt.args.p); got != tt.want {
				t.Errorf("OfPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOfAssert(t *testing.T) {
	tm := time.Now()
	// Stringer has concrete type Time
	if got := optional.OfAssert[time.Time](fmt.Stringer(tm)); got != optional.OfOk(tm) {
		t.Errorf("OfAssert() = %v, want %v", got, optional.OfOk(tm))
	}
	// Reader does not have concrete type Time
	if got := optional.OfAssert[time.Time](io.Reader(bytes.NewBufferString(""))); got != optional.OfNotOk[time.Time]() {
		t.Errorf("OfAssert() = %v, want %v", got, optional.OfNotOk[time.Time]())
	}
	// Time does not implement Reader
	if got := optional.OfAssert[io.Reader](tm); got != optional.OfNotOk[io.Reader]() {
		t.Errorf("OfAssert() = %v, want %v", got, optional.OfNotOk[io.Reader]())
	}
	// Time does implement Stringer
	if got := optional.OfAssert[fmt.Stringer](tm); got != optional.OfOk[fmt.Stringer](tm) {
		t.Errorf("OfAssert() = %v, want %v", got, optional.OfOk[fmt.Stringer](tm))
	}
	// Time (as Stringer) does implement TextMarshaler
	if got := optional.OfAssert[encoding.TextMarshaler](tm); got != optional.OfOk[encoding.TextMarshaler](tm) {
		t.Errorf("OfAssert() = %v, want %v", got, optional.OfNotOk[io.Reader]())
	}
	// nil does not implement Reader
	if got := optional.OfAssert[io.Reader, fmt.Stringer](nil); got != optional.OfNotOk[io.Reader]() {
		t.Errorf("OfAssert() = %v, want %v", got, optional.OfNotOk[io.Reader]())
	}
}

func TestOfIndex(t *testing.T) {
	type args struct {
		m map[string]int
		k string
	}
	tests := []struct {
		name string
		args args
		want optional.Value[int]
	}{
		{
			name: "nil",
			args: args{m: nil, k: "foo"},
			want: optional.OfNotOk[int](),
		},
		{
			name: "empty",
			args: args{m: make(map[string]int), k: "foo"},
			want: optional.OfNotOk[int](),
		},
		{
			name: "missing",
			args: args{m: map[string]int{"bar": 42}, k: "foo"},
			want: optional.OfNotOk[int](),
		},
		{
			name: "ok",
			args: args{m: map[string]int{"foo": 42}, k: "foo"},
			want: optional.OfOk(42),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optional.OfIndex(tt.args.m, tt.args.k); got != tt.want {
				t.Errorf("OfIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOfReceive(t *testing.T) {
	// nil
	if got := optional.OfReceive[int](nil); got != optional.OfNotOk[int]() {
		t.Errorf("OfReceive() = %v, want %v", got, optional.OfNotOk[int]())
	}
	// closed
	ch := make(chan int)
	close(ch)
	if got := optional.OfReceive(ch); got != optional.OfNotOk[int]() {
		t.Errorf("OfReceive() = %v, want %v", got, optional.OfNotOk[int]())
	}
	// ok
	ch = make(chan int, 1)
	ch <- 42
	if got := optional.OfReceive(ch); got != optional.OfOk(42) {
		t.Errorf("OfReceive() = %v, want %v", got, optional.OfOk(42))
	}
}
