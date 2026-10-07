package str

import (
	"slices"
	"testing"
)

func TestCaseConversion(t *testing.T) {
	tests := []struct {
		in                                    string
		camel, studly, snake, kebab, headline string
	}{
		{"", "", "", "", "", ""},
		{"hello world-foo_bar", "helloWorldFooBar", "HelloWorldFooBar", "hello_world_foo_bar", "hello-world-foo-bar", "Hello World Foo Bar"},
		{"HTTPServer", "httpServer", "HttpServer", "http_server", "http-server", "HTTP Server"},
		{"userID2", "userId2", "UserId2", "user_id2", "user-id2", "User ID2"},
		{"HTTP2Server", "http2Server", "Http2Server", "http2_server", "http2-server", "HTTP2 Server"},
		{"v2Beta", "v2Beta", "V2Beta", "v2_beta", "v2-beta", "V2 Beta"},
		{"steve_jobs", "steveJobs", "SteveJobs", "steve_jobs", "steve-jobs", "Steve Jobs"},
		{"EmailNotificationSent", "emailNotificationSent", "EmailNotificationSent", "email_notification_sent", "email-notification-sent", "Email Notification Sent"},
		{"FOO_BAR", "fooBar", "FooBar", "foo_bar", "foo-bar", "FOO BAR"},
		{"  --leading and trailing__  ", "leadingAndTrailing", "LeadingAndTrailing", "leading_and_trailing", "leading-and-trailing", "Leading And Trailing"},
		{"你好 世界", "你好世界", "你好世界", "你好_世界", "你好-世界", "你好 世界"},
		{"用户Name", "用户Name", "用户Name", "用户_name", "用户-name", "用户 Name"},
		{"überCool", "überCool", "ÜberCool", "über_cool", "über-cool", "Über Cool"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			checks := []struct {
				name string
				got  string
				want string
			}{
				{"Camel", Camel(tt.in), tt.camel},
				{"Studly", Studly(tt.in), tt.studly},
				{"Snake", Snake(tt.in), tt.snake},
				{"Kebab", Kebab(tt.in), tt.kebab},
				{"Headline", Headline(tt.in), tt.headline},
			}
			for _, c := range checks {
				if c.got != c.want {
					t.Errorf("%s(%q) = %q, want %q", c.name, tt.in, c.got, c.want)
				}
			}
		})
	}
}

func TestTitle(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"hello world", "Hello World"},
		{"hELLO wORLD", "Hello World"},
		{"hello-world_foo", "Hello-World_Foo"},
		{"don't stop", "Don't Stop"},
		{"1st place", "1st Place"},
		{"你好 world", "你好 World"},
		{"'quoted'", "'Quoted'"},
	}
	for _, tt := range tests {
		if got := Title(tt.in); got != tt.want {
			t.Errorf("Title(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestUcfirstLcfirst(t *testing.T) {
	tests := []struct {
		in, upper, lower string
	}{
		{"", "", ""},
		{"hello", "Hello", "hello"},
		{"Hello", "Hello", "hello"},
		{"élan", "Élan", "élan"},
		{"你好", "你好", "你好"},
		{"\xffabc", "\xffabc", "\xffabc"},
	}
	for _, tt := range tests {
		if got := Ucfirst(tt.in); got != tt.upper {
			t.Errorf("Ucfirst(%q) = %q, want %q", tt.in, got, tt.upper)
		}
		if got := Lcfirst(tt.upper); got != tt.lower {
			t.Errorf("Lcfirst(%q) = %q, want %q", tt.upper, got, tt.lower)
		}
	}
}

func TestUcSplit(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"FooBar", []string{"Foo", "Bar"}},
		{"fooBar", []string{"foo", "Bar"}},
		{"HTTP", []string{"H", "T", "T", "P"}},
		{"ÜberÉlan", []string{"Über", "Élan"}},
		{"lower", []string{"lower"}},
	}
	for _, tt := range tests {
		if got := UcSplit(tt.in); !slices.Equal(got, tt.want) {
			t.Errorf("UcSplit(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
