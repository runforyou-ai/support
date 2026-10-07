// Package data reads and writes nested JSON-like values with dot-notation
// paths, in the spirit of Laravel's data_get, data_set and Arr::dot.
//
// A tree is built from map[string]any and []any values, such as the result of
// decoding JSON into an any. Paths separate segments with ".", numeric
// segments index slices, and the "*" segment matches every key or element.
// Readers also traverse other maps with string keys, slices and arrays;
// writers traverse only map[string]any and []any.
//
// Set, Fill and Forget modify their target in place. Every other function
// leaves its input unchanged.
//
// Example:
//
//	var doc map[string]any
//	_ = json.Unmarshal(body, &doc)
//
//	name := data.GetOr(doc, "user.name", "anonymous")
//	emails, _ := data.Get(doc, "users.*.email")
//	data.Set(doc, "user.settings.theme", "dark")
package data
