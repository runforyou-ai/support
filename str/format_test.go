package str

import "testing"

func TestMask(t *testing.T) {
	tests := []struct {
		s             string
		char          rune
		index, length int
		want          string
	}{
		{"taylor@email.com", '*', 3, 0, "tay*************"},
		{"taylor@email.com", '*', 0, 6, "******@email.com"},
		{"taylor@email.com", '*', -15, 3, "t***or@email.com"},
		{"taylor@email.com", '*', -4, 0, "taylor@email****"},
		{"taylor@email.com", '*', -100, 2, "**ylor@email.com"},
		{"taylor", '*', 6, 2, "taylor"},
		{"taylor", '*', 100, 0, "taylor"},
		{"", '*', 0, 0, ""},
		{"你好世界", '＊', 1, 2, "你＊＊界"},
	}
	for _, tt := range tests {
		if got := Mask(tt.s, tt.char, tt.index, tt.length); got != tt.want {
			t.Errorf("Mask(%q, %q, %d, %d) = %q, want %q", tt.s, tt.char, tt.index, tt.length, got, tt.want)
		}
	}
}

func TestPad(t *testing.T) {
	tests := []struct {
		s                 string
		length            int
		pad               string
		left, right, both string
	}{
		{"James", 10, "-", "-----James", "James-----", "--James---"},
		{"James", 10, "-=", "-=-=-James", "James-=-=-", "-=James-=-"},
		{"James", 3, "-", "James", "James", "James"},
		{"James", 10, "", "James", "James", "James"},
		{"你好", 5, "·", "···你好", "你好···", "·你好··"},
		{"", 2, "ab", "ab", "ab", "aa"},
	}
	for _, tt := range tests {
		if got := PadLeft(tt.s, tt.length, tt.pad); got != tt.left {
			t.Errorf("PadLeft(%q, %d, %q) = %q, want %q", tt.s, tt.length, tt.pad, got, tt.left)
		}
		if got := PadRight(tt.s, tt.length, tt.pad); got != tt.right {
			t.Errorf("PadRight(%q, %d, %q) = %q, want %q", tt.s, tt.length, tt.pad, got, tt.right)
		}
		if got := PadBoth(tt.s, tt.length, tt.pad); got != tt.both {
			t.Errorf("PadBoth(%q, %d, %q) = %q, want %q", tt.s, tt.length, tt.pad, got, tt.both)
		}
	}
}

func TestStartFinish(t *testing.T) {
	tests := []struct {
		s, affix      string
		start, finish string
	}{
		{"this/string", "/", "/this/string", "this/string/"},
		{"//this/string//", "/", "/this/string//", "//this/string/"},
		{"abcabcx", "abc", "abcx", "abcabcxabc"},
		{"xabcabc", "abc", "abcxabcabc", "xabc"},
		{"value", "", "value", "value"},
		{"", "/", "/", "/"},
	}
	for _, tt := range tests {
		if got := Start(tt.s, tt.affix); got != tt.start {
			t.Errorf("Start(%q, %q) = %q, want %q", tt.s, tt.affix, got, tt.start)
		}
		if got := Finish(tt.s, tt.affix); got != tt.finish {
			t.Errorf("Finish(%q, %q) = %q, want %q", tt.s, tt.affix, got, tt.finish)
		}
	}
}

func TestWrapUnwrap(t *testing.T) {
	tests := []struct {
		s, before, after string
		wrapped          string
		unwrapped        string
	}{
		{"value", `"`, "", `"value"`, "value"},
		{"value", "[", "]", "[value]", "value"},
		{"value", "", "", "value", "value"},
		{"value", "", "!", "value!", "value"},
	}
	for _, tt := range tests {
		got := Wrap(tt.s, tt.before, tt.after)
		if got != tt.wrapped {
			t.Errorf("Wrap(%q, %q, %q) = %q, want %q", tt.s, tt.before, tt.after, got, tt.wrapped)
		}
		if back := Unwrap(got, tt.before, tt.after); back != tt.unwrapped {
			t.Errorf("Unwrap(%q, %q, %q) = %q, want %q", got, tt.before, tt.after, back, tt.unwrapped)
		}
	}

	unwrap := []struct {
		s, before, after, want string
	}{
		{"[value", "[", "]", "value"},
		{"value]", "[", "]", "value"},
		{"value", "[", "]", "value"},
		{`""`, `"`, "", ""},
	}
	for _, tt := range unwrap {
		if got := Unwrap(tt.s, tt.before, tt.after); got != tt.want {
			t.Errorf("Unwrap(%q, %q, %q) = %q, want %q", tt.s, tt.before, tt.after, got, tt.want)
		}
	}
}

func TestSquish(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"   laravel   php  framework   ", "laravel php framework"},
		{"\tfoo\n\nbar baz　qux", "foo bar baz qux"},
		{"aㅤbᅠc", "a b c"},
	}
	for _, tt := range tests {
		if got := Squish(tt.in); got != tt.want {
			t.Errorf("Squish(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"hello", "olleh"},
		{"你好世界", "界世好你"},
	}
	for _, tt := range tests {
		if got := Reverse(tt.in); got != tt.want {
			t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSlug(t *testing.T) {
	tests := []struct {
		in, sep, want string
	}{
		{"Hello World", "-", "hello-world"},
		{"  --Hello,   World!!  ", "-", "hello-world"},
		{"Laravel 11 Framework", "_", "laravel_11_framework"},
		{"你好 世界", "-", "你好-世界"},
		{"Über Café", "-", "über-café"},
		{"a b", "", "ab"},
		{"!!!", "-", ""},
		{"", "-", ""},
	}
	for _, tt := range tests {
		if got := Slug(tt.in, tt.sep); got != tt.want {
			t.Errorf("Slug(%q, %q) = %q, want %q", tt.in, tt.sep, got, tt.want)
		}
	}
}
