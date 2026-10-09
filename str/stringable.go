package str

import (
	"strings"
	"unicode"
)

// Stringable is an immutable string wrapper whose methods return new values
// for fluent chaining. The zero value wraps the empty string.
type Stringable struct {
	value string
}

// Of returns a Stringable wrapping s.
func Of(s string) Stringable {
	return Stringable{value: s}
}

// String returns the wrapped string.
func (s Stringable) String() string {
	return s.value
}

// After returns the result of After on the wrapped string.
func (s Stringable) After(search string) Stringable {
	return Of(After(s.value, search))
}

// AfterLast returns the result of AfterLast on the wrapped string.
func (s Stringable) AfterLast(search string) Stringable {
	return Of(AfterLast(s.value, search))
}

// Before returns the result of Before on the wrapped string.
func (s Stringable) Before(search string) Stringable {
	return Of(Before(s.value, search))
}

// BeforeLast returns the result of BeforeLast on the wrapped string.
func (s Stringable) BeforeLast(search string) Stringable {
	return Of(BeforeLast(s.value, search))
}

// Between returns the result of Between on the wrapped string.
func (s Stringable) Between(from, to string) Stringable {
	return Of(Between(s.value, from, to))
}

// BetweenFirst returns the result of BetweenFirst on the wrapped string.
func (s Stringable) BetweenFirst(from, to string) Stringable {
	return Of(BetweenFirst(s.value, from, to))
}

// Camel returns the result of Camel on the wrapped string.
func (s Stringable) Camel() Stringable {
	return Of(Camel(s.value))
}

// Studly returns the result of Studly on the wrapped string.
func (s Stringable) Studly() Stringable {
	return Of(Studly(s.value))
}

// Snake returns the result of Snake on the wrapped string.
func (s Stringable) Snake() Stringable {
	return Of(Snake(s.value))
}

// Kebab returns the result of Kebab on the wrapped string.
func (s Stringable) Kebab() Stringable {
	return Of(Kebab(s.value))
}

// Headline returns the result of Headline on the wrapped string.
func (s Stringable) Headline() Stringable {
	return Of(Headline(s.value))
}

// Title returns the result of Title on the wrapped string.
func (s Stringable) Title() Stringable {
	return Of(Title(s.value))
}

// Ucfirst returns the result of Ucfirst on the wrapped string.
func (s Stringable) Ucfirst() Stringable {
	return Of(Ucfirst(s.value))
}

// Lcfirst returns the result of Lcfirst on the wrapped string.
func (s Stringable) Lcfirst() Stringable {
	return Of(Lcfirst(s.value))
}

// Limit returns the result of Limit on the wrapped string.
func (s Stringable) Limit(n int, end string) Stringable {
	return Of(Limit(s.value, n, end))
}

// Words returns the result of Words on the wrapped string.
func (s Stringable) Words(n int, end string) Stringable {
	return Of(Words(s.value, n, end))
}

// Mask returns the result of Mask on the wrapped string.
func (s Stringable) Mask(char rune, index, length int) Stringable {
	return Of(Mask(s.value, char, index, length))
}

// PadLeft returns the result of PadLeft on the wrapped string.
func (s Stringable) PadLeft(length int, pad string) Stringable {
	return Of(PadLeft(s.value, length, pad))
}

// PadRight returns the result of PadRight on the wrapped string.
func (s Stringable) PadRight(length int, pad string) Stringable {
	return Of(PadRight(s.value, length, pad))
}

// PadBoth returns the result of PadBoth on the wrapped string.
func (s Stringable) PadBoth(length int, pad string) Stringable {
	return Of(PadBoth(s.value, length, pad))
}

// Start returns the result of Start on the wrapped string.
func (s Stringable) Start(prefix string) Stringable {
	return Of(Start(s.value, prefix))
}

// Finish returns the result of Finish on the wrapped string.
func (s Stringable) Finish(suffix string) Stringable {
	return Of(Finish(s.value, suffix))
}

// Replace replaces every occurrence of search with replace, like
// strings.ReplaceAll, which the package does not duplicate as a function.
// Unlike strings.ReplaceAll, an empty search returns the value unchanged
// instead of inserting replace around every rune.
func (s Stringable) Replace(search, replace string) Stringable {
	if search == "" {
		return s
	}
	return Of(strings.ReplaceAll(s.value, search, replace))
}

// ReplaceFirst returns the result of ReplaceFirst on the wrapped string.
func (s Stringable) ReplaceFirst(search, replace string) Stringable {
	return Of(ReplaceFirst(s.value, search, replace))
}

// ReplaceLast returns the result of ReplaceLast on the wrapped string.
func (s Stringable) ReplaceLast(search, replace string) Stringable {
	return Of(ReplaceLast(s.value, search, replace))
}

// ReplaceStart returns the result of ReplaceStart on the wrapped string.
func (s Stringable) ReplaceStart(search, replace string) Stringable {
	return Of(ReplaceStart(s.value, search, replace))
}

// ReplaceEnd returns the result of ReplaceEnd on the wrapped string.
func (s Stringable) ReplaceEnd(search, replace string) Stringable {
	return Of(ReplaceEnd(s.value, search, replace))
}

// ReplaceArray returns the result of ReplaceArray on the wrapped string.
func (s Stringable) ReplaceArray(search string, replaces []string) Stringable {
	return Of(ReplaceArray(s.value, search, replaces))
}

// Swap returns the result of Swap on the wrapped string.
func (s Stringable) Swap(pairs map[string]string) Stringable {
	return Of(Swap(s.value, pairs))
}

// Remove returns the result of Remove on the wrapped string.
func (s Stringable) Remove(search ...string) Stringable {
	return Of(Remove(s.value, search...))
}

// Squish returns the result of Squish on the wrapped string.
func (s Stringable) Squish() Stringable {
	return Of(Squish(s.value))
}

// Reverse returns the result of Reverse on the wrapped string.
func (s Stringable) Reverse() Stringable {
	return Of(Reverse(s.value))
}

// Substr returns the result of Substr on the wrapped string.
func (s Stringable) Substr(start, length int) Stringable {
	return Of(Substr(s.value, start, length))
}

// Slug returns the result of Slug on the wrapped string.
func (s Stringable) Slug(separator string) Stringable {
	return Of(Slug(s.value, separator))
}

// Wrap returns the result of Wrap on the wrapped string.
func (s Stringable) Wrap(before, after string) Stringable {
	return Of(Wrap(s.value, before, after))
}

// Unwrap returns the result of Unwrap on the wrapped string.
func (s Stringable) Unwrap(before, after string) Stringable {
	return Of(Unwrap(s.value, before, after))
}

// Append returns the value with parts appended in order.
func (s Stringable) Append(parts ...string) Stringable {
	return Of(s.value + strings.Join(parts, ""))
}

// Prepend returns the value with parts prepended in order.
func (s Stringable) Prepend(parts ...string) Stringable {
	return Of(strings.Join(parts, "") + s.value)
}

// Lower returns the value converted to lower case.
func (s Stringable) Lower() Stringable {
	return Of(strings.ToLower(s.value))
}

// Upper returns the value converted to upper case.
func (s Stringable) Upper() Stringable {
	return Of(strings.ToUpper(s.value))
}

// Trim removes leading and trailing whitespace, or, when cutset is given,
// every leading and trailing rune contained in any of the cutset strings. The
// cutset strings are joined into one set of runes, as in strings.Trim, so
// Trim("ab") removes any run of 'a' and 'b' rather than the substring "ab".
func (s Stringable) Trim(cutset ...string) Stringable {
	if len(cutset) == 0 {
		return Of(strings.TrimSpace(s.value))
	}
	return Of(strings.Trim(s.value, strings.Join(cutset, "")))
}

// LTrim removes leading whitespace, or, when cutset is given, every leading
// rune contained in any of the cutset strings, treated as one set of runes as
// in Trim.
func (s Stringable) LTrim(cutset ...string) Stringable {
	if len(cutset) == 0 {
		return Of(strings.TrimLeftFunc(s.value, unicode.IsSpace))
	}
	return Of(strings.TrimLeft(s.value, strings.Join(cutset, "")))
}

// RTrim removes trailing whitespace, or, when cutset is given, every trailing
// rune contained in any of the cutset strings, treated as one set of runes as
// in Trim.
func (s Stringable) RTrim(cutset ...string) Stringable {
	if len(cutset) == 0 {
		return Of(strings.TrimRightFunc(s.value, unicode.IsSpace))
	}
	return Of(strings.TrimRight(s.value, strings.Join(cutset, "")))
}

// Repeat returns the value repeated n times, or an empty value when n is zero
// or negative.
func (s Stringable) Repeat(n int) Stringable {
	if n <= 0 {
		return Of("")
	}
	return Of(strings.Repeat(s.value, n))
}

// IsEmpty reports whether the value is the empty string.
func (s Stringable) IsEmpty() bool {
	return s.value == ""
}

// IsNotEmpty reports whether the value is not the empty string.
func (s Stringable) IsNotEmpty() bool {
	return s.value != ""
}

// Contains reports whether the value contains any of the non-empty needles.
func (s Stringable) Contains(needles ...string) bool {
	return Contains(s.value, needles...)
}

// ContainsAll reports whether the value contains every non-empty needle.
func (s Stringable) ContainsAll(needles ...string) bool {
	return ContainsAll(s.value, needles...)
}

// StartsWith reports whether the value starts with any of the non-empty
// prefixes.
func (s Stringable) StartsWith(prefixes ...string) bool {
	for _, p := range prefixes {
		if p != "" && strings.HasPrefix(s.value, p) {
			return true
		}
	}
	return false
}

// EndsWith reports whether the value ends with any of the non-empty suffixes.
func (s Stringable) EndsWith(suffixes ...string) bool {
	for _, p := range suffixes {
		if p != "" && strings.HasSuffix(s.value, p) {
			return true
		}
	}
	return false
}

// Is reports whether the value matches pattern as described for Is.
func (s Stringable) Is(pattern string) bool {
	return Is(pattern, s.value)
}

// Exactly reports whether the value equals other.
func (s Stringable) Exactly(other string) bool {
	return s.value == other
}

// Length returns the number of runes in the value.
func (s Stringable) Length() int {
	return Length(s.value)
}

// Split splits the value around each occurrence of sep as strings.Split does.
func (s Stringable) Split(sep string) []string {
	return strings.Split(s.value, sep)
}

// WordCount returns the result of WordCount on the wrapped string.
func (s Stringable) WordCount() int {
	return WordCount(s.value)
}

// When returns fn applied to the value when cond is true, and the value
// unchanged otherwise.
func (s Stringable) When(cond bool, fn func(Stringable) Stringable) Stringable {
	if cond {
		return fn(s)
	}
	return s
}

// Unless returns fn applied to the value when cond is false, and the value
// unchanged otherwise.
func (s Stringable) Unless(cond bool, fn func(Stringable) Stringable) Stringable {
	return s.When(!cond, fn)
}

// WhenEmpty returns fn applied to the value when it is empty.
func (s Stringable) WhenEmpty(fn func(Stringable) Stringable) Stringable {
	return s.When(s.IsEmpty(), fn)
}

// WhenNotEmpty returns fn applied to the value when it is not empty.
func (s Stringable) WhenNotEmpty(fn func(Stringable) Stringable) Stringable {
	return s.When(s.IsNotEmpty(), fn)
}

// WhenContains returns fn applied to the value when it contains needle.
func (s Stringable) WhenContains(needle string, fn func(Stringable) Stringable) Stringable {
	return s.When(s.Contains(needle), fn)
}

// WhenStartsWith returns fn applied to the value when it starts with prefix.
func (s Stringable) WhenStartsWith(prefix string, fn func(Stringable) Stringable) Stringable {
	return s.When(s.StartsWith(prefix), fn)
}

// WhenEndsWith returns fn applied to the value when it ends with suffix.
func (s Stringable) WhenEndsWith(suffix string, fn func(Stringable) Stringable) Stringable {
	return s.When(s.EndsWith(suffix), fn)
}

// WhenExactly returns fn applied to the value when it equals other.
func (s Stringable) WhenExactly(other string, fn func(Stringable) Stringable) Stringable {
	return s.When(s.Exactly(other), fn)
}

// WhenIs returns fn applied to the value when it matches pattern as described
// for Is.
func (s Stringable) WhenIs(pattern string, fn func(Stringable) Stringable) Stringable {
	return s.When(s.Is(pattern), fn)
}

// Pipe returns the result of passing the wrapped string through fn.
func (s Stringable) Pipe(fn func(string) string) Stringable {
	return Of(fn(s.value))
}

// Tap calls fn with the value and returns the value unchanged.
func (s Stringable) Tap(fn func(Stringable)) Stringable {
	fn(s)
	return s
}
