package data

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Query encodes m as a URL query string in the style of PHP's
// http_build_query. Nested maps and slices are written with bracketed keys,
// such as a[b]=c and list[0]=x, with the brackets percent-encoded
// (a%5Bb%5D=c). Keys are sorted, spaces are encoded as "+", booleans are
// written as 1 and 0, pointers are dereferenced, and nil values and empty
// containers are omitted.
func Query(m map[string]any) string {
	var parts []string
	for _, k := range sortedKeys(m) {
		parts = appendQuery(parts, k, m[k])
	}
	return strings.Join(parts, "&")
}

// appendQuery appends the encoded pairs for v under key to parts.
func appendQuery(parts []string, key string, v any) []string {
	// Nil pointers are skipped, and other pointers without String or Error methods are dereferenced.
	for {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Pointer {
			break
		}
		if rv.IsNil() {
			return parts
		}
		if isTextual(v) {
			break
		}
		v = rv.Elem().Interface()
	}
	if !isTextual(v) {
		if keys, vals, ok := entries(v); ok {
			for i, k := range keys {
				parts = appendQuery(parts, key+"["+k+"]", vals[i])
			}
			return parts
		}
	}
	s, ok := scalar(v)
	if !ok {
		return parts
	}
	return append(parts, url.QueryEscape(key)+"="+url.QueryEscape(s))
}

// isTextual reports whether v encodes as a single value even when it is a
// container: byte slices, Stringers and errors.
func isTextual(v any) bool {
	switch v.(type) {
	case []byte, fmt.Stringer, error:
		return true
	}
	return false
}

// scalar formats v as a query value and reports false for nil values.
func scalar(v any) (string, bool) {
	switch x := v.(type) {
	case nil:
		return "", false
	case string:
		return x, true
	case []byte:
		return string(x), true
	case bool:
		if x {
			return "1", true
		}
		return "0", true
	}
	switch x := v.(type) {
	case fmt.Stringer:
		return x.String(), true
	case error:
		return x.Error(), true
	}
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), true
	case reflect.Bool:
		return scalar(rv.Bool())
	case reflect.String:
		return rv.String(), true
	}
	return fmt.Sprint(v), true
}
