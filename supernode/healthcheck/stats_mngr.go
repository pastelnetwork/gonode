package healthcheck

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	logPrefix = "statsmanager"
)

// StatsClient defines general interface monitored susbsystem should provided for stats manager
type StatsClient interface {
	// Stats return status of p2p
	Stats(ctx context.Context) (map[string]interface{}, error)
}

// StatsMngr is definitation of stats manager
type StatsMngr struct {
	mtx            sync.RWMutex
	updateDuration time.Duration
	clients        map[string]StatsClient
}

// NewStatsMngr return an instance of StatsMngr
func NewStatsMngr(duration time.Duration) *StatsMngr {
	return &StatsMngr{
		clients:        map[string]StatsClient{},
		updateDuration: duration,
	}
}

// Add adds more client to update stats
func (mngr *StatsMngr) Add(id string, client StatsClient) {
	mngr.mtx.Lock()
	defer mngr.mtx.Unlock()
	mngr.clients[id] = client
}

// Stats returns stats of all monitored clients
func (mngr *StatsMngr) Stats(ctx context.Context) (map[string]interface{}, error) {
	mngr.mtx.RLock()
	defer mngr.mtx.RUnlock()
	stats := map[string]interface{}{}
	for id, client := range mngr.clients {
		subStats, err := client.Stats(ctx)
		if err != nil {
			return nil, errors.Errorf("failed to get stats of %s: %w", id, err)
		}
		stats[id] = subStats
	}

	return stats, nil
}

// Run start update stats of system periodically
func (mngr *StatsMngr) Run(ctx context.Context) error {
	var err error
	ctx = log.ContextWithPrefix(ctx, logPrefix)
	log.WithContext(ctx).Info("StatsManager started")
	for {
		select {
		case <-ctx.Done():
			err = errors.Errorf("context done %w", ctx.Err())
			log.WithContext(ctx).WithError(err).Warnf("StatsManager stopped")
			return err
		case <-time.After(mngr.updateDuration):
			stats, subErr := mngr.Stats(ctx)
			if subErr != nil {
				log.WithContext(ctx).WithError(subErr).Warn("UpdateStatsFailed")
			} else {
				// FIXME : update local stats for fetching later
				data, subErr := json.Marshal(stats)
				if subErr != nil {
					log.WithContext(ctx).WithError(subErr).Warn("MarshalStatsFailed")
				} else {
					log.WithContext(ctx).WithField("Stats", string(data)).Warn("UpdateStatsFinished")
				}
			}
		}
	}
}
