package convert

import (
	"fmt"
	"reflect"
	"strconv"
)

// ToString returns v formatted as a string. It returns strings and byte
// slices as they are, uses the Error or String method of errors and
// fmt.Stringer values (such as time.Duration and json.Number), formats
// booleans and numbers with strconv (floats in the shortest decimal form
// without an exponent), dereferences pointers, and returns "" for nil and
// nil pointers. Any other value is formatted with fmt.Sprint.
func ToString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case []byte:
		return string(x)
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		return ""
	}
	switch x := v.(type) {
	case error:
		return x.Error()
	case fmt.Stringer:
		return x.String()
	}
	switch rv.Kind() {
	case reflect.Pointer:
		return ToString(rv.Elem().Interface())
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64)
	case reflect.String:
		return rv.String()
	}
	return fmt.Sprint(v)
}
