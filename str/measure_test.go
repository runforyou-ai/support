package str

import "testing"

func TestLength(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"hello", 5},
		{"你好", 2},
		{"héllo", 5},
	}
	for _, tt := range tests {
		if got := Length(tt.s); got != tt.want {
			t.Errorf("Length(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestSubstrCount(t *testing.T) {
	tests := []struct {
		s, sub string
		want   int
	}{
		{"hello world", "o", 2},
		{"aaaa", "aa", 2},
		{"abc", "", 0},
		{"", "a", 0},
		{"你好你好", "你", 2},
	}
	for _, tt := range tests {
		if got := SubstrCount(tt.s, tt.sub); got != tt.want {
			t.Errorf("SubstrCount(%q, %q) = %d, want %d", tt.s, tt.sub, got, tt.want)
		}
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"   ", 0},
		{"Hello, world!", 2},
		{"foo - bar", 2},
		{"  one\ttwo\nthree  ", 3},
		{"你好世界", 1},
		{"don't stop", 2},
	}
	for _, tt := range tests {
		if got := WordCount(tt.s); got != tt.want {
			t.Errorf("WordCount(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}
