package hermes

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
)

const (
	defaultBlockCountCheckForReorg = 1000
)

func (s *service) IsChainReorgDetected(ctx context.Context) (isDetected bool, ReorgStartingIndex int32, err error) {
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
