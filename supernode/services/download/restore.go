package download

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"golang.org/x/crypto/sha3"
)

const (
	requiredSymbolPercent = 9
	maxBatchSize          = 2500
)

func (task *NftDownloadingTask) restoreFileFromSymbolIDs(ctx context.Context, rqService rqnode.RaptorQ, symbolIDs []string, rqOti []byte,
	dataHash []byte, txid string) (file []byte, err error) {

	sort.Strings(symbolIDs) // Sort the keys alphabetically because we store the symbols in p2p in a sorted order

	totalSymbols := len(symbolIDs)
	requiredSymbols := totalSymbols * requiredSymbolPercent / 100

	log.WithContext(ctx).WithField("total-symbols", totalSymbols).WithField("required-symbols", requiredSymbols).WithField("txid", txid).Info("Symbols to be retrieved")

	/*
		// Divide symbols into sets
		symbolSets := make([][]string, restoreFileSymbolBatchSize)
		partSize := totalSymbols / restoreFileSymbolBatchSize
		for i := 0; i < restoreFileSymbolBatchSize; i++ {
			start := i * partSize
			end := start + partSize
			if i == (restoreFileSymbolBatchSize - 1) {
				end = totalSymbols
			}
			symbolSets[i] = symbolIDs[start:end]
		}

		// Create a map to store symbols
		symbols := make(map[string][]byte)

		// Iterate through symbol sets
		for _, symbolSet := range symbolSets {
			// If the symbol set has more than 2000 keys, split it into batches of 2000 keys
			if len(symbolSet) > maxBatchSize {
				symbolSetBatches := splitSymbolSet(symbolSet, maxBatchSize)

				log.WithContext(ctx).WithField("len-batches", len(symbolSetBatches)).Info("Symbol set split into batches")
				for i, batch := range symbolSetBatches {
					log.WithContext(ctx).WithField("batch", i).WithField("len-batch", len(batch)).Info("Batch")
				}

				var wg sync.WaitGroup
				batchResults := make(chan map[string][]byte, len(symbolSetBatches))

				// Launch a goroutine for each batch to retrieve symbols concurrently
				for _, batch := range symbolSetBatches {
					wg.Add(1)
					go func(batch []string) {
						defer wg.Done()
						retrievedSymbols, err := task.P2PClient.BatchRetrieve(ctx, batch, len(batch))
						if err != nil {
							log.WithContext(ctx).WithError(err).Error("Failed to retrieve symbols")
							return
						}
						batchResults <- retrievedSymbols
					}(batch)
				}

				// Wait for all goroutines to finish and collect results
				go func() {
					wg.Wait()
					close(batchResults)
				}()

				// Merge symbols from all batches into the symbols map
				for batchResult := range batchResults {
					for key, value := range batchResult {
						symbols[key] = value
					}
					// Check if enough symbols have been retrieved
					if len(symbols) >= requiredSymbols {
						break
					}
				}
				// Check if enough symbols have been retrieved
				if len(symbols) >= requiredSymbols {
					break
				}
			} else {
				// Retrieve symbols for the current set
				retrievedSymbols, err := task.P2PClient.BatchRetrieve(ctx, symbolSet, len(symbolSet))
				if err != nil {
					log.WithContext(ctx).WithError(err).Error("Failed to retrieve symbols")
					return nil, fmt.Errorf("failed to retrieve symbols: %w", err)
				}

				// Merge retrieved symbols into the symbols map
				for key, value := range retrievedSymbols {
					if len(value) > 0 {
						symbols[key] = value
					}
				}

				// Check if enough symbols have been retrieved
				if len(symbols) >= requiredSymbols {
					break
				}
			}
		}*/

	symbols, err := task.P2PClient.BatchRetrieve(ctx, symbolIDs, requiredSymbols)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to retrieve symbols")
		return nil, fmt.Errorf("failed to retrieve symbols: %w", err)
	}

	log.WithContext(ctx).WithField("len-symbols-found", len(symbols)).WithField("req-symbols", requiredSymbols).WithField("txid", txid).Info("Symbols retrieved")
	// Restore Nft
	var decodeInfo *rqnode.Decode
	encodeInfo := rqnode.Encode{
		Symbols: symbols,
		EncoderParam: rqnode.EncoderParameters{
			Oti: rqOti,
		},
	}

	log.WithContext(ctx).WithField("txid", txid).Info("Symbols restored successfully")
	decodeInfo, err = rqService.Decode(ctx, &encodeInfo)
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("Restore file with rqserivce")
		task.UpdateStatus(common.StatusFileDecodingFailed)
		return nil, fmt.Errorf("restore file with rqserivce: %w", err)
	}

	task.UpdateStatus(common.StatusFileDecoded)

	// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
	fileHash := sha3.Sum256(decodeInfo.File)

	if !bytes.Equal(fileHash[:], dataHash) {
		log.WithContext(ctx).Warn("hash file mismatched")
		task.UpdateStatus(common.StatusFileMismatched)
		return nil, errors.New("hash file mismatched")
	}

	return decodeInfo.File, nil
}

// Helper function to split a symbol set into batches of maxBatchSize
func splitSymbolSet(symbolSet []string, batchSize int) [][]string {
	var batches [][]string
	for i := 0; i < len(symbolSet); i += batchSize {
		end := i + batchSize
		if end > len(symbolSet) {
			end = len(symbolSet)
		}
		batches = append(batches, symbolSet[i:end])
	}

	return batches
}
func (task *NftDownloadingTask) getSymbolIDsFromMetadataFile(ctx context.Context, id string, txid string) (symbolIDs []string, err error) {
	var rqIDsData []byte
	log.WithContext(ctx).WithField("id", id).WithField("txid", txid).Debug("Retrieving symbol IDs from metadata file")
	rqIDsData, err = task.P2PClient.Retrieve(ctx, id)
	if err != nil {
		return symbolIDs, fmt.Errorf("retrieve rq metadatafile: %w", err)
	}

	if len(rqIDsData) == 0 {
		return symbolIDs, fmt.Errorf("retrieved rq metadatafile is empty: %w", err)
	}

	symbolIDs, err = task.getRQSymbolIDs(ctx, id, rqIDsData)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("txid", txid).WithField("SymbolIDsFileId", id).Warn("Parse symbol IDs failed")
		task.UpdateStatus(common.StatusSymbolFileInvalid)
		return symbolIDs, fmt.Errorf("parse symbol IDs: %w", err)
	}

	log.WithContext(ctx).WithField("len-symbol-IDs", len(symbolIDs)).WithField("txid", txid).Info("Symbol IDs retrieved")

	return symbolIDs, nil
}

// RestoreFile restores the file using the available rq-ids
func (task *NftDownloadingTask) RestoreFile(ctx context.Context, rqID []string, rqOti []byte, dataHash []byte, txid string) ([]byte, error) {
	var file []byte
	var lastErr error
	var err error

	if len(rqID) == 0 {
		task.UpdateStatus(common.StatusNftRegTicketInvalid)
		return file, errors.Errorf("ticket has empty symbol identifier files")
	}

	var rqConnection rqnode.Connection
	rqConnection, err = task.RqClient.Connect(ctx, task.NftDownloaderService.config.RaptorQServiceAddress)
	if err != nil {
		task.UpdateStatus(common.StatusRQServiceConnectionFailed)
		return file, errors.Errorf("could not connect to rqservice: %w", err)
	}
	defer rqConnection.Close()

	rqNodeConfig := &rqnode.Config{
		RqFilesDir: task.NftDownloaderService.config.RqFilesDir,
	}
	rqService := rqConnection.RaptorQ(rqNodeConfig)

	log.WithContext(ctx).WithField("txid", txid).Info("rq client connected, get symbol IDs from metadata file")
	var symbolIDs []string
	for _, id := range rqID {
		symbolIDs, err = task.getSymbolIDsFromMetadataFile(ctx, id, txid)
		if err == nil && len(symbolIDs) > 0 {
			break
		}

		log.WithContext(ctx).WithError(err).WithField("rq-metadata-file-id", id).Error("Restore file with rqserivce")
		task.UpdateStatus(common.StatusSymbolFileNotFound)
	}

	if len(symbolIDs) == 0 {
		task.UpdateStatus(common.StatusSymbolFileNotFound)
		return file, errors.Errorf("could not retrieve symbol IDs from rq metadata file")
	}

	log.WithContext(ctx).WithField("txid", txid).Info("symbol IDs retrieved, restore file from symbol IDs")
	file, err = task.restoreFileFromSymbolIDs(ctx, rqService, symbolIDs, rqOti, dataHash, txid)
	if err != nil {
		return nil, fmt.Errorf("restore file from symbol IDs: %w", err)
	}

	return file, lastErr
}
