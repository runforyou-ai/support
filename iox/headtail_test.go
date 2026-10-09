package iox

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestHeadTailBuffer(t *testing.T) {
	const marker = "<...>"
	tests := []struct {
		name          string
		head, tail    int
		writes        []string
		want          string
		wantTruncated bool
	}{
		{"empty", 3, 3, nil, "", false},
		{"fits in head", 5, 5, []string{"abc"}, "abc", false},
		{"fits in head and tail", 3, 3, []string{"abcdef"}, "abcdef", false},
		{"fits across writes", 3, 3, []string{"ab", "cd", "ef"}, "abcdef", false},
		{"truncated single write", 2, 2, []string{"abcdefg"}, "ab<...>fg", true},
		{"truncated across writes", 2, 3, []string{"ab", "cd", "efg", "h"}, "ab<...>fgh", true},
		{"head only", 3, 0, []string{"abcdef"}, "abc<...>", true},
		{"tail only", 0, 3, []string{"abcdef"}, "<...>def", true},
		{"zero limits", 0, 0, []string{"abc"}, "<...>", true},
		{"zero limits no writes", 0, 0, nil, "", false},
		{"empty write", 0, 0, []string{""}, "", false},
		{"negative limits", -1, -1, []string{"abc"}, "<...>", true},
		{"utf8 kept when contiguous", 2, 2, []string{"中"}, "中", false},
		{"utf8 split in head", 4, 3, []string{"ab中文字"}, "ab<...>字", true},
		{"utf8 split in tail", 3, 4, []string{"abc中文字"}, "abc<...>字", true},
		{"utf8 split both sides", 2, 2, []string{"中文字"}, "<...>", true},
		{"utf8 whole rune both sides", 3, 3, []string{"中文字"}, "中<...>字", true},
		{"invalid bytes replaced", 4, 4, []string{"a\xffb"}, "a�b", false},
		{"invalid bytes replaced when truncated", 2, 2, []string{"\xffbcd\xff"}, "�b<...>d�", true},
		{"exact fit then one more byte", 2, 2, []string{"abcd", "e"}, "ab<...>de", true},
		{"tail rolls over many small writes", 1, 2, []string{"a", "b", "c", "d", "e"}, "a<...>de", true},
		{"head filled by later write", 3, 1, []string{"a", "bcd", "ef"}, "abc<...>f", true},
		{"four-byte rune split in head", 2, 0, []string{"a😀"}, "a<...>", true},
		{"four-byte rune split in tail", 0, 3, []string{"😀b"}, "<...>b", true},
		{"four-byte rune kept in tail", 0, 4, []string{"a😀"}, "<...>😀", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewHeadTailBuffer(tt.head, tt.tail, marker)
			for _, w := range tt.writes {
				n, err := b.Write([]byte(w))
				if n != len(w) || err != nil {
					t.Fatalf("Write(%q) = %d, %v", w, n, err)
				}
			}
			if got := b.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
			if got := b.Truncated(); got != tt.wantTruncated {
				t.Errorf("Truncated() = %v, want %v", got, tt.wantTruncated)
			}
		})
	}
}

func TestHeadTailBufferEmptyMarker(t *testing.T) {
	b := NewHeadTailBuffer(1, 1, "")
	_, _ = b.Write([]byte("abc"))
	if got := b.String(); got != "ac" || !b.Truncated() {
		t.Errorf("String() = %q, Truncated() = %v; want \"ac\", true", got, b.Truncated())
	}
	if got := b.String(); got != "ac" {
		t.Errorf("second String() = %q, want unchanged", got)
	}
}

func TestHeadTailBufferDoesNotRetainInput(t *testing.T) {
	b := NewHeadTailBuffer(2, 2, "|")
	p := []byte("abcdef")
	_, _ = b.Write(p)
	copy(p, "zzzzzz")
	if got := b.String(); got != "ab|ef" {
		t.Errorf("String() = %q after caller reused input, want %q", got, "ab|ef")
	}
}

func TestHeadTailBufferConcurrent(t *testing.T) {
	b := NewHeadTailBuffer(100, 100, "|")
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			for j := range 100 {
				_, _ = fmt.Fprintf(b, "%d:%d\n", i, j)
				_ = b.String()
				_ = b.Truncated()
			}
		})
	}
	wg.Wait()
	got := b.String()
	if !b.Truncated() || !strings.Contains(got, "|") || len(got) != 201 {
		t.Errorf("String() = %q (len %d), want truncated 100+1+100 bytes", got, len(got))
	}
}
