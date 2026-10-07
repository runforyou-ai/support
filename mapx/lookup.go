package mapx

// GetOr returns the value stored under key, or fallback when m has no such
// key. A stored zero value is returned as is.
func GetOr[M ~map[K]V, K comparable, V any](m M, key K, fallback V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}

// Has reports whether m contains every one of keys. It returns false when no
// keys are given.
func Has[M ~map[K]V, K comparable, V any](m M, keys ...K) bool {
	if len(keys) == 0 {
		return false
	}
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

// HasAny reports whether m contains at least one of keys.
func HasAny[M ~map[K]V, K comparable, V any](m M, keys ...K) bool {
	for _, k := range keys {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}
