package data

import (
	"reflect"
	"testing"
)

type m = map[string]any

func setFixture() m {
	return m{
		"name":  "Ann",
		"nil":   nil,
		"tags":  []any{"a", m{"k": 1}},
		"users": []any{m{"name": "Ann"}, m{"name": "Bob", "age": 3}, "loose"},
		"group": m{"x": m{"v": 1}, "y": 2},
		"typed": map[string]int{"a": 1},
		"bytes": []byte("hi"),
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		value any
		check string
		want  any
	}{
		{"overwrite", "name", "Bob", "name", "Bob"},
		{"new key", "age", 3, "age", 3},
		{"create intermediates", "a.b.c", 1, "a", m{"b": m{"c": 1}}},
		{"replace scalar intermediate", "name.first", "A", "name", m{"first": "A"}},
		{"replace nil intermediate", "nil.x", 1, "nil", m{"x": 1}},
		{"slice index", "tags.0", "z", "tags", []any{"z", m{"k": 1}}},
		{"into slice element", "tags.1.k", 2, "tags", []any{"a", m{"k": 2}}},
		{"slice out of range", "tags.5", "z", "tags", []any{"a", m{"k": 1}}},
		{"slice non-numeric", "tags.x.y", "z", "tags", []any{"a", m{"k": 1}}},
		{"slice element scalar replaced", "tags.0.k", 9, "tags", []any{m{"k": 9}, m{"k": 1}}},
		{"wildcard slice field", "users.*.active", true, "users", []any{
			m{"name": "Ann", "active": true}, m{"name": "Bob", "age": 3, "active": true}, m{"active": true},
		}},
		{"wildcard slice leaf", "tags.*", 0, "tags", []any{0, 0}},
		{"wildcard map leaf", "group.*", 0, "group", m{"x": 0, "y": 0}},
		{"wildcard map field", "group.*.v", 5, "group", m{"x": m{"v": 5}, "y": m{"v": 5}}},
		{"wildcard missing", "missing.*.x", 1, "missing", nil},
		{"wildcard below new key", "a.*", 1, "a", nil},
		{"typed map untouched", "typed.a", 9, "typed", map[string]int{"a": 1}},
		{"bytes untouched", "bytes.0", 9, "bytes", []byte("hi")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := setFixture()
			Set(target, tt.path, tt.value)
			if got := target[tt.check]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("after Set(%q) %s = %#v; want %#v", tt.path, tt.check, got, tt.want)
			}
		})
	}
}

func TestSetNoop(t *testing.T) {
	Set(nil, "a", 1)
	Fill(nil, "a", 1)
	Forget(nil, "a")
	target := setFixture()
	Set(target, "", 1)
	Fill(target, "", 1)
	Forget(target, "")
	if !reflect.DeepEqual(target, setFixture()) {
		t.Errorf("empty path changed target: %v", target)
	}
}

func TestFill(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		value any
		check string
		want  any
	}{
		{"keep existing", "name", "Bob", "name", "Ann"},
		{"keep nil", "nil", 1, "nil", nil},
		{"new key", "age", 3, "age", 3},
		{"create intermediates", "a.b", 1, "a", m{"b": 1}},
		{"keep scalar intermediate", "name.first", "A", "name", "Ann"},
		{"keep nil intermediate", "nil.x", 1, "nil", nil},
		{"wildcard fills missing", "users.*.age", 1, "users", []any{
			m{"name": "Ann", "age": 1}, m{"name": "Bob", "age": 3}, "loose",
		}},
		{"wildcard leaf keeps", "tags.*", 0, "tags", []any{"a", m{"k": 1}}},
		{"wildcard map leaf keeps", "group.*", 0, "group", m{"x": m{"v": 1}, "y": 2}},
		{"slice index keeps", "tags.0", "z", "tags", []any{"a", m{"k": 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := setFixture()
			Fill(target, tt.path, tt.value)
			if got := target[tt.check]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("after Fill(%q) %s = %#v; want %#v", tt.path, tt.check, got, tt.want)
			}
		})
	}
}

func TestForget(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		check string
		want  any
	}{
		{"top level", "name", "name", nil},
		{"nested", "group.x.v", "group", m{"x": m{}, "y": 2}},
		{"missing", "group.z.v", "group", m{"x": m{"v": 1}, "y": 2}},
		{"slice element kept", "tags.0", "tags", []any{"a", m{"k": 1}}},
		{"inside slice element", "tags.1.k", "tags", []any{"a", m{}}},
		{"slice bad index", "tags.9.k", "tags", []any{"a", m{"k": 1}}},
		{"wildcard slice", "users.*.name", "users", []any{m{}, m{"age": 3}, "loose"}},
		{"wildcard map", "group.*.v", "group", m{"x": m{}, "y": 2}},
		{"trailing wildcard", "group.*", "group", m{"x": m{"v": 1}, "y": 2}},
		{"wildcard scalar", "name.*.x", "name", "Ann"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := setFixture()
			Forget(target, tt.path)
			if got := target[tt.check]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("after Forget(%q) %s = %#v; want %#v", tt.path, tt.check, got, tt.want)
			}
		})
	}
	target := setFixture()
	Forget(target, "name")
	if _, ok := target["name"]; ok {
		t.Error("Forget did not delete key")
	}
}
