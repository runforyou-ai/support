package arr

import (
	"slices"
	"strconv"
	"testing"
)

type ints []int

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		s    ints
		want ints
	}{
		{"nil", nil, nil},
		{"none", ints{1, 3}, nil},
		{"some", ints{1, 2, 3, 4}, ints{2, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.s)
			got := Filter(tt.s, isEven)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Filter() = %v; want %v", got, tt.want)
			}
			if !slices.Equal(tt.s, in) {
				t.Fatalf("Filter() mutated input: %v", tt.s)
			}
		})
	}
}

func TestReject(t *testing.T) {
	tests := []struct {
		name string
		s    ints
		want ints
	}{
		{"nil", nil, nil},
		{"all even", ints{2, 4}, nil},
		{"some", ints{1, 2, 3, 4}, ints{1, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reject(tt.s, isEven); !slices.Equal(got, tt.want) {
				t.Fatalf("Reject() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want []string
	}{
		{"nil", nil, nil},
		{"values", []int{1, 2}, []string{"1", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.s, strconv.Itoa); !slices.Equal(got, tt.want) {
				t.Fatalf("Map() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestFilterMap(t *testing.T) {
	parse := func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		return n, err == nil
	}
	tests := []struct {
		name string
		s    []string
		want []int
	}{
		{"nil", nil, nil},
		{"none", []string{"a"}, nil},
		{"some", []string{"1", "x", "3"}, []int{1, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterMap(tt.s, parse); !slices.Equal(got, tt.want) {
				t.Fatalf("FilterMap() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	concat := func(acc string, n int) string { return acc + strconv.Itoa(n) }
	tests := []struct {
		name string
		s    []int
		init string
		want string
	}{
		{"nil", nil, ">", ">"},
		{"values", []int{1, 2, 3}, ">", ">123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.s, tt.init, concat); got != tt.want {
				t.Fatalf("Reduce() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name string
		s    [][]int
		want []int
	}{
		{"nil", nil, nil},
		{"empty inner", [][]int{nil, {}}, nil},
		{"values", [][]int{{1, 2}, nil, {3}}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Flatten(tt.s); !slices.Equal(got, tt.want) {
				t.Fatalf("Flatten() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestCrossJoin(t *testing.T) {
	tests := []struct {
		name  string
		lists [][]int
		want  [][]int
	}{
		{"no lists", nil, nil},
		{"empty list", [][]int{{1, 2}, {}}, nil},
		{"single", [][]int{{1, 2}}, [][]int{{1}, {2}}},
		{"three", [][]int{{1, 2}, {3}, {4, 5}}, [][]int{{1, 3, 4}, {1, 3, 5}, {2, 3, 4}, {2, 3, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CrossJoin(tt.lists...)
			if !slices.EqualFunc(got, tt.want, slices.Equal) {
				t.Fatalf("CrossJoin() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestPad(t *testing.T) {
	tests := []struct {
		name string
		s    ints
		size int
		want ints
	}{
		{"nil zero", nil, 0, nil},
		{"nil right", nil, 2, ints{0, 0}},
		{"right", ints{1, 2}, 4, ints{1, 2, 0, 0}},
		{"left", ints{1, 2}, -4, ints{0, 0, 1, 2}},
		{"shorter", ints{1, 2, 3}, 2, ints{1, 2, 3}},
		{"shorter left", ints{1, 2, 3}, -3, ints{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.s)
			got := Pad(tt.s, tt.size, 0)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Pad() = %v; want %v", got, tt.want)
			}
			if len(got) > 0 && len(tt.s) > 0 && &got[0] == &tt.s[0] {
				t.Fatal("Pad() returned the input slice")
			}
			if !slices.Equal(tt.s, in) {
				t.Fatalf("Pad() mutated input: %v", tt.s)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name      string
		items     []string
		glue      string
		finalGlue string
		want      string
	}{
		{"nil", nil, ", ", " and ", ""},
		{"one", []string{"a"}, ", ", " and ", "a"},
		{"two", []string{"a", "b"}, ", ", " and ", "a and b"},
		{"three", []string{"a", "b", "c"}, ", ", " and ", "a, b and c"},
		{"empty final glue", []string{"a", "b", "c"}, "-", "", "a-b-c"},
		{"multibyte", []string{"甲", "乙", "丙"}, "、", "和", "甲、乙和丙"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.items, tt.glue, tt.finalGlue); got != tt.want {
				t.Fatalf("Join() = %q; want %q", got, tt.want)
			}
		})
	}
}
