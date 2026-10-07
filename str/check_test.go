package str

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		s       string
		needles []string
		any     bool
		all     bool
	}{
		{"This is my name", []string{"my"}, true, true},
		{"This is my name", []string{"my", "foo"}, true, false},
		{"This is my name", []string{"foo", "bar"}, false, false},
		{"This is my name", []string{"This", "name"}, true, true},
		{"This is my name", []string{""}, false, true},
		{"This is my name", nil, false, true},
		{"你好世界", []string{"世界"}, true, true},
	}
	for _, tt := range tests {
		if got := Contains(tt.s, tt.needles...); got != tt.any {
			t.Errorf("Contains(%q, %q) = %v, want %v", tt.s, tt.needles, got, tt.any)
		}
		if got := ContainsAll(tt.s, tt.needles...); got != tt.all {
			t.Errorf("ContainsAll(%q, %q) = %v, want %v", tt.s, tt.needles, got, tt.all)
		}
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		pattern, s string
		want       bool
	}{
		{"foo*", "foobar", true},
		{"*bar", "foobar", true},
		{"f*r", "foobar", true},
		{"*", "", true},
		{"*", "anything", true},
		{"", "", true},
		{"", "a", false},
		{"foo", "foo", true},
		{"foo", "foobar", false},
		{"Foo*", "foobar", false},
		{"a*b*c", "aXXbYYc", true},
		{"a*b*c", "aXXbYY", false},
		{"a*a", "aaa", true},
		{"*.go", "main.go.bak", false},
		{"baz*", "foobar", false},
		{"用户*", "用户名", true},
		{"a.c", "abc", false},
	}
	for _, tt := range tests {
		if got := Is(tt.pattern, tt.s); got != tt.want {
			t.Errorf("Is(%q, %q) = %v, want %v", tt.pattern, tt.s, got, tt.want)
		}
	}
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{`{"a":1}`, true},
		{`[1,2,3]`, true},
		{`"text"`, true},
		{`123`, true},
		{`{a:1}`, false},
		{``, false},
		{`text`, false},
	}
	for _, tt := range tests {
		if got := IsJSON(tt.s); got != tt.want {
			t.Errorf("IsJSON(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"https://example.com", true},
		{"http://localhost:8080/path?q=1#frag", true},
		{"ftp://user:pass@host/file", true},
		{"example.com", false},
		{"/relative/path", false},
		{"mailto:someone@example.com", false},
		{"http://", false},
		{"http://exa mple.com", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsURL(tt.s); got != tt.want {
			t.Errorf("IsURL(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"a0a2a2d2-0b87-4a18-83f2-2529882be2de", true},
		{"A0A2A2D2-0B87-4A18-83F2-2529882BE2DE", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"01890a5d-ac96-774b-bcce-b302099a8057", true},
		{"a0a2a2d2-0b87-4a18-83f2-2529882be2d", false},
		{"a0a2a2d20b874a1883f22529882be2de", false},
		{"a0a2a2d2-0b87-4a18-83f2_2529882be2de", false},
		{"g0a2a2d2-0b87-4a18-83f2-2529882be2de", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsUUID(tt.s); got != tt.want {
			t.Errorf("IsUUID(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsULID(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"01ARZ3NDEKTSV4RRFFQ69G5FAV", true},
		{"01arz3ndektsv4rrffq69g5fav", true},
		{"7ZZZZZZZZZZZZZZZZZZZZZZZZZ", true},
		{"8ZZZZZZZZZZZZZZZZZZZZZZZZZ", false},
		{"01ARZ3NDEKTSV4RRFFQ69G5FA", false},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAU", false},
		{"01ARZ3NDEKTSV4RRFFQ69G5FAI", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsULID(tt.s); got != tt.want {
			t.Errorf("IsULID(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsASCII(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"", true},
		{"Hello, World!", true},
		{"\x7f", true},
		{"héllo", false},
		{"你好", false},
	}
	for _, tt := range tests {
		if got := IsASCII(tt.s); got != tt.want {
			t.Errorf("IsASCII(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}
