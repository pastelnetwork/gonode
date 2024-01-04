package grpc

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type metrics struct {
	conn   *clientConn
	client pb.MetricsClient

	sessID string
}

// SelfHealing gets the self-healing metrics from other nodes
func (service *metrics) SelfHealing(ctx context.Context) (*types.SelfHealingMetrics, error) {
	req := pb.SelfHealingMetricsRequest{}

	rsp, err := service.client.SelfHealing(ctx, &req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("send self-healing metrics request failed")
		return nil, err
	}

	if rsp == nil {
		return nil, nil
	}

	return &types.SelfHealingMetrics{
		SentTicketsForSelfHealing:     int(rsp.SentTicketsForSelfHealing),
		EstimatedMissingKeys:          int(rsp.EstimatedMissingKeys),
		TicketsRequiredSelfHealing:    int(rsp.TicketsRequiredSelfHealing),
		SuccessfullySelfHealedTickets: int(rsp.SuccessfullySelfHealedTickets),
		SuccessfullyVerifiedTickets:   int(rsp.SuccessfullyVerifiedTickets),
	}, nil
}

func newMetricsInterface(conn *clientConn) node.MetricsInterface {
	return &metrics{
		conn:   conn,
		client: pb.NewMetricsClient(conn),
	}
}
