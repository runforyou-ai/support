package filex

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// entry describes one archive entry used by the tests.
type entry struct {
	name   string
	body   string
	mode   int64
	dir    bool
	link   string
	hard   string
	isLink bool
}

func file(name, body string) entry { return entry{name: name, body: body, mode: 0o644} }
func dir(name string) entry        { return entry{name: name, dir: true, mode: 0o755} }
func symlink(name, target string) entry {
	return entry{name: name, link: target, isLink: true, mode: 0o777}
}

// writeTarGz writes entries as a tar.gz archive and returns its path.
func writeTarGz(t *testing.T, entries ...entry) string {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		h := &tar.Header{Name: e.name, Mode: e.mode, Typeflag: tar.TypeReg, Size: int64(len(e.body))}
		switch {
		case e.dir:
			h.Typeflag, h.Size = tar.TypeDir, 0
		case e.isLink:
			h.Typeflag, h.Linkname, h.Size = tar.TypeSymlink, e.link, 0
		case e.hard != "":
			h.Typeflag, h.Linkname, h.Size = tar.TypeLink, e.hard, 0
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if h.Size > 0 {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "a.tar.gz")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// writeZip writes entries as a zip archive and returns its path.
func writeZip(t *testing.T, entries ...entry) string {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, e := range entries {
		h := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
		switch {
		case e.dir:
			h.SetMode(os.ModeDir | os.FileMode(e.mode))
		case e.isLink:
			h.SetMode(os.ModeSymlink | 0o777)
			e.body = e.link
		default:
			h.SetMode(os.FileMode(e.mode))
		}
		w, err := zw.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(e.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "a.zip")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// checkFiles verifies that every file in want exists under root with the given content.
func checkFiles(t *testing.T, root string, want map[string]string) {
	t.Helper()
	for name, body := range want {
		got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			t.Errorf("read %s: %v", name, err)
			continue
		}
		if string(got) != body {
			t.Errorf("%s = %q, want %q", name, got, body)
		}
	}
}

// checkAbsent verifies that none of the names exist under root, without following links.
func checkAbsent(t *testing.T, root string, names []string) {
	t.Helper()
	for _, name := range names {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(name))); err == nil {
			t.Errorf("%s exists, want absent", name)
		}
	}
}

func TestEntryPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root")
	sep := string(filepath.Separator)
	tests := []struct {
		root, name string
		want       string
		wantErr    bool
	}{
		{root, "a.txt", filepath.Join(root, "a.txt"), false},
		{root, "dir/sub/a.txt", filepath.Join(root, "dir", "sub", "a.txt"), false},
		{root, "./a.txt", filepath.Join(root, "a.txt"), false},
		{root, "", root, false},
		{root, ".", root, false},
		{root, "/etc/passwd", filepath.Join(root, "etc", "passwd"), false},
		{root, "a/../b", filepath.Join(root, "b"), false},
		{root + sep, "a", filepath.Join(root, "a"), false},
		{root, "..", "", true},
		{root, "../evil", "", true},
		{root, "a/../../evil", "", true},
		{root, "../root2/evil", "", true},
		{"rel", "x", filepath.Join("rel", "x"), false},
		{"rel", "../x", "", true},
	}
	for _, tt := range tests {
		got, err := EntryPath(tt.root, tt.name)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("EntryPath(%q, %q) = %q, %v; want %q, error %v", tt.root, tt.name, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestExtractTarGz(t *testing.T) {
	tests := []struct {
		name       string
		entries    []entry
		symlinks   bool
		wantErr    string
		wantFiles  map[string]string
		wantAbsent []string
	}{
		{
			name:      "files and directories",
			entries:   []entry{dir("d/"), file("d/a.txt", "a"), file("top.txt", "top"), file("x/y/z.txt", "z"), dir("empty/")},
			wantFiles: map[string]string{"d/a.txt": "a", "top.txt": "top", "x/y/z.txt": "z"},
		},
		{
			name:      "leading slash and dot entries stay inside",
			entries:   []entry{dir("./"), file("/abs.txt", "abs"), file("./dot.txt", "dot")},
			wantFiles: map[string]string{"abs.txt": "abs", "dot.txt": "dot"},
		},
		{
			name:      "empty archive",
			entries:   nil,
			wantFiles: map[string]string{},
		},
		{
			name:       "hard links are skipped",
			entries:    []entry{file("a.txt", "a"), {name: "b.txt", hard: "a.txt"}},
			wantFiles:  map[string]string{"a.txt": "a"},
			wantAbsent: []string{"b.txt"},
		},
		{
			name:      "relative symlinks inside",
			entries:   []entry{file("d/f.txt", "f"), symlink("l", "d/f.txt"), symlink("d/up", "../d/f.txt"), symlink("chain", "l"), symlink("dangling", "missing")},
			symlinks:  true,
			wantFiles: map[string]string{"l": "f", "d/up": "f", "chain": "f"},
		},
		{
			name:       "dotdot entry",
			entries:    []entry{file("../evil.txt", "x")},
			wantErr:    "escapes",
			wantAbsent: []string{"../evil.txt"},
		},
		{
			name:    "nested dotdot entry",
			entries: []entry{file("a/../../evil.txt", "x")},
			wantErr: "escapes",
		},
		{
			name:     "absolute symlink",
			entries:  []entry{symlink("l", "/etc/passwd")},
			symlinks: true,
			wantErr:  "absolute",
		},
		{
			name:     "symlink to parent",
			entries:  []entry{symlink("l", "..")},
			symlinks: true,
			wantErr:  "points outside",
		},
		{
			name:     "nested symlink escaping",
			entries:  []entry{symlink("a/b", "../../x")},
			symlinks: true,
			wantErr:  "points outside",
		},
		{
			name:     "symlink under leading slash entry escaping",
			entries:  []entry{symlink("/a/b", "../../x")},
			symlinks: true,
			wantErr:  "points outside",
		},
		{
			name:       "write through extracted symlink directory",
			entries:    []entry{dir("sub/"), symlink("l", "sub"), file("l/f.txt", "x")},
			symlinks:   true,
			wantErr:    "passes through a symlink",
			wantAbsent: []string{"sub/f.txt"},
		},
		{
			name:     "overwrite extracted symlink",
			entries:  []entry{file("f.txt", "orig"), symlink("l", "f.txt"), file("l", "evil")},
			symlinks: true,
			wantErr:  "passes through a symlink",
		},
		{
			name:     "chained symlink escaping",
			entries:  []entry{dir("d/"), symlink("d/up", ".."), symlink("esc", "d/up/..")},
			symlinks: true,
			wantErr:  "resolves outside",
		},
		{
			name:     "dangling chained symlink escaping",
			entries:  []entry{symlink("b", "."), symlink("a", "b/../outside")},
			symlinks: true,
			wantErr:  "resolves outside",
		},
		{
			name:     "symlink loop",
			entries:  []entry{symlink("x", "y"), symlink("y", "x")},
			symlinks: true,
			wantErr:  "resolves outside",
		},
		{
			name:     "symlink over existing file",
			entries:  []entry{file("a", "x"), symlink("a", "b")},
			symlinks: true,
			wantErr:  "a",
		},
		{
			name:     "symlink below file",
			entries:  []entry{file("f", "x"), symlink("f/l", "../f")},
			symlinks: true,
			wantErr:  "f",
		},
		{
			name:    "file where directory exists",
			entries: []entry{dir("d/"), file("d", "x")},
			wantErr: "d",
		},
		{
			name:    "file below file",
			entries: []entry{file("f", "x"), file("f/g", "y")},
			wantErr: "f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.symlinks && runtime.GOOS == "windows" {
				t.Skip("symbolic links need privileges on Windows")
			}
			archive := writeTarGz(t, tt.entries...)
			root := filepath.Join(t.TempDir(), "out")
			err := ExtractTarGz(archive, root)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("ExtractTarGz error = %v, want containing %q", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ExtractTarGz: %v", err)
			}
			checkFiles(t, root, tt.wantFiles)
			checkAbsent(t, root, tt.wantAbsent)
		})
	}
}

func TestExtractTarGzModes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits")
	}
	archive := writeTarGz(t, entry{name: "run.sh", body: "#!/bin/sh", mode: 0o755}, entry{name: "ro.txt", body: "r", mode: 0o444})
	root := t.TempDir()
	if err := ExtractTarGz(archive, root); err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]os.FileMode{"run.sh": 0o755, "ro.txt": 0o444} {
		info, err := os.Stat(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm() &^ 0o022; got != want&^0o022 {
			t.Errorf("%s mode = %v, want %v", name, got, want)
		}
	}
}

func TestExtractTarGzPreexistingSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symbolic links need privileges on Windows")
	}
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	err := ExtractTarGz(writeTarGz(t, file("link/evil.txt", "x")), root)
	if err == nil || !strings.Contains(err.Error(), "passes through a symlink") {
		t.Fatalf("ExtractTarGz error = %v, want passes through a symlink", err)
	}
	checkAbsent(t, outside, []string{"evil.txt"})
}

func TestExtractTarGzRootTrailingSeparator(t *testing.T) {
	root := t.TempDir() + string(filepath.Separator)
	if err := ExtractTarGz(writeTarGz(t, file("a/b.txt", "b")), root); err != nil {
		t.Fatal(err)
	}
	checkFiles(t, root, map[string]string{"a/b.txt": "b"})
}

func TestExtractTarGzRelativeRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symbolic links need privileges on Windows")
	}
	archive := writeTarGz(t, file("d/f.txt", "f"), symlink("l", "d/f.txt"))
	t.Chdir(t.TempDir())
	if err := ExtractTarGz(archive, "out"); err != nil {
		t.Fatal(err)
	}
	checkFiles(t, "out", map[string]string{"l": "f"})
}

func TestExtractTarGzInvalidArchive(t *testing.T) {
	tmp := t.TempDir()
	notGzip := filepath.Join(tmp, "plain.tar.gz")
	if err := os.WriteFile(notGzip, []byte("not gzip"), 0o600); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte("not a tar archive at all, but long enough to be read as a header"))
	_ = gz.Close()
	notTar := filepath.Join(tmp, "bad.tar.gz")
	if err := os.WriteFile(notTar, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	gz = gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "f.txt", Mode: 0o644, Size: 100, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("short"))
	_ = gz.Close()
	truncated := filepath.Join(tmp, "truncated.tar.gz")
	if err := os.WriteFile(truncated, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	valid, err := os.ReadFile(writeTarGz(t, file("a.txt", "hello")))
	if err != nil {
		t.Fatal(err)
	}
	badCRC := append([]byte(nil), valid...)
	badCRC[len(badCRC)-8] ^= 0xff
	corrupt := filepath.Join(tmp, "crc.tar.gz")
	if err := os.WriteFile(corrupt, badCRC, 0o600); err != nil {
		t.Fatal(err)
	}
	shortTrailer := filepath.Join(tmp, "trailer.tar.gz")
	if err := os.WriteFile(shortTrailer, valid[:len(valid)-4], 0o600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, path string
	}{
		{"missing", filepath.Join(tmp, "missing.tar.gz")},
		{"gzip checksum mismatch", corrupt},
		{"gzip trailer cut", shortTrailer},
		{"truncated file body", truncated},
		{"not gzip", notGzip},
		{"not tar", notTar},
	}
	for _, tt := range tests {
		if err := ExtractTarGz(tt.path, filepath.Join(tmp, "out")); err == nil {
			t.Errorf("%s: ExtractTarGz succeeded, want error", tt.name)
		}
	}
}

func TestExtractZip(t *testing.T) {
	tests := []struct {
		name       string
		entries    []entry
		wantErr    string
		wantFiles  map[string]string
		wantAbsent []string
	}{
		{
			name:      "files and directories",
			entries:   []entry{dir("d/"), file("d/a.txt", "a"), file("top.txt", "top"), file("x/y/z.txt", "z")},
			wantFiles: map[string]string{"d/a.txt": "a", "top.txt": "top", "x/y/z.txt": "z"},
		},
		{
			name:      "symlink entry is written as a regular file",
			entries:   []entry{symlink("l", "/etc/passwd")},
			wantFiles: map[string]string{"l": "/etc/passwd"},
		},
		{
			name:       "dotdot entry",
			entries:    []entry{file("../evil.txt", "x")},
			wantErr:    "escapes",
			wantAbsent: []string{"../evil.txt"},
		},
		{
			name:    "nested dotdot entry",
			entries: []entry{file("a/../../evil.txt", "x")},
			wantErr: "escapes",
		},
		{
			name:    "file where directory exists",
			entries: []entry{dir("d/"), file("d", "x")},
			wantErr: "d",
		},
		{
			name:    "directory below file",
			entries: []entry{file("f", "x"), dir("f/g/")},
			wantErr: "f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := writeZip(t, tt.entries...)
			root := filepath.Join(t.TempDir(), "out")
			err := ExtractZip(archive, root)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("ExtractZip error = %v, want containing %q", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ExtractZip: %v", err)
			}
			checkFiles(t, root, tt.wantFiles)
			checkAbsent(t, root, tt.wantAbsent)
		})
	}
}

func TestExtractZipModes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits")
	}
	root := t.TempDir()
	if err := ExtractZip(writeZip(t, entry{name: "ro.txt", body: "r", mode: 0o444}), root); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(root, "ro.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o600 != 0o600 {
		t.Errorf("mode = %v, want owner read and write", info.Mode().Perm())
	}
}

func TestExtractZipPreexistingSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symbolic links need privileges on Windows")
	}
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	err := ExtractZip(writeZip(t, file("link/evil.txt", "x")), root)
	if err == nil || !strings.Contains(err.Error(), "passes through a symlink") {
		t.Fatalf("ExtractZip error = %v, want passes through a symlink", err)
	}
	checkAbsent(t, outside, []string{"evil.txt"})
}

func TestExtractZipInvalidArchive(t *testing.T) {
	tmp := t.TempDir()
	notZip := filepath.Join(tmp, "plain.zip")
	if err := os.WriteFile(notZip, []byte("not zip"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A stored entry whose checksum does not match its content fails on read.
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.CreateRaw(&zip.FileHeader{Name: "f.txt", Method: zip.Store, CRC32: 1, CompressedSize64: 3, UncompressedSize64: 3})
	_, _ = w.Write([]byte("abc"))
	_ = zw.Close()
	badCRC := filepath.Join(tmp, "crc.zip")
	if err := os.WriteFile(badCRC, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	zw = zip.NewWriter(&buf)
	_, _ = zw.CreateRaw(&zip.FileHeader{Name: "g.txt", Method: 99})
	_ = zw.Close()
	badMethod := filepath.Join(tmp, "method.zip")
	if err := os.WriteFile(badMethod, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, path string
	}{
		{"missing", filepath.Join(tmp, "missing.zip")},
		{"not zip", notZip},
		{"bad checksum", badCRC},
		{"unknown method", badMethod},
	}
	for _, tt := range tests {
		if err := ExtractZip(tt.path, filepath.Join(tmp, "out")); err == nil {
			t.Errorf("%s: ExtractZip succeeded, want error", tt.name)
		}
	}
}

func TestExtractReadOnlyRoot(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs Unix directory permissions enforced")
	}
	tests := []struct {
		name    string
		archive func(t *testing.T) string
		extract func(archivePath, root string) error
	}{
		{"tar file in new directory", func(t *testing.T) string { return writeTarGz(t, file("d/f.txt", "x")) }, ExtractTarGz},
		{"tar symlink in new directory", func(t *testing.T) string { return writeTarGz(t, symlink("d/l", "f")) }, ExtractTarGz},
		{"zip directory", func(t *testing.T) string { return writeZip(t, dir("d/")) }, ExtractZip},
		{"zip file", func(t *testing.T) string { return writeZip(t, file("f.txt", "x")) }, ExtractZip},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := tt.archive(t)
			root := t.TempDir()
			if err := os.Chmod(root, 0o555); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(root, 0o755) })
			if err := tt.extract(archive, root); err == nil {
				t.Fatal("extraction into read-only root succeeded, want error")
			}
		})
	}
}
