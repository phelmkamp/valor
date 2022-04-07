# valor

This module provides types that optionally contain a value; hence the name valor, short for "value or".

## Installation

```bash
go get github.com/phelmkamp/valor
```

## Types

### Value

`Value` is modeled after the "comma ok" idiom. It contains a value (ok) or nothing (not ok).

```go
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
```

### Result

`Result` contains either a value or an error.

```go
// traditional
if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
    fmt.Println(res.Error())
    return
}

// try to get value, printing wrapped error if not ok
var n int
if res := result.Of(w.Write([]byte("foo"))); !res.Value().Ok(&n) {
    fmt.Println(res.Errorf("Write() failed: %w"))
    return
}

// errors.Is
if res := result.Of(w.Write([]byte("foo"))); res.ErrorIs(io.ErrShortWrite) {
    fmt.Println(res)
    return
}

// errors.As
if res := result.Of(w.Write([]byte("foo"))); res.IsError() {
    var err *fs.PathError
    if res.ErrorAs(&err) {
        fmt.Println("path=" + err.Path)
    }
    fmt.Println(res)
    return
}

// errors.Unwrap
if res := result.Of(mid(true)); res.IsError() {
    fmt.Println(res.ErrorUnwrap())
    return
}
```