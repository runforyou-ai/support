package convert

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"
)

type numCase struct {
	name    string
	in      any
	want    any
	wantErr error
}

func checkErr(t *testing.T, fn string, in any, got any, err error, tt numCase) {
	t.Helper()
	if tt.wantErr != nil {
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("%s(%#v) error = %v; want %v", fn, in, err, tt.wantErr)
		}
		return
	}
	if err != nil || got != tt.want {
		t.Errorf("%s(%#v) = %v, %v; want %v", fn, in, got, err, tt.want)
	}
}

func commonCases() []numCase {
	seven := 7
	var nilInt *int
	return []numCase{
		{"nil", nil, 0, nil},
		{"true", true, 1, nil},
		{"false", false, 0, nil},
		{"named bool", flag(true), 1, nil},
		{"int", 42, 42, nil},
		{"int8", int8(8), 8, nil},
		{"uint16", uint16(16), 16, nil},
		{"named int", level(3), 3, nil},
		{"duration", 2 * time.Nanosecond, 2, nil},
		{"float integral", 3.0, 3, nil},
		{"float truncated", 3.9, 3, nil},
		{"float32", float32(2.5), 2, nil},
		{"string", " 12 ", 12, nil},
		{"string plus", "+5", 5, nil},
		{"string float", "1.0", 1, nil},
		{"string float truncated", "2.7", 2, nil},
		{"string exponent", "2e3", 2000, nil},
		{"string empty", "", 0, nil},
		{"named string", name("9"), 9, nil},
		{"json number", json.Number("15"), 15, nil},
		{"pointer", &seven, 7, nil},
		{"nil pointer", nilInt, 0, nil},
		{"string bad", "abc", nil, strconv.ErrSyntax},
		{"string hex", "0x10", nil, strconv.ErrSyntax},
		{"string hex float", "0x1p4", nil, strconv.ErrSyntax},
		{"string underscore", "1_000", nil, strconv.ErrSyntax},
		{"string float overflow", "1e400", nil, ErrOutOfRange},
		{"float NaN", math.NaN(), nil, ErrOutOfRange},
		{"float Inf", math.Inf(1), nil, ErrOutOfRange},
		{"unsupported", []int{1}, nil, ErrUnsupported},
		{"complex", complex(1, 0), nil, ErrUnsupported},
	}
}

func TestToInt64(t *testing.T) {
	cases := append(commonCases(),
		numCase{"negative", -3, -3, nil},
		numCase{"negative float", -3.9, -3, nil},
		numCase{"min int64", int64(math.MinInt64), int64(math.MinInt64), nil},
		numCase{"min float", float64(math.MinInt64), int64(math.MinInt64), nil},
		numCase{"uint max int64", uint64(math.MaxInt64), int64(math.MaxInt64), nil},
		numCase{"uint overflow", uint64(math.MaxInt64) + 1, nil, ErrOutOfRange},
		numCase{"float overflow", 1e19, nil, ErrOutOfRange},
		numCase{"float underflow", -1e19, nil, ErrOutOfRange},
		numCase{"string max", "9223372036854775807", int64(math.MaxInt64), nil},
		numCase{"string uint overflow", "9223372036854775808", nil, ErrOutOfRange},
		numCase{"string huge", "99999999999999999999", nil, ErrOutOfRange},
	)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if w, ok := tt.want.(int); ok {
				tt.want = int64(w)
			}
			got, err := ToInt64(tt.in)
			checkErr(t, "ToInt64", tt.in, got, err, tt)
		})
	}
}

func TestToInt(t *testing.T) {
	cases := append(commonCases(),
		numCase{"negative", "-12", -12, nil},
		numCase{"uint overflow", uint64(math.MaxUint64), nil, ErrOutOfRange},
	)
	if strconv.IntSize == 32 {
		cases = append(cases, numCase{"int64 overflow", int64(math.MaxInt64), nil, ErrOutOfRange})
	} else {
		cases = append(cases, numCase{"int64 max", int64(math.MaxInt64), math.MaxInt, nil})
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.in)
			checkErr(t, "ToInt", tt.in, got, err, tt)
		})
	}
}

func TestToUint64(t *testing.T) {
	cases := append(commonCases(),
		numCase{"max", uint64(math.MaxUint64), uint64(math.MaxUint64), nil},
		numCase{"string max", "18446744073709551615", uint64(math.MaxUint64), nil},
		numCase{"string overflow", "18446744073709551616", nil, ErrOutOfRange},
		numCase{"negative int", -1, nil, ErrOutOfRange},
		numCase{"negative string", "-1", nil, ErrOutOfRange},
		numCase{"negative float", -1.5, nil, ErrOutOfRange},
		numCase{"small negative float", -0.5, uint64(0), nil},
		numCase{"float overflow", 2e19, nil, ErrOutOfRange},
	)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if w, ok := tt.want.(int); ok {
				tt.want = uint64(w)
			}
			got, err := ToUint64(tt.in)
			checkErr(t, "ToUint64", tt.in, got, err, tt)
		})
	}
}

func TestToFloat64(t *testing.T) {
	cases := []numCase{
		{"nil", nil, 0.0, nil},
		{"true", true, 1.0, nil},
		{"int", -3, -3.0, nil},
		{"uint", uint64(math.MaxUint64), float64(math.MaxUint64), nil},
		{"float", 2.5, 2.5, nil},
		{"float32", float32(0.5), 0.5, nil},
		{"string", " 1.25 ", 1.25, nil},
		{"string int", "7", 7.0, nil},
		{"string inf", "Inf", math.Inf(1), nil},
		{"json number", json.Number("-0.75"), -0.75, nil},
		{"duration", time.Microsecond, 1000.0, nil},
		{"string bad", "1.2.3", nil, strconv.ErrSyntax},
		{"string overflow", "1e400", nil, ErrOutOfRange},
		{"unsupported", map[string]int{}, nil, ErrUnsupported},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.in)
			checkErr(t, "ToFloat64", tt.in, got, err, tt)
		})
	}
	if got, err := ToFloat64("NaN"); err != nil || !math.IsNaN(got) {
		t.Errorf("ToFloat64(NaN) = %v, %v", got, err)
	}
}

func TestParseErrorWrapsNumError(t *testing.T) {
	_, err := ToInt("abc")
	var ne *strconv.NumError
	if !errors.As(err, &ne) || ne.Num != "abc" {
		t.Errorf("ToInt(abc) error = %v; want *strconv.NumError", err)
	}
}

func TestOrFallbacks(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ToIntOr ok", ToIntOr("5", 1), 5},
		{"ToIntOr fallback", ToIntOr("x", 1), 1},
		{"ToInt64Or ok", ToInt64Or(5.5, 1), int64(5)},
		{"ToInt64Or fallback", ToInt64Or(math.NaN(), 1), int64(1)},
		{"ToFloat64Or ok", ToFloat64Or("0.5", 1), 0.5},
		{"ToFloat64Or fallback", ToFloat64Or([]int{}, 1), 1.0},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %v; want %v", tt.name, tt.got, tt.want)
		}
	}
}
