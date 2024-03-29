package chainreorg

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/common"
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
			log.WithContext(ctx).Info("chain-reorg service run() has been invoked")
			if err := s.sync.WaitSynchronization(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("error syncing master-node")
				continue
			}
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

func (s *chainReorgService) checkAndDeleteTerminatedTickets(ctx context.Context) error {
	log.WithContext(ctx).Info("checking and deleting terminated tickets")
	txids, err := s.store.FetchAllTxIDs()
	if err != nil {
		return fmt.Errorf("error fetching all txids: %w", err)
	}

	log.WithContext(ctx).WithField("txids", txids).Info("fetched all txids")
	for txid := range txids {
		ticket, err := s.pastelClient.RegTicket(ctx, txid)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("error getting ticket for txid: %s", txid)
			continue
		}

		if ticket.Height < 0 {
			if err := s.store.UpdateTxIDTimestamp(txid); err != nil {
				log.WithContext(ctx).WithError(err).Errorf("error updating timestamp for txid: %s", txid)
				continue
			}

			log.WithContext(ctx).WithField("txid", txid).Info("ticket is terminated")
		}

		time.Sleep(200 * time.Millisecond) // don't spam the cnode
	}

	return nil
}

func (s *chainReorgService) run(ctx context.Context) error {
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

	if err := s.checkAndDeleteTerminatedTickets(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking and deleting terminated tickets")
		return nil
	}

	log.WithContext(ctx).WithField("last_good_block_height", lastGoodBlockHeight).Info("chain-reorg has been fixed, tickets settled - sending message to fingerprints service")
	s.blockMsgChan <- common.BlockMessage{Block: int(lastGoodBlockHeight)}
	log.WithContext(ctx).WithField("last_good_block_height", lastGoodBlockHeight).Info("chain-reorg message sent to fingerprints service")

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
