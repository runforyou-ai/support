package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ToBool converts v to a bool. Numbers are true when non-zero; json.Number
// values are parsed as numbers. Strings are trimmed and compared case
// insensitively: "1", "t", "true", "yes", "y" and "on" are true, and "0",
// "f", "false", "no", "n", "off" and "" are false; any other string returns
// an error wrapping strconv.ErrSyntax. Pointers are dereferenced, and nil
// and nil pointers are false. Other types return an error wrapping
// ErrUnsupported.
func ToBool(v any) (bool, error) {
	switch x := v.(type) {
	case nil:
		return false, nil
	case bool:
		return x, nil
	case json.Number:
		f, err := ToFloat64(x)
		return f != 0, err
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return false, nil
		}
		return ToBool(rv.Elem().Interface())
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0, nil
	case reflect.String:
		switch strings.ToLower(strings.TrimSpace(rv.String())) {
		case "1", "t", "true", "yes", "y", "on":
			return true, nil
		case "0", "f", "false", "no", "n", "off", "":
			return false, nil
		}
		return false, fmt.Errorf("convert: parse bool %q: %w", rv.String(), strconv.ErrSyntax)
	}
	return false, fmt.Errorf("%w: %T to bool", ErrUnsupported, v)
}

// ToBoolOr converts v to a bool like ToBool, returning fallback when the
// conversion fails.
func ToBoolOr(v any, fallback bool) bool {
	b, err := ToBool(v)
	if err != nil {
		return fallback
	}
	return b
}
