package data

import (
	"reflect"
	"testing"
)

func sample() map[string]any {
	return map[string]any{
		"user": map[string]any{
			"name":  "Ann",
			"tags":  []any{"a", "b"},
			"empty": nil,
		},
		"users": []any{
			map[string]any{"name": "Ann", "pets": []any{map[string]any{"name": "Rex"}}},
			map[string]any{"name": "Bob", "pets": []any{map[string]any{"name": "Tom"}, map[string]any{"name": "Kit"}}},
			map[string]any{"email": "c@example.com"},
		},
		"scores": map[string]int{"math": 90},
		"grades": map[string]int{"b": 2, "a": 1},
		"ptr":    &map[string]any{"k": []any{"v"}},
		"ids":    []int{7, 8},
		"pair":   [2]string{"x", "y"},
		"levels": map[level]any{"low": 1},
		"byInt":  map[int]any{1: "one"},
		"count":  3,
	}
}

type level string

func TestGet(t *testing.T) {
	m := sample()
	tests := []struct {
		name   string
		target any
		path   string
		want   any
		ok     bool
	}{
		{"empty path", 5, "", 5, true},
		{"top level", m, "count", 3, true},
		{"nested", m, "user.name", "Ann", true},
		{"slice index", m, "user.tags.1", "b", true},
		{"present nil", m, "user.empty", nil, true},
		{"missing", m, "user.age", nil, false},
		{"index out of range", m, "user.tags.2", nil, false},
		{"non-numeric index", m, "user.tags.x", nil, false},
		{"signed index", m, "user.tags.+1", nil, false},
		{"empty index", m, "user.tags.", nil, false},
		{"huge index", m, "user.tags.99999999999999999999", nil, false},
		{"scalar parent", m, "count.x", nil, false},
		{"nil target", nil, "a", nil, false},
		{"typed map", m, "scores.math", 90, true},
		{"typed map missing", m, "scores.art", nil, false},
		{"named key map", m, "levels.low", 1, true},
		{"non-string key map", m, "byInt.1", nil, false},
		{"typed slice", m, "ids.1", 8, true},
		{"typed slice out of range", m, "ids.5", nil, false},
		{"array", m, "pair.0", "x", true},
		{"wildcard values", m, "user.tags.*", []any{"a", "b"}, true},
		{"wildcard field", m, "users.*.name", []any{"Ann", "Bob"}, true},
		{"wildcard flatten", m, "users.*.pets.*.name", []any{"Rex", "Tom", "Kit"}, true},
		{"wildcard map sorted", map[string]any{"b": 2, "a": 1}, "*", []any{1, 2}, true},
		{"wildcard typed map", m, "scores.*", []any{90}, true},
		{"wildcard typed map sorted", m, "grades.*", []any{1, 2}, true},
		{"through pointer", m, "ptr.k.0", "v", true},
		{"wildcard through pointer", m, "ptr.*", []any{[]any{"v"}}, true},
		{"nil pointer", map[string]any{"p": (*map[string]any)(nil)}, "p.k", nil, false},
		{"wildcard nil pointer", map[string]any{"p": (*map[string]any)(nil)}, "p.*", nil, false},
		{"wildcard typed slice", m, "ids.*", []any{7, 8}, true},
		{"wildcard non-string key map", m, "byInt.*", nil, false},
		{"wildcard no match", m, "users.*.age", nil, false},
		{"wildcard scalar", m, "count.*", nil, false},
		{"wildcard nil", nil, "*", nil, false},
		{"wildcard empty", map[string]any{"a": []any{}}, "a.*", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Get(tt.target, tt.path)
			if ok != tt.ok || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get(%q) = %v, %v; want %v, %v", tt.path, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestGetOr(t *testing.T) {
	m := sample()
	tests := []struct {
		path     string
		fallback any
		want     any
	}{
		{"user.name", "x", "Ann"},
		{"user.age", 30, 30},
		{"user.empty", "x", nil},
	}
	for _, tt := range tests {
		if got := GetOr(m, tt.path, tt.fallback); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetOr(%q) = %v; want %v", tt.path, got, tt.want)
		}
	}
}

func TestGetAs(t *testing.T) {
	m := sample()
	if got, ok := GetAs[string](m, "user.name"); !ok || got != "Ann" {
		t.Errorf("GetAs[string] = %q, %v", got, ok)
	}
	tests := []struct {
		name string
		path string
	}{
		{"wrong type", "count"},
		{"missing", "user.age"},
		{"nil value", "user.empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := GetAs[string](m, tt.path); ok || got != "" {
				t.Errorf("GetAs[string](%q) = %q, %v; want zero, false", tt.path, got, ok)
			}
		})
	}
	if got, ok := GetAs[[]any](m, "users.*.name"); !ok || len(got) != 2 {
		t.Errorf("GetAs[[]any] wildcard = %v, %v", got, ok)
	}
}

func TestHas(t *testing.T) {
	m := sample()
	tests := []struct {
		path string
		want bool
	}{
		{"", true},
		{"user", true},
		{"user.empty", true},
		{"user.tags.0", true},
		{"user.tags.9", false},
		{"users.*.email", true},
		{"users.*.phone", false},
		{"nope", false},
	}
	for _, tt := range tests {
		if got := Has(m, tt.path); got != tt.want {
			t.Errorf("Has(%q) = %v; want %v", tt.path, got, tt.want)
		}
	}
}
