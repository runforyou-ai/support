// Package arr provides slice helpers that complement the standard slices
// package: predicate queries, transformations, grouping, set operations,
// random selection and numeric aggregation.
//
// Functions never modify their input; they return new slices and maps. Nil
// and empty slices are valid input everywhere, and functions that build a
// slice return nil when the result is empty; OrEmpty turns a nil slice into an
// empty one where that difference matters.
//
// Example:
//
//	active := arr.Filter(users, func(u User) bool { return u.Active })
//	names := arr.Map(active, func(u User) string { return u.Name })
//	fmt.Println(arr.Join(names, ", ", " and "))
//	byTeam := arr.GroupBy(users, func(u User) string { return u.Team })
package arr
