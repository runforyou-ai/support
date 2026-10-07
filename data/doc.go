// Package data reads and writes nested JSON-like values with dot-notation
// paths, in the spirit of Laravel's data_get, data_set and Arr::dot.
//
// A tree is built from map[string]any and []any values, such as the result of
// decoding JSON into an any. Paths separate segments with ".", numeric
// segments index slices, and the "*" segment matches every key or element.
// Readers and Query also traverse other maps with string keys, slices,
// arrays and non-nil pointers to them; writers traverse only map[string]any
// and []any. Values must not contain reference cycles.
//
// No function modifies its input: Set, Fill and Forget return a new map in
// which the maps and slices along the path are copied.
//
// Example:
//
//	var doc map[string]any
//	_ = json.Unmarshal(body, &doc)
//
//	name := data.GetOr(doc, "user.name", "anonymous")
//	emails, _ := data.Get(doc, "users.*.email")
//	doc = data.Set(doc, "user.settings.theme", "dark")
package data
