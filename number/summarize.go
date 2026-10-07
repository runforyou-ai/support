package number

import (
	"math"
	"strconv"
	"strings"
)

// abbreviateUnits are the suffixes used by Abbreviate, indexed by power of
// one thousand.
var abbreviateUnits = []string{"", "K", "M", "B", "T", "Q"}

// humanUnits are the suffixes used by ForHumans, indexed by power of one
// thousand.
var humanUnits = []string{"", " thousand", " million", " billion", " trillion", " quadrillion"}

// Abbreviate formats v with a short scale suffix (K, M, B, T, Q), keeping
// at most precision digits after the decimal point and trimming trailing
// zeros: Abbreviate(1200000, 2) is "1.2M" and Abbreviate(1000, 2) is "1K".
// A value that rounds up to the next unit uses it, so 999999 is "1M" at
// precision 0. Values of a thousand quadrillion and above stack suffixes as
// Laravel does, for example "1KQ". A negative precision uses the shortest
// representation.
func Abbreviate(v float64, precision int) string {
	return summarize(v, precision, abbreviateUnits)
}

// ForHumans is like Abbreviate with spelled-out scale words: ForHumans(1500000,
// 1) is "1.5 million".
func ForHumans(v float64, precision int) string {
	return summarize(v, precision, humanUnits)
}

// summarize scales v by powers of one thousand and appends the matching unit.
func summarize(v float64, precision int, units []string) string {
	switch {
	case math.IsNaN(v) || math.IsInf(v, 0):
		return strconv.FormatFloat(v, 'f', -1, 64)
	case v < 0:
		s := summarize(-v, precision, units)
		if s == "0" {
			return s
		}
		return "-" + s
	case v >= 1e15:
		return summarize(v/1e15, precision, units) + units[len(units)-1]
	}
	unit := 0
	for unit < len(units)-1 && v >= math.Pow10(3*(unit+1)) {
		unit++
	}
	intPart, frac := roundDecimal(v/math.Pow10(3*unit), precision)
	// Rounding up to one thousand moves to the next unit.
	if len(intPart) > 3 && unit < len(units)-1 {
		unit++
		intPart, frac = roundDecimal(v/math.Pow10(3*unit), precision)
	}
	frac = strings.TrimRight(frac, "0")
	if frac != "" {
		intPart += "." + frac
	}
	return intPart + units[unit]
}
