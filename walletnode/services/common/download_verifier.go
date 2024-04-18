package common

import (
	"context"
	"time"
)

const (
	//RetryTime is the time WN will retry to download if it gets an error during data validation
	RetryTime = 10
)

// DownloadVerifier represents an interface with the p2p data validation download function
type DownloadVerifier interface {
	Download(ctx context.Context, hashOnly bool) error
}

// DownloadWithRetry represents a func that will retry downloading data from p2p for one minute
func DownloadWithRetry(ctx context.Context, fn DownloadVerifier, retryN, maxRetry time.Time, hashOnly bool) error {
	err := fn.Download(ctx, hashOnly)
	if err != nil {
		if retryN.After(maxRetry) {
			return err
		}

		time.Sleep(15 * time.Second)
		return DownloadWithRetry(ctx, fn, time.Now().UTC(), maxRetry, hashOnly)
	}

	return nil
}
