package data

import (
	"maps"
	"reflect"
	"strconv"
)

// Dot flattens the nested map[string]any and []any values of m into a
// single-level map whose keys are dot-notation paths, such as "user.name"
// and "tags.0". Empty nested maps and slices are kept as values under their
// own path. Other values, including other map and slice types, are stored
// as they are.
func Dot(m map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		flatten(out, k, v)
	}
	return out
}

// flatten stores v in out under prefix, expanding nested containers.
func flatten(out map[string]any, prefix string, v any) {
	switch n := v.(type) {
	case map[string]any:
		if len(n) == 0 {
			out[prefix] = map[string]any{}
		}
		for k, c := range n {
			flatten(out, prefix+"."+k, c)
		}
	case []any:
		if len(n) == 0 {
			out[prefix] = []any{}
		}
		for i, c := range n {
			flatten(out, prefix+"."+strconv.Itoa(i), c)
		}
	default:
		out[prefix] = v
	}
}

// Undot expands the dot-notation keys of m into nested map[string]any
// values. It is the inverse of Dot except that numeric segments become map
// keys such as "0"; slices are not rebuilt. Keys are applied in ascending
// order, so when a key is also the prefix of another, such as "a" and
// "a.b", the nested value replaces a non-map value. Nested maps taken from m
// are copied before they are extended.
func Undot(m map[string]any) map[string]any {
	out := make(map[string]any)
	owned := map[uintptr]bool{}
	for _, key := range sortedKeys(m) {
		segs := splitPath(key)
		node := out
		for _, seg := range segs[:len(segs)-1] {
			child, ok := node[seg].(map[string]any)
			if !ok {
				child = map[string]any{}
			} else if !owned[reflect.ValueOf(child).Pointer()] {
				child = maps.Clone(child)
			}
			owned[reflect.ValueOf(child).Pointer()] = true
			node[seg] = child
			node = child
		}
		node[segs[len(segs)-1]] = m[key]
	}
	return out
}
