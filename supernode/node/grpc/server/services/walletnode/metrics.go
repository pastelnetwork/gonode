package walletnode

import (
	"context"
	"google.golang.org/grpc"

	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/metrics"
)

// Metrics represents grpc service for metrics.
type Metrics struct {
	pb.UnimplementedMetricsServer
	*common.Metrics
}

// SelfHealing returns the sum of metrics stored in sn history db for self-healing
func (service *Metrics) SelfHealing(ctx context.Context, _ *pb.SelfHealingMetricsRequest) (*pb.SelfHealingMetricsReply, error) {
	res, err := service.Metrics.GetSelfHealingMetrics(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		return nil, err
	}

	return &pb.SelfHealingMetricsReply{
		SentTicketsForSelfHealing:     int64(res.SentTicketsForSelfHealing),
		EstimatedMissingKeys:          int64(res.EstimatedMissingKeys),
		TicketsRequiredSelfHealing:    int64(res.TicketsRequiredSelfHealing),
		SuccessfullySelfHealedTickets: int64(res.SuccessfullySelfHealedTickets),
		SuccessfullyVerifiedTickets:   int64(res.SuccessfullyVerifiedTickets),
	}, nil
}

// Desc returns a description of the service.
func (service *Metrics) Desc() *grpc.ServiceDesc {
	return &pb.Metrics_ServiceDesc
}

// NewMetricsService returns a new MetricService instance.
func NewMetricsService(service *metrics.FetchMetricsService) *Metrics {
	return &Metrics{
		Metrics: common.NewMetricsService(service),
	}
}
