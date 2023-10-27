// nolint: wrapcheck
package utils

import (
	"context"
	"time"
)

// Sleep is a time.Sleep, but it can be interrupted by context.
// If interrupted, Sleep returns ctx.Err().
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	// do something every second until do something else finish.
//	go func() {
//		for {
//			doSomething()
//			if err := Sleep(ctx, 1 * time.Second); err != nil {
//				return // interrupted
//			}
//		}
//	}()
//
//	doSomethingElse()
//	cancel() // interrupt the goroutine
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
