//go:generate mockery --name=Client

package pastel

import "context"

// Client represents pastel RPC API client.
type Client interface {
	// MasterNodesTop returns 10 top masternodes for the current or n-th block.
	// Command `masternode top`.
	MasterNodesTop(ctx context.Context) (MasterNodes, error)

	// MasterNodeStatus returns masternode status information.
	// Command `masternode status`.
	MasterNodeStatus(ctx context.Context) (*MasterNodeStatus, error)

	// MasterNodeConfig returns settings from masternode.conf.
	// Command `masternode list-conf`.
	MasterNodeConfig(ctx context.Context) (*MasterNodeConfig, error)

	// StorageNetworkFee returns network median storage fee.
	// Command `storagefee getnetworkfee`.
	StorageNetworkFee(ctx context.Context) (float64, error)

	// IDTickets returns masternode PastelIDs tickets.
	// Command `tickets list id`.
	IDTickets(ctx context.Context, idType IDTicketType) (IDTickets, error)

	// TicketOwnership returns ownership by the given pastelID and passphrase, if successful returns owership.
	// Command `tickets tools validateownership <txid> <pastelid> <passphrase>`.
	TicketOwnership(ctx context.Context, txID, pastelID, passphrase string) (string, error)

	// GetTicket returns the Art Register ticket by the given transaction ID.
	// Command `tickets get <txid>`.
	GetTicket(ctx context.Context, txID string) (*RegisterTicket, error)

	// ListAvailableTradeTickets returns list available trade tickets by the given pastelID.
	// Command `tickets list trade available <pastelID>`.
	ListAvailableTradeTickets(ctx context.Context, pastelID string) ([]TradeTicket, error)

	// Sign signs data by the given pastelID and passphrase, if successful returns signature.
	// Command `pastelid sign "text" "PastelID" "passphrase"`.
	Sign(ctx context.Context, data []byte, pastelID, passphrase string) (signature []byte, err error)

	// Verify verifies signed data by the given its signature and pastelID, if successful returns true.
	// Command `pastelid verify "text" "signature" "PastelID"`.
	Verify(ctx context.Context, data []byte, signature, pastelID string) (ok bool, err error)
}
