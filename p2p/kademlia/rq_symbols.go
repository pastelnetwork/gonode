package kademlia

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
)

const (
	defaultSoreSymbolsInterval = 30 * time.Second
	loadSymbolsBatchSize       = 1000
)

func (s *DHT) startStoreSymbolsWorker(ctx context.Context) {
	log.P2P().WithContext(ctx).Info("start delete data worker")

	for {
		select {
		case <-time.After(defaultSoreSymbolsInterval):
			if err := s.storeSymbols(ctx); err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("store symbols")
			}
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing store symbols worker")
			return
		}
	}
}

func (s *DHT) storeSymbols(ctx context.Context) error {
	dirs, err := s.rqstore.GetToDoStoreSymbolDirs()
	if err != nil {
		return fmt.Errorf("get to do store symbol dirs: %w", err)
	}

	for _, dir := range dirs {
		log.WithContext(ctx).WithField("dir", dir).WithField("txid", dir.TXID).Info("rq_symbols worker: start scanning dir & storing raptorQ symbols")
		if err := s.scanDirAndStoreSymbols(ctx, dir.Dir, dir.TXID); err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("scan and store symbols")
		}

		log.WithContext(ctx).WithField("dir", dir).WithField("txid", dir.TXID).Info("rq_symbols worker: scanned dir & stored raptorQ symbols")
	}

	return nil
}

func (s *DHT) scanDirAndStoreSymbols(ctx context.Context, dir, txid string) error {
	keys, err := utils.ReadDirFilenames(dir)
	if err != nil {
		return fmt.Errorf("read dir filenames: %w", err)
	}

	// Iterate over sorted keys in batches
	batchKeys := make(map[string][]byte)
	count := 0

	log.WithContext(ctx).WithField("total-count", len(keys)).WithField("txid", txid).WithField("dir", dir).Info("read dir & found keys")
	for key := range keys {
		batchKeys[key] = nil
		count++
		if count%loadSymbolsBatchSize == 0 {
			if err := s.storeSymbolsInP2P(ctx, dir, batchKeys); err != nil {
				return err
			}
			batchKeys = make(map[string][]byte) // Reset batchKeys after storing
		}
	}

	// Store any remaining symbols in the last batch
	if len(batchKeys) > 0 {
		if err := s.storeSymbolsInP2P(ctx, dir, batchKeys); err != nil {
			return fmt.Errorf("scanDirAndStoreSymbols: store symbols in p2p: %w", err)
		}

		log.WithContext(ctx).WithField("dir", dir).WithField("txid", txid).Info("rq_symbols worker: stored raptorQ symbols in p2p")
	}

	if err := s.rqstore.SetIsCompleted(txid); err != nil {
		return fmt.Errorf("error updating first batch stored flag in rq DB: %w", err)
	}

	return nil
}
func (s *DHT) storeSymbolsInP2P(ctx context.Context, dir string, batchKeys map[string][]byte) error {
	loadedSymbols, err := utils.LoadSymbols(dir, batchKeys)
	if err != nil {
		return fmt.Errorf("p2p worker: load batch symbols from db: %w", err)
	}
	// Prepare batch for P2P storage
	result := make([][]byte, len(loadedSymbols))
	i := 0
	for _, value := range loadedSymbols {
		result[i] = value
		i++
	}

	// Store the loaded symbols in P2P
	if err := s.StoreBatch(ctx, result, 1, dir); err != nil {
		return fmt.Errorf("p2p worker: store batch raptorq symbols in p2p: %w", err)
	}

	if err := utils.DeleteSymbols(ctx, dir, loadedSymbols); err != nil {
		return fmt.Errorf("p2p worker: delete batch symbols from db: %w", err)
	}

	return nil
}
