package support

import (
	"reflect"
	"strings"
)

// Blank reports whether v is nil, a nil pointer or interface, a string that
// is empty or contains only whitespace, or an empty slice, map, array or
// channel. Numbers and booleans are never blank, including 0 and false.
func Blank(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return strings.TrimSpace(rv.String()) == ""
	case reflect.Slice, reflect.Map, reflect.Chan:
		return rv.IsNil() || rv.Len() == 0
	case reflect.Array:
		return rv.Len() == 0
	case reflect.Pointer, reflect.Interface, reflect.Func:
		return rv.IsNil()
	}
	return false
}

// Filled reports whether v is not blank. It is the inverse of Blank.
func Filled(v any) bool {
	return !Blank(v)
}

// NilIfZero returns nil when v is the zero value of T, and a pointer to a
// copy of v otherwise.
func NilIfZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// MapPtr returns nil when p is nil, and a pointer to fn(*p) otherwise. fn is
// not called for a nil pointer.
func MapPtr[T, R any](p *T, fn func(T) R) *R {
	if p == nil {
		return nil
	}
	r := fn(*p)
	return &r
}

// Deref returns the value p points to, or the zero value of T when p is nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOr returns the value p points to, or fallback when p is nil.
func DerefOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// Tap calls fn with v and returns v.
func Tap[T any](v T, fn func(T)) T {
	fn(v)
	return v
}

// With returns the result of calling fn with v.
func With[T, R any](v T, fn func(T) R) R {
	return fn(v)
}

// Transform returns fn(v) when v is filled, and fallback when v is blank.
func Transform[T, R any](v T, fn func(T) R, fallback R) R {
	if Blank(v) {
		return fallback
	}
	return fn(v)
}
