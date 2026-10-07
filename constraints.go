package support

// Signed is satisfied by all signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is satisfied by all unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is satisfied by all integer types.
type Integer interface {
	Signed | Unsigned
}

// Float is satisfied by all floating-point types.
type Float interface {
	~float32 | ~float64
}

// Number is satisfied by all integer and floating-point types.
type Number interface {
	Integer | Float
}
