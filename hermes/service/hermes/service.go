package hermes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/store"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer"

	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/errgroup"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	svc "github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	synchronizationIntervalSec = 5
	synchronizationTimeoutSec  = 60
	runTaskInterval            = 2 * time.Minute
	masterNodeSuccessfulStatus = "Masternode successfully started"
)

type service struct {
	config       *Config
	pastelClient pastel.Client
	p2p          node.HermesP2PInterface
	sn           node.SNClientInterface
	store        store.DDStore

	scorer             *scorer.Scorer
	isMasterNodeSynced bool

	currentBlockCount           int32
	restartPastelDExecutionTime time.Time

	checkChainReorgTime time.Time

	// used in cleanup inactive tickets task
	currentNFTBlock    int
	currentActionBlock int

	// used in fingerprints store task
	latestNFTBlockHeight        int
	latestSenseBlockHeight      int
	latestCollectionBlockHeight int
}

func toFloat64Array(data []float32) []float64 {
	ret := make([]float64, len(data))
	for idx, value := range data {
		ret[idx] = float64(value)
	}

	return ret
}

// Run starts task
func (s *service) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			if err := s.run(ctx); err != nil {
				if utils.IsContextErr(err) {
					log.WithContext(ctx).WithError(err).Error("closing hermes due to context err")
					return err
				}

				log.WithContext(ctx).WithError(err).Error("failed to run hermes, retrying.")
			} else {
				return nil
			}
		}
	}
}

func (s *service) run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, "hermes")
	if _, err := os.Stat(s.config.DataFile); os.IsNotExist(err) {
		return errors.Errorf("dataFile dd service not found: %w", err)
	}

	if err := s.waitSynchronization(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to initial wait synchronization")
	} else {
		s.isMasterNodeSynced = true
	}

	snAddr := fmt.Sprintf("%s:%d", s.config.SNHost, s.config.SNPort)
	log.WithContext(ctx).WithField("sn-addr", snAddr).Info("connecting with SN-Service")

	conn, err := s.sn.Connect(ctx, snAddr)
	if err != nil {
		return errors.Errorf("unable to connect with SN service: %w", err)
	}
	s.p2p = conn.HermesP2P()
	log.WithContext(ctx).Info("connection established with SN-Service")

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.runStoreFingerprintsTask(gctx)
	})

	group.Go(func() error {
		return s.runStoreCollectionsTask(gctx)
	})

	group.Go(func() error {
		return s.processPastelBlock(ctx)
	})

	group.Go(func() error {
		return s.CleanupInactiveTickets(gctx)
	})

	group.Go(func() error {
		return s.runChainReorgTask(gctx)
	})

	/*group.Go(func() error {
		return s.scorer.Start(gctx)
	})*/

	group.Go(func() error {
		return s.restartPastelID(gctx)
	})

	if err := group.Wait(); err != nil {
		log.WithContext(gctx).WithError(err).Errorf("First runTask() failed")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(runTaskInterval):
			// Check if node is synchronized or not
			if !s.isMasterNodeSynced {
				if err := s.checkSynchronized(ctx); err != nil {
					log.WithContext(ctx).WithError(err).Debug("Failed to check synced status from master node")
					continue
				}

				log.WithContext(ctx).Debug("Done for waiting synchronization status")
				s.isMasterNodeSynced = true
			}

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				return s.runStoreFingerprintsTask(gctx)
			})

			group.Go(func() error {
				return s.processPastelBlock(ctx)
			})

			group.Go(func() error {
				return s.CleanupInactiveTickets(gctx)
			})

			group.Go(func() error {
				return s.restartPastelID(gctx)
			})

			group.Go(func() error {
				return s.runChainReorgTask(gctx)
			})

			/*group.Go(func() error {
				return s.scorer.Start(gctx)
			})*/

			if err := group.Wait(); err != nil {
				log.WithContext(gctx).WithError(err).Errorf("run task failed")
			}
		}
	}
}

func (s *service) checkSynchronized(ctx context.Context) error {
	st, err := s.pastelClient.MasterNodeStatus(ctx)
	if err != nil {
		return errors.Errorf("getMasterNodeStatus: %w", err)
	}

	if st == nil {
		return errors.New("empty status")
	}

	if st.Status == masterNodeSuccessfulStatus {
		return nil
	}

	return errors.Errorf("node not synced, status is %s", st.Status)
}

func (s *service) waitSynchronization(ctx context.Context) error {
	checkTimeout := func(checked chan<- struct{}) {
		time.Sleep(synchronizationTimeoutSec * time.Second)
		close(checked)
	}

	timeoutCh := make(chan struct{})
	go checkTimeout(timeoutCh)

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(synchronizationIntervalSec * time.Second):
			err := s.checkSynchronized(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Debug("Failed to check synced status from master node")
			} else {
				log.WithContext(ctx).Info("Done for waiting synchronization status")
				return nil
			}
		case <-timeoutCh:
			return errors.New("timeout expired")
		}
	}
}

// Utility function to get dd and fp file from an id hash, where the file should be stored
func (s *service) tryToGetFingerprintFileFromHash(ctx context.Context, hash string) (*pastel.DDAndFingerprints, error) {
	rawFile, err := s.p2p.Retrieve(ctx, hash)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("retrieve err")
		return nil, errors.Errorf("Error finding dd and fp file: %w", err)
	}

	decData, err := zstd.Decompress(nil, rawFile)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("decompress err")
		return nil, errors.Errorf("decompress: %w", err)
	}

	splits := bytes.Split(decData, []byte{pastel.SeparatorByte})
	if (len(splits)) < 2 {
		log.WithContext(ctx).WithError(err).Error("incorrecrt split err")
		return nil, errors.Errorf("error separating file by separator bytes, separator not found")
	}
	file, err := utils.B64Decode(splits[0])
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("b64 decode err")
		return nil, errors.Errorf("decode file: %w", err)
	}

	ddFingerprint := &pastel.DDAndFingerprints{}
	if err := json.Unmarshal(file, ddFingerprint); err != nil {
		log.WithContext(ctx).WithError(err).Error("unmarshal err")
		return nil, errors.Errorf("unmarshal json: %w", err)
	}

	return ddFingerprint, nil
}

func (s *service) runStoreFingerprintsTask(ctx context.Context) error {
	log.WithContext(ctx).Info("getting Activation tickets, checking non seed records.")
	nonseed, err := s.store.CheckNonSeedRecord(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get nonseed record")
	} else if !nonseed {
		log.WithContext(ctx).Info("No NonSeed Record, set latestBlockHeight to 0")
		s.latestNFTBlockHeight = 0
		s.latestSenseBlockHeight = 0
	}

	if err := s.parseSenseTickets(ctx); err != nil {
		return err
	}

	return s.parseNFTTickets(ctx)
}

func (s *service) runStoreCollectionsTask(ctx context.Context) error {
	return s.parseCollectionTickets(ctx)
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
func NewService(config *Config, pastelClient pastel.Client, sn node.SNClientInterface) (svc.SvcInterface, error) {
	store, err := store.NewSQLiteStore(config.DataFile)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise database: %w", err)
	}

	scorerConfig := &scorer.Config{
		NumberSuperNodes:          config.NumberSuperNodes,
		ConnectToNodeTimeout:      config.ConnectToNodeTimeout,
		ConnectToNextNodeDelay:    config.ConnectToNextNodeDelay,
		AcceptNodesTimeout:        config.AcceptNodesTimeout,
		CreatorPastelID:           config.CreatorPastelID,
		CreatorPastelIDPassphrase: config.CreatorPastelIDPassphrase,
	}

	return &service{
		config:                      config,
		pastelClient:                pastelClient,
		store:                       store,
		sn:                          sn,
		currentNFTBlock:             1,
		currentActionBlock:          1,
		scorer:                      scorer.New(scorerConfig, pastelClient, sn, store),
		restartPastelDExecutionTime: time.Now(),
		checkChainReorgTime:         time.Now(),
	}, nil
}

func typeMapper(val string) string {
	if val == "action-reg" {
		return "SENSE"
	}
	if val == "nft-reg" {
		return "NFT"
	}

	return val
}
