package pastel

import "context"

// Client represents pastel RPC API client.
type Client interface {
	// MasterNodesTop returns a result of the `masternode top`.
	MasterNodesTop(ctx context.Context) (MasterNodes, error)

	// MasterNodeStatus returns a result of the `masternode status`.
	MasterNodeStatus(ctx context.Context) (*MasterNodeStatus, error)

	// MasterNodeConfig returns a result of the `masternode list-conf`.
	MasterNodeConfig(ctx context.Context) (*MasterNodeConfig, error)

	// StorageFee returns a result of the `storagefee getnetworkfee`.
	StorageFee(ctx context.Context) (*StorageFee, error)

	// IDTickets returns a result of the `tickets list id`.
	IDTickets(ctx context.Context, idType IDTicketType) (IDTickets, error)
}
