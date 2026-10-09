package str

import "strings"

// After returns the part of s after the first occurrence of search. It
// returns s unchanged when search is empty or not found.
func After(s, search string) string {
	if search == "" {
		return s
	}
	if i := strings.Index(s, search); i >= 0 {
		return s[i+len(search):]
	}
	return s
}

// AfterLast returns the part of s after the last occurrence of search. It
// returns s unchanged when search is empty or not found.
func AfterLast(s, search string) string {
	if search == "" {
		return s
	}
	if i := strings.LastIndex(s, search); i >= 0 {
		return s[i+len(search):]
	}
	return s
}

// Before returns the part of s before the first occurrence of search. It
// returns s unchanged when search is empty or not found.
func Before(s, search string) string {
	if search == "" {
		return s
	}
	if i := strings.Index(s, search); i >= 0 {
		return s[:i]
	}
	return s
}

// BeforeLast returns the part of s before the last occurrence of search. It
// returns s unchanged when search is empty or not found.
func BeforeLast(s, search string) string {
	if search == "" {
		return s
	}
	if i := strings.LastIndex(s, search); i >= 0 {
		return s[:i]
	}
	return s
}

// Between returns the part of s after the first occurrence of from and before
// the last occurrence of to. It returns s unchanged when from or to is empty;
// a delimiter that is not found leaves that side untouched.
func Between(s, from, to string) string {
	if from == "" || to == "" {
		return s
	}
	return BeforeLast(After(s, from), to)
}

// BetweenFirst returns the part of s after the first occurrence of from and
// before the first occurrence of to that follows it. It returns s unchanged
// when from or to is empty; a delimiter that is not found leaves that side
// untouched.
func BetweenFirst(s, from, to string) string {
	if from == "" || to == "" {
		return s
	}
	return Before(After(s, from), to)
}

// CharAt returns the character at the given rune index of s. A negative index
// counts from the end, so -1 is the last character. It returns false when the
// index is out of range. A character is a single rune, so a letter with a
// combining mark or an emoji sequence joined by U+200D spans several indexes.
func CharAt(s string, index int) (string, bool) {
	runes := []rune(s)
	if index < 0 {
		index += len(runes)
	}
	if index < 0 || index >= len(runes) {
		return "", false
	}
	return string(runes[index]), true
}

// Substr returns the part of s starting at rune offset start and spanning
// length runes, following PHP's mb_substr. A negative start counts from the
// end of s and is clamped to the beginning. A negative length omits that many
// runes from the end of s, and a length of zero returns an empty string; pass
// Length(s) to take everything up to the end. Out-of-range values return an
// empty string. Offsets count runes, not grapheme clusters, so a cut may
// separate a combining mark or split an emoji sequence joined by U+200D.
func Substr(s string, start, length int) string {
	runes := []rune(s)
	n := len(runes)
	if start < 0 {
		start = max(n+start, 0)
	}
	if start >= n {
		return ""
	}
	end := n
	if length < 0 {
		end = n + length
	} else if length < n-start {
		end = start + length
	}
	if end <= start {
		return ""
	}
	return string(runes[start:end])
}
