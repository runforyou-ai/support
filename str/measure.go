package str

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Length returns the number of runes in s.
func Length(s string) int {
	return utf8.RuneCountInString(s)
}

// SubstrCount returns the number of non-overlapping occurrences of sub in s.
// It returns 0 when sub is empty.
func SubstrCount(s, sub string) int {
	if sub == "" {
		return 0
	}
	return strings.Count(s, sub)
}

// WordCount returns the number of whitespace-separated words in s, ignoring
// fields that contain no letter or digit such as a lone dash. Text without
// spaces, such as a Chinese sentence, counts as a single word.
func WordCount(s string) int {
	count := 0
	for _, field := range strings.Fields(s) {
		if strings.IndexFunc(field, isWordRune) >= 0 {
			count++
		}
	}
	return count
}

// isWordRune reports whether r is a letter or digit.
func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
