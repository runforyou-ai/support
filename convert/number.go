package convert

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// numKind identifies which field of a number holds its value.
type numKind int

const (
	kindInt numKind = iota
	kindUint
	kindFloat
)

// number is a parsed numeric value.
type number struct {
	kind numKind
	i    int64
	u    uint64
	f    float64
	// wide marks an integer literal outside the int64 and uint64 ranges,
	// held approximately in f.
	wide bool
}

// parseNumber extracts the numeric value of v.
func parseNumber(v any) (number, error) {
	switch x := v.(type) {
	case nil:
		return number{}, nil
	case bool:
		if x {
			return number{i: 1}, nil
		}
		return number{}, nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return number{}, nil
		}
		return parseNumber(rv.Elem().Interface())
	case reflect.Bool:
		return parseNumber(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return number{kind: kindInt, i: rv.Int()}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return number{kind: kindUint, u: rv.Uint()}, nil
	case reflect.Float32, reflect.Float64:
		return number{kind: kindFloat, f: rv.Float()}, nil
	case reflect.String:
		return parseString(rv.String())
	}
	return number{}, fmt.Errorf("%w: %T to number", ErrUnsupported, v)
}

// parseString parses a trimmed decimal integer or float string.
func parseString(s string) (number, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return number{}, nil
	}
	i, intErr := strconv.ParseInt(s, 10, 64)
	if intErr == nil {
		return number{kind: kindInt, i: i}, nil
	}
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return number{kind: kindUint, u: u}, nil
	}
	// An integer literal that fits neither int64 nor uint64 is only exact as a float.
	wide := errors.Is(intErr, strconv.ErrRange)
	// Hexadecimal literals and digit separators are not decimal input.
	if strings.ContainsAny(s, "xX_") {
		return number{}, fmt.Errorf("convert: %w", &strconv.NumError{Func: "ParseFloat", Num: s, Err: strconv.ErrSyntax})
	}
	f, err := strconv.ParseFloat(s, 64)
	if errors.Is(err, strconv.ErrRange) {
		return number{}, fmt.Errorf("%w: %w", ErrOutOfRange, err)
	}
	if err != nil {
		return number{}, fmt.Errorf("convert: %w", err)
	}
	return number{kind: kindFloat, f: f, wide: wide}, nil
}

// ToInt64 converts v to an int64. It accepts every integer and float kind,
// including named types such as time.Duration (as nanoseconds), booleans
// (1 or 0), json.Number and decimal strings, which are trimmed and may be
// written as floats such as "1.0" or "2e3". Floats and float strings are
// truncated toward zero. Pointers are dereferenced, and nil, nil pointers and
// empty strings convert to 0. Values that do not fit, including NaN and
// infinities, return an error wrapping ErrOutOfRange; malformed strings
// return an error wrapping a *strconv.NumError; other types return an error
// wrapping ErrUnsupported.
func ToInt64(v any) (int64, error) {
	n, err := parseNumber(v)
	if err != nil {
		return 0, err
	}
	switch n.kind {
	case kindUint:
		if n.u > math.MaxInt64 {
			return 0, fmt.Errorf("%w: %v overflows int64", ErrOutOfRange, v)
		}
		return int64(n.u), nil
	case kindFloat:
		f := math.Trunc(n.f)
		if n.wide || math.IsNaN(f) || f < math.MinInt64 || f >= math.MaxInt64 {
			return 0, fmt.Errorf("%w: %v overflows int64", ErrOutOfRange, v)
		}
		return int64(f), nil
	}
	return n.i, nil
}

// ToInt converts v to an int following the rules of ToInt64, returning an
// error wrapping ErrOutOfRange when the value does not fit in an int.
func ToInt(v any) (int, error) {
	i, err := ToInt64(v)
	if err != nil {
		return 0, err
	}
	if int64(int(i)) != i {
		return 0, fmt.Errorf("%w: %v overflows int", ErrOutOfRange, v)
	}
	return int(i), nil
}

// ToUint64 converts v to a uint64 following the rules of ToInt64. Negative
// values return an error wrapping ErrOutOfRange.
func ToUint64(v any) (uint64, error) {
	n, err := parseNumber(v)
	if err != nil {
		return 0, err
	}
	switch n.kind {
	case kindInt:
		if n.i < 0 {
			return 0, fmt.Errorf("%w: %v is negative", ErrOutOfRange, v)
		}
		return uint64(n.i), nil
	case kindFloat:
		f := math.Trunc(n.f)
		if n.wide || math.IsNaN(f) || f < 0 || f >= math.MaxUint64 {
			return 0, fmt.Errorf("%w: %v overflows uint64", ErrOutOfRange, v)
		}
		return uint64(f), nil
	}
	return n.u, nil
}

// ToFloat64 converts v to a float64 following the input rules of ToInt64.
// Integers beyond 2^53 may lose precision, float32 values are widened
// exactly, and strings such as "NaN" and "Inf" are accepted. Strings whose
// magnitude exceeds float64 return an error wrapping ErrOutOfRange.
func ToFloat64(v any) (float64, error) {
	n, err := parseNumber(v)
	if err != nil {
		return 0, err
	}
	switch n.kind {
	case kindUint:
		return float64(n.u), nil
	case kindFloat:
		return n.f, nil
	}
	return float64(n.i), nil
}

// ToIntOr converts v to an int like ToInt, returning fallback when the
// conversion fails.
func ToIntOr(v any, fallback int) int {
	i, err := ToInt(v)
	if err != nil {
		return fallback
	}
	return i
}

// ToInt64Or converts v to an int64 like ToInt64, returning fallback when the
// conversion fails.
func ToInt64Or(v any, fallback int64) int64 {
	i, err := ToInt64(v)
	if err != nil {
		return fallback
	}
	return i
}

// ToFloat64Or converts v to a float64 like ToFloat64, returning fallback
// when the conversion fails.
func ToFloat64Or(v any, fallback float64) float64 {
	f, err := ToFloat64(v)
	if err != nil {
		return fallback
	}
	return f
}
