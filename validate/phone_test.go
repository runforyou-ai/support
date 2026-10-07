package validate

import "testing"

func TestE164(t *testing.T) {
	tests := []struct {
		s    string
		want string
		ok   bool
	}{
		{"+15550109999", "+15550109999", true},
		{"+1 (555) 010-9999", "+15550109999", true},
		{"  +86 138 0013 8000\n", "+8613800138000", true},
		{"+44 20　7946 0958", "+442079460958", true},
		{"+123456", "+123456", true},
		{"+123456789012345", "+123456789012345", true},
		{"+12345", "", false},
		{"+1234567890123456", "", false},
		{"15550109999", "", false},
		{"++15550109999", "", false},
		{"+1.555.010.9999", "", false},
		{"+1555010999x", "", false},
		{"+１５５５０１０９９９９", "", false},
		{"+", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, ok := E164(tt.s)
		if got != tt.want || ok != tt.ok {
			t.Errorf("E164(%q) = %q, %v, want %q, %v", tt.s, got, ok, tt.want, tt.ok)
		}
	}
}
