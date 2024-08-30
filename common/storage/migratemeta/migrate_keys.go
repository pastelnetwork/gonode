package migratemeta

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/sqlite"
)

func MigrateKeys(ctx context.Context, p2p p2p.P2P, ph *mixins.PastelHandler, metaMigratorStore *sqlite.MigrationMetaStore, txid string) error {
	time.Sleep(5 * time.Minute)

	logger := log.WithContext(ctx).WithField("txid", txid)
	batchSize := 1000

	RQIDs, err := ph.GetRQIDs(ctx, txid)
	if err != nil {
		logger.WithError(err).Error("error retrieving RQIDs")
		return err
	}
	logger.Info("RQIDs retrieved by TxID")

	migrationID, err := metaMigratorStore.CreateNewMigration(ctx)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	logger.Info("migration has been created")

	var keys []string
	for _, RQID := range RQIDs {
		data, err := p2p.Retrieve(ctx, RQID, true)
		if err != nil {
			logger.WithError(err).Error("error retrieving key from p2p")
		}

		if len(data) != 0 {
			decodedRQID := base58.Decode(RQID)
			dbKey := hex.EncodeToString(decodedRQID)
			keys = append(keys, dbKey)
		} else {
			continue
		}

		if len(keys) >= batchSize {
			if err := metaMigratorStore.InsertMetaMigrationData(ctx, migrationID, keys); err != nil {
				logger.WithError(err).Error("error inserting batch keys to meta-migration")
			}

			logger.Info("keys added to meta-migration")
			keys = nil
		}
	}

	if len(keys) > 0 {
		if err := metaMigratorStore.InsertMetaMigrationData(ctx, migrationID, keys); err != nil {
			logger.WithError(err).Error("error inserting batch of stale data to migration-meta")
			return fmt.Errorf("failed to insert stale data for migration %d: %w", migrationID, err)
		}

		logger.Info("keys added to meta-migration")
	}

	if err := metaMigratorStore.ProcessMigrationInBatches(ctx, sqlite.Migration{
		ID: migrationID,
	}); err != nil {
		logger.WithError(err).Error("error processing migration in batches")
	}
	logger.Info("keys have been migrated")

	return nil
}
