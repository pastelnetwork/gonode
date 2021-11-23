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
	defaultRpcConnectTimeout           = 15 * time.Second
	// Update duration in case last update was success
	defaultSuccessUpdateDuration = 10 * time.Second
	// Update duration in case last update was failed - prevent too much call to pasteld
	defaultFailedUpdateDuration = 5 * time.Second
)

type PastelClient interface {
	GetBlockCount(ctx context.Context) (int32, error)
}

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
		ctx, cancel := context.WithTimeout(context.Background(), defaultRpcConnectTimeout)
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
