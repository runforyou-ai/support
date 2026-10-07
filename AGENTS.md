# support

General-purpose helpers for Go, organized by domain in the spirit of Laravel's `Illuminate\Support`.

## Scope

- Each package covers one domain (`str`, `arr`, `mapx`, `set`, `data`, `number`, `convert`, `random`, `validate`, `filex`, `iox`); the root package `support` holds generic value helpers and shared type constraints.
- The module targets the Go version in `go.mod`. Add a function only when that standard library (`strings`, `slices`, `maps`, `cmp`, `iter`, `sync`, `uuid`, builtin `new(expr)`, …) does not already provide it in a direct form.
- Framework concerns (paths, config, routing, logging) and thin wrappers around third-party libraries are out of scope.

## Rules

- The module has no third-party dependencies, including in tests.
- Subpackages may import the root package `support`; the root package imports no subpackage.
- String helpers operate on runes: lengths, offsets and truncation are measured in characters, never bytes.
- Functions do not panic on ordinary input (empty values, out-of-range offsets, nil maps); they return zero values or the input unchanged. Lookups that may miss return `(value, ok)`.
- Functions never mutate their arguments; they return new slices and maps. Methods of container types, such as `set.Set`, may modify their receiver.
- Every exported identifier has an English doc comment starting with its name; every package has a `doc.go` with a package comment and a short usage example.
- Every exported function has table-driven tests in `<file>_test.go`; every package-level function also has a runnable `Example` with an `// Output:` block in `example_test.go`.
- Tests compile and pass on 32-bit platforms; values such as `math.MinInt64` go in `int64` fields, not `int`.
- Tests use only the standard `testing` package.

## Commands

```bash
go test -race ./...
go vet ./...
golangci-lint run
```
