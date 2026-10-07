// Package set provides Set, a generic unordered collection of distinct
// comparable values backed by a map.
//
// A Set is a map type, so len, range, clear and the maps package work on it
// directly. The zero value is an empty set ready to use: reading methods
// accept a nil Set, and Add allocates the map on first use through its pointer
// receiver. Union, Intersect, Diff and Clone return new sets and never modify
// their receiver or argument; Add and Delete modify the receiver.
//
// Example:
//
//	var seen set.Set[string]
//	for _, name := range names {
//		if !seen.Add(name) {
//			fmt.Println("duplicate:", name)
//		}
//	}
//	admins := set.Of("ann", "bob")
//	fmt.Println(set.Sorted(seen.Intersect(admins)))
package set
