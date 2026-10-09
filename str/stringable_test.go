package str

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

var _ fmt.Stringer = Stringable{}

func TestStringableTransforms(t *testing.T) {
	tests := []struct {
		name string
		got  Stringable
		want string
	}{
		{"Of", Of("value"), "value"},
		{"zero", Stringable{}, ""},
		{"After", Of("a/b/c").After("/"), "b/c"},
		{"AfterLast", Of("a/b/c").AfterLast("/"), "c"},
		{"Before", Of("a/b/c").Before("/"), "a"},
		{"BeforeLast", Of("a/b/c").BeforeLast("/"), "a/b"},
		{"Between", Of("[a] [b]").Between("[", "]"), "a] [b"},
		{"BetweenFirst", Of("[a] [b]").BetweenFirst("[", "]"), "a"},
		{"Camel", Of("foo_bar").Camel(), "fooBar"},
		{"Studly", Of("foo_bar").Studly(), "FooBar"},
		{"Snake", Of("FooBar").Snake(), "foo_bar"},
		{"Kebab", Of("FooBar").Kebab(), "foo-bar"},
		{"Headline", Of("foo_bar").Headline(), "Foo Bar"},
		{"Title", Of("foo bar").Title(), "Foo Bar"},
		{"Ucfirst", Of("foo").Ucfirst(), "Foo"},
		{"Lcfirst", Of("Foo").Lcfirst(), "foo"},
		{"Limit", Of("hello world").Limit(5, "..."), "hello..."},
		{"Words", Of("one two three").Words(2, "..."), "one two..."},
		{"Mask", Of("secret").Mask('*', 2, 0), "se****"},
		{"PadLeft", Of("7").PadLeft(3, "0"), "007"},
		{"PadRight", Of("7").PadRight(3, "0"), "700"},
		{"PadBoth", Of("7").PadBoth(3, "0"), "070"},
		{"Start", Of("path").Start("/"), "/path"},
		{"Finish", Of("path").Finish("/"), "path/"},
		{"Replace", Of("a-b-c").Replace("-", "+"), "a+b+c"},
		{"Replace empty", Of("abc").Replace("", "+"), "abc"},
		{"ReplaceFirst", Of("a-b-c").ReplaceFirst("-", "+"), "a+b-c"},
		{"ReplaceLast", Of("a-b-c").ReplaceLast("-", "+"), "a-b+c"},
		{"ReplaceStart", Of("a-b").ReplaceStart("a", "x"), "x-b"},
		{"ReplaceEnd", Of("a-b").ReplaceEnd("b", "x"), "a-x"},
		{"ReplaceArray", Of("? ?").ReplaceArray("?", []string{"1", "2"}), "1 2"},
		{"Swap", Of("ab").Swap(map[string]string{"a": "b", "b": "a"}), "ba"},
		{"Remove", Of("a-b-c").Remove("-"), "abc"},
		{"Squish", Of("  a   b ").Squish(), "a b"},
		{"Reverse", Of("abc").Reverse(), "cba"},
		{"Substr", Of("abcdef").Substr(1, 3), "bcd"},
		{"Slug", Of("Hello World").Slug("-"), "hello-world"},
		{"Wrap", Of("x").Wrap("(", ")"), "(x)"},
		{"Unwrap", Of("(x)").Unwrap("(", ")"), "x"},
		{"Append", Of("a").Append("b", "c"), "abc"},
		{"Append none", Of("a").Append(), "a"},
		{"Prepend", Of("c").Prepend("a", "b"), "abc"},
		{"Lower", Of("ABC").Lower(), "abc"},
		{"Upper", Of("abc").Upper(), "ABC"},
		{"Trim", Of(" \tabc\n ").Trim(), "abc"},
		{"Trim cutset", Of("--abc-/").Trim("-", "/"), "abc"},
		{"Trim cutset is a rune set", Of("abcba").Trim("ab"), "c"},
		{"LTrim", Of("  abc  ").LTrim(), "abc  "},
		{"LTrim cutset", Of("--abc--").LTrim("-"), "abc--"},
		{"RTrim", Of("  abc  ").RTrim(), "  abc"},
		{"RTrim cutset", Of("--abc--").RTrim("-"), "--abc"},
		{"Repeat", Of("ab").Repeat(3), "ababab"},
		{"Repeat zero", Of("ab").Repeat(0), ""},
		{"Repeat negative", Of("ab").Repeat(-1), ""},
		{"Pipe", Of("abc").Pipe(strings.ToUpper), "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringablePredicates(t *testing.T) {
	s := Of("Hello World")
	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{"IsEmpty", s.IsEmpty(), false},
		{"IsEmpty zero", Of("").IsEmpty(), true},
		{"IsNotEmpty", s.IsNotEmpty(), true},
		{"IsNotEmpty zero", Of("").IsNotEmpty(), false},
		{"Contains", s.Contains("xyz", "World"), true},
		{"Contains none", s.Contains("xyz"), false},
		{"ContainsAll", s.ContainsAll("Hello", "World"), true},
		{"ContainsAll missing", s.ContainsAll("Hello", "xyz"), false},
		{"StartsWith", s.StartsWith("x", "Hello"), true},
		{"StartsWith none", s.StartsWith("World"), false},
		{"StartsWith empty", s.StartsWith(""), false},
		{"EndsWith", s.EndsWith("x", "World"), true},
		{"EndsWith none", s.EndsWith("Hello"), false},
		{"EndsWith empty", s.EndsWith(""), false},
		{"Is", s.Is("Hello*"), true},
		{"Is mismatch", s.Is("World*"), false},
		{"Exactly", s.Exactly("Hello World"), true},
		{"Exactly mismatch", s.Exactly("hello world"), false},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestStringableTerminals(t *testing.T) {
	if got := Of("你好 world").Length(); got != 8 {
		t.Errorf("Length = %d, want 8", got)
	}
	if got := Of("one two three").WordCount(); got != 3 {
		t.Errorf("WordCount = %d, want 3", got)
	}
	splits := []struct {
		s, sep string
		want   []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"abc", "", []string{"a", "b", "c"}},
		{"", ",", []string{""}},
	}
	for _, tt := range splits {
		if got := Of(tt.s).Split(tt.sep); !slices.Equal(got, tt.want) {
			t.Errorf("Split(%q, %q) = %q, want %q", tt.s, tt.sep, got, tt.want)
		}
	}
	if got := fmt.Sprint(Of("printed")); got != "printed" {
		t.Errorf("fmt.Sprint = %q, want %q", got, "printed")
	}
}

func TestStringableConditionals(t *testing.T) {
	mark := func(s Stringable) Stringable { return s.Append("!") }
	tests := []struct {
		name string
		got  Stringable
		want string
	}{
		{"When true", Of("a").When(true, mark), "a!"},
		{"When false", Of("a").When(false, mark), "a"},
		{"Unless true", Of("a").Unless(true, mark), "a"},
		{"Unless false", Of("a").Unless(false, mark), "a!"},
		{"WhenEmpty empty", Of("").WhenEmpty(mark), "!"},
		{"WhenEmpty filled", Of("a").WhenEmpty(mark), "a"},
		{"WhenNotEmpty filled", Of("a").WhenNotEmpty(mark), "a!"},
		{"WhenNotEmpty empty", Of("").WhenNotEmpty(mark), ""},
		{"WhenContains hit", Of("abc").WhenContains("b", mark), "abc!"},
		{"WhenContains miss", Of("abc").WhenContains("x", mark), "abc"},
		{"WhenStartsWith hit", Of("abc").WhenStartsWith("a", mark), "abc!"},
		{"WhenStartsWith miss", Of("abc").WhenStartsWith("c", mark), "abc"},
		{"WhenEndsWith hit", Of("abc").WhenEndsWith("c", mark), "abc!"},
		{"WhenEndsWith miss", Of("abc").WhenEndsWith("a", mark), "abc"},
		{"WhenExactly hit", Of("abc").WhenExactly("abc", mark), "abc!"},
		{"WhenExactly miss", Of("abc").WhenExactly("ab", mark), "abc"},
		{"WhenIs hit", Of("abc").WhenIs("a*", mark), "abc!"},
		{"WhenIs miss", Of("abc").WhenIs("b*", mark), "abc"},
	}
	for _, tt := range tests {
		if got := tt.got.String(); got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestStringableTap(t *testing.T) {
	var seen string
	got := Of("abc").Tap(func(s Stringable) { seen = s.String() }).Upper()
	if seen != "abc" || got.String() != "ABC" {
		t.Errorf("Tap saw %q and returned %q", seen, got)
	}
}

func TestStringableImmutable(t *testing.T) {
	s := Of("abc")
	_ = s.Upper().Append("x")
	if s.String() != "abc" {
		t.Errorf("original changed to %q", s)
	}
}
