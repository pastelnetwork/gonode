package metrics

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// FetchMetricsService represent metrics service.
type FetchMetricsService struct {
	*common.SuperNodeService
	historyDB storage.LocalStoreInterface
}

// GetSelfHealingMetrics return the self-healing metrics from DB
func (service *FetchMetricsService) GetSelfHealingMetrics(ctx context.Context) (metrics *types.SelfHealingMetrics, err error) {
	res, err := service.historyDB.GetSelfHealingMetrics()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving self-healing metrics from DB")
		return nil, err
	}

	return &res, nil
}

// NewService returns a new Service instance.
func NewService(historyDB storage.LocalStoreInterface) *FetchMetricsService {
	return &FetchMetricsService{
		historyDB: historyDB,
	}
}
