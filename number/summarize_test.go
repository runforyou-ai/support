package number

import (
	"math"
	"testing"
)

func TestAbbreviate(t *testing.T) {
	tests := []struct {
		v         float64
		precision int
		want      string
	}{
		{0, 0, "0"},
		{0, 2, "0"},
		{1, 0, "1"},
		{0.5, 1, "0.5"},
		{0.5, 0, "1"},
		{999, 0, "999"},
		{999.5, 0, "1K"},
		{1000, 0, "1K"},
		{1000, 2, "1K"},
		{1200, 1, "1.2K"},
		{1234, 2, "1.23K"},
		{1200000, 0, "1M"},
		{1200000, 1, "1.2M"},
		{1200000, 2, "1.2M"},
		{1234567, -1, "1.234567M"},
		{999999, 0, "1M"},
		{999999, 3, "999.999K"},
		{1e9, 0, "1B"},
		{1.5e12, 1, "1.5T"},
		{999.9999e12, 2, "1Q"},
		{1e15, 0, "1Q"},
		{2.5e15, 1, "2.5Q"},
		{1e18, 0, "1KQ"},
		{-1234, 1, "-1.2K"},
		{-0.1, 0, "0"},
		{math.NaN(), 0, "NaN"},
		{math.Inf(-1), 0, "-Inf"},
	}
	for _, tt := range tests {
		if got := Abbreviate(tt.v, tt.precision); got != tt.want {
			t.Errorf("Abbreviate(%v, %d) = %q, want %q", tt.v, tt.precision, got, tt.want)
		}
	}
}

func TestForHumans(t *testing.T) {
	tests := []struct {
		v         float64
		precision int
		want      string
	}{
		{0, 0, "0"},
		{123, 0, "123"},
		{1000, 0, "1 thousand"},
		{1500000, 1, "1.5 million"},
		{1e9, 0, "1 billion"},
		{1.25e12, 2, "1.25 trillion"},
		{1e15, 0, "1 quadrillion"},
		{1e18, 0, "1 thousand quadrillion"},
		{-1500, 1, "-1.5 thousand"},
	}
	for _, tt := range tests {
		if got := ForHumans(tt.v, tt.precision); got != tt.want {
			t.Errorf("ForHumans(%v, %d) = %q, want %q", tt.v, tt.precision, got, tt.want)
		}
	}
}
