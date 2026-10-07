package set

import (
	"maps"
	"slices"
	"testing"
)

func check(t *testing.T, fn string, got Set[int], want []int) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s() = nil; want a non-nil set", fn)
	}
	if !slices.Equal(Sorted(got), want) {
		t.Fatalf("%s() = %v; want %v", fn, Sorted(got), want)
	}
}

func TestOfAndCollect(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"nil", nil, nil},
		{"empty", []int{}, nil},
		{"duplicates", []int{3, 1, 3, 2, 1}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.values)
			check(t, "Of", Of(tt.values...), tt.want)
			check(t, "Collect", Collect(tt.values), tt.want)
			if !slices.Equal(tt.values, in) {
				t.Fatalf("Collect() mutated input: %v", tt.values)
			}
		})
	}
}

func TestCollectBy(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want []int
	}{
		{"nil", nil, nil},
		{"lengths", []string{"a", "bb", "c", "ddd"}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check(t, "CollectBy", CollectBy(tt.s, func(s string) int { return len(s) }), tt.want)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name  string
		s     Set[int]
		v     int
		added bool
		want  []int
	}{
		{"nil set", nil, 1, true, []int{1}},
		{"absent", Of(1), 2, true, []int{1, 2}},
		{"present", Of(1, 2), 2, false, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.s
			if got := s.Add(tt.v); got != tt.added {
				t.Fatalf("Add(%d) = %v; want %v", tt.v, got, tt.added)
			}
			check(t, "Add", s, tt.want)
		})
	}
	var holder struct{ ids Set[string] }
	if !holder.ids.Add("a") || holder.ids.Add("a") || !holder.ids.Has("a") {
		t.Fatalf("Add on a zero struct field: %v", holder.ids)
	}
}

func TestHasAndLen(t *testing.T) {
	tests := []struct {
		name string
		s    Set[int]
		v    int
		has  bool
		len  int
	}{
		{"nil", nil, 1, false, 0},
		{"absent", Of(2, 3), 1, false, 2},
		{"present", Of(1, 2, 3), 1, true, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Has(tt.v); got != tt.has {
				t.Errorf("Has(%d) = %v; want %v", tt.v, got, tt.has)
			}
			if got := tt.s.Len(); got != tt.len {
				t.Errorf("Len() = %d; want %d", got, tt.len)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name string
		s    Set[int]
		v    int
		want []int
	}{
		{"nil", nil, 1, nil},
		{"absent", Of(2), 1, []int{2}},
		{"present", Of(1, 2), 1, []int{2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Delete(tt.v)
			if !slices.Equal(Sorted(tt.s), tt.want) {
				t.Fatalf("Delete(%d) left %v; want %v", tt.v, Sorted(tt.s), tt.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		name string
		s    Set[int]
		want []int
	}{
		{"nil", nil, nil},
		{"values", Of(1, 2), []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Clone()
			check(t, "Clone", got, tt.want)
			got.Add(99)
			if tt.s.Has(99) {
				t.Fatal("Clone() shares storage with the receiver")
			}
		})
	}
}

func TestSetOperations(t *testing.T) {
	tests := []struct {
		name      string
		s, o      Set[int]
		union     []int
		intersect []int
		diff      []int
	}{
		{"both nil", nil, nil, nil, nil, nil},
		{"nil receiver", nil, Of(1), []int{1}, nil, nil},
		{"nil argument", Of(1), nil, []int{1}, nil, []int{1}},
		{"overlap", Of(1, 2, 3), Of(2, 3, 4, 5), []int{1, 2, 3, 4, 5}, []int{2, 3}, []int{1}},
		{"larger receiver", Of(1, 2, 3, 4), Of(4, 5), []int{1, 2, 3, 4, 5}, []int{4}, []int{1, 2, 3}},
		{"disjoint", Of(1), Of(2), []int{1, 2}, nil, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, o := maps.Clone(tt.s), maps.Clone(tt.o)
			check(t, "Union", tt.s.Union(tt.o), tt.union)
			check(t, "Intersect", tt.s.Intersect(tt.o), tt.intersect)
			check(t, "Diff", tt.s.Diff(tt.o), tt.diff)
			if !maps.Equal(tt.s, s) || !maps.Equal(tt.o, o) {
				t.Fatalf("operations mutated their operands: %v, %v", tt.s, tt.o)
			}
		})
	}
}

func TestValuesAndSorted(t *testing.T) {
	tests := []struct {
		name string
		s    Set[string]
		want []string
	}{
		{"nil", nil, nil},
		{"empty", Set[string]{}, nil},
		{"values", Of("c", "a", "b"), []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := tt.s.Values()
			if (values == nil) != (tt.want == nil) {
				t.Fatalf("Values() = %#v; want nil=%v", values, tt.want == nil)
			}
			slices.Sort(values)
			if !slices.Equal(values, tt.want) {
				t.Fatalf("Values() = %v; want %v", values, tt.want)
			}
			if got := Sorted(tt.s); !slices.Equal(got, tt.want) || (got == nil) != (tt.want == nil) {
				t.Fatalf("Sorted() = %#v; want %#v", got, tt.want)
			}
		})
	}
}
