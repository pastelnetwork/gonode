package pastelblock

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/domain"
)

const (
	runTaskInterval = 2 * time.Minute
)

// Run stores the latest block hash and height to DB if not stored already
func (s *pastelBlockService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			// Check if node is synchronized or not
			if !s.sync.GetSyncStatus() {
				if err := s.sync.CheckSynchronized(ctx); err != nil {
					log.WithContext(ctx).WithError(err).Debug("Failed to check synced status from master node")
					continue
				}

				log.WithContext(ctx).Debug("Done for waiting synchronization status")
				s.sync.SetSyncStatus(true)
			}

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				return s.run(gctx)
			})

			if err := group.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("run task failed")
			}
		}
	}
}

func (s *pastelBlockService) Stats(_ context.Context) (map[string]interface{}, error) {
	//pastel-block service stats can be implemented here
	return nil, nil
}

func (s pastelBlockService) run(ctx context.Context) error {
	log.WithContext(ctx).Info("pastel block service run() has been invoked")
	pastelBlock, err := s.store.GetLatestPastelBlock(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving latest pastel block from DB")
		return nil
	}
	log.WithContext(ctx).WithField("db_block_count", pastelBlock.BlockHeight).Info("latest block count retrieved from DB")

	blockCountfromPastel, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting block count from pastel")
		return nil
	}
	log.WithContext(ctx).WithField("pastel_block_count", blockCountfromPastel).Info("pastel block count has been retrieved")

	if pastelBlock.BlockHeight >= blockCountfromPastel {
		log.WithContext(ctx).Info("no block with new height found to be inserted")
		return nil
	}

	for i := pastelBlock.BlockHeight + 1; i <= blockCountfromPastel; i++ {
		blockHash, err := s.pastelClient.GetBlockHash(ctx, i)
		if err != nil {
			log.WithContext(ctx).WithField("block_count", i).WithError(err).Error("error getting block hash from pastel")
			return nil
		}

		if err := s.store.StorePastelBlock(ctx, domain.PastelBlock{
			BlockHeight:        i,
			BlockHash:          blockHash,
			DatetimeBlockAdded: time.Now().Format(time.RFC3339),
		}); err != nil {
			log.WithContext(ctx).WithError(err).Error("error storing pastel block to DB")
			return nil
		}
	}

	return nil
}
