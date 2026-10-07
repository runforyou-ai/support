package number

import (
	"math"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		decimals int
		want     string
	}{
		{"zero", 0, 0, "0"},
		{"zero decimals", 0, 2, "0.00"},
		{"integer", 1234567, 0, "1,234,567"},
		{"three digits", 123, 0, "123"},
		{"four digits", 1000, 0, "1,000"},
		{"pad decimals", 1234.5, 2, "1,234.50"},
		{"round down", 1234.564, 2, "1,234.56"},
		{"round half up", 2.675, 2, "2.68"},
		{"round half exact", 0.125, 2, "0.13"},
		{"round half to integer", 2.5, 0, "3"},
		{"round half odd", 1.5, 0, "2"},
		{"carry", 999.995, 2, "1,000.00"},
		{"carry all nines", 9.99, 1, "10.0"},
		{"negative", -1234.567, 2, "-1,234.57"},
		{"negative half away from zero", -2.5, 0, "-3"},
		{"negative rounds to zero", -0.001, 2, "0.00"},
		{"negative zero", math.Copysign(0, -1), 1, "0.0"},
		{"shortest", 1234.5678, -1, "1,234.5678"},
		{"shortest integer", 1e6, -1, "1,000,000"},
		{"small", 0.0000001, 7, "0.0000001"},
		{"large", 1e21, 0, "1,000,000,000,000,000,000,000"},
		{"nan", math.NaN(), 2, "NaN"},
		{"positive infinity", math.Inf(1), 2, "+Inf"},
		{"negative infinity", math.Inf(-1), 2, "-Inf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.v, tt.decimals); got != tt.want {
				t.Errorf("Format(%v, %d) = %q, want %q", tt.v, tt.decimals, got, tt.want)
			}
		})
	}
}

func TestFormatWith(t *testing.T) {
	tests := []struct {
		name         string
		v            float64
		decimals     int
		decimalSep   string
		thousandsSep string
		want         string
	}{
		{"european", 1234567.891, 2, ",", ".", "1.234.567,89"},
		{"space thousands", 1234567.891, 1, ".", " ", "1 234 567.9"},
		{"no thousands", 1234567, 0, ".", "", "1234567"},
		{"multibyte separators", -9876543.21, 2, "·", "’", "-9’876’543·21"},
		{"no fraction", 1000, 0, ",", ".", "1.000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatWith(tt.v, tt.decimals, tt.decimalSep, tt.thousandsSep); got != tt.want {
				t.Errorf("FormatWith() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		v        float64
		decimals int
		want     string
	}{
		{12.5, 2, "12.50%"},
		{12.5, 0, "13%"},
		{0, 0, "0%"},
		{-3.333, 1, "-3.3%"},
		{1250, 1, "1,250.0%"},
		{33.3333, -1, "33.3333%"},
	}
	for _, tt := range tests {
		if got := Percentage(tt.v, tt.decimals); got != tt.want {
			t.Errorf("Percentage(%v, %d) = %q, want %q", tt.v, tt.decimals, got, tt.want)
		}
	}
}

func TestFileSize(t *testing.T) {
	tests := []struct {
		bytes     int64
		precision int
		want      string
	}{
		{0, 0, "0 B"},
		{1, 0, "1 B"},
		{921, 0, "921 B"},
		{1000, 2, "0.98 KB"},
		{1024, 0, "1 KB"},
		{1536, 2, "1.50 KB"},
		{2048, 0, "2 KB"},
		{1024 * 1024, 0, "1 MB"},
		{5 * 1024 * 1024 * 1024, 1, "5.0 GB"},
		{1 << 40, 0, "1 TB"},
		{1 << 50, 0, "1 PB"},
		{1 << 60, 0, "1 EB"},
		{math.MaxInt64, 2, "8.00 EB"},
		{-1536, 2, "-1.50 KB"},
		{-2048, 0, "-2 KB"},
		{math.MinInt64, 0, "-8 EB"},
	}
	for _, tt := range tests {
		if got := FileSize(tt.bytes, tt.precision); got != tt.want {
			t.Errorf("FileSize(%d, %d) = %q, want %q", tt.bytes, tt.precision, got, tt.want)
		}
	}
}
