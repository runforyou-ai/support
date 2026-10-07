package arr

import (
	"maps"
	"slices"
	"testing"
)

func parity(n int) string {
	if n%2 == 0 {
		return "even"
	}
	return "odd"
}

func TestKeyBy(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want map[string]int
	}{
		{"nil", nil, map[string]int{}},
		{"later wins", []int{1, 2, 3, 4}, map[string]int{"odd": 3, "even": 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeyBy(tt.s, parity)
			if got == nil || !maps.Equal(got, tt.want) {
				t.Fatalf("KeyBy() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want map[string][]int
	}{
		{"nil", nil, map[string][]int{}},
		{"ordered groups", []int{1, 2, 3, 4, 5}, map[string][]int{"odd": {1, 3, 5}, "even": {2, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupBy(tt.s, parity)
			if got == nil || !maps.EqualFunc(got, tt.want, slices.Equal) {
				t.Fatalf("GroupBy() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name        string
		s           ints
		wantMatched ints
		wantRest    ints
	}{
		{"nil", nil, nil, nil},
		{"mixed", ints{1, 2, 3, 4}, ints{2, 4}, ints{1, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, rest := Partition(tt.s, isEven)
			if !slices.Equal(matched, tt.wantMatched) || !slices.Equal(rest, tt.wantRest) {
				t.Fatalf("Partition() = %v, %v; want %v, %v", matched, rest, tt.wantMatched, tt.wantRest)
			}
		})
	}
}
