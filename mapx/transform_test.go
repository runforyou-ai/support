package mapx

import (
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name  string
		m     scores
		key   string
		value int
		want  scores
	}{
		{"nil map", nil, "a", 1, scores{"a": 1}},
		{"absent", scores{"a": 1}, "b", 2, scores{"a": 1, "b": 2}},
		{"present", scores{"a": 1}, "a", 5, scores{"a": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := maps.Clone(tt.m)
			checkResult(t, "Add", Add(tt.m, tt.key, tt.value), tt.want, tt.m, orig)
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name string
		ms   []scores
		want scores
	}{
		{"no maps", nil, scores{}},
		{"nil maps", []scores{nil, nil}, scores{}},
		{"later wins", []scores{{"a": 1, "b": 2}, nil, {"b": 3, "c": 4}}, scores{"a": 1, "b": 3, "c": 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var first, orig scores
			if len(tt.ms) > 0 {
				first, orig = tt.ms[0], maps.Clone(tt.ms[0])
			}
			checkResult(t, "Merge", Merge(tt.ms...), tt.want, first, orig)
		})
	}
}

func TestMapValues(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want map[string]string
	}{
		{"nil map", nil, map[string]string{}},
		{"values", scores{"a": 1, "b": 2}, map[string]string{"a": "a=1", "b": "b=2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapValues(tt.m, func(k string, v int) string { return k + "=" + strconv.Itoa(v) })
			if got == nil || !maps.Equal(got, tt.want) {
				t.Fatalf("MapValues() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want scores
	}{
		{"nil map", nil, scores{}},
		{"values", scores{"a": 1, "b": 2}, scores{"A": 1, "B": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapKeys(tt.m, func(k string, _ int) string { return strings.ToUpper(k) })
			if got == nil || !maps.Equal(got, map[string]int(tt.want)) {
				t.Fatalf("MapKeys() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestMapKeysCollision(t *testing.T) {
	got := MapKeys(scores{"a": 1, "A": 2}, func(k string, _ int) string { return strings.ToLower(k) })
	if len(got) != 1 || (got["a"] != 1 && got["a"] != 2) {
		t.Fatalf("MapKeys() = %v; want one of the colliding values", got)
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want map[int]string
	}{
		{"nil map", nil, map[int]string{}},
		{"values", scores{"a": 1, "b": 2}, map[int]string{1: "a", 2: "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Invert(tt.m)
			if got == nil || !maps.Equal(got, tt.want) {
				t.Fatalf("Invert() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestSortedKeys(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want []string
	}{
		{"nil map", nil, nil},
		{"values", scores{"c": 3, "a": 1, "b": 2}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortedKeys(tt.m); !slices.Equal(got, tt.want) {
				t.Fatalf("SortedKeys() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestSortedEntries(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want []Entry[string, int]
	}{
		{"nil map", nil, nil},
		{"values", scores{"c": 3, "a": 1, "b": 2}, []Entry[string, int]{{"a", 1}, {"b", 2}, {"c", 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortedEntries(tt.m); !slices.Equal(got, tt.want) {
				t.Fatalf("SortedEntries() = %v; want %v", got, tt.want)
			}
		})
	}
}
