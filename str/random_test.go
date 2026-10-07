package str

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	tests := []struct {
		n, want int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{16, 16},
		{1000, 1000},
	}
	for _, tt := range tests {
		got := Random(tt.n)
		if len(got) != tt.want {
			t.Errorf("len(Random(%d)) = %d, want %d", tt.n, len(got), tt.want)
		}
		for _, r := range got {
			if !strings.ContainsRune(alphanumeric, r) {
				t.Errorf("Random(%d) contains %q", tt.n, r)
			}
		}
	}

	seen := map[string]bool{}
	for range 1000 {
		s := Random(16)
		if seen[s] {
			t.Fatalf("Random(16) repeated %q", s)
		}
		seen[s] = true
	}

	counts := map[rune]int{}
	for _, r := range Random(62 * 1000) {
		counts[r]++
	}
	if len(counts) != len(alphanumeric) {
		t.Errorf("Random used %d distinct characters, want %d", len(counts), len(alphanumeric))
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		name    string
		fn      func() string
		version byte
	}{
		{"UUID", UUID, '4'},
		{"UUIDv7", UUIDv7, '7'},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := map[string]bool{}
			for range 1000 {
				id := tt.fn()
				if !IsUUID(id) || id != strings.ToLower(id) {
					t.Fatalf("%s() = %q, not a canonical lowercase UUID", tt.name, id)
				}
				if id[14] != tt.version {
					t.Fatalf("%s() = %q, version %c, want %c", tt.name, id, id[14], tt.version)
				}
				if !strings.ContainsRune("89ab", rune(id[19])) {
					t.Fatalf("%s() = %q, wrong variant", tt.name, id)
				}
				if seen[id] {
					t.Fatalf("%s() repeated %q", tt.name, id)
				}
				seen[id] = true
			}
		})
	}
}

func TestUUIDv7Timestamp(t *testing.T) {
	before := time.Now().UnixMilli()
	id := UUIDv7()
	after := time.Now().UnixMilli()
	ms, err := strconv.ParseInt(strings.ReplaceAll(id[:13], "-", ""), 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	if ms < before || ms > after {
		t.Errorf("UUIDv7 timestamp %d outside [%d, %d]", ms, before, after)
	}
}

func TestULID(t *testing.T) {
	seen := map[string]bool{}
	for range 1000 {
		id := ULID()
		if !IsULID(id) || id != strings.ToUpper(id) {
			t.Fatalf("ULID() = %q, not an uppercase ULID", id)
		}
		if seen[id] {
			t.Fatalf("ULID() repeated %q", id)
		}
		seen[id] = true
	}

	before := time.Now().UnixMilli()
	id := ULID()
	after := time.Now().UnixMilli()
	var ms int64
	for _, c := range id[:10] {
		ms = ms<<5 | int64(strings.IndexRune(crockford, c))
	}
	if ms < before || ms > after {
		t.Errorf("ULID timestamp %d outside [%d, %d]", ms, before, after)
	}
}
