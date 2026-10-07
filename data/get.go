package data

import "slices"

// Get returns the value at path in target and reports whether it exists.
// An empty path returns target itself. A present key holding nil counts as
// found.
//
// A "*" segment matches every element of a slice or every value of a map
// (in ascending key order) and Get returns the matches as a []any; each
// further "*" in the path flattens the result by one level. Elements that
// do not contain the rest of the path are skipped, and ok is false when
// nothing matches.
func Get(target any, path string) (any, bool) {
	if path == "" {
		return target, true
	}
	return get(target, splitPath(path))
}

// get resolves segs against node.
func get(node any, segs []string) (any, bool) {
	for i, seg := range segs {
		if seg != wildcard {
			var ok bool
			if node, ok = lookup(node, seg); !ok {
				return nil, false
			}
			continue
		}
		_, vals, ok := entries(node)
		if !ok {
			return nil, false
		}
		rest := segs[i+1:]
		flatten := slices.Contains(rest, wildcard)
		out := []any{}
		for _, v := range vals {
			r, ok := get(v, rest)
			switch {
			case !ok:
			case flatten:
				out = append(out, r.([]any)...)
			default:
				out = append(out, r)
			}
		}
		if len(out) == 0 {
			return nil, false
		}
		return out, true
	}
	return node, true
}

// GetOr returns the value at path in target, or fallback when it does not
// exist.
func GetOr(target any, path string, fallback any) any {
	if v, ok := Get(target, path); ok {
		return v
	}
	return fallback
}

// GetAs returns the value at path in target asserted to T. It reports false
// when the value does not exist or is not a T; no conversion is attempted
// and a nil value never satisfies T.
func GetAs[T any](target any, path string) (T, bool) {
	v, _ := Get(target, path)
	t, ok := v.(T)
	return t, ok
}

// Has reports whether a value exists at path in target, following the same
// rules as Get.
func Has(target any, path string) bool {
	_, ok := Get(target, path)
	return ok
}
