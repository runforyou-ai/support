package support

import (
	"errors"
	"testing"
)

func TestRescue(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (int, error)
		want int
	}{
		{"success", func() (int, error) { return 1, nil }, 1},
		{"success with zero", func() (int, error) { return 0, nil }, 0},
		{"error", func() (int, error) { return 1, errors.New("boom") }, -1},
		{"panic", func() (int, error) { panic("boom") }, -1},
		{"panic with error", func() (int, error) { panic(errors.New("boom")) }, -1},
		{"panic with nil", func() (int, error) { panic(nil) }, -1}, //nolint:govet // panic(nil) is the case under test.
		{"runtime panic", func() (int, error) {
			var m map[string]int
			m["a"] = 1 //nolint:staticcheck // The nil map write is the panic under test.
			return 1, nil
		}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rescue(tt.fn, -1); got != tt.want {
				t.Errorf("Rescue() = %d, want %d", got, tt.want)
			}
		})
	}
}
