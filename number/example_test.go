package number_test

import (
	"fmt"

	"github.com/runforyou-ai/support/number"
)

func ExampleFormat() {
	fmt.Println(number.Format(1234567.891, 2))
	fmt.Println(number.Format(2.5, 0))
	fmt.Println(number.Format(-1234.5, -1))
	// Output:
	// 1,234,567.89
	// 3
	// -1,234.5
}

func ExampleFormatWith() {
	fmt.Println(number.FormatWith(1234567.891, 2, ",", "."))
	// Output: 1.234.567,89
}

func ExamplePercentage() {
	fmt.Println(number.Percentage(12.5, 2))
	// Output: 12.50%
}

func ExampleFileSize() {
	fmt.Println(number.FileSize(1536, 2))
	fmt.Println(number.FileSize(5*1024*1024*1024, 0))
	// Output:
	// 1.50 KB
	// 5 GB
}

func ExampleAbbreviate() {
	fmt.Println(number.Abbreviate(1000, 2))
	fmt.Println(number.Abbreviate(1234567, 2))
	// Output:
	// 1K
	// 1.23M
}

func ExampleForHumans() {
	fmt.Println(number.ForHumans(1500000, 1))
	// Output: 1.5 million
}

func ExampleOrdinal() {
	fmt.Println(number.Ordinal(1), number.Ordinal(12), number.Ordinal(23), number.Ordinal(112))
	// Output: 1st 12th 23rd 112th
}

func ExampleClamp() {
	fmt.Println(number.Clamp(15, 1, 10))
	fmt.Println(number.Clamp(0.5, 1.0, 0.0))
	// Output:
	// 10
	// 0.5
}

func ExamplePairs() {
	fmt.Println(number.Pairs(25, 10, 0))
	fmt.Println(number.Pairs(25, 10, 1))
	// Output:
	// [[0 9] [10 19] [20 25]]
	// [[1 10] [11 20] [21 25]]
}
