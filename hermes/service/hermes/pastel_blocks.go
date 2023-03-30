package hermes

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
)

// processPastelBlock stores the latest block hash and height to DB if not stored already
func (s *service) processPastelBlock(ctx context.Context) error {
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
