// Package mapx provides generic map helpers that complement the standard maps
// package: key selection, predicate filtering, lookups with fallbacks, key and
// value transformation, merging and ordered iteration.
//
// Functions never modify their input and always return a new, non-nil map.
// Nil maps are valid input everywhere.
//
// Example:
//
//	public := mapx.Except(settings, "password", "token")
//	timeout := mapx.GetOr(options, "timeout", 30)
//	merged := mapx.Merge(defaults, overrides)
//	for _, e := range mapx.SortedEntries(merged) {
//		fmt.Println(e.Key, e.Value)
//	}
package mapx
