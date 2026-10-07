package support

import (
	"context"
	"time"
)

// RetryOption configures Retry and RetryValue.
type RetryOption func(*retryConfig)

type retryConfig struct {
	backoff func(attempt int) time.Duration
	when    func(err error) bool
}

// WithDelay waits d between attempts.
func WithDelay(d time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.backoff = func(int) time.Duration { return d }
	}
}

// WithBackoff waits backoff(attempt) after the given failed attempt, where
// attempt starts at 1.
func WithBackoff(backoff func(attempt int) time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.backoff = backoff
	}
}

// WithExponentialBackoff waits base, then twice the previous delay after each
// failed attempt, capped at limit.
func WithExponentialBackoff(base, limit time.Duration) RetryOption {
	return WithBackoff(func(attempt int) time.Duration {
		d := base
		for i := 1; i < attempt && d < limit; i++ {
			if d > limit/2 {
				return limit
			}
			d *= 2
		}
		return min(d, limit)
	})
}

// WithRetryIf retries only while when returns true for the returned error.
func WithRetryIf(when func(err error) bool) RetryOption {
	return func(c *retryConfig) {
		c.when = when
	}
}

// Retry calls fn up to times times until it returns nil, passing the attempt
// number starting at 1. fn always runs at least once, even when times is
// zero or negative. It returns nil on success, the last error once the
// attempts are exhausted or the retry condition rejects an error, or the
// context error when ctx is done while waiting between attempts.
func Retry(ctx context.Context, times int, fn func(attempt int) error, opts ...RetryOption) error {
	_, err := RetryValue(ctx, times, func(attempt int) (struct{}, error) {
		return struct{}{}, fn(attempt)
	}, opts...)
	return err
}

// RetryValue is like Retry for functions that return a value. It returns the
// value of the first successful attempt.
func RetryValue[T any](ctx context.Context, times int, fn func(attempt int) (T, error), opts ...RetryOption) (T, error) {
	var cfg retryConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	for attempt := 1; ; attempt++ {
		v, err := fn(attempt)
		if err == nil || attempt >= times || (cfg.when != nil && !cfg.when(err)) {
			return v, err
		}
		var delay time.Duration
		if cfg.backoff != nil {
			delay = cfg.backoff(attempt)
		}
		if delay <= 0 {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return v, ctxErr
			}
			continue
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return v, ctx.Err()
		case <-timer.C:
		}
	}
}
