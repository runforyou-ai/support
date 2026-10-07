package filex

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWithin(t *testing.T) {
	sep := string(filepath.Separator)
	abs := filepath.Join(t.TempDir(), "root")
	tests := []struct {
		name       string
		root, path string
		want       bool
	}{
		{"same", abs, abs, true},
		{"same with trailing separator", abs + sep, abs, true},
		{"child", abs, filepath.Join(abs, "a", "b"), true},
		{"unclean child", abs, abs + sep + "a" + sep + ".." + sep + "b", true},
		{"parent", abs, filepath.Dir(abs), false},
		{"escape via dotdot", abs, abs + sep + ".." + sep + "x", false},
		{"sibling with shared prefix", abs, abs + "2", false},
		{"dotdot-prefixed name", abs, filepath.Join(abs, "..x"), true},
		{"relative child", "dir", filepath.Join("dir", "f"), true},
		{"relative escape", "dir", filepath.Join("dir", "..", "f"), false},
		{"dot root", ".", "f", true},
		{"dot root parent", ".", "..", false},
		{"empty root", "", "f", true},
		{"empty both", "", "", true},
		{"relative path absolute root", abs, "f", false},
		{"absolute path relative root", "dir", abs, false},
	}
	for _, tt := range tests {
		if got := Within(tt.root, tt.path); got != tt.want {
			t.Errorf("%s: Within(%q, %q) = %v, want %v", tt.name, tt.root, tt.path, got, tt.want)
		}
	}
}

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"nil", nil, false},
		{"empty", []byte{}, false},
		{"text", []byte("hello\nworld"), false},
		{"utf8", []byte("你好"), false},
		{"nul", []byte("a\x00b"), true},
		{"nul first", []byte{0}, true},
		{"nul at last sniffed byte", []byte(strings.Repeat("a", 7999) + "\x00"), true},
		{"nul after sniff window", []byte(strings.Repeat("a", 8000) + "\x00"), false},
	}
	for _, tt := range tests {
		if got := IsBinary(tt.data); got != tt.want {
			t.Errorf("%s: IsBinary = %v, want %v", tt.name, got, tt.want)
		}
	}
}
