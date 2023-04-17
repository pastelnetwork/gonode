package hermes

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	service2 "github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/chainreorg"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/cleaner"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/collection"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/fingerprint"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/pastelblock"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/pasteldstarter"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/pastelnetwork/gonode/pastel"
	"os"
	"reflect"
)

type service struct {
	config       *Config
	pastelClient pastel.Client
	p2p          node.HermesP2PInterface
	sn           node.SNClientInterface
	store        *store.SQLiteStore
}

func (s *service) Run(ctx context.Context) error {
	collectionService, err := collection.NewCollectionService(s.store, s.pastelClient)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to initialize collection service")
	}

	restartPastelDService, err := pasteldstarter.NewRestartPastelDService(s.pastelClient)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to initialize restart pasteld service")
	}

	fingerprintService, err := fingerprint.NewFingerprintService(s.store, s.pastelClient, s.p2p)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to initialize fingerprint service")
	}

	chainReorgService, err := chainreorg.NewChainReorgService(s.store, s.pastelClient)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to initialize chain reorg service")
	}

	cleanerService, err := cleaner.NewCleanupService(s.pastelClient, s.p2p)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to initialize chain reorg service")
	}

	pastelBlockService, err := pastelblock.NewPastelBlocksService(s.store, s.pastelClient)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to initialize restart pastel-d service")
	}

	return runServices(ctx, chainReorgService, cleanerService, collectionService, fingerprintService, pastelBlockService, restartPastelDService)
}

func runServices(ctx context.Context, services ...service2.SvcInterface) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, service := range services {
		service := service

		group.Go(func() error {
			err := service.Run(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("service %s stopped", reflect.TypeOf(service))
			} else {
				log.WithContext(ctx).Warnf("service %s stopped", reflect.TypeOf(service))
			}
			return err
		})
	}

	return group.Wait()
}

// Stats return status of dupe detection
func (s *service) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Get last inserted item
	lastItem, err := s.store.GetLatestFingerprints(ctx)
	if err != nil {
		return nil, errors.Errorf("getLatestFingerprint: %w", err)
	}
	stats["last_insert_time"] = lastItem.DatetimeFingerprintAddedToDatabase

	// Get total of records
	recordCount, err := s.store.GetFingerprintsCount(ctx)
	if err != nil {
		return nil, errors.Errorf("getRecordCount: %w", err)
	}
	stats["record_count"] = recordCount

	fi, err := os.Stat(s.config.DataFile)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get dd db size")
	} else {
		stats["db_size"] = utils.BytesToMB(uint64(fi.Size()))
	}

	return stats, nil
}

// NewService returns a new ddscan service
func NewService(config *Config, pastelClient pastel.Client, sn node.SNClientInterface) (service2.SvcInterface, error) {
	store, err := store.NewSQLiteStore(config.DataFile)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise database: %w", err)
	}

	snAddr := fmt.Sprintf("%s:%d", config.SNHost, config.SNPort)
	log.WithField("sn-addr", snAddr).Info("connecting with SN-Service")

	conn, err := sn.Connect(context.Background(), snAddr)
	if err != nil {
		return nil, errors.Errorf("unable to connect with SN service: %w", err)
	}
	log.Info("connection established with SN-Service")

	p2p := conn.HermesP2P()
	log.Info("connection established with SN-P2P-Service")

	return &service{
		config:       config,
		pastelClient: pastelClient,
		store:        store,
		sn:           sn,
		p2p:          p2p,
	}, nil
}
