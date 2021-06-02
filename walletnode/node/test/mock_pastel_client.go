//go:generate mockery --name=PastelClient

package test

import (
	"context"

	"github.com/pastelnetwork/gonode/pastel"
)

// PastelClient mock pastel.Client. redefine pastel.Client interface
// TODO : go:generate does not work with accros module on current CircleCI build
// not sure how to fix it
type PastelClient interface {
	// MasterNodesTop returns a result of the `masternode top`.
	MasterNodesTop(ctx context.Context) (pastel.MasterNodes, error)

	// MasterNodeStatus returns a result of the `masternode status`.
	MasterNodeStatus(ctx context.Context) (*pastel.MasterNodeStatus, error)

	// MasterNodeConfig returns a result of the `masternode list-conf`.
	MasterNodeConfig(ctx context.Context) (*pastel.MasterNodeConfig, error)

	// StorageFee returns a result of the `storagefee getnetworkfee`.
	StorageFee(ctx context.Context) (*pastel.StorageFee, error)

	// IDTickets returns a result of the `tickets list id`.
	IDTickets(ctx context.Context, idType pastel.IDTicketType) (pastel.IDTickets, error)
}
