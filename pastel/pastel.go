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

	// Sign signs data by the given pastelID and passphrase, if successful returns signature.
	// Command `pastelid sign "text" "PastelID" "passphrase"`.
	Sign(ctx context.Context, data []byte, pastelID, passphrase string) (signature []byte, err error)

	// Verify verifies signed data by the given its signature and pastelID, if successful returns true.
	// Command `pastelid verify "text" "signature" "PastelID"`.
	Verify(ctx context.Context, data []byte, signature, pastelID string) (ok bool, err error)

	// Do an transaction by the given address to sent to and ammount to send, if successful return id of transaction.
	// Command `sendtoaddress  "pastelId" "amount"`.
	SendToAddress(ctx context.Context, pastelId string, amount int64) (txId TxIdType, error error)

	// ActTickets returns activated art tickets.
	// Command `tickets list act`.
	ActTickets(ctx context.Context, actType ActTicketType, minHeight int) (ActTickets, error)

	// ActTickets returns art registration tickets.
	// Command `tickets get <txid>`.
	RegTicket(ctx context.Context, regTxid string) (RegTicket, error)

	// GetBlockCount returns the number of blocks in the best valid block chain
	// Command `getblockcount `
	GetBlockCount(ctx context.Context) (int64, error)
}
