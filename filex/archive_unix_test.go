//go:build unix

package filex

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestExtractTarGzModes(t *testing.T) {
	tests := []struct {
		name  string
		umask int
		want  map[string]os.FileMode
	}{
		{"default umask", 0o022, map[string]os.FileMode{"run.sh": 0o755, "ro.txt": 0o644, "none.txt": 0o600}},
		{"umask removes owner bits", 0o277, map[string]os.FileMode{"run.sh": 0o700, "ro.txt": 0o600, "none.txt": 0o600}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := writeTarGz(t,
				entry{name: "run.sh", body: "#!/bin/sh", mode: 0o755},
				entry{name: "ro.txt", body: "r", mode: 0o444},
				entry{name: "none.txt", body: "first", mode: 0},
				entry{name: "none.txt", body: "second", mode: 0},
			)
			root := t.TempDir()
			old := syscall.Umask(tt.umask)
			err := ExtractTarGz(archive, root)
			syscall.Umask(old)
			if err != nil {
				t.Fatal(err)
			}
			checkFiles(t, root, map[string]string{"none.txt": "second"})
			for name, want := range tt.want {
				info, err := os.Stat(filepath.Join(root, name))
				if err != nil {
					t.Fatal(err)
				}
				if got := info.Mode().Perm(); got != want {
					t.Errorf("%s mode = %v, want %v", name, got, want)
				}
			}
		})
	}
}
