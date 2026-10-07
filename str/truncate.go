package str

import (
	"strings"
	"unicode"
)

// Limit truncates s to at most n runes and appends end when s was truncated.
// Trailing whitespace of the kept part is removed before end is appended. A
// negative n is treated as zero.
func Limit(s string, n int, end string) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	n = max(n, 0)
	return strings.TrimRightFunc(string(runes[:n]), unicode.IsSpace) + end
}

// Words truncates s to its first n whitespace-separated words and appends end
// when s was truncated. Trailing whitespace of the kept part is removed before
// end is appended. A negative n is treated as zero.
func Words(s string, n int, end string) string {
	n = max(n, 0)
	count := 0
	inWord := false
	for i, r := range s {
		space := unicode.IsSpace(r)
		if !space && !inWord {
			if count == n {
				return strings.TrimRightFunc(s[:i], unicode.IsSpace) + end
			}
			count++
		}
		inWord = !space
	}
	return s
}

// Excerpt extracts the first occurrence of phrase in text together with up to
// radius runes on each side, matching phrase case-insensitively. Each side is
// trimmed of surrounding whitespace, and omission is added to a side that was
// cut short. It returns false when phrase is not found. A negative radius is
// treated as zero.
func Excerpt(text, phrase string, radius int, omission string) (string, bool) {
	radius = max(radius, 0)
	runes := []rune(text)
	n := len([]rune(phrase))
	at := -1
	for i := 0; i+n <= len(runes); i++ {
		if strings.EqualFold(string(runes[i:i+n]), phrase) {
			at = i
			break
		}
	}
	if at < 0 {
		return "", false
	}

	before := []rune(strings.TrimLeftFunc(string(runes[:at]), unicode.IsSpace))
	start := strings.TrimLeftFunc(string(before[max(len(before)-radius, 0):]), unicode.IsSpace)
	if start != string(before) {
		start = omission + start
	}

	after := []rune(strings.TrimRightFunc(string(runes[at+n:]), unicode.IsSpace))
	end := strings.TrimRightFunc(string(after[:min(radius, len(after))]), unicode.IsSpace)
	if end != string(after) {
		end += omission
	}

	return start + string(runes[at:at+n]) + end, true
}
