# valor

[![Go Reference](https://pkg.go.dev/badge/github.com/phelmkamp/valor.svg)](https://pkg.go.dev/github.com/phelmkamp/valor)
[![Go Report Card](https://goreportcard.com/badge/github.com/phelmkamp/valor)](https://goreportcard.com/report/github.com/phelmkamp/valor)
[![codecov](https://codecov.io/gh/phelmkamp/valor/branch/main/graph/badge.svg?token=GH8IYR78VD)](https://codecov.io/gh/phelmkamp/valor)

This module provides option and result types that optionally contain a value; hence the name valor, short for "value or".

This is not an attempt to make Go code look less like Go.
Instead, the goal is to codify the ["comma ok"](https://blog.toshima.ru/2019/07/21/go-comma-ok-idiom.html) and
["errors are values"](https://go.dev/blog/errors-are-values) principles that Go already encourages.

## Installation

```bash
go get github.com/phelmkamp/valor
```

## Types

### Value

[`optional.Value`](https://pkg.go.dev/github.com/phelmkamp/valor/optional) is modeled after the "comma ok" idiom.
It contains a value (ok) or nothing (not ok).

```go
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
```

### Result

[`result.Result`](https://pkg.go.dev/github.com/phelmkamp/valor/result) contains either a value or an error.

```go
// traditional
if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
    fmt.Println(res.Error())
    return
}

// try to get value, printing wrapped error if not ok
// note: only relevant values are in-scope after handling
var n int
if res := result.Of(w.Write([]byte("foo"))); !res.Value().Ok(&n) {
    fmt.Println(res.Errorf("Write() failed: %w").Error())
    return
}

// same as above with multiple values
var s string
var b bool
if res := two.TupleResultOf(multi(false)); !res.Value().Do(
    func(t two.Tuple[string, bool]) { s, b = t.Values() },
).IsOk() {
    fmt.Println(res.Errorf("multi() failed: %w").Error())
    return
}

// errors.Is
if res := result.Of(w.Write([]byte("foo"))); res.ErrorIs(io.ErrShortWrite) {
    fmt.Println(res.Error())
    return
}

// errors.As
if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
    var err *fs.PathError
    if res.ErrorAs(&err) {
        fmt.Println("path=" + err.Path)
    }
    fmt.Println(res.Error())
    return
}

// errors.Unwrap
if res := result.Of(mid(true)); res.IsError() {
    fmt.Println(res.ErrorUnwrap().Error())
    return
}
```

### Tuples

[`unit.Type`](https://pkg.go.dev/github.com/phelmkamp/valor/tuple/unit), [`singleton.Set`](https://pkg.go.dev/github.com/phelmkamp/valor/tuple/singleton),
[`two.Tuple`](https://pkg.go.dev/github.com/phelmkamp/valor/tuple/two), [`three.Tuple`](https://pkg.go.dev/github.com/phelmkamp/valor/tuple/three), and
[`four.Tuple`](https://pkg.go.dev/github.com/phelmkamp/valor/tuple/four) contain zero through four values respectively.
Among other things, they enable `Value` and `Result` to contain a variable number of values.

{% raw %}
```go
get := func(string, int, bool) {
    return "a", 1, true
}
val := two.TupleValueOf(get())
fmt.Println(val) // {{a 1} true}
```
{% endraw %}

### Enum

[`enum.Enum`](https://pkg.go.dev/github.com/phelmkamp/valor/enum) is an enumerated type.
It's initialized with a set of allowed values and then each "copy" optionally contains a currently selected value.

```go
const (
	Clubs    = "clubs"
	Diamonds = "diamonds"
	Hearts   = "hearts"
	Spades   = "spades"
)
var Suit = enum.OfString(Clubs, Diamonds, Hearts, Spades)
func main() {
    fmt.Println(Suit.Values())          // [clubs diamonds hearts spades]
    fmt.Println(Suit.ValueOf("Foo"))    // { false}
    fmt.Println(Suit.ValueOf(Hearts))   // {hearts true}
}
```

## Similar concepts in other languages

### Rust

`Value` is like Rust's [`Option`](https://doc.rust-lang.org/std/option/enum.Option.html).
`Result` is like Rust's [`Result`](https://doc.rust-lang.org/std/result/enum.Result.html).

A `switch` statement provides similar semantics to Rust's `match`:

```go
var foo int
switch val {
case val.OfOk():
    foo = val.MustOk()
    fmt.Println("Ok", foo)
case optional.OfNotOk[int]():
    fmt.Println("Not Ok")
    return
}

var n int
switch res {
case res.OfOk():
    n = res.Value().MustOk()
    fmt.Println("Ok", n)
case res.OfError():
    fmt.Println("Error", res.Error())
    return
}
```

### Java

`Value` is like Java's [`Optional`](https://docs.oracle.com/en/java/javase/11/docs/api/java.base/java/util/Optional.html).

### More

 * [https://en.wikipedia.org/wiki/Result_type](https://en.wikipedia.org/wiki/Result_type)
 * [https://en.wikipedia.org/wiki/Option_type](https://en.wikipedia.org/wiki/Option_type)

## Releases

This module is currently at v0 but every effort will be made to avoid breaking changes.
Instead, functionality will be deprecated as needed with plans to remove in v1.

## Linter

[valorcheck](https://github.com/phelmkamp/valor/tree/main/valorcheck#readme) is a linter to check that access to an optional value is guarded against the case where the value is not present.
