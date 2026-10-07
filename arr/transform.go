package arr

// Filter returns the elements of s for which pred returns true, in order.
func Filter[S ~[]E, E any](s S, pred func(E) bool) S {
	var out S
	for _, v := range s {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reject returns the elements of s for which pred returns false, in order. It
// is the inverse of Filter.
func Reject[S ~[]E, E any](s S, pred func(E) bool) S {
	return Filter(s, func(v E) bool { return !pred(v) })
}

// Map returns a slice holding fn applied to each element of s, in order.
func Map[E, R any](s []E, fn func(E) R) []R {
	if len(s) == 0 {
		return nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}

// MapErr returns a slice holding fn applied to each element of s, in order.
// It stops at the first error and returns nil and that error. It returns nil
// and a nil error for an empty s.
func MapErr[E, R any](s []E, fn func(E) (R, error)) ([]R, error) {
	if len(s) == 0 {
		return nil, nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		r, err := fn(v)
		if err != nil {
			return nil, err
		}
		out[i] = r
	}
	return out, nil
}

// OrEmpty returns s itself when it is not nil, and a non-nil empty slice
// otherwise. It is useful where nil and empty encode differently, such as
// JSON null and [].
func OrEmpty[S ~[]E, E any](s S) S {
	if s == nil {
		return S{}
	}
	return s
}

// FilterMap applies fn to each element of s and returns the results for which
// fn reports true, in order.
func FilterMap[E, R any](s []E, fn func(E) (R, bool)) []R {
	var out []R
	for _, v := range s {
		if r, ok := fn(v); ok {
			out = append(out, r)
		}
	}
	return out
}

// Reduce folds s from left to right, starting with init and replacing the
// accumulator with fn(accumulator, element) for each element.
func Reduce[E, A any](s []E, init A, fn func(A, E) A) A {
	acc := init
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// Flatten concatenates the inner slices of s into a single slice. It flattens
// one level only.
func Flatten[E any](s [][]E) []E {
	n := 0
	for _, inner := range s {
		n += len(inner)
	}
	if n == 0 {
		return nil
	}
	out := make([]E, 0, n)
	for _, inner := range s {
		out = append(out, inner...)
	}
	return out
}

// CrossJoin returns the cartesian product of lists: every combination that
// takes one element from each list, ordered with the last list varying
// fastest. It returns nil when no lists are given or any list is empty. The
// result holds the product of the list lengths, so it is meant for small
// inputs.
func CrossJoin[E any](lists ...[]E) [][]E {
	if len(lists) == 0 {
		return nil
	}
	total := 1
	for _, l := range lists {
		if len(l) == 0 {
			return nil
		}
		total *= len(l)
	}
	out := make([][]E, total)
	for i := range out {
		row := make([]E, len(lists))
		rem := i
		for j := len(lists) - 1; j >= 0; j-- {
			row[j] = lists[j][rem%len(lists[j])]
			rem /= len(lists[j])
		}
		out[i] = row
	}
	return out
}

// Pad returns a copy of s extended to |size| elements with value. A positive
// size pads on the right and a negative size pads on the left. When |size| is
// not greater than len(s), it returns a copy of s.
func Pad[S ~[]E, E any](s S, size int, value E) S {
	left := size < 0
	if left {
		size = -size
	}
	if size <= len(s) {
		return append(S(nil), s...)
	}
	out := make(S, 0, size)
	fill := size - len(s)
	if !left {
		out = append(out, s...)
	}
	for range fill {
		out = append(out, value)
	}
	if left {
		out = append(out, s...)
	}
	return out
}

// Join concatenates items with glue, using finalGlue between the last two
// items, for example "a, b and c". An empty finalGlue means glue.
func Join(items []string, glue, finalGlue string) string {
	if finalGlue == "" {
		finalGlue = glue
	}
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	}
	n := len(items) - 1
	size := len(finalGlue) + len(glue)*(n-1)
	for _, item := range items {
		size += len(item)
	}
	buf := make([]byte, 0, size)
	for i, item := range items {
		switch {
		case i == n:
			buf = append(buf, finalGlue...)
		case i > 0:
			buf = append(buf, glue...)
		}
		buf = append(buf, item...)
	}
	return string(buf)
}
