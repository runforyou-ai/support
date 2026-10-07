package filex

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWriteAtomic(t *testing.T) {
	unix := runtime.GOOS != "windows"
	tests := []struct {
		name     string
		setup    func(t *testing.T, path string)
		path     func(dir string) string
		perm     os.FileMode
		data     string
		wantErr  bool
		wantMode os.FileMode
		unixOnly bool
	}{
		{
			name:     "new file gets perm exactly",
			perm:     0o666,
			data:     "new",
			wantMode: 0o666,
		},
		{
			name:     "new file ignores non-permission bits",
			perm:     os.ModeSetuid | 0o640,
			data:     "new",
			wantMode: 0o640,
		},
		{
			name: "existing file keeps its mode",
			setup: func(t *testing.T, path string) {
				writeOrFail(t, path, "old", 0o600)
			},
			perm:     0o644,
			data:     "replaced",
			wantMode: 0o600,
		},
		{
			name: "existing file is truncated",
			setup: func(t *testing.T, path string) {
				writeOrFail(t, path, "a much longer old content", 0o644)
			},
			perm:     0o644,
			data:     "",
			wantMode: 0o644,
		},
		{
			name: "symlink is replaced, target untouched",
			setup: func(t *testing.T, path string) {
				target := filepath.Join(filepath.Dir(path), "target")
				writeOrFail(t, target, "target", 0o600)
				if err := os.Symlink(target, path); err != nil {
					t.Fatal(err)
				}
			},
			perm:     0o640,
			data:     "replaced",
			wantMode: 0o640,
			unixOnly: true,
		},
		{
			name: "directory is an error",
			setup: func(t *testing.T, path string) {
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatal(err)
				}
			},
			perm:    0o644,
			wantErr: true,
		},
		{
			name:    "missing parent directory",
			path:    func(dir string) string { return filepath.Join(dir, "missing", "f") },
			perm:    0o644,
			wantErr: true,
		},
		{
			name: "parent is a file",
			path: func(dir string) string {
				parent := filepath.Join(dir, "parent")
				_ = os.WriteFile(parent, nil, 0o644)
				return filepath.Join(parent, "f")
			},
			perm:    0o644,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.unixOnly && !unix {
				t.Skip("symbolic links need privileges on Windows")
			}
			dir := t.TempDir()
			path := filepath.Join(dir, "f.txt")
			if tt.path != nil {
				path = tt.path(dir)
			}
			if tt.setup != nil {
				tt.setup(t, path)
			}
			err := WriteAtomic(path, []byte(tt.data), tt.perm)
			if tt.wantErr {
				if err == nil {
					t.Fatal("WriteAtomic succeeded, want error")
				}
			} else if err != nil {
				t.Fatalf("WriteAtomic: %v", err)
			}
			entries, _ := os.ReadDir(filepath.Dir(path))
			for _, e := range entries {
				if filepath.Ext(e.Name()) != "" && e.Name()[0] == '.' {
					t.Errorf("temporary file %s left behind", e.Name())
				}
			}
			if tt.wantErr {
				return
			}
			info, err := os.Lstat(path)
			if err != nil {
				t.Fatal(err)
			}
			if !info.Mode().IsRegular() {
				t.Fatalf("mode = %v, want regular file", info.Mode())
			}
			if unix && info.Mode().Perm() != tt.wantMode {
				t.Errorf("mode = %v, want %v", info.Mode().Perm(), tt.wantMode)
			}
			if got, _ := os.ReadFile(path); string(got) != tt.data {
				t.Errorf("content = %q, want %q", got, tt.data)
			}
			if tt.unixOnly {
				if got, _ := os.ReadFile(filepath.Join(dir, "target")); string(got) != "target" {
					t.Errorf("symlink target content = %q, want unchanged", got)
				}
			}
		})
	}
}

func TestWriteAtomicReadOnlyDir(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs Unix directory permissions enforced")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "f")
	writeOrFail(t, path, "old", 0o644)
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	if err := WriteAtomic(path, []byte("new"), 0o644); err == nil {
		t.Fatal("WriteAtomic succeeded in read-only directory, want error")
	}
	if got, _ := os.ReadFile(path); string(got) != "old" {
		t.Errorf("content = %q, want unchanged", got)
	}
}

func TestWriteAtomicDoesNotMutateData(t *testing.T) {
	data := []byte("data")
	if err := WriteAtomic(filepath.Join(t.TempDir(), "f"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if string(data) != "data" {
		t.Errorf("data = %q, want unchanged", data)
	}
}

// writeOrFail writes content to path with mode, failing the test on error.
func writeOrFail(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func TestWriteAtomicDirSyncError(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs Unix directory permissions enforced")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "f")
	writeOrFail(t, path, "old", 0o644)
	// Write and search permission without read permission lets the rename succeed but not the directory open.
	if err := os.Chmod(dir, 0o300); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	if err := WriteAtomic(path, []byte("new"), 0o644); err == nil {
		t.Fatal("WriteAtomic succeeded without directory read permission, want error")
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(path); string(got) != "new" {
		t.Errorf("content = %q, want new", got)
	}
}
