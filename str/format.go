package str

import (
	"strings"
	"unicode"
)

// Mask replaces length runes of s starting at rune offset index with char. A
// negative index counts from the end of s and is clamped to the beginning. A
// length of zero or less masks everything up to the end. It returns s
// unchanged when the range is empty.
func Mask(s string, char rune, index, length int) string {
	runes := []rune(s)
	n := len(runes)
	if index < 0 {
		index = max(n+index, 0)
	}
	end := n
	if length > 0 && length < n-index {
		end = index + length
	}
	if index >= end {
		return s
	}
	out := make([]rune, n)
	copy(out, runes)
	for i := index; i < end; i++ {
		out[i] = char
	}
	return string(out)
}

// PadLeft pads the left side of s with pad until it is length runes long. A
// multi-character pad is repeated and cut to fit. It returns s unchanged when
// pad is empty or s is already long enough.
func PadLeft(s string, length int, pad string) string {
	return padding(pad, length-Length(s)) + s
}

// PadRight pads the right side of s with pad until it is length runes long. A
// multi-character pad is repeated and cut to fit. It returns s unchanged when
// pad is empty or s is already long enough.
func PadRight(s string, length int, pad string) string {
	return s + padding(pad, length-Length(s))
}

// PadBoth pads both sides of s with pad until it is length runes long, putting
// the extra rune on the right when the padding is odd. A multi-character pad
// is repeated and cut to fit. It returns s unchanged when pad is empty or s is
// already long enough.
func PadBoth(s string, length int, pad string) string {
	total := length - Length(s)
	if total <= 0 {
		return s
	}
	return padding(pad, total/2) + s + padding(pad, total-total/2)
}

// padding returns pad repeated and cut to exactly n runes.
func padding(pad string, n int) string {
	if n <= 0 || pad == "" {
		return ""
	}
	runes := []rune(pad)
	out := make([]rune, n)
	for i := range out {
		out[i] = runes[i%len(runes)]
	}
	return string(out)
}

// Start returns s with a single leading prefix, collapsing any repeated
// prefixes. It returns s unchanged when prefix is empty.
func Start(s, prefix string) string {
	if prefix == "" {
		return s
	}
	for strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}
	return prefix + s
}

// Finish returns s with a single trailing suffix, collapsing any repeated
// suffixes. It returns s unchanged when suffix is empty.
func Finish(s, suffix string) string {
	if suffix == "" {
		return s
	}
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s + suffix
}

// Wrap surrounds s with before and after. An empty after means before is used
// on both sides.
func Wrap(s, before, after string) string {
	if after == "" {
		after = before
	}
	return before + s + after
}

// Unwrap removes before from the start of s and after from the end of s,
// each only when present. An empty after means before is used on both sides.
func Unwrap(s, before, after string) string {
	if after == "" {
		after = before
	}
	if before != "" {
		s = strings.TrimPrefix(s, before)
	}
	if after != "" {
		s = strings.TrimSuffix(s, after)
	}
	return s
}

// Squish removes leading and trailing whitespace from s and collapses every
// inner run of Unicode whitespace, including the Hangul fillers U+1160 and
// U+3164, into a single space.
func Squish(s string) string {
	return strings.Join(strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == 'ᅠ' || r == 'ㅤ'
	}), " ")
}

// Reverse returns s with its runes in reverse order.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Slug converts s to a URL-friendly slug. Letters are lowercased, Unicode
// letters, digits and marks are kept (including Chinese), and every other run
// of runes is replaced by a single separator. Leading and trailing separators
// are removed. Letters are not transliterated.
func Slug(s, separator string) string {
	var b strings.Builder
	pending := false
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsMark(r) {
			pending = true
			continue
		}
		if pending && b.Len() > 0 {
			b.WriteString(separator)
		}
		pending = false
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
