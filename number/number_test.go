package number

import (
	"math"
	"reflect"
	"testing"
)

func TestOrdinal(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0th"},
		{1, "1st"},
		{2, "2nd"},
		{3, "3rd"},
		{4, "4th"},
		{10, "10th"},
		{11, "11th"},
		{12, "12th"},
		{13, "13th"},
		{21, "21st"},
		{22, "22nd"},
		{23, "23rd"},
		{101, "101st"},
		{111, "111th"},
		{112, "112th"},
		{1003, "1003rd"},
		{-1, "-1st"},
		{-12, "-12th"},
		{-22, "-22nd"},
		{math.MinInt64, "-9223372036854775808th"},
	}
	for _, tt := range tests {
		if got := Ordinal(tt.n); got != tt.want {
			t.Errorf("Ordinal(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestClamp(t *testing.T) {
	ints := []struct {
		v, lo, hi, want int
	}{
		{5, 1, 10, 5},
		{-5, 1, 10, 1},
		{15, 1, 10, 10},
		{1, 1, 10, 1},
		{10, 1, 10, 10},
		{15, 10, 1, 10},
		{0, 10, 1, 1},
		{3, 3, 3, 3},
	}
	for _, tt := range ints {
		if got := Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
	floats := []struct {
		v, lo, hi, want float64
	}{
		{0.5, 0, 1, 0.5},
		{1.5, 0, 1, 1},
		{-0.5, 1, 0, 0},
	}
	for _, tt := range floats {
		if got := Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%v, %v, %v) = %v, want %v", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
	strs := []struct {
		v, lo, hi, want string
	}{
		{"m", "b", "y", "m"},
		{"a", "b", "y", "b"},
		{"z", "y", "b", "y"},
	}
	for _, tt := range strs {
		if got := Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%q, %q, %q) = %q, want %q", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
}

func TestPairs(t *testing.T) {
	tests := []struct {
		name          string
		to, by, start int
		want          [][2]int
	}{
		{"zero based", 25, 10, 0, [][2]int{{0, 9}, {10, 19}, {20, 25}}},
		{"one based", 25, 10, 1, [][2]int{{1, 10}, {11, 20}, {21, 25}}},
		{"exact one based", 20, 10, 1, [][2]int{{1, 10}, {11, 20}}},
		{"exact zero based", 20, 10, 0, [][2]int{{0, 9}, {10, 19}}},
		{"last single", 21, 10, 0, [][2]int{{0, 9}, {10, 19}, {20, 21}}},
		{"by one", 3, 1, 0, [][2]int{{0, 0}, {1, 1}, {2, 2}}},
		{"by larger than range", 5, 100, 0, [][2]int{{0, 5}}},
		{"negative start", 5, 5, -5, [][2]int{{-5, -1}, {0, 4}}},
		{"start equals to", 10, 5, 10, nil},
		{"start after to", 10, 5, 20, nil},
		{"zero by", 10, 0, 0, nil},
		{"negative by", 10, -1, 0, nil},
		{"near max int", math.MaxInt, 5, math.MaxInt - 7, [][2]int{{math.MaxInt - 7, math.MaxInt - 3}, {math.MaxInt - 2, math.MaxInt}}},
		{"full int range", math.MaxInt, math.MaxInt, math.MinInt, [][2]int{{math.MinInt, -2}, {-1, math.MaxInt - 2}, {math.MaxInt - 1, math.MaxInt}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pairs(tt.to, tt.by, tt.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pairs(%d, %d, %d) = %v, want %v", tt.to, tt.by, tt.start, got, tt.want)
			}
		})
	}
}
