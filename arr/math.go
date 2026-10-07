package arr

import "github.com/runforyou-ai/support"

// Sum returns the sum of the elements of s, or 0 when s is empty.
func Sum[E support.Number](s []E) E {
	var total E
	for _, v := range s {
		total += v
	}
	return total
}

// SumBy returns the sum of fn(element) over s, or 0 when s is empty.
func SumBy[E any, N support.Number](s []E, fn func(E) N) N {
	var total N
	for _, v := range s {
		total += fn(v)
	}
	return total
}

// Avg returns the arithmetic mean of the elements of s, or 0 when s is empty.
// The sum is accumulated in float64.
func Avg[E support.Number](s []E) float64 {
	return AvgBy(s, func(v E) E { return v })
}

// AvgBy returns the arithmetic mean of fn(element) over s, or 0 when s is
// empty. The sum is accumulated in float64.
func AvgBy[E any, N support.Number](s []E, fn func(E) N) float64 {
	if len(s) == 0 {
		return 0
	}
	var total float64
	for _, v := range s {
		total += float64(fn(v))
	}
	return total / float64(len(s))
}
