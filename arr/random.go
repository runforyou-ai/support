package arr

import "math/rand/v2"

// Random returns a randomly chosen element of s. The boolean is false when s
// is empty.
func Random[E any](s []E) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[rand.IntN(len(s))], true
}

// Sample returns n elements of s taken from distinct positions, in random
// order. n is clamped to len(s); it returns nil when n is not positive.
func Sample[S ~[]E, E any](s S, n int) S {
	n = min(n, len(s))
	if n <= 0 {
		return nil
	}
	out := make(S, n)
	for i, p := range rand.Perm(len(s))[:n] {
		out[i] = s[p]
	}
	return out
}

// Shuffle returns a copy of s with its elements in random order.
func Shuffle[S ~[]E, E any](s S) S {
	return Sample(s, len(s))
}
