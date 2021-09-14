package healthcheck

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
)

// PastelStatsClient is definitation of stats manager
type PastelStatsClient struct {
	pastelClient pastel.Client
}

// NewPastelStatsClient return an instance of PastelStatsClient
func NewPastelStatsClient(pastelClient pastel.Client) *PastelStatsClient {
	return &PastelStatsClient{
		pastelClient: pastelClient,
	}
}

// Stats returns stats of all monitored clients
func (client *PastelStatsClient) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	// Get master top nodes
	topNodes, err := client.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to MasterNodesTop(): %w", err)
	}
	stats["master_top_nodes"] = topNodes

	// Get masternode status
	masternodeStatus, err := client.pastelClient.MasterNodeStatus(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to MasterNodeStatus(): %w", err)
	}
	stats["master_node_status"] = masternodeStatus

	// get current block height
	blockHeight, err := client.pastelClient.GetBlockCount(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to GetBlockCount(): %w", err)
	}
	stats["block_height"] = blockHeight

	return stats, nil
}
