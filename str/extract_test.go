package str

import "testing"

func TestAfterBefore(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(string, string) string
		s, sep string
		want   string
	}{
		{"After", After, "hannah", "han", "nah"},
		{"After repeated", After, "a/b/c", "/", "b/c"},
		{"After missing", After, "hannah", "xxx", "hannah"},
		{"After empty search", After, "hannah", "", "hannah"},
		{"After unicode", After, "你好，世界", "，", "世界"},
		{"AfterLast", AfterLast, "a/b/c", "/", "c"},
		{"AfterLast missing", AfterLast, "abc", "/", "abc"},
		{"AfterLast empty search", AfterLast, "abc", "", "abc"},
		{"Before", Before, "hannah", "nah", "han"},
		{"Before repeated", Before, "a/b/c", "/", "a"},
		{"Before missing", Before, "hannah", "xxx", "hannah"},
		{"Before empty search", Before, "hannah", "", "hannah"},
		{"BeforeLast", BeforeLast, "a/b/c", "/", "a/b"},
		{"BeforeLast missing", BeforeLast, "abc", "/", "abc"},
		{"BeforeLast empty search", BeforeLast, "abc", "", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.s, tt.sep); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		s, from, to string
		want        string
		wantFirst   string
	}{
		{"[a] [b]", "[", "]", "a] [b", "a"},
		{"foofoobar", "foo", "bar", "foo", "foo"},
		{"abc", "", "c", "abc", "abc"},
		{"abc", "a", "", "abc", "abc"},
		{"abc", "x", "c", "ab", "ab"},
		{"abc", "a", "x", "bc", "bc"},
		{"《你好》", "《", "》", "你好", "你好"},
	}
	for _, tt := range tests {
		if got := Between(tt.s, tt.from, tt.to); got != tt.want {
			t.Errorf("Between(%q, %q, %q) = %q, want %q", tt.s, tt.from, tt.to, got, tt.want)
		}
		if got := BetweenFirst(tt.s, tt.from, tt.to); got != tt.wantFirst {
			t.Errorf("BetweenFirst(%q, %q, %q) = %q, want %q", tt.s, tt.from, tt.to, got, tt.wantFirst)
		}
	}
}

func TestCharAt(t *testing.T) {
	tests := []struct {
		s      string
		index  int
		want   string
		wantOK bool
	}{
		{"hello", 0, "h", true},
		{"hello", 4, "o", true},
		{"hello", -1, "o", true},
		{"hello", -5, "h", true},
		{"hello", 5, "", false},
		{"hello", -6, "", false},
		{"", 0, "", false},
		{"你好世界", 2, "世", true},
	}
	for _, tt := range tests {
		got, ok := CharAt(tt.s, tt.index)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("CharAt(%q, %d) = %q, %v, want %q, %v", tt.s, tt.index, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestSubstr(t *testing.T) {
	tests := []struct {
		s             string
		start, length int
		want          string
	}{
		{"hello world", 0, 5, "hello"},
		{"hello world", 6, 5, "world"},
		{"hello world", 6, 100, "world"},
		{"hello world", -5, 3, "wor"},
		{"hello world", -100, 5, "hello"},
		{"hello world", 0, -6, "hello"},
		{"hello world", 2, -2, "llo wor"},
		{"hello world", 0, 0, ""},
		{"hello world", 11, 1, ""},
		{"hello world", 5, -7, ""},
		{"hello", 1, Length("hello"), "ello"},
		{"你好世界", 1, 2, "好世"},
		{"", 0, 1, ""},
	}
	for _, tt := range tests {
		if got := Substr(tt.s, tt.start, tt.length); got != tt.want {
			t.Errorf("Substr(%q, %d, %d) = %q, want %q", tt.s, tt.start, tt.length, got, tt.want)
		}
	}
}
