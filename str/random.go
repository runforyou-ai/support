package str

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// alphanumeric is the alphabet used by Random.
const alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// crockford is Crockford's base32 alphabet used by ULID.
const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// Random returns a cryptographically secure random string of n characters
// drawn uniformly from [A-Za-z0-9]. It returns an empty string when n is zero
// or negative.
func Random(n int) string {
	if n <= 0 {
		return ""
	}
	// 248 is the largest multiple of 62 not above 256; larger bytes are
	// rejected so that every character is equally likely.
	const limit = 256 - 256%len(alphanumeric)
	out := make([]byte, 0, n)
	buf := make([]byte, n+n/4+8)
	for len(out) < n {
		randomBytes(buf)
		for _, c := range buf {
			if int(c) < limit {
				out = append(out, alphanumeric[int(c)%len(alphanumeric)])
				if len(out) == n {
					break
				}
			}
		}
	}
	return string(out)
}

// UUID returns a random version 4 UUID in canonical lowercase form.
func UUID() string {
	var b [16]byte
	randomBytes(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return formatUUID(b)
}

// UUIDv7 returns a version 7 UUID in canonical lowercase form, combining the
// current Unix time in milliseconds with random bits so that values sort by
// creation time at millisecond precision.
func UUIDv7() string {
	var b [16]byte
	randomBytes(b[6:])
	putMillis(b[:6], time.Now().UnixMilli())
	b[6] = b[6]&0x0f | 0x70
	b[8] = b[8]&0x3f | 0x80
	return formatUUID(b)
}

// ULID returns a ULID combining the current Unix time in milliseconds with 80
// random bits, encoded as 26 uppercase Crockford base32 characters. ULIDs
// generated within the same millisecond are not guaranteed to be ordered.
func ULID() string {
	var b [16]byte
	randomBytes(b[6:])
	putMillis(b[:6], time.Now().UnixMilli())
	var hi, lo uint64
	for i := 0; i < 8; i++ {
		hi = hi<<8 | uint64(b[i])
		lo = lo<<8 | uint64(b[i+8])
	}
	var out [26]byte
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = crockford[lo&31]
		lo = lo>>5 | hi<<59
		hi >>= 5
	}
	return string(out[:])
}

// randomBytes fills b with cryptographically secure random bytes.
func randomBytes(b []byte) {
	// crypto/rand.Read never returns an error since Go 1.24.
	_, _ = rand.Read(b)
}

// putMillis writes the low 48 bits of ms into b in big-endian order.
func putMillis(b []byte, ms int64) {
	for i := 5; i >= 0; i-- {
		b[i] = byte(ms)
		ms >>= 8
	}
}

// formatUUID formats b in the canonical 8-4-4-4-12 form.
func formatUUID(b [16]byte) string {
	var out [36]byte
	hex.Encode(out[0:8], b[0:4])
	out[8] = '-'
	hex.Encode(out[9:13], b[4:6])
	out[13] = '-'
	hex.Encode(out[14:18], b[6:8])
	out[18] = '-'
	hex.Encode(out[19:23], b[8:10])
	out[23] = '-'
	hex.Encode(out[24:], b[10:])
	return string(out[:])
}
