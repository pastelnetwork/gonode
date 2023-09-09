package blocktracker

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	defaultRetries                     = 3
	defaultDelayDurationBetweenRetries = 5 * time.Second
	defaultRPCConnectTimeout           = 15 * time.Second
	// Update duration in case last update was success
	defaultSuccessUpdateDuration = 10 * time.Second
	// Update duration in case last update was failed - prevent too much call to pasteld
	defaultFailedUpdateDuration = 5 * time.Second
	defaultNextBlockTimeout     = 30 * time.Minute
)

// PastelClient defines interface functions BlockCntTracker expects from pastel
type PastelClient interface {
	// GetBlockCount returns block height of blockchain
	GetBlockCount(ctx context.Context) (int32, error)
}

// BlockCntTracker defines a block tracker - that will keep current block height
type BlockCntTracker struct {
	mtx                 sync.Mutex
	pastelClient        PastelClient
	curBlockCnt         int32
	lastSuccess         time.Time
	lastRetried         time.Time
	lastErr             error
	delayBetweenRetries time.Duration
	retries             int
}

// New returns an instance of BlockCntTracker
func New(pastelClient PastelClient) *BlockCntTracker {
	return &BlockCntTracker{
		pastelClient:        pastelClient,
		curBlockCnt:         0,
		delayBetweenRetries: defaultDelayDurationBetweenRetries,
		retries:             defaultRetries,
	}
}

func (tracker *BlockCntTracker) refreshBlockCount(retries int) {
	tracker.lastRetried = time.Now()
	for i := 0; i < retries; i = i + 1 {
		ctx, cancel := context.WithTimeout(context.Background(), defaultRPCConnectTimeout)
		blockCnt, err := tracker.pastelClient.GetBlockCount(ctx)
		if err == nil {
			tracker.curBlockCnt = blockCnt
			tracker.lastSuccess = time.Now()
			cancel()
			tracker.lastErr = nil
			return
		}
		cancel()

		tracker.lastErr = err
		// delay between retries
		time.Sleep(tracker.delayBetweenRetries)
	}

}

// GetBlockCount return current block count
// it will get from cache if last refresh is small than defaultSuccessUpdateDuration
// or will refresh it by call from pastel daemon to get the latest one if defaultSuccessUpdateDuration expired
func (tracker *BlockCntTracker) GetBlockCount() (int32, error) {
	tracker.mtx.Lock()
	defer tracker.mtx.Unlock()

	shouldRefresh := false

	if tracker.lastSuccess.After(tracker.lastRetried) {
		if time.Now().After(tracker.lastSuccess.Add(defaultSuccessUpdateDuration)) {
			shouldRefresh = true
		}
	} else {
		// prevent update too much
		if time.Now().After(tracker.lastRetried.Add(defaultFailedUpdateDuration)) {
			shouldRefresh = true
		}
	}

	if shouldRefresh {
		tracker.refreshBlockCount(tracker.retries)
	}

	if tracker.curBlockCnt == 0 {
		return 0, errors.Errorf("failed to get blockcount: %w", tracker.lastErr)
	}

	return tracker.curBlockCnt, nil
}

func (tracker *BlockCntTracker) WaitTillNextBlock(ctx context.Context, blockCnt int32) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(defaultNextBlockTimeout):
			return errors.Errorf("timeout waiting for next block")
		case <-time.After(defaultSuccessUpdateDuration):
			curBlockCnt, err := tracker.GetBlockCount()
			if err != nil {
				return errors.Errorf("failed to get blockcount: %w", err)
			}

			if curBlockCnt > blockCnt {
				return nil
			}
		}
	}
}
