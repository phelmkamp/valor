package result_test

import (
	"errors"
	"fmt"
	"github.com/phelmkamp/valor/result"
	"github.com/phelmkamp/valor/value"
	"io"
	"io/fs"
	"strconv"
	"strings"
)

func Example() {
	var w strings.Builder
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
	fmt.Println(n)

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
	// Output: 3
	// fail
}

func mid(fail bool) (string, error) {
	var s string
	if res := result.Of(leaf(fail)); !value.Map(res.Value(), strconv.Itoa).Ok(&s) {
		return "", res.Errorf("leaf() failed: %w").Error()
	}
	return "", nil
}

func leaf(fail bool) (int, error) {
	if fail {
		return 0, errors.New("fail")
	}
	return 1, nil
}
