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
		{"asia/shanghai", false},
		{"ASIA/SHANGHAI", false},
		{"Asia/shanghai", false},
		{"utc", false},
		{"./UTC", false},
		{"Asia//Shanghai", false},
		{"Etc/GMT+8", true},
		{"America/Argentina/Buenos_Aires", true},
	}
	for _, tt := range tests {
		if got := Timezone(tt.name); got != tt.want {
			t.Errorf("Timezone(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTimezoneZoneinfoWithoutZone(t *testing.T) {
	// A $ZONEINFO directory without the zone defers to the next source.
	t.Setenv("ZONEINFO", t.TempDir())
	if !Timezone("Asia/Shanghai") {
		t.Error(`Timezone("Asia/Shanghai") = false with an empty $ZONEINFO, want true`)
	}
}
