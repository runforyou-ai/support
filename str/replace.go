package str

import (
	"cmp"
	"slices"
	"strings"
	"unicode/utf8"
)

// ReplaceFirst replaces the first occurrence of search in s with replace. It
// returns s unchanged when search is empty or not found.
func ReplaceFirst(s, search, replace string) string {
	if search == "" {
		return s
	}
	return strings.Replace(s, search, replace, 1)
}

// ReplaceLast replaces the last occurrence of search in s with replace. It
// returns s unchanged when search is empty or not found.
func ReplaceLast(s, search, replace string) string {
	if search == "" {
		return s
	}
	i := strings.LastIndex(s, search)
	if i < 0 {
		return s
	}
	return s[:i] + replace + s[i+len(search):]
}

// ReplaceStart replaces search with replace only when s starts with search.
// It returns s unchanged when search is empty.
func ReplaceStart(s, search, replace string) string {
	if search == "" || !strings.HasPrefix(s, search) {
		return s
	}
	return replace + s[len(search):]
}

// ReplaceEnd replaces search with replace only when s ends with search. It
// returns s unchanged when search is empty.
func ReplaceEnd(s, search, replace string) string {
	if search == "" || !strings.HasSuffix(s, search) {
		return s
	}
	return s[:len(s)-len(search)] + replace
}

// ReplaceArray replaces the occurrences of search in s one by one with the
// successive elements of replaces. Occurrences beyond the end of replaces are
// left as they are. It returns s unchanged when search is empty.
func ReplaceArray(s, search string, replaces []string) string {
	if search == "" {
		return s
	}
	segments := strings.Split(s, search)
	var b strings.Builder
	b.WriteString(segments[0])
	for i, segment := range segments[1:] {
		if i < len(replaces) {
			b.WriteString(replaces[i])
		} else {
			b.WriteString(search)
		}
		b.WriteString(segment)
	}
	return b.String()
}

// Swap replaces every key of pairs found in s with its value in a single pass,
// like PHP's strtr: replaced text is never searched again, and at each
// position the longest matching key wins. Keys of equal length are tried in
// lexical order. Empty keys are ignored.
func Swap(s string, pairs map[string]string) string {
	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		if k != "" {
			keys = append(keys, k)
		}
	}
	if len(keys) == 0 {
		return s
	}
	slices.SortFunc(keys, func(a, b string) int {
		if c := cmp.Compare(len(b), len(a)); c != 0 {
			return c
		}
		return strings.Compare(a, b)
	})

	var b strings.Builder
	for i := 0; i < len(s); {
		matched := false
		for _, k := range keys {
			if strings.HasPrefix(s[i:], k) {
				b.WriteString(pairs[k])
				i += len(k)
				matched = true
				break
			}
		}
		if !matched {
			_, size := utf8.DecodeRuneInString(s[i:])
			b.WriteString(s[i : i+size])
			i += size
		}
	}
	return b.String()
}

// Remove removes every occurrence of each search string from s, applying the
// searches in order. Empty search strings are ignored.
func Remove(s string, search ...string) string {
	for _, v := range search {
		if v != "" {
			s = strings.ReplaceAll(s, v, "")
		}
	}
	return s
}
