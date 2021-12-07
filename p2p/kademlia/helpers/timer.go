package helpers

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

// StartTimer news a timer for doing the function
func StartTimer(ctx context.Context, name string, done chan struct{}, interval time.Duration, fn func() error) {
	go func() {
		// new a timer for doing the function periodly
		ticker := time.NewTimer(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				// do the function
				if err := fn(); err != nil {
					log.P2P().WithContext(ctx).WithError(err).Errorf("run worker(%s) failed", name)
				}

				// reset the timer
				ticker.Reset(interval)
			}
		}
	}()
}
