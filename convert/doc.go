// Package convert converts loosely typed values, such as decoded JSON,
// configuration entries or form input, to basic Go types.
//
// Conversions accept the built-in scalar types and named types derived from
// them, json.Number, time.Duration, and pointers to any of these. Failures
// wrap ErrUnsupported, ErrOutOfRange or a *strconv.NumError so callers can
// test them with errors.Is and errors.As.
//
// Example:
//
//	port, err := convert.ToInt(cfg["port"])
//	debug := convert.ToBoolOr(os.Getenv("DEBUG"), false)
//	label := convert.ToString(cfg["name"])
package convert
