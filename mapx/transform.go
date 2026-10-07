package mapx

import (
	"cmp"
	"maps"
	"slices"
)

// Entry is a key-value pair of a map.
type Entry[K, V any] struct {
	Key   K
	Value V
}

// Add returns a copy of m with key set to value when m has no such key, or a
// plain copy of m otherwise.
func Add[M ~map[K]V, K comparable, V any](m M, key K, value V) M {
	out := make(M, len(m)+1)
	maps.Copy(out, m)
	if _, ok := out[key]; !ok {
		out[key] = value
	}
	return out
}

// Merge returns a new map holding the entries of every map in ms. When a key
// appears in several maps, the value from the later map wins.
func Merge[M ~map[K]V, K comparable, V any](ms ...M) M {
	n := 0
	for _, m := range ms {
		n = max(n, len(m))
	}
	out := make(M, n)
	for _, m := range ms {
		maps.Copy(out, m)
	}
	return out
}

// MapValues returns a new map with the keys of m and fn(key, value) as
// values.
func MapValues[M ~map[K]V, K comparable, V, R any](m M, fn func(K, V) R) map[K]R {
	out := make(map[K]R, len(m))
	for k, v := range m {
		out[k] = fn(k, v)
	}
	return out
}

// MapKeys returns a new map with fn(key, value) as keys and the values of m.
// When fn maps several entries to the same key, which of their values is kept
// is unspecified, because map iteration order is random.
func MapKeys[M ~map[K]V, K comparable, V any, R comparable](m M, fn func(K, V) R) map[R]V {
	out := make(map[R]V, len(m))
	for k, v := range m {
		out[fn(k, v)] = v
	}
	return out
}

// Invert returns a new map with the keys and values of m swapped. When
// several keys share a value, which of those keys is kept is unspecified,
// because map iteration order is random.
func Invert[M ~map[K]V, K, V comparable](m M) map[V]K {
	out := make(map[V]K, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

// SortedKeys returns the keys of m in ascending order.
func SortedKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	return slices.Sorted(maps.Keys(m))
}

// SortedEntries returns the entries of m ordered by ascending key, with NaN
// keys first as ordered by cmp.Compare.
func SortedEntries[M ~map[K]V, K cmp.Ordered, V any](m M) []Entry[K, V] {
	if len(m) == 0 {
		return nil
	}
	out := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		out = append(out, Entry[K, V]{Key: k, Value: v})
	}
	slices.SortFunc(out, func(a, b Entry[K, V]) int { return cmp.Compare(a.Key, b.Key) })
	return out
}
