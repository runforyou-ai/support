package str

import (
	"encoding/json"
	"net/url"
	"strings"
)

// Contains reports whether s contains any of the non-empty needles. It
// returns false when no non-empty needle is given.
func Contains(s string, needles ...string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(s, n) {
			return true
		}
	}
	return false
}

// ContainsAll reports whether s contains every non-empty needle. Empty needles
// are ignored, so it returns true when no non-empty needle is given.
func ContainsAll(s string, needles ...string) bool {
	for _, n := range needles {
		if n != "" && !strings.Contains(s, n) {
			return false
		}
	}
	return true
}

// Is reports whether the whole of s matches pattern, where each * matches any
// sequence of characters, including an empty one, and every other character
// matches itself. The match is case-sensitive.
func Is(pattern, s string) bool {
	p, i := 0, 0
	star, mark := -1, 0
	for i < len(s) {
		switch {
		case p < len(pattern) && pattern[p] == '*':
			star, mark = p, i
			p++
		case p < len(pattern) && pattern[p] == s[i]:
			p++
			i++
		case star >= 0:
			p = star + 1
			mark++
			i = mark
		default:
			return false
		}
	}
	for p < len(pattern) && pattern[p] == '*' {
		p++
	}
	return p == len(pattern)
}

// IsJSON reports whether s is a valid JSON value.
func IsJSON(s string) bool {
	return json.Valid([]byte(s))
}

// IsURL reports whether s is an absolute URL with both a scheme and a host.
func IsURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsUUID reports whether s is a UUID of any version in the canonical
// 8-4-4-4-12 hexadecimal form, ignoring case.
func IsUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !isHex(c) {
			return false
		}
	}
	return true
}

// NormalizeUUID trims surrounding whitespace from s and returns it in
// lowercase when the rest is a UUID accepted by IsUUID. It returns "" and
// false otherwise, so braced and "urn:uuid:" forms are rejected.
func NormalizeUUID(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if !IsUUID(s) {
		return "", false
	}
	return strings.ToLower(s), true
}

// isHex reports whether c is an ASCII hexadecimal digit.
func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// IsULID reports whether s is a ULID: 26 characters of Crockford's base32
// alphabet, ignoring case, whose first character is between 0 and 7.
func IsULID(s string) bool {
	if len(s) != 26 || s[0] < '0' || s[0] > '7' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !strings.ContainsRune(crockford, rune(upperASCII(s[i]))) {
			return false
		}
	}
	return true
}

// upperASCII converts an ASCII lowercase letter to upper case.
func upperASCII(c byte) byte {
	if 'a' <= c && c <= 'z' {
		return c - 'a' + 'A'
	}
	return c
}

// IsASCII reports whether s contains only ASCII characters.
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}
