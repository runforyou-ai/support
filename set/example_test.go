package set_test

import (
	"fmt"

	"github.com/runforyou-ai/support/set"
)

func ExampleOf() {
	s := set.Of("b", "a", "b")
	fmt.Println(s.Len(), set.Sorted(s))
	// Output: 2 [a b]
}

func ExampleCollect() {
	s := set.Collect([]int{3, 1, 3, 2})
	fmt.Println(set.Sorted(s))
	// Output: [1 2 3]
}

func ExampleCollectBy() {
	type user struct {
		Name string
		Team string
	}
	users := []user{{"Ann", "core"}, {"Bob", "web"}, {"Cid", "core"}}
	teams := set.CollectBy(users, func(u user) string { return u.Team })
	fmt.Println(set.Sorted(teams))
	// Output: [core web]
}

func ExampleSorted() {
	fmt.Println(set.Sorted(set.Of(3, 1, 2)))
	// Output: [1 2 3]
}

func ExampleSet_Add() {
	var seen set.Set[string]
	for _, name := range []string{"ann", "bob", "ann"} {
		if !seen.Add(name) {
			fmt.Println("duplicate:", name)
		}
	}
	fmt.Println(seen.Len())
	// Output:
	// duplicate: ann
	// 2
}

func ExampleSet_Has() {
	s := set.Of("read", "write")
	fmt.Println(s.Has("read"), s.Has("admin"))
	// Output: true false
}

func ExampleSet_Delete() {
	s := set.Of(1, 2, 3)
	s.Delete(2)
	fmt.Println(set.Sorted(s))
	// Output: [1 3]
}

func ExampleSet_Union() {
	fmt.Println(set.Sorted(set.Of(1, 2).Union(set.Of(2, 3))))
	// Output: [1 2 3]
}

func ExampleSet_Intersect() {
	fmt.Println(set.Sorted(set.Of(1, 2, 3).Intersect(set.Of(2, 3, 4))))
	// Output: [2 3]
}

func ExampleSet_Diff() {
	fmt.Println(set.Sorted(set.Of(1, 2, 3).Diff(set.Of(2))))
	// Output: [1 3]
}
