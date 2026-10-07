package support

import (
	"strings"
	"testing"
)

type name string

type point struct{ X, Y int }

func TestBlankAndFilled(t *testing.T) {
	var (
		nilPtr   *int
		nilSlice []int
		nilMap   map[string]int
		nilChan  chan int
		nilFunc  func()
		nilError error
	)
	tests := []struct {
		name  string
		v     any
		blank bool
	}{
		{"nil", nil, true},
		{"nil error interface", nilError, true},
		{"typed nil pointer", nilPtr, true},
		{"nil slice", nilSlice, true},
		{"nil map", nilMap, true},
		{"nil chan", nilChan, true},
		{"nil func", nilFunc, true},
		{"empty string", "", true},
		{"whitespace string", " \t\n", true},
		{"string", " a ", false},
		{"empty named string", name("  "), true},
		{"named string", name("bob"), false},
		{"empty slice", []int{}, true},
		{"slice", []int{0}, false},
		{"empty map", map[string]int{}, true},
		{"map", map[string]int{"a": 0}, false},
		{"empty chan", make(chan int, 1), true},
		{"empty array", [0]int{}, true},
		{"array", [2]int{}, false},
		{"zero int", 0, false},
		{"int", 7, false},
		{"zero float", 0.0, false},
		{"false", false, false},
		{"true", true, false},
		{"zero struct", point{}, false},
		{"struct", point{1, 2}, false},
		{"pointer", Ptr(0), false},
		{"func", func() {}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Blank(tt.v); got != tt.blank {
				t.Errorf("Blank(%#v) = %v, want %v", tt.v, got, tt.blank)
			}
			if got := Filled(tt.v); got == tt.blank {
				t.Errorf("Filled(%#v) = %v, want %v", tt.v, got, !tt.blank)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	strs := []struct {
		name   string
		values []string
		want   string
	}{
		{"no values", nil, ""},
		{"all zero", []string{"", ""}, ""},
		{"first", []string{"a", "b"}, "a"},
		{"skips zero", []string{"", "b", "c"}, "b"},
		{"whitespace is not zero", []string{" ", "b"}, " "},
	}
	for _, tt := range strs {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(tt.values...); got != tt.want {
				t.Errorf("Default(%q) = %q, want %q", tt.values, got, tt.want)
			}
		})
	}
	ints := []struct {
		values []int
		want   int
	}{
		{[]int{0, 0, 3}, 3},
		{[]int{-1, 2}, -1},
		{[]int{0}, 0},
	}
	for _, tt := range ints {
		if got := Default(tt.values...); got != tt.want {
			t.Errorf("Default(%v) = %d, want %d", tt.values, got, tt.want)
		}
	}
	if got := Default(point{}, point{1, 2}); got != (point{1, 2}) {
		t.Errorf("Default(structs) = %v, want {1 2}", got)
	}
}

func TestPtr(t *testing.T) {
	v := 5
	p := Ptr(v)
	if p == &v || *p != 5 {
		t.Fatalf("Ptr(5) = %p (%d), want a new pointer to 5", p, *p)
	}
	*p = 6
	if v != 5 {
		t.Errorf("Ptr mutated its argument: v = %d", v)
	}
	if got := *Ptr("s"); got != "s" {
		t.Errorf("*Ptr(%q) = %q", "s", got)
	}
}

func TestDeref(t *testing.T) {
	tests := []struct {
		name string
		p    *int
		want int
	}{
		{"nil", nil, 0},
		{"zero", Ptr(0), 0},
		{"value", Ptr(42), 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Deref(tt.p); got != tt.want {
				t.Errorf("Deref() = %d, want %d", got, tt.want)
			}
		})
	}
	if got := Deref[[]int](nil); got != nil {
		t.Errorf("Deref[[]int](nil) = %v, want nil", got)
	}
}

func TestDerefOr(t *testing.T) {
	tests := []struct {
		name     string
		p        *string
		fallback string
		want     string
	}{
		{"nil", nil, "fallback", "fallback"},
		{"empty value is kept", Ptr(""), "fallback", ""},
		{"value", Ptr("v"), "fallback", "v"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DerefOr(tt.p, tt.fallback); got != tt.want {
				t.Errorf("DerefOr() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTap(t *testing.T) {
	tests := []struct {
		name string
		v    []int
	}{
		{"nil", nil},
		{"values", []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			var seen []int
			got := Tap(tt.v, func(v []int) {
				calls++
				seen = v
			})
			if calls != 1 || len(seen) != len(tt.v) || len(got) != len(tt.v) {
				t.Errorf("Tap() = %v after %d calls with %v, want %v once", got, calls, seen, tt.v)
			}
		})
	}
	m := map[string]int{}
	got := Tap(m, func(m map[string]int) { m["a"] = 1 })
	if got["a"] != 1 {
		t.Errorf("Tap returned %v, want the tapped map", got)
	}
}

func TestWith(t *testing.T) {
	tests := []struct {
		v    string
		want int
	}{
		{"", 0},
		{"abc", 3},
	}
	for _, tt := range tests {
		if got := With(tt.v, func(s string) int { return len(s) }); got != tt.want {
			t.Errorf("With(%q, len) = %d, want %d", tt.v, got, tt.want)
		}
	}
}

func TestTransform(t *testing.T) {
	upper := func(s string) string { return strings.ToUpper(s) }
	tests := []struct {
		name string
		v    string
		want string
	}{
		{"filled", "abc", "ABC"},
		{"empty", "", "none"},
		{"whitespace", "  ", "none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Transform(tt.v, upper, "none"); got != tt.want {
				t.Errorf("Transform(%q) = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
	called := false
	if got := Transform([]int(nil), func([]int) int { called = true; return 1 }, -1); got != -1 || called {
		t.Errorf("Transform(nil slice) = %d (called %v), want -1 without calling fn", got, called)
	}
	if got := Transform(0, func(n int) int { return n + 1 }, -1); got != 1 {
		t.Errorf("Transform(0) = %d, want 1", got)
	}
}
