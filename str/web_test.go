package str

import "testing"

func TestIsHTTPURL(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"https://example.com", true},
		{"http://example.com/docs?page=2#intro", true},
		{"HTTPS://example.com", true},
		{"https://example.com:8080/", true},
		{"http://:8080", false},
		{"https://[::1]:80", true},
		{"https://user@example.com", false},
		{"https://user:pass@example.com", false},
		{"ftp://example.com", false},
		{"https:///path", false},
		{"https:example.com", false},
		{"//example.com", false},
		{"/path", false},
		{" https://example.com", false},
		{"https://example.com:abc", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsHTTPURL(tt.s); got != tt.want {
			t.Errorf("IsHTTPURL(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsHTTPOrigin(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"https://example.com", true},
		{"http://localhost:8080", true},
		{"https://[::1]:443", true},
		{"https://:443", false},
		{"HTTPS://example.com", false},
		{"https://example.com/", false},
		{"https://example.com/app", false},
		{"https://example.com?", false},
		{"https://example.com#top", false},
		{"https://user@example.com", false},
		{"https://exa mple.com", false},
		{"ftp://example.com", false},
		{"https://", false},
		{"example.com", false},
		{"https://example.com:abc", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsHTTPOrigin(tt.s); got != tt.want {
			t.Errorf("IsHTTPOrigin(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsHTTPBaseURL(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"https://api.example.com", true},
		{"http://localhost:8080/v1/", true},
		{"http://:8080/v1", false},
		{"HTTPS://api.example.com", true},
		{"https://api.example.com/v1?", true},
		{"https://api.example.com/#", true},
		{"https://api.example.com/v1?key=1", false},
		{"https://api.example.com#top", false},
		{"https://user@api.example.com", false},
		{"ftp://api.example.com", false},
		{"https:///v1", false},
		{"/v1", false},
		{" https://api.example.com", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsHTTPBaseURL(tt.s); got != tt.want {
			t.Errorf("IsHTTPBaseURL(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"taylor@example.com", true},
		{"Taylor@Example.COM", true},
		{"first.last+tag@example.co", true},
		{"taylor@localhost", true},
		{"taylor@[192.0.2.1]", true},
		{"张三@example.com", true},
		{" taylor@example.com", false},
		{"taylor@example.com ", false},
		{"<taylor@example.com>", false},
		{"Taylor <taylor@example.com>", false},
		{"taylor@example.com (Taylor)", false},
		{`"taylor"@example.com`, false},
		{"taylor@@example.com", false},
		{"taylor.@example.com", false},
		{"taylor@example..com", false},
		{"taylor", false},
		{"@example.com", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsEmail(tt.s); got != tt.want {
			t.Errorf("IsEmail(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		s, want string
	}{
		{"taylor@example.com", "t***@example.com"},
		{"t@example.com", "t***@example.com"},
		{"张三丰@example.com", "张***@example.com"},
		{"a@b@c", "a***@b@c"},
		{"taylor@", "t***@"},
		{"\xffbad@example.com", "�***@example.com"},
		{"@example.com", "@example.com"},
		{"taylor", "taylor"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := MaskEmail(tt.s); got != tt.want {
			t.Errorf("MaskEmail(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}
