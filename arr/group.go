package arr

// KeyBy returns a map from fn(element) to element. When several elements
// share a key, the later element wins. The result is never nil.
func KeyBy[E any, K comparable](s []E, fn func(E) K) map[K]E {
	out := make(map[K]E, len(s))
	for _, v := range s {
		out[fn(v)] = v
	}
	return out
}

// GroupBy groups the elements of s by fn(element). Each group keeps the order
// of s. The result is never nil.
func GroupBy[E any, K comparable](s []E, fn func(E) K) map[K][]E {
	out := make(map[K][]E)
	for _, v := range s {
		k := fn(v)
		out[k] = append(out[k], v)
	}
	return out
}

// Partition splits s into the elements for which pred returns true and the
// rest, both in order.
func Partition[S ~[]E, E any](s S, pred func(E) bool) (matched, rest S) {
	for _, v := range s {
		if pred(v) {
			matched = append(matched, v)
		} else {
			rest = append(rest, v)
		}
	}
	return matched, rest
}
