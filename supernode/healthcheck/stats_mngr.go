package healthcheck

import (
	"context"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

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
	mtx          sync.RWMutex
	clients      map[string]StatsClient
	config       *Config
	currentStats map[string]interface{}
}

// NewStatsMngr return an instance of StatsMngr
func NewStatsMngr(config *Config) *StatsMngr {
	return &StatsMngr{
		clients:      map[string]StatsClient{},
		config:       config,
		currentStats: map[string]interface{}{"time_stamp": ""},
	}
}

// Add adds more client to update stats
func (mngr *StatsMngr) Add(id string, client StatsClient) {
	mngr.mtx.Lock()
	defer mngr.mtx.Unlock()
	mngr.clients[id] = client
}

// Stats returns cached stats of all monitored clients
func (mngr *StatsMngr) Stats(_ context.Context) (map[string]interface{}, error) {
	mngr.mtx.RLock()
	defer mngr.mtx.RUnlock()
	stats := mngr.currentStats
	return stats, nil
}

// updateStats returns stats of all monitored clients
func (mngr *StatsMngr) updateStats(ctx context.Context) (map[string]interface{}, error) {
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

	stats["time_stamp"] = time.Now()
	return stats, nil
}

func (mngr *StatsMngr) run(ctx context.Context) error {
	var err error
	ctx = log.ContextWithPrefix(ctx, logPrefix)
	log.WithContext(ctx).Info("StatsManager started")
	for {
		select {
		case <-ctx.Done():
			err = errors.Errorf("context done %w", ctx.Err())
			log.WithContext(ctx).WithError(err).Warnf("StatsManager stopped")
			return err
		case <-time.After(mngr.config.UpdateInterval):
			stats, subErr := mngr.updateStats(ctx)
			if subErr != nil {
				log.WithContext(ctx).WithError(subErr).Warn("UpdateStatsFailed")
			} else {
				mngr.mtx.Lock()
				mngr.currentStats = stats
				mngr.mtx.Unlock()
			}
		}
	}
}

// Run start update stats of system periodically
func (mngr *StatsMngr) Run(ctx context.Context) error {
	for {
		if err := mngr.run(ctx); err != nil {
			if utils.IsContextErr(err) {
				return err
			}
			log.WithContext(ctx).WithError(err).Error("failed to run stats manager, retrying.")
		} else {
			return nil
		}
	}
}
