package str

import (
	"strings"
	"testing"
)

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
		{"HTTP://Example.COM/Path", true},
		{"http://127.0.0.1:65535/x", true},
		{"https://[2001:DB8::1]/", true},
		{"http://a.com:/x", false},
		{"http://a.com:0/", false},
		{"http://a.com:65536/", false},
		{"http://a.com;x/", false},
		{"http://a.com%40b.com/", false},
		{"http://exa%2fmple.com/", false},
		{"http://a.com\\@b.com/", false},
		{"http://a..b/", false},
		{"http://-a.com/", false},
		{"http://example.com./", false},
		{"http://\u212a.com/", false},
		{"http://例子.com/", false},
		{"http://1.2.3/", false},
		{"http://0X7F000001/", false},
		{"http://[0:0:0:0:0:0:0:1]/", true},
		{"https://api.example.com/v1\\", false},
		{"https://example.com/a\\b", false},
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
		{"https://example.com:", false},
		{"https://example.com:0", false},
		{"https://example.com:080", false},
		{"https://example.com:65535", true},
		{"https://example.com:65536", false},
		{"https://example.com:99999", false},
		{"https://example.com:8080:1", false},
		{"http://example.com;x", false},
		{"https://example.com.", false},
		{"https://EXAMPLE.com", false},
		{"https://a..b", false},
		{"https://-a.com", false},
		{"https://a-.com", false},
		{"https://a-b.example-1.com", true},
		{"https://xn--fsqu00a.com", true},
		{"https://例子.com", false},
		{"https://_x.com", false},
		{"https://exa%2fmple.com", false},
		{"https://exa\\mple.com", false},
		{"https://example.com\n", false},
		{"https://" + strings.Repeat("a", 63) + ".com", true},
		{"https://" + strings.Repeat("a", 64) + ".com", false},
		{"https://" + strings.Repeat("a.", 126) + "com", false},
		{"http://192.0.2.1:8080", true},
		{"http://192.0.2.256", false},
		{"http://01.2.3.4", false},
		{"http://1.2.3", false},
		{"http://a.123", false},
		{"https://[::1]", true},
		{"https://[::1", false},
		{"https://[::1]x", false},
		{"https://[::1]:", false},
		{"https://[2001:DB8::1]", false},
		{"https://[fe80::1%25en0]", false},
		{"https://[192.0.2.1]", false},
		{"https://0x7f000001", false},
		{"https://0x7f.0x0.0x0.0x1", false},
		{"https://0x", false},
		{"https://a.0xg", true},
		{"https://[0:0:0:0:0:0:0:1]", false},
		{"https://[2001:0db8::1]", false},
		{"https://[2001:db8::0:1]", false},
		{"https://[2001:db8::1]", true},
		{"https://[::ffff:c000:201]", false},
		{"https://[::ffff:192.0.2.1]", false},
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
		{"https://api.example.com/v1?", false},
		{"https://api.example.com/#", false},
		{"https://api.example.com/v1?key=1", false},
		{"https://api.example.com#top", false},
		{"https://user@api.example.com", false},
		{"ftp://api.example.com", false},
		{"https:///v1", false},
		{"/v1", false},
		{" https://api.example.com", false},
		{"https://api.example.com:/v1", false},
		{"https://api.example.com:0/v1", false},
		{"https://api.example.com;x/v1", false},
		{"https://api.example.com/v1\\", false},
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
		{"o'brien+x_y=z@mail.example.co.uk", true},
		{"a@b.c", true},
		{"taylor@localhost", false},
		{"taylor@[192.0.2.1]", false},
		{"taylor@192.0.2.1", false},
		{"张三@example.com", false},
		{"测试@例子.cn", false},
		{"taylor@例子.cn", false},
		{"taylor@-example.com", false},
		{"taylor@example-.com", false},
		{"taylor@example.com.", false},
		{"taylor@exa_mple.com", false},
		{".taylor@example.com", false},
		{"tay..lor@example.com", false},
		{"tay lor@example.com", false},
		{"taylor@example.com\r\nBcc: x@y.com", false},
		{"taylor@b@example.com", false},
		{strings.Repeat("a", 64) + "@example.com", true},
		{"a@" + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 63) + "." + strings.Repeat("e", 56) + ".com", true},
		{strings.Repeat("a", 65) + "@example.com", false},
		{"a@" + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 63) + "." + strings.Repeat("e", 57) + ".com", false},
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
