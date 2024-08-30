package metamigrator

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/log"
	"time"
)

const (
	batchSize = 5000
)

var (
	staleTime = time.Now().AddDate(0, -3, 0)
)

func (task *MetaMigratorTask) IdentifyMigrationData(ctx context.Context) (err error) {
	var migrationID int

	migrationID, err = task.service.metaStore.GetPendingMigrationID(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving pending migration")
		return fmt.Errorf("failed to get pending migration ID: %w", err)
	}

	if migrationID == 0 {
		log.WithContext(ctx).Info("creating new migration")

		migrationID, err = task.service.metaStore.CreateNewMigration(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error creating new migration")
			return fmt.Errorf("failed to create new migration: %w", err)
		}
	}
	log.WithContext(ctx).WithField("migration_id", migrationID).Info("migration info has been sorted")

	totalCount, err := task.service.metaStore.GetCountOfStaleData(ctx, staleTime)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving stale data count")
		return fmt.Errorf("failed to get count of stale data: %w", err)
	}

	if totalCount == 0 {
		log.WithContext(ctx).Info("no stale data found to migrate")
		return nil
	}
	log.WithContext(ctx).WithField("total_keys", totalCount).Info("total-data that needs to migrate has been identified")

	numOfBatches := getNumOfBatches(totalCount)
	log.WithContext(ctx).WithField("no_of_batches", numOfBatches).Info("batches required to store migration-meta has been calculated")

	for batchNo := 0; batchNo < numOfBatches; batchNo++ {
		staleData, err := task.service.metaStore.GetStaleDataInBatches(ctx, batchSize, batchNo, staleTime)
		if err != nil {
			log.WithContext(ctx).Error("error retrieving batch of stale data")
			return fmt.Errorf("failed to get stale data in batch %d: %w", batchNo, err)
		}

		if err := task.service.metaStore.InsertMetaMigrationData(ctx, migrationID, staleData); err != nil {
			log.WithContext(ctx).Error("error inserting batch of stale data to migration-meta")
			return fmt.Errorf("failed to insert stale data for migration %d: %w", migrationID, err)
		}

		log.WithContext(ctx).WithField("batch", batchNo).Debug("data added to migration-meta for migration")
	}

	return nil
}

func getNumOfBatches(totalCount int) int {
	numBatches := totalCount / batchSize
	if totalCount%batchSize != 0 {
		numBatches++
	}

	return numBatches
}
