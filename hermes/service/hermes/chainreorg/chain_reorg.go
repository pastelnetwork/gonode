package chainreorg

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/domain"
)

const (
	defaultBlockCountCheckForReorg = 1000
	runTaskInterval                = 60 * time.Minute
)

func (s *chainReorgService) Run(ctx context.Context) error {
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

func (s *chainReorgService) run(ctx context.Context) error {
	log.WithContext(ctx).Info("chain-reorg service run() has been invoked")
	isDetected, lastGoodBlockHeight, err := s.IsChainReorgDetected(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error detecting chain reorg in hermes")
		return nil
	}

	if !isDetected {
		log.WithContext(ctx).Info("chain reorg is not detected, no action needed")
		return nil
	}
	log.WithContext(ctx).WithField("last_good_block_height", lastGoodBlockHeight).
		Info("chain-reorg is detected, going to fix now..")

	if err := s.FixChainReorg(ctx, lastGoodBlockHeight); err != nil {
		log.WithContext(ctx).WithError(err).Error("error fixing the chain reorg in db")
		return nil
	}

	log.WithContext(ctx).Info("chain-reorg has been fixed")
	return nil
}

func (s chainReorgService) Stats(_ context.Context) (map[string]interface{}, error) {
	//chain-reorg stats can be implemented here
	return nil, nil
}

func (s *chainReorgService) IsChainReorgDetected(ctx context.Context) (isDetected bool, ReorgStartingIndex int32, err error) {
	latestBlockCount, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting block count")
		return false, 0, nil
	}

	for blockCount, iterationCount := latestBlockCount, 0; blockCount > 0; blockCount, iterationCount = blockCount-1, iterationCount+1 {
		// to run for last 1000 blocks if no reorg has been detected
		if iterationCount == defaultBlockCountCheckForReorg && !isDetected {
			return false, 0, nil
		}

		blockHash, err := s.pastelClient.GetBlockHash(ctx, blockCount)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error getting block hash")
			continue
		}

		pastelBlock, err := s.store.GetPastelBlockByHash(ctx, blockHash)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error getting block hash")
			continue
		}

		//if the blockHash and height does not match with the one stored in DB,
		//isDetected should update to true, but the loop should not break until it finds the last good block count
		//after reorg detection
		if !pastelBlock.IsValid(blockCount, blockHash) {
			isDetected = true
		}

		//finding the last good block count, in case if the reorg detected
		if pastelBlock.IsValid(blockCount, blockHash) && isDetected && blockCount > 0 {
			return isDetected, blockCount, nil
		}

	}

	return false, 0, nil
}

func (s *chainReorgService) FixChainReorg(ctx context.Context, lastKnownGoodBlockCount int32) error {
	startingBlockToFixChainReorg := lastKnownGoodBlockCount + 1

	latestBlockCount, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting block count")
		return nil
	}

	for blockCount := startingBlockToFixChainReorg; blockCount <= latestBlockCount; blockCount++ {
		blockHash, err := s.pastelClient.GetBlockHash(ctx, blockCount)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error getting block hash through pastel client")
			continue
		}

		pastelBlock, err := s.store.GetPastelBlockByHeight(ctx, blockCount)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error getting block by height")
			continue
		}

		if !pastelBlock.IsValid(blockCount, blockHash) {
			if err := s.store.UpdatePastelBlock(ctx, domain.PastelBlock{
				BlockHeight: blockCount,
				BlockHash:   blockHash,
			}); err != nil {
				log.WithContext(ctx).WithField("block_count", blockCount).WithError(err).
					Error("error updating block hash")
				continue
			}
		}
	}

	return nil
}
