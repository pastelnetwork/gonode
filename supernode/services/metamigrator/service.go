package metamigrator

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/sqlite"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix                         = "MetaMigrator"
	defaultMetaMigratorDataIdentifier = 120 * time.Minute
	lowSpaceThresholdGB               = 50 // in GB

)

// MetaMigratorService represents the MetaMigrator service.
type MetaMigratorService struct {
	*common.SuperNodeService
	metaStore *sqlite.MigrationMetaStore
}

// Run starts the MetaMigrator service task
func (service *MetaMigratorService) Run(ctx context.Context) error {
	for {
		select {
		case <-time.After(defaultMetaMigratorDataIdentifier):
			log.WithContext(ctx).Info("meta-migrator-worker run() has been invoked")

			isLow, err := utils.CheckDiskSpace(lowSpaceThresholdGB)
			if err != nil {
				log.WithContext(ctx).WithField("method", "MetaMigratorService").WithError(err).Error("check disk space failed")
				continue
			}

			if !isLow {
				log.WithContext(ctx).Info("disk space is not lower than threshold, not proceeding further")
				continue
			}

			task := NewMetaMigratorTask(service)
			if err := task.IdentifyMigrationData(ctx); err != nil {
				log.WithContext(ctx).WithError(err).Error("failed to identify migration data")
			}
		case <-ctx.Done():
			log.Println("Context done being called in IdentifyMigrationData loop in service.go")
			return nil
		}
	}
}

// NewService returns a new MetaMigratorService instance.
func NewService(metaStore *sqlite.MigrationMetaStore) *MetaMigratorService {
	return &MetaMigratorService{
		SuperNodeService: common.NewSuperNodeService(nil, nil, nil),
		metaStore:        metaStore,
	}
}
