package data_test

import (
	"fmt"

	"github.com/runforyou-ai/support/data"
)

func ExampleGet() {
	doc := map[string]any{
		"user": map[string]any{"name": "Ann", "tags": []any{"admin", "dev"}},
		"users": []any{
			map[string]any{"name": "Ann"},
			map[string]any{"name": "Bob"},
		},
	}
	name, ok := data.Get(doc, "user.name")
	fmt.Println(name, ok)
	tag, _ := data.Get(doc, "user.tags.1")
	fmt.Println(tag)
	names, _ := data.Get(doc, "users.*.name")
	fmt.Println(names)
	_, ok = data.Get(doc, "user.age")
	fmt.Println(ok)
	// Output:
	// Ann true
	// dev
	// [Ann Bob]
	// false
}

func ExampleGetOr() {
	doc := map[string]any{"user": map[string]any{"name": "Ann"}}
	fmt.Println(data.GetOr(doc, "user.name", "anonymous"))
	fmt.Println(data.GetOr(doc, "user.nickname", "anonymous"))
	// Output:
	// Ann
	// anonymous
}

func ExampleGetAs() {
	doc := map[string]any{"port": 8080.0, "host": "localhost"}
	host, ok := data.GetAs[string](doc, "host")
	fmt.Println(host, ok)
	_, ok = data.GetAs[string](doc, "port")
	fmt.Println(ok)
	// Output:
	// localhost true
	// false
}

func ExampleHas() {
	doc := map[string]any{"user": map[string]any{"email": nil}}
	fmt.Println(data.Has(doc, "user.email"))
	fmt.Println(data.Has(doc, "user.phone"))
	// Output:
	// true
	// false
}

func ExampleSet() {
	doc := map[string]any{
		"users": []any{map[string]any{"name": "Ann"}, map[string]any{"name": "Bob"}},
	}
	doc = data.Set(doc, "settings.theme", "dark")
	doc = data.Set(doc, "users.*.active", true)
	doc = data.Set(doc, "users.1.name", "Bobby")
	fmt.Println(doc["settings"])
	fmt.Println(doc["users"])
	// Output:
	// map[theme:dark]
	// [map[active:true name:Ann] map[active:true name:Bobby]]
}

func ExampleFill() {
	doc := map[string]any{"theme": "light"}
	doc = data.Fill(doc, "theme", "dark")
	doc = data.Fill(doc, "lang", "en")
	fmt.Println(doc)
	// Output:
	// map[lang:en theme:light]
}

func ExampleForget() {
	doc := map[string]any{
		"user":  map[string]any{"name": "Ann", "password": "secret"},
		"users": []any{map[string]any{"name": "Bob", "password": "x"}},
	}
	doc = data.Forget(doc, "user.password")
	doc = data.Forget(doc, "users.*.password")
	fmt.Println(doc)
	// Output:
	// map[user:map[name:Ann] users:[map[name:Bob]]]
}

func ExampleDot() {
	flat := data.Dot(map[string]any{
		"db":   map[string]any{"host": "localhost", "port": 5432},
		"tags": []any{"a", "b"},
	})
	fmt.Println(flat)
	// Output:
	// map[db.host:localhost db.port:5432 tags.0:a tags.1:b]
}

func ExampleUndot() {
	nested := data.Undot(map[string]any{
		"db.host": "localhost",
		"db.port": 5432,
		"tags.0":  "a",
	})
	fmt.Println(nested)
	// Output:
	// map[db:map[host:localhost port:5432] tags:map[0:a]]
}

func ExampleQuery() {
	fmt.Println(data.Query(map[string]any{
		"q":      "go helpers",
		"page":   2,
		"filter": map[string]any{"active": true},
		"ids":    []any{3, 5},
	}))
	// Output:
	// filter%5Bactive%5D=1&ids%5B0%5D=3&ids%5B1%5D=5&page=2&q=go+helpers
}
