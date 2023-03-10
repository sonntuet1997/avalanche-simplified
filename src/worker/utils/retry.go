package utils

import (
	"context"
	"github.com/avast/retry-go"
	"time"
)

var (
	defaultType               = retry.FixedDelay
	defaultDelayInterval      = 3 * time.Second
	defaultAttempts      uint = 3
)

func Retry(ctx context.Context, fn func() error) error {
	return _retry(ctx, fn, defaultType, defaultAttempts, defaultDelayInterval)
}

func _retry(ctx context.Context, fn func() error, delayType retry.DelayTypeFunc, attempts uint, delayInterval time.Duration) error {
	return retry.Do(
		func() error {
			return fn()
		},
		retry.Context(ctx),
		retry.DelayType(delayType),
		retry.Attempts(attempts),
		retry.Delay(delayInterval),
		retry.LastErrorOnly(true),
		retry.RetryIf(retry.IsRecoverable),
	)
}
