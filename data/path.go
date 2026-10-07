package data

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

const wildcard = "*"

// splitPath splits a dot-notation path into its segments.
func splitPath(path string) []string {
	return strings.Split(path, ".")
}

// index parses seg as a slice index below n.
func index(seg string, n int) (int, bool) {
	if seg == "" {
		return 0, false
	}
	for i := 0; i < len(seg); i++ {
		if seg[i] < '0' || seg[i] > '9' {
			return 0, false
		}
	}
	i, err := strconv.Atoi(seg)
	if err != nil || i >= n {
		return 0, false
	}
	return i, true
}

// lookup returns the child of node named by seg.
func lookup(node any, seg string) (any, bool) {
	switch n := node.(type) {
	case map[string]any:
		v, ok := n[seg]
		return v, ok
	case []any:
		i, ok := index(seg, len(n))
		if !ok {
			return nil, false
		}
		return n[i], true
	case nil:
		return nil, false
	}
	rv := reflect.ValueOf(node)
	switch rv.Kind() {
	case reflect.Map:
		kt := rv.Type().Key()
		if kt.Kind() != reflect.String {
			return nil, false
		}
		e := rv.MapIndex(reflect.ValueOf(seg).Convert(kt))
		if !e.IsValid() {
			return nil, false
		}
		return e.Interface(), true
	case reflect.Slice, reflect.Array:
		i, ok := index(seg, rv.Len())
		if !ok {
			return nil, false
		}
		return rv.Index(i).Interface(), true
	}
	return nil, false
}

// entries returns the keys and values of a map or slice node, with map keys
// sorted and slice keys as decimal indexes.
func entries(node any) ([]string, []any, bool) {
	switch n := node.(type) {
	case map[string]any:
		keys := sortedKeys(n)
		vals := make([]any, len(keys))
		for i, k := range keys {
			vals[i] = n[k]
		}
		return keys, vals, true
	case []any:
		keys := make([]string, len(n))
		for i := range n {
			keys[i] = strconv.Itoa(i)
		}
		return keys, slices.Clone(n), true
	case nil:
		return nil, nil, false
	}
	rv := reflect.ValueOf(node)
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, nil, false
		}
		mk := rv.MapKeys()
		slices.SortFunc(mk, func(a, b reflect.Value) int { return strings.Compare(a.String(), b.String()) })
		keys := make([]string, len(mk))
		vals := make([]any, len(mk))
		for i, k := range mk {
			keys[i] = k.String()
			vals[i] = rv.MapIndex(k).Interface()
		}
		return keys, vals, true
	case reflect.Slice, reflect.Array:
		keys := make([]string, rv.Len())
		vals := make([]any, rv.Len())
		for i := range keys {
			keys[i] = strconv.Itoa(i)
			vals[i] = rv.Index(i).Interface()
		}
		return keys, vals, true
	}
	return nil, nil, false
}

// sortedKeys returns the keys of m in ascending order.
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
