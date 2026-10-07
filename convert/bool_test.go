package convert

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"
)

type flag bool

func TestToBool(t *testing.T) {
	yes := "yes"
	var nilStr *string
	tests := []struct {
		name    string
		in      any
		want    bool
		wantErr error
	}{
		{"nil", nil, false, nil},
		{"true", true, true, nil},
		{"named bool", flag(true), true, nil},
		{"int zero", 0, false, nil},
		{"int nonzero", -2, true, nil},
		{"uint", uint8(1), true, nil},
		{"float", 0.5, true, nil},
		{"float zero", 0.0, false, nil},
		{"json number", json.Number("2"), true, nil},
		{"json number zero", json.Number("0.0"), false, nil},
		{"json number bad", json.Number("x"), false, strconv.ErrSyntax},
		{"string true", " TRUE ", true, nil},
		{"string on", "On", true, nil},
		{"string y", "y", true, nil},
		{"string off", "off", false, nil},
		{"string empty", "  ", false, nil},
		{"string n", "N", false, nil},
		{"named string", name("1"), true, nil},
		{"string bad", "maybe", false, strconv.ErrSyntax},
		{"pointer", &yes, true, nil},
		{"nil pointer", nilStr, false, nil},
		{"unsupported", []int{1}, false, ErrUnsupported},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.in)
			if got != tt.want || !errors.Is(err, tt.wantErr) || (tt.wantErr == nil) != (err == nil) {
				t.Errorf("ToBool(%#v) = %v, %v; want %v, %v", tt.in, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestToBoolOr(t *testing.T) {
	tests := []struct {
		in       any
		fallback bool
		want     bool
	}{
		{"yes", false, true},
		{"no", true, false},
		{"maybe", true, true},
		{struct{}{}, false, false},
	}
	for _, tt := range tests {
		if got := ToBoolOr(tt.in, tt.fallback); got != tt.want {
			t.Errorf("ToBoolOr(%#v, %v) = %v; want %v", tt.in, tt.fallback, got, tt.want)
		}
	}
}
