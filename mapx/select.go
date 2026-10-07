package mapx

import "maps"

// Only returns a new map holding the entries of m whose key is one of keys.
// Keys missing from m are skipped.
func Only[M ~map[K]V, K comparable, V any](m M, keys ...K) M {
	out := make(M, min(len(keys), len(m)))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			out[k] = v
		}
	}
	return out
}

// Except returns a new map holding the entries of m whose key is not one of
// keys.
func Except[M ~map[K]V, K comparable, V any](m M, keys ...K) M {
	out := make(M, len(m))
	maps.Copy(out, m)
	for _, k := range keys {
		delete(out, k)
	}
	return out
}

// Filter returns a new map holding the entries of m for which pred returns
// true.
func Filter[M ~map[K]V, K comparable, V any](m M, pred func(K, V) bool) M {
	out := make(M)
	for k, v := range m {
		if pred(k, v) {
			out[k] = v
		}
	}
	return out
}

// Reject returns a new map holding the entries of m for which pred returns
// false. It is the inverse of Filter.
func Reject[M ~map[K]V, K comparable, V any](m M, pred func(K, V) bool) M {
	return Filter(m, func(k K, v V) bool { return !pred(k, v) })
}
