package common

import (
	"context"
	"time"
)

// DownloadVerifier represents an interface with the p2p data validation download function
type DownloadVerifier interface {
	Download(ctx context.Context) error
}

// DownloadWithRetry represents a func that will retry downloading data from p2p for one minute
func DownloadWithRetry(ctx context.Context, fn DownloadVerifier, retryN, maxRetry time.Time) error {
	err := fn.Download(ctx)
	if err != nil {
		if retryN.After(maxRetry) {
			return err
		}

		time.Sleep(20 * time.Millisecond)
		return DownloadWithRetry(ctx, fn, time.Now(), maxRetry)
	}

	return nil
}
