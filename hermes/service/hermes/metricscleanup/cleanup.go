package metricscleanup

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	runTaskInterval = 24 * time.Hour
)

var (
	metricsCleanupThreshold = time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02")
)

// Run cleans up inactive tickets

// Run stores the latest block hash and height to DB if not stored already
func (s *metricsCleanupService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				err := s.run(gctx)
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("error executing metrics cleanup service")
				}

				return nil
			})

			if err := group.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("run task failed")
			}
		}
	}

}

func (s metricsCleanupService) Stats(_ context.Context) (map[string]interface{}, error) {
	//cleanup service stats can be implemented here
	return nil, nil
}

func (s *metricsCleanupService) run(ctx context.Context) error {
	err := s.historyDB.RemoveStorageChallengeStaleData(ctx, metricsCleanupThreshold)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error removing storage challenge stale data")
	}
	log.WithContext(ctx).Info("storage challenge stale data has been removed")

	err = s.historyDB.RemoveSelfHealingStaleData(ctx, metricsCleanupThreshold)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error removing self-healing stale data")
	}
	log.WithContext(ctx).Info("self-healing stale data has been removed")

	return nil
}
