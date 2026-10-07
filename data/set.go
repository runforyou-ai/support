package data

import (
	"reflect"
	"slices"
)

// Set modifies target in place, storing value at path and overwriting any
// existing value, like Laravel's data_set.
//
// Missing intermediate segments are created as map[string]any, and existing
// intermediate values that are not maps or slices are replaced by one.
// Numeric segments index existing []any elements; an index that is out of
// range, or a non-numeric segment applied to a slice, leaves target
// unchanged because slices cannot grow in place. A "*" segment applies the
// rest of the path to every existing value of a map or element of a slice.
// Only map[string]any and []any are traversed; other map and slice types are
// left untouched. A nil target or an empty path does nothing.
func Set(target map[string]any, path string, value any) {
	if target == nil || path == "" {
		return
	}
	assign(target, splitPath(path), value, true)
}

// Fill modifies target in place, storing value at path only when nothing
// exists there yet. It creates missing intermediate maps like Set but never
// replaces existing values, including nil values and intermediate values
// that are not containers.
func Fill(target map[string]any, path string, value any) {
	if target == nil || path == "" {
		return
	}
	assign(target, splitPath(path), value, false)
}

// assign stores value at segs below the container node.
func assign(node any, segs []string, value any, overwrite bool) {
	seg, rest := segs[0], segs[1:]
	switch n := node.(type) {
	case map[string]any:
		if seg == wildcard {
			for k, v := range n {
				if len(rest) == 0 {
					if overwrite {
						n[k] = value
					}
				} else if c, store := descend(v, true, rest, value, overwrite); store {
					n[k] = c
				}
			}
			return
		}
		v, exists := n[seg]
		if len(rest) == 0 {
			if overwrite || !exists {
				n[seg] = value
			}
		} else if c, store := descend(v, exists, rest, value, overwrite); store {
			n[seg] = c
		}
	case []any:
		for i, v := range n {
			if seg != wildcard {
				if j, ok := index(seg, len(n)); !ok || j != i {
					continue
				}
			}
			if len(rest) == 0 {
				if overwrite {
					n[i] = value
				}
			} else if c, store := descend(v, true, rest, value, overwrite); store {
				n[i] = c
			}
		}
	}
}

// descend applies rest below child and returns a replacement for child when
// the slot must be updated.
func descend(child any, exists bool, rest []string, value any, overwrite bool) (any, bool) {
	switch child.(type) {
	case map[string]any, []any:
		assign(child, rest, value, overwrite)
		return nil, false
	}
	if exists && (!overwrite || isContainer(child)) {
		return nil, false
	}
	// A wildcard has nothing to match inside a newly created map.
	if slices.Contains(rest, wildcard) {
		return nil, false
	}
	m := map[string]any{}
	assign(m, rest, value, overwrite)
	return m, true
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

// Forget modifies target in place, deleting the map key at path, like
// Laravel's data_forget. A "*" segment applies the rest of the path to every
// value of a map or element of a slice; a trailing "*" deletes nothing.
// Slice elements cannot be removed in place, so a path ending at a slice
// index does nothing. A nil target or an empty path does nothing.
func Forget(target map[string]any, path string) {
	if target == nil || path == "" {
		return
	}
	forget(target, splitPath(path))
}

// forget deletes segs below node.
func forget(node any, segs []string) {
	seg, rest := segs[0], segs[1:]
	if seg == wildcard {
		if len(rest) == 0 {
			return
		}
		switch n := node.(type) {
		case map[string]any:
			for _, v := range n {
				forget(v, rest)
			}
		case []any:
			for _, v := range n {
				forget(v, rest)
			}
		}
		return
	}
	switch n := node.(type) {
	case map[string]any:
		if len(rest) == 0 {
			delete(n, seg)
		} else if v, ok := n[seg]; ok {
			forget(v, rest)
		}
	case []any:
		if i, ok := index(seg, len(n)); ok && len(rest) > 0 {
			forget(n[i], rest)
		}
	}
}
