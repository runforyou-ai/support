package str

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Camel converts s to camelCase. Words are split as described for Snake and
// lowercased before the first rune of every word but the first is
// capitalized, so "HTTPServer" becomes "httpServer".
func Camel(s string) string {
	words := splitWords(s)
	var b strings.Builder
	for i, w := range words {
		w = strings.ToLower(w)
		if i > 0 {
			w = Ucfirst(w)
		}
		b.WriteString(w)
	}
	return b.String()
}

// Studly converts s to StudlyCase (PascalCase). Words are split as described
// for Snake and lowercased before their first rune is capitalized, so
// "HTTP_SERVER" becomes "HttpServer".
func Studly(s string) string {
	var b strings.Builder
	for _, w := range splitWords(s) {
		b.WriteString(Ucfirst(strings.ToLower(w)))
	}
	return b.String()
}

// Snake converts s to snake_case. Words are separated by any rune that is not
// a letter, digit or mark, by a lowercase letter, digit or uncased letter
// followed by an uppercase letter, and before the last uppercase letter of an
// acronym followed by a lowercase letter. Digits stay attached to the
// preceding word, so "HTTPServer" becomes "http_server" and "userID2" becomes
// "user_id2". Letters without case, such as Chinese, are kept as they are.
func Snake(s string) string {
	return joinLower(s, "_")
}

// Kebab converts s to kebab-case, splitting words as described for Snake.
func Kebab(s string) string {
	return joinLower(s, "-")
}

// Headline converts s to space-separated words with the first rune of each
// word capitalized, splitting words as described for Snake. The remaining
// runes keep their case, so "EmailNotificationSent" becomes
// "Email Notification Sent" and "HTTPServer" becomes "HTTP Server".
func Headline(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = Ucfirst(w)
	}
	return strings.Join(words, " ")
}

// Title capitalizes the first letter of each word in s and lowercases the
// remaining letters. A word is a run of letters, digits, marks and
// apostrophes; every other rune is kept in place as a boundary.
func Title(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inWord := false
	for _, r := range s {
		isWord := unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) ||
			(inWord && (r == '\'' || r == '’'))
		switch {
		case isWord && !inWord:
			b.WriteRune(unicode.ToTitle(r))
		case isWord:
			b.WriteRune(unicode.ToLower(r))
		default:
			b.WriteRune(r)
		}
		inWord = isWord
	}
	return b.String()
}

// Ucfirst returns s with its first rune converted to upper case.
func Ucfirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 || r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

// Lcfirst returns s with its first rune converted to lower case.
func Lcfirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 || r == utf8.RuneError {
		return s
	}
	return string(unicode.ToLower(r)) + s[size:]
}

// UcSplit splits s before every uppercase letter, so "FooBar" becomes
// ["Foo", "Bar"] and "HTTP" becomes ["H", "T", "T", "P"]. It returns nil when
// s is empty.
func UcSplit(s string) []string {
	var parts []string
	start := 0
	for i, r := range s {
		if i > start && unicode.IsUpper(r) {
			parts = append(parts, s[start:i])
			start = i
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

// joinLower splits s into words, lowercases them and joins them with sep.
func joinLower(s, sep string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, sep)
}

// splitWords splits s into words for case conversion as described for Snake.
func splitWords(s string) []string {
	runes := []rune(s)
	var words []string
	var current []rune
	flush := func() {
		if len(current) > 0 {
			words = append(words, string(current))
			current = current[:0]
		}
	}
	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsMark(r) {
			flush()
			continue
		}
		if len(current) > 0 && unicode.IsUpper(r) {
			prev := current[len(current)-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if !unicode.IsUpper(prev) || nextLower {
				flush()
			}
		}
		current = append(current, r)
	}
	flush()
	return words
}
