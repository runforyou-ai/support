package support

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// PanicError is an error built from a recovered panic. It keeps the panic
// value and the stack of the goroutine at the moment the panic was
// recovered.
type PanicError struct {
	// Value is the value passed to panic.
	Value any
	// Stack is the formatted stack trace, as returned by debug.Stack.
	Stack []byte
	pcs   []uintptr
}

// NewPanicError returns a PanicError for value with the current stack. Call
// it from the deferred function that recovered the panic, so that the stack
// still includes the frames that panicked.
func NewPanicError(value any) *PanicError {
	pcs := make([]uintptr, 64)
	return &PanicError{Value: value, Stack: debug.Stack(), pcs: pcs[:runtime.Callers(2, pcs)]}
}

// Error returns the panic value followed by the stack trace.
func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v\n%s", e.Value, e.Stack)
}

// Unwrap returns the panic value when it is an error, and nil otherwise.
func (e *PanicError) Unwrap() error {
	err, _ := e.Value.(error)
	return err
}

// StackTrace returns the program counters of the stack captured by
// NewPanicError, innermost frame first. Pass them to runtime.CallersFrames to
// resolve function names, files and lines.
func (e *PanicError) StackTrace() []uintptr {
	return e.pcs
}

// Catch calls fn and returns its error. When fn panics, Catch recovers and
// returns a *PanicError holding the panic value and stack instead.
func Catch(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewPanicError(r)
		}
	}()
	return fn()
}
