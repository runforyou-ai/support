package mapx_test

import (
	"fmt"
	"strings"

	"github.com/runforyou-ai/support/mapx"
)

// printSorted prints m in ascending key order.
func printSorted[V any](m map[string]V) {
	var parts []string
	for _, e := range mapx.SortedEntries(m) {
		parts = append(parts, fmt.Sprintf("%s=%v", e.Key, e.Value))
	}
	fmt.Println(strings.Join(parts, " "))
}

func ExampleOnly() {
	user := map[string]string{"name": "Ann", "email": "ann@example.com", "password": "secret"}
	printSorted(mapx.Only(user, "name", "email"))
	// Output: email=ann@example.com name=Ann
}

func ExampleExcept() {
	user := map[string]string{"name": "Ann", "password": "secret"}
	printSorted(mapx.Except(user, "password"))
	// Output: name=Ann
}

func ExampleFilter() {
	stock := map[string]int{"apple": 3, "pear": 0, "plum": 5}
	printSorted(mapx.Filter(stock, func(_ string, n int) bool { return n > 0 }))
	// Output: apple=3 plum=5
}

func ExampleReject() {
	stock := map[string]int{"apple": 3, "pear": 0, "plum": 5}
	printSorted(mapx.Reject(stock, func(_ string, n int) bool { return n > 0 }))
	// Output: pear=0
}

func ExampleGetOr() {
	opts := map[string]int{"retries": 0}
	fmt.Println(mapx.GetOr(opts, "retries", 3))
	fmt.Println(mapx.GetOr(opts, "timeout", 30))
	// Output:
	// 0
	// 30
}

func ExampleHas() {
	m := map[string]int{"a": 1, "b": 2}
	fmt.Println(mapx.Has(m, "a", "b"))
	fmt.Println(mapx.Has(m, "a", "c"))
	// Output:
	// true
	// false
}

func ExampleHasAny() {
	m := map[string]int{"a": 1, "b": 2}
	fmt.Println(mapx.HasAny(m, "c", "b"))
	fmt.Println(mapx.HasAny(m, "c", "d"))
	// Output:
	// true
	// false
}

func ExampleAdd() {
	m := map[string]int{"a": 1}
	printSorted(mapx.Add(m, "a", 9))
	printSorted(mapx.Add(m, "b", 2))
	// Output:
	// a=1
	// a=1 b=2
}

func ExampleOrEmpty() {
	var labels map[string]string
	fmt.Println(labels == nil, mapx.OrEmpty(labels) == nil, len(mapx.OrEmpty(labels)))
	// Output: true false 0
}

func ExampleMerge() {
	defaults := map[string]int{"port": 80, "workers": 4}
	overrides := map[string]int{"port": 8080}
	printSorted(mapx.Merge(defaults, overrides))
	// Output: port=8080 workers=4
}

func ExampleMapValues() {
	prices := map[string]int{"apple": 120, "plum": 85}
	printSorted(mapx.MapValues(prices, func(_ string, cents int) float64 { return float64(cents) / 100 }))
	// Output: apple=1.2 plum=0.85
}

func ExampleMapKeys() {
	m := map[string]int{"a": 1, "b": 2}
	printSorted(mapx.MapKeys(m, func(k string, _ int) string { return strings.ToUpper(k) }))
	// Output: A=1 B=2
}

func ExampleInvert() {
	codes := map[string]int{"ok": 200, "not_found": 404}
	inverted := mapx.Invert(codes)
	fmt.Println(inverted[200], inverted[404])
	// Output: ok not_found
}

func ExampleSortedKeys() {
	fmt.Println(mapx.SortedKeys(map[string]int{"b": 2, "c": 3, "a": 1}))
	// Output: [a b c]
}

func ExampleSortedEntries() {
	for _, e := range mapx.SortedEntries(map[int]string{3: "c", 1: "a", 2: "b"}) {
		fmt.Println(e.Key, e.Value)
	}
	// Output:
	// 1 a
	// 2 b
	// 3 c
}
