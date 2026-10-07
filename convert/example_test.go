package convert_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/runforyou-ai/support/convert"
)

func ExampleToString() {
	n := 7
	fmt.Println(convert.ToString(3.5))
	fmt.Println(convert.ToString(true))
	fmt.Println(convert.ToString(&n))
	fmt.Println(convert.ToString(90 * time.Second))
	fmt.Printf("%q\n", convert.ToString(nil))
	// Output:
	// 3.5
	// true
	// 7
	// 1m30s
	// ""
}

func ExampleToBool() {
	for _, v := range []any{"Yes", "off", 2, "maybe"} {
		b, err := convert.ToBool(v)
		fmt.Println(b, err != nil)
	}
	// Output:
	// true false
	// false false
	// true false
	// false true
}

func ExampleToBoolOr() {
	fmt.Println(convert.ToBoolOr("on", false))
	fmt.Println(convert.ToBoolOr("maybe", true))
	// Output:
	// true
	// true
}

func ExampleToInt() {
	for _, v := range []any{" 42 ", 3.9, json.Number("15"), "1.0", "abc"} {
		i, err := convert.ToInt(v)
		fmt.Println(i, err)
	}
	// Output:
	// 42 <nil>
	// 3 <nil>
	// 15 <nil>
	// 1 <nil>
	// 0 convert: strconv.ParseFloat: parsing "abc": invalid syntax
}

func ExampleToInt64() {
	d, _ := convert.ToInt64(time.Millisecond)
	fmt.Println(d)
	_, err := convert.ToInt64(1e19)
	fmt.Println(errors.Is(err, convert.ErrOutOfRange))
	// Output:
	// 1000000
	// true
}

func ExampleToUint64() {
	u, _ := convert.ToUint64("18446744073709551615")
	fmt.Println(u)
	_, err := convert.ToUint64(-1)
	fmt.Println(errors.Is(err, convert.ErrOutOfRange))
	// Output:
	// 18446744073709551615
	// true
}

func ExampleToFloat64() {
	f, _ := convert.ToFloat64("2.5")
	fmt.Println(f)
	_, err := convert.ToFloat64([]int{1})
	fmt.Println(errors.Is(err, convert.ErrUnsupported))
	// Output:
	// 2.5
	// true
}

func ExampleToIntOr() {
	fmt.Println(convert.ToIntOr("8080", 80))
	fmt.Println(convert.ToIntOr("auto", 80))
	// Output:
	// 8080
	// 80
}

func ExampleToInt64Or() {
	fmt.Println(convert.ToInt64Or(int32(5), -1))
	fmt.Println(convert.ToInt64Or(struct{}{}, -1))
	// Output:
	// 5
	// -1
}

func ExampleToFloat64Or() {
	fmt.Println(convert.ToFloat64Or("0.25", 1))
	fmt.Println(convert.ToFloat64Or("n/a", 1))
	// Output:
	// 0.25
	// 1
}
