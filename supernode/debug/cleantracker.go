package debug

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p"
)

type trackItem struct {
	key     string
	expires time.Time
}

// NewCleanTracker creates a tracker
func NewCleanTracker(p2pClient p2p.Client) *CleanTracker {
	return &CleanTracker{
		expireTrackMap: map[time.Time]*trackItem{},
		p2pClient:      p2pClient,
	}
}

// CleanTracker will periodically check to remove expired temporary p2p's (key, value)
type CleanTracker struct {
	expireTrackMap map[time.Time]*trackItem
	p2pClient      p2p.Client
	mtx            sync.Mutex
}

// Track add key into tracking list
func (tracker *CleanTracker) Track(key string, expires time.Time) {
	tracker.mtx.Lock()
	defer tracker.mtx.Unlock()
	tracker.expireTrackMap[time.Now().UTC()] = &trackItem{
		key:     key,
		expires: expires,
	}
}

// CheckExpires will check all items to remove expired items
func (tracker *CleanTracker) CheckExpires(ctx context.Context) {
	deleteItems := []*trackItem{}

	// determines which key need to be delete
	tracker.mtx.Lock()
	now := time.Now().UTC()
	for key, trackItem := range tracker.expireTrackMap {
		if now.After(trackItem.expires) {
			deleteItems = append(deleteItems, trackItem)
			delete(tracker.expireTrackMap, key)
		}
	}
	tracker.mtx.Unlock()

	// try to delete
	for _, item := range deleteItems {
		err := tracker.p2pClient.Delete(ctx, item.key)
		log.WithContext(ctx).WithError(err).WithField("key", item.key).Warn("Remove p2p test key")
	}
}
