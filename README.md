# support

[![Go Reference](https://pkg.go.dev/badge/github.com/runforyou-ai/support.svg)](https://pkg.go.dev/github.com/runforyou-ai/support)

General-purpose helpers for Go, organized by domain in the spirit of Laravel's `Illuminate\Support`. Zero third-party dependencies.

[简体中文](README.zh-CN.md)

```bash
go get github.com/runforyou-ai/support
```

Requires Go 1.27 or later.

## Packages

| Package | Contents |
| --- | --- |
| [`support`](https://pkg.go.dev/github.com/runforyou-ai/support) | Value helpers (`Blank`, `Filled`, `Deref`, `DerefOr`, `NilIfZero`, `MapPtr`, `Tap`, `With`, `Transform`), `Retry` with backoff, `Rescue`, `Catch` and `PanicError`, numeric type constraints |
| [`str`](https://pkg.go.dev/github.com/runforyou-ai/support/str) | String helpers and the fluent `Stringable`: extraction, case conversion, truncation, padding, replacement, slugs, checks (`IsUUID`, `NormalizeUUID`, `IsHTTPURL`, `IsEmail`, …), `MaskEmail`, random strings and ULID |
| [`arr`](https://pkg.go.dev/github.com/runforyou-ai/support/arr) | Slice helpers: `Filter`, `Map`, `MapErr`, `Reduce`, `GroupBy`, `KeyBy`, `Associate`, `Partition`, `Unique`, `Diff`, `Sum`, `Sample`, `OrEmpty` and more |
| [`mapx`](https://pkg.go.dev/github.com/runforyou-ai/support/mapx) | Generic map helpers: `Only`, `Except`, `Filter`, `Merge`, `MapValues`, `Invert`, `SortedKeys`, `OrEmpty` and more |
| [`set`](https://pkg.go.dev/github.com/runforyou-ai/support/set) | Generic `Set` type with `Add`, `Has`, `Union`, `Intersect`, `Diff` |
| [`data`](https://pkg.go.dev/github.com/runforyou-ai/support/data) | Dot-notation access to nested `map[string]any` / `[]any`: `Get`, `Set`, `Fill`, `Forget`, `Dot`, `Undot`, `Query` |
| [`number`](https://pkg.go.dev/github.com/runforyou-ai/support/number) | Number formatting: `Format`, `Percentage`, `FileSize`, `Abbreviate`, `ForHumans`, `Ordinal`, `Clamp`, `Pairs` |
| [`convert`](https://pkg.go.dev/github.com/runforyou-ai/support/convert) | Conversion from `any`: `ToString`, `ToBool`, `ToInt`, `ToInt64`, `ToUint64`, `ToFloat64` |
| [`random`](https://pkg.go.dev/github.com/runforyou-ai/support/random) | Secure random `Bytes`, `Hex`, `Base64URL`, `Token` with hash, `SHA256Hex` |
| [`validate`](https://pkg.go.dev/github.com/runforyou-ai/support/validate) | `E164` phone numbers and IANA `Timezone` names |
| [`filex`](https://pkg.go.dev/github.com/runforyou-ai/support/filex) | Safe `ExtractTarGz` / `ExtractZip`, `OpenRotating`, `WriteAtomic`, `Within`, `IsBinary` |
| [`iox`](https://pkg.go.dev/github.com/runforyou-ai/support/iox) | `HeadTailBuffer`, a writer that keeps the beginning and end of long output |

## Examples

```go
import (
	"github.com/runforyou-ai/support"
	"github.com/runforyou-ai/support/arr"
	"github.com/runforyou-ai/support/data"
	"github.com/runforyou-ai/support/number"
	"github.com/runforyou-ai/support/random"
	"github.com/runforyou-ai/support/set"
	"github.com/runforyou-ai/support/str"
)

str.Of("  Hello World  ").Squish().Snake().String()           // "hello_world"
str.Of(title).When(draft, func(s str.Stringable) str.Stringable {
	return s.Prepend("[Draft] ")
}).Limit(20, "...").String()
str.Slug("Hello 世界!", "-")                                   // "hello-世界"

adults := arr.Filter(users, func(u User) bool { return u.Age >= 18 })
byTeam := arr.GroupBy(users, func(u User) string { return u.Team })

name, _ := data.GetAs[string](payload, "user.profile.name")
number.FileSize(1536, 2)                                       // "1.50 KB"

limit := support.DerefOr(req.Limit, 20)
seen := set.Of(ids...)
token, hash := random.Token(32)
err := support.Retry(ctx, 3, send, support.WithExponentialBackoff(100*time.Millisecond, 2*time.Second))
```

Every package-level function has a runnable example on [pkg.go.dev](https://pkg.go.dev/github.com/runforyou-ai/support).

## Design

- Only fills gaps in the standard library; it does not wrap `strings`, `slices`, `maps` or `cmp`.
- String helpers work on runes, so lengths, offsets and truncation are correct for multibyte text.
- Functions do not panic on ordinary input and never mutate their arguments; writers such as `data.Set` return a new value.
- Framework concerns (paths, config, routing, logging) are out of scope.

## License

[MIT](LICENSE)
