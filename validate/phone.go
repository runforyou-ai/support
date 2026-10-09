package validate

import (
	"strings"
	"unicode"
)

// E164 normalizes an international phone number to E.164 form. It trims s,
// which must then start with '+' directly followed by a country calling code
// digit from 1 to 9. It removes Unicode whitespace, '-', '(' and ')' from the
// rest and requires 6 to 15 ASCII digits in total and nothing else. It returns
// the normalized number and true, or "" and false when s does not match. The
// country calling code is not checked against any numbering plan.
func E164(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '+' || s[1] < '1' || s[1] > '9' {
		return "", false
	}
	n := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '-' || r == '(' || r == ')' {
			return -1
		}
		return r
	}, s)
	digits := n[1:]
	if len(digits) < 6 || len(digits) > 15 {
		return "", false
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return "", false
		}
	}
	return n, true
}
