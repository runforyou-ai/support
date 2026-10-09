package validate

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
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

// zoneinfoCases lists, for each kind of $ZONEINFO source, the names checked in
// a process that uses it.
var zoneinfoCases = map[string][]struct {
	name string
	want bool
}{
	"dir": {
		// An invalid file whose name differs in case is skipped, as by time.LoadLocation.
		{"Asia/Shanghai", true},
		{"asia/shanghai", false},
		{"Test/Zone", true},
		{"test/zone", false},
		{"Test/zone", false},
		// UTC is built in and never read, while the stored "utc" file matches itself.
		{"UTC", true},
		{"utc", true},
	},
	"zip": {
		{"Test/Zone", true},
		{"test/zone", false},
		{"Asia/Shanghai", true},
		// The zip entry is the source time.LoadLocation reads, and it matches exactly.
		{"asia/shanghai", true},
		{"UTC", true},
	},
	// time.LoadLocation skips compressed entries and archives with a comment,
	// so the system zone data decides.
	"deflate": {
		{"Test/Zone", false},
		{"Asia/Shanghai", true},
		{"asia/shanghai", false},
	},
	"comment": {
		{"Test/Zone", false},
		{"Asia/Shanghai", true},
		{"asia/shanghai", false},
	},
}

func TestTimezoneZoneinfo(t *testing.T) {
	if kind := os.Getenv("TIMEZONE_HELPER"); kind != "" {
		for _, tt := range zoneinfoCases[kind] {
			if got := Timezone(tt.name); got != tt.want {
				t.Errorf("%s: Timezone(%q) = %v, want %v", kind, tt.name, got, tt.want)
			}
		}
		return
	}
	// time and Timezone read $ZONEINFO once, so each source is checked in a new
	// process that sees it from its first lookup.
	dir := t.TempDir()
	for name, data := range map[string][]byte{"asia/shanghai": []byte("not a zone"), "Test/Zone": utcTZif, "utc": utcTZif} {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	archives := map[string]string{}
	for kind, write := range map[string]func(*zip.Writer, string) (io.Writer, error){
		"zip": func(zw *zip.Writer, name string) (io.Writer, error) {
			return zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
		},
		"deflate": func(zw *zip.Writer, name string) (io.Writer, error) {
			return zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
		},
		"comment": func(zw *zip.Writer, name string) (io.Writer, error) {
			if err := zw.SetComment("zones"); err != nil {
				return nil, err
			}
			return zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
		},
	} {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		for _, name := range []string{"Test/Zone", "asia/shanghai"} {
			w, err := write(zw, name)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := w.Write(utcTZif); err != nil {
				t.Fatal(err)
			}
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		archives[kind] = filepath.Join(t.TempDir(), "zoneinfo.zip")
		if err := os.WriteFile(archives[kind], buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	archives["dir"] = dir
	env := slices.DeleteFunc(os.Environ(), func(v string) bool { return strings.HasPrefix(v, "ZONEINFO=") })
	for kind, source := range archives {
		cmd := exec.Command(os.Args[0], "-test.run=^TestTimezoneZoneinfo$")
		cmd.Env = append(slices.Clone(env), "TIMEZONE_HELPER="+kind, "ZONEINFO="+source)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("%s helper process: %v\n%s", kind, err, out)
		}
	}
}
