package validate

import (
	"strings"
	"unicode"
)

// E164 normalizes an international phone number to E.164 form. It trims s,
// removes Unicode whitespace, '-', '(' and ')', and then requires a leading
// '+' followed by 6 to 15 ASCII digits and nothing else. It returns the
// normalized number and true, or "" and false when s does not match. The
// country calling code is not checked against any numbering plan.
func E164(s string) (string, bool) {
	n := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '-' || r == '(' || r == ')' {
			return -1
		}
		return r
	}, strings.TrimSpace(s))
	digits, found := strings.CutPrefix(n, "+")
	if !found || len(digits) < 6 || len(digits) > 15 {
		return "", false
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return "", false
		}
	}
	return n, true
}
