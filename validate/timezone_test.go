package validate

import (
	"os"
	"os/exec"
	"path/filepath"
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

// utcTZif is a minimal TZif file describing a single zone at UTC.
var utcTZif = append(append([]byte("TZif"), make([]byte, 16+4*4)...), 0, 0, 0, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 'U', 'T', 'C', 0)

func TestTimezoneZoneinfo(t *testing.T) {
	if os.Getenv("TIMEZONE_HELPER") == "" {
		// time and Timezone read $ZONEINFO once, so the check runs in a new
		// process that sees the prepared directory from its first lookup.
		dir := t.TempDir()
		files := map[string][]byte{
			"asia/shanghai": []byte("not a zone"),
			"Test/Zone":     utcTZif,
		}
		for name, data := range files {
			path := filepath.Join(dir, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		cmd := exec.Command(os.Args[0], "-test.run=^TestTimezoneZoneinfo$")
		cmd.Env = append(os.Environ(), "TIMEZONE_HELPER=1", "ZONEINFO="+dir)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("helper process: %v\n%s", err, out)
		}
		return
	}
	tests := []struct {
		name string
		want bool
	}{
		// An invalid file whose name differs in case is skipped, as by time.LoadLocation.
		{"Asia/Shanghai", true},
		{"asia/shanghai", false},
		{"Test/Zone", true},
		{"test/zone", false},
		{"Test/zone", false},
	}
	for _, tt := range tests {
		if got := Timezone(tt.name); got != tt.want {
			t.Errorf("Timezone(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}
