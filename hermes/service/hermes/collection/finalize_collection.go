package collection

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
)

func (s *collectionService) finalizeCollections(ctx context.Context) error {
	block, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting block count")
		return err
	}

	collections, err := s.store.GetAllInProcessCollections(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting all in-process collections")
		return err
	}

	for i := 0; i < len(collections); i++ {
		isFinalized, err := s.finalizeCollection(ctx, collections[i].CollectionTicketTXID, block)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error finalizing collection")
		} else if isFinalized {
			log.WithContext(ctx).WithField("collection_txid", collections[i].CollectionTicketTXID).Info("collection has been finalized")
		}
	}

	return nil
}

func (s *collectionService) finalizeCollection(ctx context.Context, txid string, block int32) (bool, error) {
	finalize, err := s.shouldCollectionBeFinalized(ctx, txid, block)
	if err != nil {
		return false, fmt.Errorf("error checking if collection should be finalized: %w", err)
	}

	if !finalize {
		return false, nil
	}

	// finalize collection
	err = s.store.FinalizeCollectionState(ctx, txid)
	if err != nil {
		return false, fmt.Errorf("error finalizing collection state: %w", err)
	}

	return true, nil
}

func (s *collectionService) shouldCollectionBeFinalized(_ context.Context, _ string, _ int32) (bool, error) {
	// implement logic to check if collection should be finalized

	return false, nil
}
