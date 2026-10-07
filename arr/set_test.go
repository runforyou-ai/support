package arr

import (
	"slices"
	"strings"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name string
		s    ints
		want ints
	}{
		{"nil", nil, nil},
		{"no duplicates", ints{3, 1, 2}, ints{3, 1, 2}},
		{"duplicates", ints{3, 1, 3, 2, 1}, ints{3, 1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.s); !slices.Equal(got, tt.want) {
				t.Fatalf("Unique() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueBy(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []string
	}{
		{"nil", nil, nil},
		{"case insensitive", []string{"Go", "go", "Rust", "GO", "rust"}, []string{"Go", "Rust"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueBy(tt.s, strings.ToLower); !slices.Equal(got, tt.want) {
				t.Fatalf("UniqueBy() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name   string
		s      ints
		others [][]int
		want   ints
	}{
		{"nil", nil, [][]int{{1}}, nil},
		{"no others", ints{1, 2}, nil, ints{1, 2}},
		{"several others", ints{1, 2, 3, 2, 4}, [][]int{{3}, {4, 5}}, ints{1, 2, 2}},
		{"all removed", ints{1}, [][]int{{1}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.s, tt.others...); !slices.Equal(got, tt.want) {
				t.Fatalf("Diff() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name   string
		s      ints
		others [][]int
		want   ints
	}{
		{"nil", nil, [][]int{{1}}, nil},
		{"no others", ints{1, 2}, nil, ints{1, 2}},
		{"several others", ints{1, 2, 3, 2, 4}, [][]int{{2, 3, 4}, {4, 2}}, ints{2, 2, 4}},
		{"empty other", ints{1, 2}, [][]int{nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersect(tt.s, tt.others...); !slices.Equal(got, tt.want) {
				t.Fatalf("Intersect() = %v; want %v", got, tt.want)
			}
		})
	}
}
