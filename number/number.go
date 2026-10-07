package number

import (
	"cmp"
	"strconv"
)

// Ordinal returns n followed by its English ordinal suffix: "1st", "2nd",
// "3rd", "4th", "11th", "21st". Negative numbers keep their sign, so -2 is
// "-2nd".
func Ordinal(n int) string {
	m := n % 100
	if m < 0 {
		m = -m
	}
	suffix := "th"
	if m < 11 || m > 13 {
		switch m % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return strconv.Itoa(n) + suffix
}

// Clamp returns v limited to the range [lo, hi]. When lo is greater than hi
// the bounds are swapped.
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if lo > hi {
		lo, hi = hi, lo
	}
	return min(max(v, lo), hi)
}

// Pairs splits the range from start up to to into consecutive inclusive
// [lower, upper] pairs of width by, following Laravel's Number::pairs with
// its offset fixed at 1: each lower bound starts at start and advances by
// by, each upper bound is lower+by-1 capped at to, and pairs are produced
// while lower is less than to. Pairs(25, 10, 0) is [[0 9] [10 19] [20 25]]
// and Pairs(25, 10, 1) is [[1 10] [11 20] [21 25]]. It returns nil when by
// is not positive or start is not less than to. As in Laravel, a lower
// bound equal to to starts no pair, so Pairs(20, 10, 0) is [[0 9] [10 19]].
func Pairs(to, by, start int) [][2]int {
	if by <= 0 || start >= to {
		return nil
	}
	var pairs [][2]int
	for lower := start; ; lower += by {
		// The unsigned difference is exact because lower < to.
		rem := uint64(to) - uint64(lower)
		if uint64(by-1) >= rem {
			return append(pairs, [2]int{lower, to})
		}
		pairs = append(pairs, [2]int{lower, lower + by - 1})
		if rem == uint64(by) {
			return pairs
		}
	}
}
