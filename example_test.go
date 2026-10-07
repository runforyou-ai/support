package support_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/runforyou-ai/support"
)

func ExampleBlank() {
	var p *int
	fmt.Println(support.Blank(""), support.Blank("  "), support.Blank(p), support.Blank([]int{}))
	fmt.Println(support.Blank("a"), support.Blank(0), support.Blank(false))
	// Output:
	// true true true true
	// false false false
}

func ExampleFilled() {
	fmt.Println(support.Filled("a"), support.Filled(map[string]int{}))
	// Output: true false
}

func ExampleDefault() {
	fmt.Println(support.Default("", "nickname", "name"))
	fmt.Println(support.Default(0, 0, 20))
	// Output:
	// nickname
	// 20
}

func ExamplePtr() {
	p := support.Ptr(42)
	fmt.Println(*p)
	// Output: 42
}

func ExampleDeref() {
	var limit *int
	fmt.Println(support.Deref(limit), support.Deref(support.Ptr(5)))
	// Output: 0 5
}

func ExampleDerefOr() {
	var limit *int
	fmt.Println(support.DerefOr(limit, 20), support.DerefOr(support.Ptr(5), 20))
	// Output: 20 5
}

func ExampleTap() {
	var b strings.Builder
	s := support.Tap("hello", func(s string) { b.WriteString("saw " + s) })
	fmt.Println(s, "/", b.String())
	// Output: hello / saw hello
}

func ExampleWith() {
	n := support.With("hello", func(s string) int { return len(s) })
	fmt.Println(n)
	// Output: 5
}

func ExampleTransform() {
	upper := func(s string) string { return strings.ToUpper(s) }
	fmt.Println(support.Transform("go", upper, "n/a"))
	fmt.Println(support.Transform("  ", upper, "n/a"))
	// Output:
	// GO
	// n/a
}

func ExampleRetry() {
	err := support.Retry(context.Background(), 3, func(attempt int) error {
		fmt.Println("attempt", attempt)
		if attempt < 2 {
			return errors.New("temporary")
		}
		return nil
	}, support.WithDelay(time.Millisecond))
	fmt.Println(err)
	// Output:
	// attempt 1
	// attempt 2
	// <nil>
}

func ExampleRetryValue() {
	v, err := support.RetryValue(context.Background(), 3, func(attempt int) (string, error) {
		if attempt < 3 {
			return "", errors.New("temporary")
		}
		return "ok on " + strconv.Itoa(attempt), nil
	})
	fmt.Println(v, err)
	// Output: ok on 3 <nil>
}

func ExampleWithDelay() {
	err := support.Retry(context.Background(), 2, func(int) error {
		return errors.New("down")
	}, support.WithDelay(time.Millisecond))
	fmt.Println(err)
	// Output: down
}

func ExampleWithBackoff() {
	err := support.Retry(context.Background(), 3, func(attempt int) error {
		fmt.Println("attempt", attempt)
		return errors.New("down")
	}, support.WithBackoff(func(attempt int) time.Duration {
		return time.Duration(attempt) * time.Millisecond
	}))
	fmt.Println(err)
	// Output:
	// attempt 1
	// attempt 2
	// attempt 3
	// down
}

func ExampleWithExponentialBackoff() {
	// Waits 1ms, 2ms, then 3ms (capped) between attempts.
	err := support.Retry(context.Background(), 4, func(attempt int) error {
		if attempt < 4 {
			return errors.New("temporary")
		}
		return nil
	}, support.WithExponentialBackoff(time.Millisecond, 3*time.Millisecond))
	fmt.Println(err)
	// Output: <nil>
}

func ExampleWithRetryIf() {
	errNotFound := errors.New("not found")
	err := support.Retry(context.Background(), 5, func(attempt int) error {
		fmt.Println("attempt", attempt)
		return errNotFound
	}, support.WithRetryIf(func(err error) bool {
		return !errors.Is(err, errNotFound)
	}))
	fmt.Println(err)
	// Output:
	// attempt 1
	// not found
}

func ExampleRescue() {
	parse := func(s string) func() (int, error) {
		return func() (int, error) { return strconv.Atoi(s) }
	}
	fmt.Println(support.Rescue(parse("42"), -1))
	fmt.Println(support.Rescue(parse("x"), -1))
	fmt.Println(support.Rescue(func() (int, error) { panic("boom") }, -1))
	// Output:
	// 42
	// -1
	// -1
}
