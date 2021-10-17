package healthcheck

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p"
)

// P2PTracker will periodically check to remove expired temporary p2p's (key, value)
type P2PTracker struct {
	pollDuration  time.Duration
	expireTimeMap map[string]time.Time
	p2p           p2p.P2P
	mtx           sync.Mutex
}

// NewP2PTracker creates a tracker
func NewP2PTracker(p2p p2p.P2P, pollDuration time.Duration) *P2PTracker {
	return &P2PTracker{
		p2p:           p2p,
		expireTimeMap: map[string]time.Time{},
		pollDuration:  pollDuration,
	}
}

// Track add key into tracking list
func (tracker *P2PTracker) Track(key string, expires time.Time) {
	tracker.mtx.Lock()
	defer tracker.mtx.Unlock()
	tracker.expireTimeMap[key] = expires
}

// Run executes main delete function in background
func (tracker *P2PTracker) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, "p2ptracker")
	log.WithContext(ctx).Info("p2ptracker started")
	defer log.WithContext(ctx).Info("p2ptracker stopped")

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(tracker.pollDuration):
			tracker.checkExpires(ctx)
		}
	}
}

func (tracker *P2PTracker) checkExpires(ctx context.Context) {
	deleteKeys := []string{}

	// determines which key need to be delete
	tracker.mtx.Lock()
	now := time.Now()
	for key, expires := range tracker.expireTimeMap {
		if now.After(expires) {
			deleteKeys = append(deleteKeys, key)
			delete(tracker.expireTimeMap, key)
		}
	}
	tracker.mtx.Unlock()

	// try to delete
	for _, key := range deleteKeys {
		err := tracker.p2p.Delete(ctx, key)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("key", key).Error("RemoveKeyFailed")
		} else {
			log.WithContext(ctx).WithField("key", key).Info("RemoveKeySuccessfully")
		}
	}
}
