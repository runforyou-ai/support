package validate

import (
	"testing"
	_ "time/tzdata"
)

func TestTimezone(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"UTC", true},
		{"Asia/Shanghai", true},
		{"America/New_York", true},
		{"Europe/London", true},
		{"", false},
		{"Local", false},
		{"Mars/Olympus_Mons", false},
		{"asia/shanghai_x", false},
		{"../etc/passwd", false},
		{"/Asia/Shanghai", false},
	}
	for _, tt := range tests {
		if got := Timezone(tt.name); got != tt.want {
			t.Errorf("Timezone(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}
