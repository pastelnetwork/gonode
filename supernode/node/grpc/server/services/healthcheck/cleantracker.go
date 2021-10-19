package healthcheck

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

type CleanFunc func(ctx context.Context, key string)

type trackItem struct {
	key       string
	expires   time.Time
	cleanFunc CleanFunc
}

// CleanTracker is interface of lifetime tracker that will periodically check to remove expired temporary p2p's (key, value)
type CleanTracker interface {
	// Track monitor to delete temporary key
	Track(key string, cleanFunc CleanFunc, expires time.Time)

	// Run executes main delete function in background
	Run(ctx context.Context) error
}

// NewCleanTracker creates a tracker
func NewCleanTracker(pollDuration time.Duration) CleanTracker {
	return &cleanTrackerImpl{
		expireTrackMap: map[time.Time]*trackItem{},
		pollDuration:   pollDuration,
	}
}

// cleanTrackerImpl will periodically check to remove expired temporary p2p's (key, value)
type cleanTrackerImpl struct {
	pollDuration   time.Duration
	expireTrackMap map[time.Time]*trackItem
	mtx            sync.Mutex
}

// Track add key into tracking list
func (tracker *cleanTrackerImpl) Track(key string, cleanFunc CleanFunc, expires time.Time) {
	tracker.mtx.Lock()
	defer tracker.mtx.Unlock()
	tracker.expireTrackMap[time.Now()] = &trackItem{
		key:       key,
		expires:   expires,
		cleanFunc: cleanFunc,
	}
}

// Run executes main delete function in background
func (tracker *cleanTrackerImpl) Run(ctx context.Context) error {
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

func (tracker *cleanTrackerImpl) checkExpires(ctx context.Context) {
	deleteItems := []*trackItem{}

	// determines which key need to be delete
	tracker.mtx.Lock()
	now := time.Now()
	for key, trackItem := range tracker.expireTrackMap {
		if now.After(trackItem.expires) {
			deleteItems = append(deleteItems, trackItem)
			delete(tracker.expireTrackMap, key)
		}
	}
	tracker.mtx.Unlock()

	// try to delete
	for _, item := range deleteItems {
		item.cleanFunc(ctx, item.key)

	}
}
