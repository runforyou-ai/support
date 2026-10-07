package arr

import (
	"slices"
	"testing"
)

func TestRandom(t *testing.T) {
	tests := []struct {
		name   string
		s      []int
		wantOK bool
	}{
		{"nil", nil, false},
		{"one", []int{7}, true},
		{"many", []int{1, 2, 3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Random(tt.s)
			if ok != tt.wantOK {
				t.Fatalf("Random() ok = %v; want %v", ok, tt.wantOK)
			}
			if ok && !slices.Contains(tt.s, got) {
				t.Fatalf("Random() = %v; not in %v", got, tt.s)
			}
			if !ok && got != 0 {
				t.Fatalf("Random() = %v; want zero", got)
			}
		})
	}
}

func TestSample(t *testing.T) {
	tests := []struct {
		name    string
		s       ints
		n       int
		wantLen int
	}{
		{"nil", nil, 3, 0},
		{"zero", ints{1, 2}, 0, 0},
		{"negative", ints{1, 2}, -1, 0},
		{"some", ints{1, 2, 3, 4, 5}, 3, 3},
		{"clamped", ints{1, 2, 3}, 10, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.s)
			got := Sample(tt.s, tt.n)
			if len(got) != tt.wantLen {
				t.Fatalf("Sample() len = %d; want %d", len(got), tt.wantLen)
			}
			if len(Unique(got)) != len(got) {
				t.Fatalf("Sample() = %v; positions not distinct", got)
			}
			for _, v := range got {
				if !slices.Contains(tt.s, v) {
					t.Fatalf("Sample() = %v; %v not in input", got, v)
				}
			}
			if !slices.Equal(tt.s, in) {
				t.Fatalf("Sample() mutated input: %v", tt.s)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name string
		s    ints
	}{
		{"nil", nil},
		{"many", ints{1, 2, 3, 4, 5, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.s)
			got := Shuffle(tt.s)
			if !slices.Equal(tt.s, in) {
				t.Fatalf("Shuffle() mutated input: %v", tt.s)
			}
			sorted := slices.Sorted(slices.Values(got))
			if !slices.Equal(sorted, slices.Sorted(slices.Values(in))) {
				t.Fatalf("Shuffle() = %v; not a permutation of %v", got, in)
			}
		})
	}
}
