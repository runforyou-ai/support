package convert

import "errors"

// ErrUnsupported is returned when a value's type cannot be converted to the
// requested type.
var ErrUnsupported = errors.New("convert: unsupported type")

// ErrOutOfRange is returned when a value does not fit in the requested type,
// including NaN and infinite floats converted to integers and negative values
// converted to unsigned integers.
var ErrOutOfRange = errors.New("convert: value out of range")
