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
	// Command `sendtoaddress  "pastelID" "amount"`.
	SendToAddress(ctx context.Context, pastelID string, amount int64) (txID TxIDType, error error)

	// ActTickets returns activated art tickets.
	// Command `tickets list act`.
	ActTickets(ctx context.Context, actType ActTicketType, minHeight int) (ActTickets, error)

	// ActTickets returns art registration tickets.
	// Command `tickets get <txid>`.
	RegTicket(ctx context.Context, regTxid string) (RegTicket, error)

	// GetBlockVerbose1 Return block info with verbose is 1
	// Command `getblock height 1`
	GetBlockVerbose1(ctx context.Context, blkHeight int32) (*GetBlockVerbose1Result, error)

	// GetBlockCount returns the number of blocks in the best valid block chain
	// Command `getblockcount `
	GetBlockCount(ctx context.Context) (int32, error)

	// GetBlockHash returns the hash of block
	// Command `getblockhash <blkIndex> `
	GetBlockHash(ctx context.Context, blkIndex int32) (string, error)

	// GetInfo returns the general info of wallet server
	// Command `getinfo `
	GetInfo(ctx context.Context) (*GetInfoResult, error)

	// GetTransaction returns details of transaction
	// Command `gettransaction  <txid>`
	GetTransaction(ctx context.Context, txID TxIDType) (*GetTransactionResult, error)

	// GetNetworkFeePerMB return network storage fee
	// Command `storagefee  getnetworkfee`
	GetNetworkFeePerMB(ctx context.Context) (int64, error)

	// GetArtTicketFeePerKB return network ticket fee
	// Command `storagefee  getartticketfee`
	GetArtTicketFeePerKB(ctx context.Context) (int64, error)
}
