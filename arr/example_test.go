package arr_test

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/runforyou-ai/support/arr"
)

func isEven(n int) bool { return n%2 == 0 }

func ExampleFirst() {
	fmt.Println(arr.First([]int{1, 2, 3, 4}, isEven))
	fmt.Println(arr.First([]int{5, 6}, nil))
	fmt.Println(arr.First([]int{1, 3}, isEven))
	// Output:
	// 2 true
	// 5 true
	// 0 false
}

func ExampleLast() {
	fmt.Println(arr.Last([]int{1, 2, 3, 4, 5}, isEven))
	fmt.Println(arr.Last([]int{5, 6}, nil))
	// Output:
	// 4 true
	// 6 true
}

func ExampleEvery() {
	fmt.Println(arr.Every([]int{2, 4}, isEven))
	fmt.Println(arr.Every([]int{2, 3}, isEven))
	fmt.Println(arr.Every(nil, isEven))
	// Output:
	// true
	// false
	// true
}

func ExampleCount() {
	fmt.Println(arr.Count([]int{1, 2, 3, 4}, isEven))
	// Output: 2
}

func ExampleFilter() {
	fmt.Println(arr.Filter([]int{1, 2, 3, 4}, isEven))
	// Output: [2 4]
}

func ExampleReject() {
	fmt.Println(arr.Reject([]int{1, 2, 3, 4}, isEven))
	// Output: [1 3]
}

func ExampleMap() {
	fmt.Printf("%q\n", arr.Map([]int{1, 2, 3}, strconv.Itoa))
	// Output: ["1" "2" "3"]
}

func ExampleFilterMap() {
	nums := arr.FilterMap([]string{"1", "x", "3"}, func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		return n, err == nil
	})
	fmt.Println(nums)
	// Output: [1 3]
}

func ExampleReduce() {
	product := arr.Reduce([]int{2, 3, 4}, 1, func(acc, n int) int { return acc * n })
	fmt.Println(product)
	// Output: 24
}

func ExampleFlatten() {
	fmt.Println(arr.Flatten([][]int{{1, 2}, {3}, nil, {4}}))
	// Output: [1 2 3 4]
}

func ExampleCrossJoin() {
	for _, row := range arr.CrossJoin([]string{"S", "M"}, []string{"red", "blue"}) {
		fmt.Println(row)
	}
	// Output:
	// [S red]
	// [S blue]
	// [M red]
	// [M blue]
}

func ExamplePad() {
	fmt.Println(arr.Pad([]int{1, 2}, 4, 0))
	fmt.Println(arr.Pad([]int{1, 2}, -4, 0))
	// Output:
	// [1 2 0 0]
	// [0 0 1 2]
}

func ExampleJoin() {
	fmt.Println(arr.Join([]string{"a", "b", "c"}, ", ", " and "))
	fmt.Println(arr.Join([]string{"a", "b"}, ", ", " or "))
	fmt.Println(arr.Join([]string{"a", "b", "c"}, ", ", ""))
	// Output:
	// a, b and c
	// a or b
	// a, b, c
}

func ExampleKeyBy() {
	type user struct {
		ID   int
		Name string
	}
	byID := arr.KeyBy([]user{{1, "Ann"}, {2, "Bob"}}, func(u user) int { return u.ID })
	fmt.Println(byID[2].Name)
	// Output: Bob
}

func ExampleGroupBy() {
	groups := arr.GroupBy([]string{"apple", "avocado", "banana"}, func(s string) byte { return s[0] })
	fmt.Println(groups['a'])
	fmt.Println(groups['b'])
	// Output:
	// [apple avocado]
	// [banana]
}

func ExamplePartition() {
	even, odd := arr.Partition([]int{1, 2, 3, 4, 5}, isEven)
	fmt.Println(even, odd)
	// Output: [2 4] [1 3 5]
}

func ExampleUnique() {
	fmt.Println(arr.Unique([]int{3, 1, 3, 2, 1}))
	// Output: [3 1 2]
}

func ExampleUniqueBy() {
	fmt.Println(arr.UniqueBy([]string{"Go", "go", "Rust"}, strings.ToLower))
	// Output: [Go Rust]
}

func ExampleDiff() {
	fmt.Println(arr.Diff([]int{1, 2, 3, 4, 5}, []int{2, 4}, []int{5}))
	// Output: [1 3]
}

func ExampleIntersect() {
	fmt.Println(arr.Intersect([]int{1, 2, 3, 4}, []int{2, 3, 4}, []int{4, 3}))
	// Output: [3 4]
}

func ExampleRandom() {
	v, ok := arr.Random([]string{"only"})
	fmt.Println(v, ok)
	_, ok = arr.Random([]string(nil))
	fmt.Println(ok)
	// Output:
	// only true
	// false
}

func ExampleSample() {
	picked := arr.Sample([]int{1, 2, 3, 4, 5}, 3)
	fmt.Println(len(picked))
	fmt.Println(len(arr.Sample([]int{1, 2}, 10)))
	// Output:
	// 3
	// 2
}

func ExampleShuffle() {
	s := []int{1, 2, 3, 4}
	shuffled := arr.Shuffle(s)
	slices.Sort(shuffled)
	fmt.Println(shuffled, s)
	// Output: [1 2 3 4] [1 2 3 4]
}

func ExampleSum() {
	fmt.Println(arr.Sum([]float64{1.5, 2.5, 3}))
	// Output: 7
}

func ExampleSumBy() {
	type line struct{ Qty int }
	fmt.Println(arr.SumBy([]line{{2}, {3}}, func(l line) int { return l.Qty }))
	// Output: 5
}

func ExampleAvg() {
	fmt.Println(arr.Avg([]int{1, 2, 3, 4}))
	fmt.Println(arr.Avg([]int(nil)))
	// Output:
	// 2.5
	// 0
}

func ExampleAvgBy() {
	type score struct{ Points int }
	fmt.Println(arr.AvgBy([]score{{80}, {95}}, func(s score) int { return s.Points }))
	// Output: 87.5
}
