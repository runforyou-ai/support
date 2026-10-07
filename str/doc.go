// Package str provides string helpers in the spirit of Laravel's Str, plus
// Stringable, an immutable value type for fluent chaining.
//
// All helpers operate on runes: lengths, offsets and truncation are measured
// in characters, never bytes. Functions never mutate their inputs and do not
// panic on ordinary input such as empty strings or out-of-range offsets.
//
// Example:
//
//	str.Snake("HTTPServerError")        // "http_server_error"
//	str.Limit("The quick brown fox", 9, "...") // "The quick..."
//	str.Of("  hello   world ").Squish().Studly().String() // "HelloWorld"
package str
