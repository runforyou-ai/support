package str_test

import (
	"fmt"
	"strings"

	"github.com/runforyou-ai/support/str"
)

func ExampleAfter() {
	fmt.Println(str.After("https://example.com/path", "://"))
	// Output: example.com/path
}

func ExampleAfterLast() {
	fmt.Println(str.AfterLast("app/Http/Controllers/UserController", "/"))
	// Output: UserController
}

func ExampleBefore() {
	fmt.Println(str.Before("user@example.com", "@"))
	// Output: user
}

func ExampleBeforeLast() {
	fmt.Println(str.BeforeLast("archive.tar.gz", "."))
	// Output: archive.tar
}

func ExampleBetween() {
	fmt.Println(str.Between("This is my name", "This", "name"))
	// Output: is my
}

func ExampleBetweenFirst() {
	fmt.Println(str.BetweenFirst("[a] bc [d]", "[", "]"))
	// Output: a
}

func ExampleCharAt() {
	c, ok := str.CharAt("你好世界", -1)
	fmt.Println(c, ok)
	_, ok = str.CharAt("abc", 10)
	fmt.Println(ok)
	// Output:
	// 界 true
	// false
}

func ExampleSubstr() {
	fmt.Println(str.Substr("Hello World", 6, 5))
	fmt.Println(str.Substr("Hello World", -5, 3))
	fmt.Println(str.Substr("Hello World", 0, -6))
	// Output:
	// World
	// Wor
	// Hello
}

func ExampleLength() {
	fmt.Println(str.Length("你好, world"))
	// Output: 9
}

func ExampleSubstrCount() {
	fmt.Println(str.SubstrCount("I like bananas and banana bread", "banana"))
	// Output: 2
}

func ExampleWordCount() {
	fmt.Println(str.WordCount("Hello, world - again!"))
	// Output: 3
}

func ExampleCamel() {
	fmt.Println(str.Camel("hello world-foo_bar"))
	// Output: helloWorldFooBar
}

func ExampleStudly() {
	fmt.Println(str.Studly("http_server_error"))
	// Output: HttpServerError
}

func ExampleSnake() {
	fmt.Println(str.Snake("HTTPServer"))
	fmt.Println(str.Snake("userID2"))
	// Output:
	// http_server
	// user_id2
}

func ExampleKebab() {
	fmt.Println(str.Kebab("fooBarBaz"))
	// Output: foo-bar-baz
}

func ExampleHeadline() {
	fmt.Println(str.Headline("steve_jobs"))
	fmt.Println(str.Headline("EmailNotificationSent"))
	// Output:
	// Steve Jobs
	// Email Notification Sent
}

func ExampleTitle() {
	fmt.Println(str.Title("a nice tITLE, isn't it?"))
	// Output: A Nice Title, Isn't It?
}

func ExampleUcfirst() {
	fmt.Println(str.Ucfirst("élan"))
	// Output: Élan
}

func ExampleLcfirst() {
	fmt.Println(str.Lcfirst("Hello"))
	// Output: hello
}

func ExampleUcSplit() {
	fmt.Printf("%q\n", str.UcSplit("FooBar"))
	// Output: ["Foo" "Bar"]
}

func ExampleLimit() {
	fmt.Println(str.Limit("The quick brown fox", 10, "..."))
	// Output: The quick...
}

func ExampleWords() {
	fmt.Println(str.Words("Perfectly balanced, as all things should be.", 3, " >>>"))
	// Output: Perfectly balanced, as >>>
}

func ExampleExcerpt() {
	s, ok := str.Excerpt("This is my name", "my", 3, "...")
	fmt.Println(s, ok)
	_, ok = str.Excerpt("This is my name", "your", 3, "...")
	fmt.Println(ok)
	// Output:
	// ...is my na... true
	// false
}

func ExampleMask() {
	fmt.Println(str.Mask("taylor@example.com", '*', 3, 0))
	fmt.Println(str.Mask("13812345678", '*', 3, 4))
	fmt.Println(str.Mask("4111111111111111", '*', 0, 12))
	// Output:
	// tay***************
	// 138****5678
	// ************1111
}

func ExamplePadLeft() {
	fmt.Println(str.PadLeft("42", 5, "0"))
	// Output: 00042
}

func ExamplePadRight() {
	fmt.Println(str.PadRight("James", 10, "-="))
	// Output: James-=-=-
}

func ExamplePadBoth() {
	fmt.Println(str.PadBoth("James", 10, "_"))
	// Output: __James___
}

func ExampleStart() {
	fmt.Println(str.Start("//path", "/"))
	// Output: /path
}

func ExampleFinish() {
	fmt.Println(str.Finish("path", "/"))
	// Output: path/
}

func ExampleWrap() {
	fmt.Println(str.Wrap("value", `"`, ""))
	fmt.Println(str.Wrap("value", "[", "]"))
	// Output:
	// "value"
	// [value]
}

func ExampleUnwrap() {
	fmt.Println(str.Unwrap("[value]", "[", "]"))
	// Output: value
}

func ExampleSquish() {
	fmt.Println(str.Squish("   laravel \t  php \n framework   "))
	// Output: laravel php framework
}

func ExampleReverse() {
	fmt.Println(str.Reverse("你好, Go"))
	// Output: oG ,好你
}

func ExampleSlug() {
	fmt.Println(str.Slug("Hello, World! 你好", "-"))
	// Output: hello-world-你好
}

func ExampleReplaceFirst() {
	fmt.Println(str.ReplaceFirst("the cat and the dog", "the", "a"))
	// Output: a cat and the dog
}

func ExampleReplaceLast() {
	fmt.Println(str.ReplaceLast("the cat and the dog", "the", "a"))
	// Output: the cat and a dog
}

func ExampleReplaceStart() {
	fmt.Println(str.ReplaceStart("Hello World", "Hello", "Goodbye"))
	fmt.Println(str.ReplaceStart("Hello World", "World", "Goodbye"))
	// Output:
	// Goodbye World
	// Hello World
}

func ExampleReplaceEnd() {
	fmt.Println(str.ReplaceEnd("Hello World", "World", "Gophers"))
	// Output: Hello Gophers
}

func ExampleReplaceArray() {
	fmt.Println(str.ReplaceArray("between ? and ?", "?", []string{"8:30", "9:00"}))
	// Output: between 8:30 and 9:00
}

func ExampleSwap() {
	fmt.Println(str.Swap("Tacos are great!", map[string]string{
		"Tacos": "Burritos",
		"great": "fantastic",
	}))
	// Output: Burritos are fantastic!
}

func ExampleRemove() {
	fmt.Println(str.Remove("Peter Piper picked a peck", "e", "P"))
	// Output: tr ipr pickd a pck
}

func ExampleContains() {
	fmt.Println(str.Contains("This is my name", "your", "my"))
	// Output: true
}

func ExampleContainsAll() {
	fmt.Println(str.ContainsAll("This is my name", "my", "name"))
	fmt.Println(str.ContainsAll("This is my name", "my", "your"))
	// Output:
	// true
	// false
}

func ExampleIs() {
	fmt.Println(str.Is("foo*", "foobar"))
	fmt.Println(str.Is("baz*", "foobar"))
	// Output:
	// true
	// false
}

func ExampleIsJSON() {
	fmt.Println(str.IsJSON(`{"first": "John"}`))
	fmt.Println(str.IsJSON(`{first: "John"}`))
	// Output:
	// true
	// false
}

func ExampleIsURL() {
	fmt.Println(str.IsURL("https://example.com/docs"))
	fmt.Println(str.IsURL("example.com"))
	// Output:
	// true
	// false
}

func ExampleIsUUID() {
	fmt.Println(str.IsUUID("a0a2a2d2-0b87-4a18-83f2-2529882be2de"))
	fmt.Println(str.IsUUID("laravel"))
	// Output:
	// true
	// false
}

func ExampleIsULID() {
	fmt.Println(str.IsULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
	fmt.Println(str.IsULID("laravel"))
	// Output:
	// true
	// false
}

func ExampleIsASCII() {
	fmt.Println(str.IsASCII("Taylor"))
	fmt.Println(str.IsASCII("ü"))
	// Output:
	// true
	// false
}

func ExampleRandom() {
	token := str.Random(40)
	fmt.Println(len(token))
	// Output: 40
}

func ExampleUUID() {
	id := str.UUID()
	fmt.Println(str.IsUUID(id), id[14:15])
	// Output: true 4
}

func ExampleUUIDv7() {
	id := str.UUIDv7()
	fmt.Println(str.IsUUID(id), id[14:15])
	// Output: true 7
}

func ExampleULID() {
	id := str.ULID()
	fmt.Println(str.IsULID(id), len(id))
	// Output: true 26
}

func ExampleOf() {
	s := str.Of("  hello   world  ").
		Squish().
		Studly().
		Append("Controller")
	fmt.Println(s)
	// Output: HelloWorldController
}

func ExampleStringable_When() {
	path := func(p string) string {
		return str.Of(p).
			Trim("/").
			WhenNotEmpty(func(s str.Stringable) str.Stringable { return s.Start("/") }).
			WhenEmpty(func(s str.Stringable) str.Stringable { return s.Append("/") }).
			String()
	}
	fmt.Println(path("users/42/"))
	fmt.Println(path("///"))
	// Output:
	// /users/42
	// /
}

func ExampleStringable_Pipe() {
	s := str.Of("Hello World").
		Pipe(strings.ToLower).
		Slug("_").
		When(true, func(s str.Stringable) str.Stringable { return s.Upper() })
	fmt.Println(s.String(), s.Length(), s.StartsWith("HELLO"))
	// Output: HELLO_WORLD 11 true
}

func ExampleStringable_Tap() {
	s := str.Of("draft").
		Tap(func(s str.Stringable) { fmt.Println("before:", s) }).
		Prepend("final ").
		Headline()
	fmt.Println(s)
	// Output:
	// before: draft
	// Final Draft
}
