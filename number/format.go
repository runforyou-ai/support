package number

import (
	"math"
	"strconv"
	"strings"
)

// Format formats v with decimals digits after the decimal point, using ","
// as the thousands separator and "." as the decimal separator. Values are
// rounded half away from zero. A negative decimals uses the shortest
// representation that round-trips v.
func Format(v float64, decimals int) string {
	return FormatWith(v, decimals, ".", ",")
}

// FormatWith is like Format with custom decimal and thousands separators.
func FormatWith(v float64, decimals int, decimalSep, thousandsSep string) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	intPart, frac := roundDecimal(math.Abs(v), decimals)
	var b strings.Builder
	// A value that rounds to zero carries no sign.
	if v < 0 && strings.Trim(intPart+frac, "0") != "" {
		b.WriteByte('-')
	}
	for i := range len(intPart) {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteString(thousandsSep)
		}
		b.WriteByte(intPart[i])
	}
	if frac != "" {
		b.WriteString(decimalSep)
		b.WriteString(frac)
	}
	return b.String()
}

// Percentage formats v, which is already a percentage value, with decimals
// digits and a trailing "%": Percentage(12.5, 2) is "12.50%".
func Percentage(v float64, decimals int) string {
	return Format(v, decimals) + "%"
}

// fileSizeUnits are the binary units used by FileSize.
var fileSizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// FileSize formats a byte count in base-1024 units (B, KB, MB, GB, TB, PB,
// EB) with precision digits after the decimal point. Like Laravel, it moves
// to the next unit once the value exceeds 0.9 of it, so 1000 bytes is
// "0.98 KB" at precision 2. Negative sizes keep their sign.
func FileSize(bytes int64, precision int) string {
	size := math.Abs(float64(bytes))
	unit := 0
	for size/1024 > 0.9 && unit < len(fileSizeUnits)-1 {
		size /= 1024
		unit++
	}
	if bytes < 0 {
		size = -size
	}
	return Format(size, precision) + " " + fileSizeUnits[unit]
}

// roundDecimal returns the integer and fractional digits of the
// non-negative value abs rounded half away from zero to decimals digits.
// A negative decimals keeps the shortest representation.
func roundDecimal(abs float64, decimals int) (intPart, frac string) {
	intPart, frac, _ = strings.Cut(strconv.FormatFloat(abs, 'f', -1, 64), ".")
	if decimals < 0 {
		return intPart, frac
	}
	if len(frac) <= decimals {
		return intPart, frac + strings.Repeat("0", decimals-len(frac))
	}
	digits := []byte(intPart + frac[:decimals])
	if frac[decimals] >= '5' {
		i := len(digits) - 1
		for ; i >= 0 && digits[i] == '9'; i-- {
			digits[i] = '0'
		}
		if i < 0 {
			digits = append([]byte{'1'}, digits...)
		} else {
			digits[i]++
		}
	}
	n := len(digits) - decimals
	return string(digits[:n]), string(digits[n:])
}
