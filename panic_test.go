package support

import (
	"errors"
	"io"
	"runtime"
	"strings"
	"testing"
)

func TestCatch(t *testing.T) {
	errFailed := errors.New("failed")
	tests := []struct {
		name      string
		fn        func() error
		wantErr   error
		wantPanic bool
		wantValue string
	}{
		{"success", func() error { return nil }, nil, false, ""},
		{"error", func() error { return errFailed }, errFailed, false, ""},
		{"panic with string", func() error { panic("boom") }, nil, true, "boom"},
		{"panic with error", func() error { panic(io.EOF) }, io.EOF, true, "EOF"},
		{"panic with nil", func() error { panic(nil) }, nil, true, "runtime error: panic called with nil argument"}, //nolint:govet // panic(nil) is the case under test.
		{"runtime panic", func() error {
			var m map[string]int
			m["a"] = 1 //nolint:staticcheck // The nil map write is the panic under test.
			return nil
		}, nil, true, "assignment to entry in nil map"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Catch(tt.fn)
			var pe *PanicError
			isPanic := errors.As(err, &pe)
			if isPanic != tt.wantPanic {
				t.Fatalf("Catch() = %v, want panic=%v", err, tt.wantPanic)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("Catch() = %v, want errors.Is %v", err, tt.wantErr)
			}
			if !tt.wantPanic {
				if err != tt.wantErr { //nolint:errorlint // Catch returns fn's error unchanged.
					t.Fatalf("Catch() = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if got := errorString(pe.Value); got != tt.wantValue {
				t.Errorf("PanicError.Value = %q, want %q", got, tt.wantValue)
			}
			msg := pe.Error()
			if !strings.HasPrefix(msg, "panic: "+tt.wantValue+"\n") || !strings.Contains(msg, "goroutine") {
				t.Errorf("PanicError.Error() = %q, want panic value and stack", msg)
			}
			if len(pe.Stack) == 0 {
				t.Error("PanicError.Stack is empty")
			}
			if !framesContain(pe.StackTrace(), "TestCatch") {
				t.Error("PanicError.StackTrace() does not include the panicking frames")
			}
		})
	}
}

func TestPanicErrorUnwrap(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  error
	}{
		{"error value", io.EOF, io.EOF},
		{"string value", "boom", nil},
		{"nil value", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPanicError(tt.value).Unwrap(); got != tt.want { //nolint:errorlint // Unwrap returns the value unchanged.
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func errorString(v any) string {
	if err, ok := v.(error); ok {
		return err.Error()
	}
	s, _ := v.(string)
	return s
}

func framesContain(pcs []uintptr, name string) bool {
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		if strings.Contains(f.Function, name) {
			return true
		}
		if !more {
			return false
		}
	}
}
