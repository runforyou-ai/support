package mapx

import "testing"

func TestGetOr(t *testing.T) {
	tests := []struct {
		name     string
		m        scores
		key      string
		fallback int
		want     int
	}{
		{"nil map", nil, "a", 9, 9},
		{"present", scores{"a": 1}, "a", 9, 1},
		{"present zero", scores{"a": 0}, "a", 9, 0},
		{"missing", scores{"a": 1}, "b", 9, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOr(tt.m, tt.key, tt.fallback); got != tt.want {
				t.Fatalf("GetOr() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestHas(t *testing.T) {
	m := scores{"a": 1, "b": 2}
	tests := []struct {
		name string
		m    scores
		keys []string
		want bool
	}{
		{"nil map", nil, []string{"a"}, false},
		{"no keys", m, nil, false},
		{"all present", m, []string{"a", "b"}, true},
		{"one missing", m, []string{"a", "x"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Has(tt.m, tt.keys...); got != tt.want {
				t.Fatalf("Has() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestHasAny(t *testing.T) {
	m := scores{"a": 1, "b": 2}
	tests := []struct {
		name string
		m    scores
		keys []string
		want bool
	}{
		{"nil map", nil, []string{"a"}, false},
		{"no keys", m, nil, false},
		{"one present", m, []string{"x", "b"}, true},
		{"none present", m, []string{"x", "y"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAny(tt.m, tt.keys...); got != tt.want {
				t.Fatalf("HasAny() = %v; want %v", got, tt.want)
			}
		})
	}
}
