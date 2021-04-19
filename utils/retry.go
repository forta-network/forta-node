package utils

import (
	"context"
	"time"
)

func sliceContains(needle string, haystack []string) bool {
	for _, str := range haystack {
		if str == needle {
			return true
		}
	}
	return false
}

func RetryForErrors(ctx context.Context, handler func() error, errs []string, interval time.Duration, timeout time.Duration) error {
	timeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ticker := time.NewTicker(interval)
	for {
		if timeCtx.Err() != nil {
			return ctx.Err()
		}
		err := handler()
		if err == nil {
			return nil
		}
		if !sliceContains(err.Error(), errs) {
			return err
		}
		<-ticker.C
	}

}
