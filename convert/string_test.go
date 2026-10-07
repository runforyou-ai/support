package convert

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type level int

type name string

type point struct{ X, Y int }

type stringerPtr struct{}

func (*stringerPtr) String() string { return "ptr" }

func TestToString(t *testing.T) {
	s := "hi"
	n := 42
	var nilInt *int
	var nilStringer *stringerPtr
	var nilTime *time.Time
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"nil", nil, ""},
		{"string", "abc", "abc"},
		{"bytes", []byte("abc"), "abc"},
		{"error", errors.New("boom"), "boom"},
		{"stringer", time.Duration(1500) * time.Millisecond, "1.5s"},
		{"json number", json.Number("12.50"), "12.50"},
		{"bool", true, "true"},
		{"int", -7, "-7"},
		{"int8", int8(8), "8"},
		{"uint64", uint64(18446744073709551615), "18446744073709551615"},
		{"float64", 1.5, "1.5"},
		{"float64 large", 1e21, "1000000000000000000000"},
		{"float32", float32(0.1), "0.1"},
		{"named int", level(3), "3"},
		{"named string", name("ann"), "ann"},
		{"pointer", &s, "hi"},
		{"pointer to int", &n, "42"},
		{"nil pointer", nilInt, ""},
		{"nil stringer pointer", nilStringer, ""},
		{"nil value-receiver pointer", nilTime, ""},
		{"pointer stringer", &stringerPtr{}, "ptr"},
		{"fallback", point{1, 2}, "{1 2}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.in); got != tt.want {
				t.Errorf("ToString(%#v) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}
