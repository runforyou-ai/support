package support

import (
	"context"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"
)

var (
	errTemporary = errors.New("temporary")
	errFatal     = errors.New("fatal")
)

// failUntil returns a function that fails with errTemporary until the given
// attempt and counts its calls.
func failUntil(success int, calls *int) func(attempt int) (int, error) {
	return func(attempt int) (int, error) {
		*calls++
		if attempt < success {
			return attempt, errTemporary
		}
		return attempt * 10, nil
	}
}

func TestRetryValue(t *testing.T) {
	tests := []struct {
		name      string
		times     int
		success   int
		opts      []RetryOption
		wantValue int
		wantErr   error
		wantCalls int
	}{
		{"first attempt", 3, 1, nil, 10, nil, 1},
		{"third attempt", 3, 3, nil, 30, nil, 3},
		{"exhausted returns last error", 3, 10, nil, 3, errTemporary, 3},
		{"zero times runs once", 0, 10, nil, 1, errTemporary, 1},
		{"negative times runs once", -5, 1, nil, 10, nil, 1},
		{"one time", 1, 2, nil, 1, errTemporary, 1},
		{"retry if allows", 5, 3, []RetryOption{WithRetryIf(func(err error) bool {
			return errors.Is(err, errTemporary)
		})}, 30, nil, 3},
		{"retry if stops", 5, 3, []RetryOption{WithRetryIf(func(error) bool { return false })}, 1, errTemporary, 1},
		{"with delay", 3, 3, []RetryOption{WithDelay(time.Microsecond)}, 30, nil, 3},
		{"with zero delay", 3, 3, []RetryOption{WithDelay(0)}, 30, nil, 3},
		{"with backoff", 3, 3, []RetryOption{WithBackoff(func(attempt int) time.Duration {
			return time.Duration(attempt) * time.Microsecond
		})}, 30, nil, 3},
		{"with exponential backoff", 4, 4, []RetryOption{WithExponentialBackoff(time.Microsecond, 3*time.Microsecond)}, 40, nil, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			v, err := RetryValue(context.Background(), tt.times, failUntil(tt.success, &calls), tt.opts...)
			if v != tt.wantValue || !errors.Is(err, tt.wantErr) || calls != tt.wantCalls {
				t.Errorf("RetryValue() = (%d, %v) after %d calls, want (%d, %v) after %d calls",
					v, err, calls, tt.wantValue, tt.wantErr, tt.wantCalls)
			}
		})
	}
}

func TestRetry(t *testing.T) {
	tests := []struct {
		name      string
		times     int
		errs      []error
		opts      []RetryOption
		wantErr   error
		wantCalls int
	}{
		{"success", 3, []error{nil}, nil, nil, 1},
		{"second attempt", 3, []error{errTemporary, nil}, nil, nil, 2},
		{"exhausted", 2, []error{errTemporary, errFatal, nil}, nil, errFatal, 2},
		{"zero times", 0, []error{errTemporary, nil}, nil, errTemporary, 1},
		{"retry if stops on fatal", 5, []error{errTemporary, errFatal, nil}, []RetryOption{WithRetryIf(func(err error) bool {
			return !errors.Is(err, errFatal)
		})}, errFatal, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attempts []int
			err := Retry(context.Background(), tt.times, func(attempt int) error {
				attempts = append(attempts, attempt)
				return tt.errs[attempt-1]
			}, tt.opts...)
			if !errors.Is(err, tt.wantErr) || len(attempts) != tt.wantCalls {
				t.Errorf("Retry() = %v after attempts %v, want %v after %d calls", err, attempts, tt.wantErr, tt.wantCalls)
			}
			for i, a := range attempts {
				if a != i+1 {
					t.Errorf("attempt numbers = %v, want 1..n", attempts)
					break
				}
			}
		})
	}
}

func TestBackoffOptions(t *testing.T) {
	tests := []struct {
		name string
		opt  RetryOption
		want []time.Duration
	}{
		{"delay", WithDelay(5 * time.Millisecond), []time.Duration{5 * time.Millisecond, 5 * time.Millisecond, 5 * time.Millisecond}},
		{"backoff", WithBackoff(func(attempt int) time.Duration {
			return time.Duration(attempt) * time.Second
		}), []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}},
		{"exponential", WithExponentialBackoff(100*time.Millisecond, time.Second), []time.Duration{
			100 * time.Millisecond, 200 * time.Millisecond, 400 * time.Millisecond, 800 * time.Millisecond, time.Second, time.Second,
		}},
		{"exponential base above limit", WithExponentialBackoff(2*time.Second, time.Second), []time.Duration{time.Second, time.Second}},
		{"exponential zero base", WithExponentialBackoff(0, time.Second), []time.Duration{0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg retryConfig
			tt.opt(&cfg)
			for i, want := range tt.want {
				if got := cfg.backoff(i + 1); got != want {
					t.Errorf("backoff(%d) = %v, want %v", i+1, got, want)
				}
			}
		})
	}
}

func TestExponentialBackoffNoOverflow(t *testing.T) {
	var cfg retryConfig
	WithExponentialBackoff(time.Second, math.MaxInt64)(&cfg)
	prev := time.Duration(0)
	for attempt := 1; attempt <= 100; attempt++ {
		d := cfg.backoff(attempt)
		if d < prev {
			t.Fatalf("backoff(%d) = %v, less than backoff(%d) = %v", attempt, d, attempt-1, prev)
		}
		prev = d
	}
	if prev != math.MaxInt64 {
		t.Errorf("backoff(100) = %v, want the limit", prev)
	}
}

func TestRetryContext(t *testing.T) {
	tests := []struct {
		name      string
		opts      []RetryOption
		cancel    func(attempt int, cancel context.CancelFunc)
		wantCalls int
	}{
		// The first wait is short, so the cancellation lands in the second, long wait.
		{"canceled during wait", []RetryOption{WithBackoff(func(attempt int) time.Duration {
			if attempt == 1 {
				return time.Microsecond
			}
			return time.Hour
		})}, func(attempt int, cancel context.CancelFunc) {
			if attempt == 2 {
				time.AfterFunc(5*time.Millisecond, cancel)
			}
		}, 2},
		{"canceled before wait", []RetryOption{WithDelay(time.Hour)}, func(attempt int, cancel context.CancelFunc) {
			if attempt == 1 {
				cancel()
			}
		}, 1},
		{"canceled without delay", nil, func(attempt int, cancel context.CancelFunc) {
			if attempt == 1 {
				cancel()
			}
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			v, err := RetryValue(ctx, 5, func(attempt int) (string, error) {
				calls++
				tt.cancel(attempt, cancel)
				return strconv.Itoa(attempt), errTemporary
			}, tt.opts...)
			if !errors.Is(err, context.Canceled) || calls != tt.wantCalls || v != strconv.Itoa(tt.wantCalls) {
				t.Errorf("RetryValue() = (%q, %v) after %d calls, want (%q, context.Canceled) after %d calls",
					v, err, calls, strconv.Itoa(tt.wantCalls), tt.wantCalls)
			}
		})
	}
}

func TestRetryContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := Retry(ctx, 3, func(int) error { return errTemporary }, WithDelay(time.Hour))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Retry() = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > time.Minute {
		t.Errorf("Retry() waited %v, want it to stop at the deadline", elapsed)
	}
}
