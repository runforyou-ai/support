package arr

// First returns the first element of s for which pred returns true. When pred
// is nil it returns the first element. The boolean is false when no element
// matches.
func First[E any](s []E, pred func(E) bool) (E, bool) {
	for _, v := range s {
		if pred == nil || pred(v) {
			return v, true
		}
	}
	var zero E
	return zero, false
}

// Last returns the last element of s for which pred returns true. When pred
// is nil it returns the last element. The boolean is false when no element
// matches.
func Last[E any](s []E, pred func(E) bool) (E, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if pred == nil || pred(s[i]) {
			return s[i], true
		}
	}
	var zero E
	return zero, false
}

// Every reports whether pred returns true for every element of s. It returns
// true for an empty slice.
func Every[E any](s []E, pred func(E) bool) bool {
	for _, v := range s {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Count returns the number of elements of s for which pred returns true.
func Count[E any](s []E, pred func(E) bool) int {
	n := 0
	for _, v := range s {
		if pred(v) {
			n++
		}
	}
	return n
}
