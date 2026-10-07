package support

// Rescue returns the result of fn, or fallback when fn returns an error or
// panics. The error or panic is discarded; use Catch to keep it.
func Rescue[T any](fn func() (T, error), fallback T) (result T) {
	defer func() {
		if recover() != nil {
			result = fallback
		}
	}()
	v, err := fn()
	if err != nil {
		return fallback
	}
	return v
}
