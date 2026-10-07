package mapx

import (
	"maps"
	"testing"
)

type scores map[string]int

func checkResult(t *testing.T, fn string, got, want, in, orig scores) {
	t.Helper()
	if got == nil || !maps.Equal(got, want) {
		t.Fatalf("%s() = %v; want %v", fn, got, want)
	}
	if !maps.Equal(in, orig) {
		t.Fatalf("%s() mutated input: %v", fn, in)
	}
}

func TestOnly(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		keys []string
		want scores
	}{
		{"nil map", nil, []string{"a"}, scores{}},
		{"no keys", scores{"a": 1}, nil, scores{}},
		{"some keys", scores{"a": 1, "b": 2, "c": 3}, []string{"a", "c", "x"}, scores{"a": 1, "c": 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := maps.Clone(tt.m)
			checkResult(t, "Only", Only(tt.m, tt.keys...), tt.want, tt.m, orig)
		})
	}
}

func TestExcept(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		keys []string
		want scores
	}{
		{"nil map", nil, []string{"a"}, scores{}},
		{"no keys", scores{"a": 1}, nil, scores{"a": 1}},
		{"some keys", scores{"a": 1, "b": 2, "c": 3}, []string{"a", "x"}, scores{"b": 2, "c": 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := maps.Clone(tt.m)
			checkResult(t, "Except", Except(tt.m, tt.keys...), tt.want, tt.m, orig)
		})
	}
}

func positive(_ string, v int) bool { return v > 0 }

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want scores
	}{
		{"nil map", nil, scores{}},
		{"mixed", scores{"a": 1, "b": -2, "c": 3}, scores{"a": 1, "c": 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := maps.Clone(tt.m)
			checkResult(t, "Filter", Filter(tt.m, positive), tt.want, tt.m, orig)
		})
	}
}

func TestReject(t *testing.T) {
	tests := []struct {
		name string
		m    scores
		want scores
	}{
		{"nil map", nil, scores{}},
		{"mixed", scores{"a": 1, "b": -2, "c": 3}, scores{"b": -2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := maps.Clone(tt.m)
			checkResult(t, "Reject", Reject(tt.m, positive), tt.want, tt.m, orig)
		})
	}
}
