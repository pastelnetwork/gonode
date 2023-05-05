package collection

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
)

func (s *collectionService) finalizeCollections(ctx context.Context) error {
	collections, err := s.store.GetAllInProcessCollections(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error getting all in_process collections")
		return err
	}

	for i := 0; i < len(collections); i++ {
		isFinalized, err := s.finalizeCollection(ctx, collections[i].CollectionTicketTXID)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error finalizing collection")
		} else if isFinalized {
			log.WithContext(ctx).WithField("collection_txid", collections[i].CollectionTicketTXID).Info("collection has been finalized")
		}
	}

	return nil
}

func (s *collectionService) finalizeCollection(ctx context.Context, txid string) (bool, error) {
	finalize, err := s.shouldCollectionBeFinalized(ctx, txid)
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

func (s *collectionService) shouldCollectionBeFinalized(ctx context.Context, txID string) (bool, error) {
	t, err := s.pastelClient.CollectionActTicket(ctx, txID)
	if err != nil {
		return false, fmt.Errorf("error getting collection act ticket: %w", err)
	}

	if t.CollectionActTicketData.IsExpiredByHeight {
		return true, nil
	}

	if t.CollectionActTicketData.IsFull {
		return true, nil
	}

	return false, nil
}
