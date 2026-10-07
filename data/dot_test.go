package data

import (
	"reflect"
	"testing"
)

func TestDot(t *testing.T) {
	tests := []struct {
		name string
		in   m
		want m
	}{
		{"nil", nil, m{}},
		{"flat", m{"a": 1}, m{"a": 1}},
		{"nested", m{"a": m{"b": m{"c": 1}}, "d": 2}, m{"a.b.c": 1, "d": 2}},
		{"slice", m{"list": []any{"x", m{"y": 1}}}, m{"list.0": "x", "list.1.y": 1}},
		{"empty containers", m{"a": m{}, "b": []any{}}, m{"a": m{}, "b": []any{}}},
		{"typed kept", m{"t": []int{1}}, m{"t": []int{1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dot(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dot = %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestUndot(t *testing.T) {
	tests := []struct {
		name string
		in   m
		want m
	}{
		{"nil", nil, m{}},
		{"flat", m{"a": 1}, m{"a": 1}},
		{"nested", m{"a.b.c": 1, "a.d": 2, "e": 3}, m{"a": m{"b": m{"c": 1}, "d": 2}, "e": 3}},
		{"numeric segments", m{"list.0": "x", "list.1": "y"}, m{"list": m{"0": "x", "1": "y"}}},
		{"prefix conflict", m{"a": 1, "a.b": 2}, m{"a": m{"b": 2}}},
		{"merge into map value", m{"a": m{"x": 1}, "a.y": 2}, m{"a": m{"x": 1, "y": 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Undot(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Undot = %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestUndotDoesNotMutate(t *testing.T) {
	inner := m{"x": 1}
	in := m{"a": inner, "a.y": 2, "a.z.w": 3}
	Undot(in)
	if !reflect.DeepEqual(inner, m{"x": 1}) {
		t.Errorf("Undot mutated nested input map: %v", inner)
	}
}

func TestDotRoundTrip(t *testing.T) {
	in := m{"a": m{"b": 1, "c": m{"d": "x"}}, "e": true}
	if got := Undot(Dot(in)); !reflect.DeepEqual(got, in) {
		t.Errorf("Undot(Dot) = %#v; want %#v", got, in)
	}
}
