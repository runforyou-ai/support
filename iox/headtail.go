package iox

import (
	"strings"
	"sync"
	"unicode/utf8"
)

// HeadTailBuffer is an io.Writer that keeps the first head bytes and the most
// recent tail bytes written to it and discards everything in between. It is
// safe for concurrent use, so the same buffer can collect several streams.
// Create it with NewHeadTailBuffer.
type HeadTailBuffer struct {
	mu        sync.Mutex
	headLimit int
	tailLimit int
	marker    string
	head      []byte
	tail      []byte
	dropped   bool
}

// NewHeadTailBuffer returns a buffer that keeps the first head bytes and the
// last tail bytes written to it. marker is inserted by String between the two
// parts when bytes were dropped. Negative limits are treated as zero.
func NewHeadTailBuffer(head, tail int, marker string) *HeadTailBuffer {
	return &HeadTailBuffer{headLimit: max(head, 0), tailLimit: max(tail, 0), marker: marker}
}

// Write appends p: it first fills the head, then appends the rest to the tail
// and discards the oldest tail bytes beyond the tail limit. It always reports
// len(p) bytes written and a nil error, and never retains p.
func (b *HeadTailBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	written := len(p)
	if room := b.headLimit - len(b.head); room > 0 {
		n := min(room, len(p))
		b.head = append(b.head, p[:n]...)
		p = p[n:]
	}
	// Bytes that cannot fit in the tail are dropped before being copied.
	if len(p) > b.tailLimit {
		p = p[len(p)-b.tailLimit:]
		b.dropped = true
	}
	b.tail = append(b.tail, p...)
	if excess := len(b.tail) - b.tailLimit; excess > 0 {
		b.tail = append(b.tail[:0], b.tail[excess:]...)
		b.dropped = true
	}
	return written, nil
}

// String returns the retained output as valid UTF-8. When bytes were dropped
// it returns head + marker + tail; a multi-byte character split by a cut is
// removed from either side, and any other invalid byte sequence is replaced
// with U+FFFD.
func (b *HeadTailBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.dropped {
		return strings.ToValidUTF8(string(b.head)+string(b.tail), "�")
	}
	head := b.head
	// Drop a trailing character whose remaining bytes were discarded.
	for i := len(head) - 1; i >= 0 && i >= len(head)-utf8.UTFMax; i-- {
		if utf8.RuneStart(head[i]) {
			if !utf8.FullRune(head[i:]) {
				head = head[:i]
			}
			break
		}
	}
	tail := b.tail
	// Drop leading continuation bytes whose first byte was discarded.
	for n := 0; n < utf8.UTFMax-1 && len(tail) > 0 && !utf8.RuneStart(tail[0]); n++ {
		tail = tail[1:]
	}
	return strings.ToValidUTF8(string(head), "�") + b.marker + strings.ToValidUTF8(string(tail), "�")
}

// Truncated reports whether any bytes were dropped between head and tail.
func (b *HeadTailBuffer) Truncated() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dropped
}
