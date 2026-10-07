package data

import (
	"maps"
	"reflect"
	"slices"
)

// Set returns a copy of target with value stored at path, overwriting any
// existing value, like Laravel's data_set. Maps and slices along the path
// are copied; other nested values are shared with target, which is never
// modified.
//
// Missing or nil intermediate segments are created as map[string]any, and
// existing intermediate values that are not maps or slices are replaced by
// one. Numeric segments index existing []any elements; an index that is out
// of range, or a non-numeric segment applied to a slice, changes nothing. A
// "*" segment applies the rest of the path to every existing value of a map
// or element of a slice. Only map[string]any and []any are traversed; other
// map and slice types are left as they are. A nil target is treated as an
// empty map, and an empty path returns a shallow copy of target.
func Set(target map[string]any, path string, value any) map[string]any {
	return write(target, path, value, true)
}

// Fill returns a copy of target with value stored at path only when nothing
// exists there yet. It creates missing intermediate maps like Set but never
// replaces existing values, including nil values and intermediate values
// that are not containers. target is never modified.
func Fill(target map[string]any, path string, value any) map[string]any {
	return write(target, path, value, false)
}

// write implements Set and Fill.
func write(target map[string]any, path string, value any, overwrite bool) map[string]any {
	if path == "" {
		return maps.Clone(target)
	}
	if target == nil {
		target = map[string]any{}
	}
	return assign(target, splitPath(path), value, overwrite).(map[string]any)
}

// assign returns a copy of node, a map[string]any or []any, with value stored at segs.
func assign(node any, segs []string, value any, overwrite bool) any {
	seg, rest := segs[0], segs[1:]
	if n, ok := node.(map[string]any); ok {
		out := make(map[string]any, len(n)+1)
		maps.Copy(out, n)
		if seg == wildcard {
			for k, v := range n {
				if c, store := place(v, true, rest, value, overwrite); store {
					out[k] = c
				}
			}
			return out
		}
		v, exists := n[seg]
		if c, store := place(v, exists, rest, value, overwrite); store {
			out[seg] = c
		}
		return out
	}
	n := node.([]any)
	out := slices.Clone(n)
	for i, v := range n {
		if seg != wildcard {
			if j, ok := index(seg, len(n)); !ok || j != i {
				continue
			}
		}
		if c, store := place(v, true, rest, value, overwrite); store {
			out[i] = c
		}
	}
	return out
}

// place returns the new value for a slot holding child when rest is applied
// below it, and reports whether the slot must be updated.
func place(child any, exists bool, rest []string, value any, overwrite bool) (any, bool) {
	if len(rest) == 0 {
		return value, overwrite || !exists
	}
	switch child.(type) {
	case map[string]any, []any:
		return assign(child, rest, value, overwrite), true
	}
	if exists && (!overwrite || isContainer(child)) {
		return nil, false
	}
	// A wildcard has nothing to match inside a newly created map.
	if slices.Contains(rest, wildcard) {
		return nil, false
	}
	return assign(map[string]any{}, rest, value, overwrite), true
}

// isContainer reports whether v is a map or slice of a type writers do not
// traverse.
func isContainer(v any) bool {
	if v == nil {
		return false
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// Forget returns a copy of target with the map key at path deleted, like
// Laravel's data_forget. Maps and slices along the path are copied; target
// is never modified. A "*" segment applies the rest of the path to every
// value of a map or element of a slice; a trailing "*" deletes nothing.
// Slice elements are not removed, so a path ending at a slice index changes
// nothing. An empty path returns a shallow copy of target, and a nil target
// returns nil.
func Forget(target map[string]any, path string) map[string]any {
	if target == nil || path == "" {
		return maps.Clone(target)
	}
	return forget(target, splitPath(path)).(map[string]any)
}

// forget returns node with segs deleted, copying the containers it changes.
func forget(node any, segs []string) any {
	seg, rest := segs[0], segs[1:]
	switch n := node.(type) {
	case map[string]any:
		out := maps.Clone(n)
		switch {
		case seg == wildcard && len(rest) > 0:
			for k, v := range n {
				out[k] = forget(v, rest)
			}
		case seg == wildcard:
		case len(rest) == 0:
			delete(out, seg)
		default:
			if v, ok := n[seg]; ok {
				out[seg] = forget(v, rest)
			}
		}
		return out
	case []any:
		if len(rest) == 0 {
			return node
		}
		out := slices.Clone(n)
		for i, v := range n {
			if seg != wildcard {
				if j, ok := index(seg, len(n)); !ok || j != i {
					continue
				}
			}
			out[i] = forget(v, rest)
		}
		return out
	}
	return node
}
