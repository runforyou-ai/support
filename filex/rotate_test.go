package filex

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// readFiles returns the content of path and its numbered backups that exist.
func readFiles(t *testing.T, path string, upTo int) map[string]string {
	t.Helper()
	got := map[string]string{}
	names := []string{path}
	for i := 1; i <= upTo; i++ {
		names = append(names, fmt.Sprintf("%s.%d", path, i))
	}
	for _, name := range names {
		data, err := os.ReadFile(name)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		got[strings.TrimPrefix(name, path)] = string(data)
	}
	return got
}

func TestOpenRotating(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		maxSize  int64
		backups  int
		writes   []string
		want     map[string]string
	}{
		{
			name:    "below limit",
			maxSize: 10,
			backups: 2,
			writes:  []string{"abc", "def"},
			want:    map[string]string{"": "abcdef"},
		},
		{
			name:    "exactly at limit does not rotate",
			maxSize: 6,
			backups: 2,
			writes:  []string{"abc", "def"},
			want:    map[string]string{"": "abcdef"},
		},
		{
			name:    "rotates after exceeding limit",
			maxSize: 5,
			backups: 2,
			writes:  []string{"abc", "def", "g"},
			want:    map[string]string{"": "g", ".1": "abcdef"},
		},
		{
			name:    "write is never split",
			maxSize: 2,
			backups: 1,
			writes:  []string{"abcdefgh"},
			want:    map[string]string{"": "", ".1": "abcdefgh"},
		},
		{
			name:    "backups shift and oldest is removed",
			maxSize: 1,
			backups: 2,
			writes:  []string{"aa", "bb", "cc", "d"},
			want:    map[string]string{"": "d", ".1": "cc", ".2": "bb"},
		},
		{
			name:    "no backups truncates",
			maxSize: 3,
			backups: 0,
			writes:  []string{"abcd", "ef"},
			want:    map[string]string{"": "ef"},
		},
		{
			name:    "negative backups truncates",
			maxSize: 3,
			backups: -1,
			writes:  []string{"abcd", "ef"},
			want:    map[string]string{"": "ef"},
		},
		{
			name:     "existing size counts toward limit",
			existing: "12345",
			maxSize:  6,
			backups:  1,
			writes:   []string{"67", "8"},
			want:     map[string]string{"": "8", ".1": "1234567"},
		},
		{
			name:    "zero limit rotates after every write",
			maxSize: 0,
			backups: 3,
			writes:  []string{"a", "b"},
			want:    map[string]string{"": "", ".1": "b", ".2": "a"},
		},
		{
			name:    "empty writes do not rotate",
			maxSize: 0,
			backups: 1,
			writes:  []string{"", ""},
			want:    map[string]string{"": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "logs", "app.log")
			if tt.existing != "" {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(tt.existing), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			f, err := OpenRotating(path, tt.maxSize, tt.backups)
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range tt.writes {
				if n, err := f.Write([]byte(w)); n != len(w) || err != nil {
					t.Fatalf("Write(%q) = %d, %v", w, n, err)
				}
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			got := readFiles(t, path, 5)
			if fmt.Sprint(got) != fmt.Sprint(tt.want) {
				t.Errorf("files = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenRotatingErrors(t *testing.T) {
	tmp := t.TempDir()
	blocker := filepath.Join(tmp, "file")
	if err := os.WriteFile(blocker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, path string
	}{
		{"parent is a file", filepath.Join(blocker, "app.log")},
		{"path is a directory", tmp},
	}
	for _, tt := range tests {
		if f, err := OpenRotating(tt.path, 10, 1); err == nil || f != nil {
			t.Errorf("%s: OpenRotating = %v, %v; want nil, error", tt.name, f, err)
		}
	}
}

func TestRotatingFileClose(t *testing.T) {
	f, err := OpenRotating(filepath.Join(t.TempDir(), "app.log"), 10, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close = %v", err)
	}
	if err := f.Close(); err != nil {
		t.Errorf("second Close = %v, want nil", err)
	}
	if n, err := f.Write([]byte("x")); n != 0 || !errors.Is(err, os.ErrClosed) {
		t.Errorf("Write after Close = %d, %v; want 0, os.ErrClosed", n, err)
	}
}

func TestRotatingFileRotateError(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "logs")
	path := filepath.Join(dir, "app.log")
	f, err := OpenRotating(path, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	// A directory occupying the backup name makes the rename fail.
	if err := os.MkdirAll(filepath.Join(path+".1", "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	if n, err := f.Write([]byte("ab")); n != 2 || err == nil {
		t.Errorf("Write = %d, %v; want 2, rotate error", n, err)
	}
	if n, err := f.Write([]byte("c")); n != 1 || err == nil || errors.Is(err, os.ErrClosed) {
		t.Errorf("Write after failed rotation = %d, %v; want 1, rotate error", n, err)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "abc" {
		t.Fatalf("current file = %q, %v; want abc", data, err)
	}
	if err := os.RemoveAll(path + ".1"); err != nil {
		t.Fatal(err)
	}
	if n, err := f.Write([]byte("d")); n != 1 || err != nil {
		t.Fatalf("Write after obstacle removed = %d, %v; want 1, nil", n, err)
	}
	if got := readFiles(t, path, 1); fmt.Sprint(got) != fmt.Sprint(map[string]string{"": "", ".1": "abcd"}) {
		t.Errorf("files = %v", got)
	}
}

func TestRotatingFileRemoveError(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs Unix directory permissions enforced")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	f, err := OpenRotating(path, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	if n, err := f.Write([]byte("ab")); n != 2 || err == nil {
		t.Errorf("Write = %d, %v; want 2, rotate error", n, err)
	}
	if n, err := f.Write([]byte("c")); n != 1 || err == nil {
		t.Errorf("second Write = %d, %v; want 1, rotate error", n, err)
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if n, err := f.Write([]byte("d")); n != 1 || err != nil {
		t.Fatalf("Write after permissions restored = %d, %v; want 1, nil", n, err)
	}
	if data, err := os.ReadFile(path); err != nil || len(data) != 0 {
		t.Errorf("current file = %q, %v; want empty", data, err)
	}
}

func TestRotatingFileShiftStopsAtMissingBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	// A missing .2, as left by a partly completed rotation, is filled first.
	for name, body := range map[string]string{".1": "one", ".3": "three"} {
		if err := os.WriteFile(path+name, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	f, err := OpenRotating(path, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Write([]byte("ab")); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"": "", ".1": "ab", ".2": "one", ".3": "three"}
	if got := readFiles(t, path, 3); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("files = %v, want %v", got, want)
	}
}

func TestRotatingFileReopensMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	f, err := OpenRotating(path, 100, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	// Simulate a rotation that moved the file away and could not reopen it.
	_ = f.file.Close()
	f.file = nil
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if n, err := f.Write([]byte("x")); n != 1 || err != nil {
		t.Fatalf("Write = %d, %v; want 1, nil", n, err)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "x" {
		t.Errorf("current file = %q, %v; want x", data, err)
	}
}

func TestRotatingFileReopenError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	f, err := OpenRotating(path, 100, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	_ = f.file.Close()
	f.file = nil
	// A directory in place of the current file cannot be opened for writing.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if n, err := f.Write([]byte("x")); n != 0 || err == nil || errors.Is(err, os.ErrClosed) {
		t.Errorf("Write = %d, %v; want 0, open error", n, err)
	}
}

func TestRotatingFileConcurrent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	f, err := OpenRotating(path, 100, 50)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			for j := range 50 {
				_, _ = fmt.Fprintf(f, "%d-%02d\n", i, j)
			}
		})
	}
	wg.Wait()
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, content := range readFiles(t, path, 50) {
		total += strings.Count(content, "\n")
	}
	if total != 400 {
		t.Errorf("lines across files = %d, want 400", total)
	}
}

func TestRotatingFileKeepsBackupOnShiftError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	if err := os.WriteFile(path+".1", []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A non-empty directory in place of the oldest backup cannot be removed.
	if err := os.MkdirAll(filepath.Join(path+".2", "keep"), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := OpenRotating(path, 4, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	if _, err := f.Write([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("defg")); err == nil {
		t.Fatal("Write succeeded, want rotation error")
	}
	if data, err := os.ReadFile(path + ".1"); err != nil || string(data) != "old" {
		t.Fatalf("backup .1 = %q, %v; want old", data, err)
	}
}
