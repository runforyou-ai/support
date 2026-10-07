package random_test

import (
	"fmt"

	"github.com/runforyou-ai/support/random"
)

func ExampleBytes() {
	fmt.Println(len(random.Bytes(16)))
	fmt.Println(len(random.Bytes(0)))
	// Output:
	// 16
	// 0
}

func ExampleHex() {
	fmt.Println(len(random.Hex(16)))
	// Output: 32
}

func ExampleBase64URL() {
	fmt.Println(len(random.Base64URL(32)))
	// Output: 43
}

func ExampleSHA256Hex() {
	fmt.Println(random.SHA256Hex([]byte("abc")))
	// Output: ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
}

func ExampleToken() {
	token, hash := random.Token(32)
	fmt.Println(len(token), len(hash))
	fmt.Println(random.VerifyToken(token, hash))
	// Output:
	// 43 64
	// true
}

func ExampleHashToken() {
	fmt.Println(random.HashToken("abc"))
	// Output: ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad
}

func ExampleVerifyToken() {
	token, hash := random.Token(32)
	fmt.Println(random.VerifyToken(token, hash))
	fmt.Println(random.VerifyToken("guess", hash))
	// Output:
	// true
	// false
}
