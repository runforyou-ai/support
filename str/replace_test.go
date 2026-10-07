package str

import "testing"

func TestReplaceFirstLast(t *testing.T) {
	tests := []struct {
		s, search, replace string
		first, last        string
	}{
		{"the quick brown fox jumps over the lazy dog", "the", "a", "a quick brown fox jumps over the lazy dog", "the quick brown fox jumps over a lazy dog"},
		{"abc", "x", "y", "abc", "abc"},
		{"abc", "", "y", "abc", "abc"},
		{"你好你好", "你", "您", "您好你好", "你好您好"},
	}
	for _, tt := range tests {
		if got := ReplaceFirst(tt.s, tt.search, tt.replace); got != tt.first {
			t.Errorf("ReplaceFirst(%q, %q, %q) = %q, want %q", tt.s, tt.search, tt.replace, got, tt.first)
		}
		if got := ReplaceLast(tt.s, tt.search, tt.replace); got != tt.last {
			t.Errorf("ReplaceLast(%q, %q, %q) = %q, want %q", tt.s, tt.search, tt.replace, got, tt.last)
		}
	}
}

func TestReplaceStartEnd(t *testing.T) {
	tests := []struct {
		s, search, replace string
		start, end         string
	}{
		{"Hello World", "Hello", "Laravel", "Laravel World", "Hello World"},
		{"Hello World", "World", "Laravel", "Hello World", "Hello Laravel"},
		{"aXa", "a", "b", "bXa", "aXb"},
		{"abc", "", "x", "abc", "abc"},
	}
	for _, tt := range tests {
		if got := ReplaceStart(tt.s, tt.search, tt.replace); got != tt.start {
			t.Errorf("ReplaceStart(%q, %q, %q) = %q, want %q", tt.s, tt.search, tt.replace, got, tt.start)
		}
		if got := ReplaceEnd(tt.s, tt.search, tt.replace); got != tt.end {
			t.Errorf("ReplaceEnd(%q, %q, %q) = %q, want %q", tt.s, tt.search, tt.replace, got, tt.end)
		}
	}
}

func TestReplaceArray(t *testing.T) {
	tests := []struct {
		s, search string
		replaces  []string
		want      string
	}{
		{"The event will take place between ? and ?", "?", []string{"8:30", "9:00"}, "The event will take place between 8:30 and 9:00"},
		{"? ? ?", "?", []string{"a"}, "a ? ?"},
		{"? ?", "?", []string{"a", "b", "c"}, "a b"},
		{"? ?", "?", nil, "? ?"},
		{"no match", "?", []string{"a"}, "no match"},
		{"abc", "", []string{"x"}, "abc"},
		{"", "?", []string{"x"}, ""},
	}
	for _, tt := range tests {
		replaces := append([]string(nil), tt.replaces...)
		if got := ReplaceArray(tt.s, tt.search, tt.replaces); got != tt.want {
			t.Errorf("ReplaceArray(%q, %q, %q) = %q, want %q", tt.s, tt.search, tt.replaces, got, tt.want)
		}
		for i := range replaces {
			if replaces[i] != tt.replaces[i] {
				t.Errorf("ReplaceArray mutated replaces")
			}
		}
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		s     string
		pairs map[string]string
		want  string
	}{
		{"Tacos are great!", map[string]string{"Tacos": "Burritos", "great": "fantastic"}, "Burritos are fantastic!"},
		{"ab", map[string]string{"a": "b", "b": "a"}, "ba"},
		{"abc", map[string]string{"a": "1", "ab": "2"}, "2c"},
		{"abc", map[string]string{"": "x", "c": "3"}, "ab3"},
		{"abc", nil, "abc"},
		{"你好世界", map[string]string{"你好": "再见"}, "再见世界"},
		{"", map[string]string{"a": "b"}, ""},
	}
	for _, tt := range tests {
		if got := Swap(tt.s, tt.pairs); got != tt.want {
			t.Errorf("Swap(%q, %v) = %q, want %q", tt.s, tt.pairs, got, tt.want)
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		s      string
		search []string
		want   string
	}{
		{"Peter Piper picked a peck", []string{"e"}, "Ptr Pipr pickd a pck"},
		{"Peter Piper", []string{"e", "P"}, "tr ipr"},
		{"abc", []string{""}, "abc"},
		{"abc", nil, "abc"},
		{"你好世界", []string{"好"}, "你世界"},
	}
	for _, tt := range tests {
		if got := Remove(tt.s, tt.search...); got != tt.want {
			t.Errorf("Remove(%q, %q) = %q, want %q", tt.s, tt.search, got, tt.want)
		}
	}
}
