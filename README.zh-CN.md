# support

[![Go Reference](https://pkg.go.dev/badge/github.com/runforyou-ai/support.svg)](https://pkg.go.dev/github.com/runforyou-ai/support)

Go 通用辅助函数库，参照 Laravel `Illuminate\Support` 按领域分包，不依赖任何第三方库。

[English](README.md)

```bash
go get github.com/runforyou-ai/support
```

需要 Go 1.27 及以上版本。

## 包

| 包 | 内容 |
| --- | --- |
| [`support`](https://pkg.go.dev/github.com/runforyou-ai/support) | 值辅助函数（`Blank`、`Filled`、`Deref`、`DerefOr`、`NilIfZero`、`MapPtr`、`Tap`、`With`、`Transform`）、带退避的 `Retry`、`Rescue`、`Catch` 与 `PanicError`、数值类型约束 |
| [`str`](https://pkg.go.dev/github.com/runforyou-ai/support/str) | 字符串函数与链式 `Stringable`：截取、大小写转换、截断、填充、替换、Slug、格式校验（`IsUUID`、`NormalizeUUID`、`IsHTTPURL`、`IsEmail` 等）、`MaskEmail`、随机字符串与 ULID |
| [`arr`](https://pkg.go.dev/github.com/runforyou-ai/support/arr) | 切片函数：`Filter`、`Map`、`MapErr`、`Reduce`、`GroupBy`、`KeyBy`、`Associate`、`Partition`、`Unique`、`Diff`、`Sum`、`Sample`、`OrEmpty` 等 |
| [`mapx`](https://pkg.go.dev/github.com/runforyou-ai/support/mapx) | 泛型 map 函数：`Only`、`Except`、`Filter`、`Merge`、`MapValues`、`Invert`、`SortedKeys`、`OrEmpty` 等 |
| [`set`](https://pkg.go.dev/github.com/runforyou-ai/support/set) | 泛型集合 `Set`：`Add`、`Has`、`Union`、`Intersect`、`Diff` |
| [`data`](https://pkg.go.dev/github.com/runforyou-ai/support/data) | 以点号路径读写嵌套的 `map[string]any` / `[]any`：`Get`、`Set`、`Fill`、`Forget`、`Dot`、`Undot`、`Query` |
| [`number`](https://pkg.go.dev/github.com/runforyou-ai/support/number) | 数字格式化：`Format`、`Percentage`、`FileSize`、`Abbreviate`、`ForHumans`、`Ordinal`、`Clamp`、`Pairs` |
| [`convert`](https://pkg.go.dev/github.com/runforyou-ai/support/convert) | 从 `any` 转换：`ToString`、`ToBool`、`ToInt`、`ToInt64`、`ToUint64`、`ToFloat64` |
| [`random`](https://pkg.go.dev/github.com/runforyou-ai/support/random) | 安全随机数：`Bytes`、`Hex`、`Base64URL`、带哈希的 `Token`、`SHA256Hex` |
| [`validate`](https://pkg.go.dev/github.com/runforyou-ai/support/validate) | `E164` 手机号与 IANA `Timezone` 时区名校验 |
| [`filex`](https://pkg.go.dev/github.com/runforyou-ai/support/filex) | 防路径穿越的 `ExtractTarGz` / `ExtractZip`、按大小轮转的 `OpenRotating`、`WriteAtomic`、`Within`、`IsBinary` |
| [`iox`](https://pkg.go.dev/github.com/runforyou-ai/support/iox) | `HeadTailBuffer`：保留长输出首尾内容的 Writer |

## 示例

```go
str.Of("  Hello World  ").Squish().Snake().String()           // "hello_world"
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

每个包级函数在 [pkg.go.dev](https://pkg.go.dev/github.com/runforyou-ai/support) 上都有可运行的示例。

## 设计原则

- 只补标准库缺少的能力，不包装 `strings`、`slices`、`maps`、`cmp` 已有的函数。
- 字符串函数按字符（rune）计算长度、位置和截断，多字节文本结果正确。
- 普通输入不会 panic，也不修改传入的参数；`data.Set` 等写入函数返回新值。
- 不包含路径、配置、路由等框架相关功能。

## 许可证

[MIT](LICENSE)
