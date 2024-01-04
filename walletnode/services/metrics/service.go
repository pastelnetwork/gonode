package metrics

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// MetricsService represents a service for Self-Healing & Storage Challenge Metrics
type MetricsService struct {
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface
	historyDB     storage.LocalStoreInterface
}

// GetSelfHealingMetrics sum up the self-healing metrics after retrieving them from every node
func (ms *MetricsService) GetSelfHealingMetrics(ctx context.Context, pid, passphrase string) (types.SelfHealingMetrics, error) {
	listOfMasternodes, err := ms.pastelHandler.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving list of masternodes extra")
	}

	if len(listOfMasternodes) == 0 {
		err := errors.New("no masternode found in masternode list extra")
		log.WithContext(ctx).WithError(err)
		return types.SelfHealingMetrics{}, err
	}

	secInfo := &alts.SecInfo{
		PastelID:   pid,
		PassPhrase: passphrase,
		Algorithm:  pastel.SignAlgorithmED448,
	}

	var (
		metrics types.SelfHealingMetrics
		wg      sync.WaitGroup
		mu      sync.Mutex
	)
	sem := make(chan struct{}, 5) // Limit set to 10

	for _, node := range listOfMasternodes {
		node := node

		if node.ExtAddress == "" || node.ExtKey == "" {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			// Acquire a token from the semaphore channel
			sem <- struct{}{}
			defer func() { <-sem }() // Release the token back into the channel

			res, err := ms.getSelfHealingMetricsFromSN(ctx, secInfo, node.ExtAddress)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("error retrieving metrics")
				return
			}

			if res == nil {
				return
			}
			log.WithContext(ctx).
				WithField("metrics", res).
				WithField("processing_sn_address", node.ExtAddress).
				Info("metrics response received")

			mu.Lock()
			metrics.SentTicketsForSelfHealing = metrics.SentTicketsForSelfHealing + res.SentTicketsForSelfHealing
			metrics.EstimatedMissingKeys = metrics.EstimatedMissingKeys + res.EstimatedMissingKeys
			metrics.TicketsRequiredSelfHealing = metrics.TicketsRequiredSelfHealing + res.TicketsRequiredSelfHealing
			metrics.SuccessfullySelfHealedTickets = metrics.SuccessfullySelfHealedTickets + res.SuccessfullySelfHealedTickets
			metrics.SuccessfullyVerifiedTickets = metrics.SuccessfullyVerifiedTickets + res.SuccessfullyVerifiedTickets
			mu.Unlock()
		}()
	}
	wg.Wait()

	return metrics, nil
}

func (ms *MetricsService) getSelfHealingMetricsFromSN(ctx context.Context, secInfo *alts.SecInfo, processingSupernodeAddr string) (metrics *types.SelfHealingMetrics, err error) {
	logger := log.WithContext(ctx).WithField("processing_sn_address", processingSupernodeAddr)

	logger.Info("sending get self-healing-metric request to sn")

	// Create a context with timeout
	connectCtx, cancel := context.WithTimeout(ctx, time.Second*10) // 10 seconds timeout
	defer cancel()
	//Connect over grpc
	nodeClientConn, err := ms.nodeClient.Connect(connectCtx, processingSupernodeAddr, secInfo, "")
	if err != nil {
		err = fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
		logger.Warn(err.Error())
		return nil, err
	}
	defer func() {
		if nodeClientConn != nil {
			nodeClientConn.Close()
		}
	}()

	metricsIF := nodeClientConn.MetricsInterface()

	res, err := metricsIF.SelfHealing(ctx)
	if err != nil {
		logger.WithError(err).Error("error getting metrics")
		return nil, err
	}

	return res, nil
}

// NewService returns a new Metrics Service instance
func NewService(pastelClient pastel.Client, nodeClient node.ClientInterface,
	historyDB storage.LocalStoreInterface,
) *MetricsService {
	return &MetricsService{
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		historyDB:     historyDB,
	}
}
