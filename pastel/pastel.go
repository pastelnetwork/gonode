//go:generate mockery --name=Client

package pastel

import "context"

// Client represents pastel RPC API client.
type Client interface {
	// MasterNodesTop returns 10 top masternodes for the current or n-th block.
	// Command `masternode top`.
	MasterNodesTop(ctx context.Context) (MasterNodes, error)

	// MasterNodesList returns all masternodes.
	// Command `masternode list full`.
	MasterNodesList(ctx context.Context) (MasterNodes, error)

	// MasterNodesExtra returns all masternodes.
	// Command `masternode list extra`.
	MasterNodesExtra(ctx context.Context) (MasterNodes, error)

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

	// ListAvailableTradeTickets returns list available trade tickets by the given pastelID.
	// Command `tickets list trade available`.
	ListAvailableTradeTickets(ctx context.Context) ([]TradeTicket, error)

	// Sign signs data by the given pastelID and passphrase, if successful returns signature.
	// Command `pastelid sign "text" "PastelID" "passphrase"` "algorithm"
	// Algorithm is ed448[default] or legroast
	Sign(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error)

	// Verify verifies signed data by the given its signature and pastelID, if successful returns true.
	// Command `pastelid verify "text" "signature" "PastelID"`.
	Verify(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error)

	// Do an transaction by the given address to sent to and ammount to send, if successful return id of transaction.
	// input account is default account
	// Command `sendtoaddress  "pastelID" "amount"`.
	SendToAddress(ctx context.Context, pastelID string, amount int64) (txID string, error error)

	// Do an transaction by the given address to sent to and ammount to send, if successful return id of transaction.
	// Command `sendmany  "pastelID" "{...}"`.
	SendFromAddress(ctx context.Context, fromID string, toID string, amount float64) (txID string, error error)

	// ActTickets returns activated NFT tickets.
	// Command `tickets list act`.
	ActTickets(ctx context.Context, actType ActTicketType, minHeight int) (ActTickets, error)

	// RegTicket returns NFT registration tickets.
	// Command `tickets get <txid>`.
	RegTicket(ctx context.Context, regTxid string) (RegTicket, error)

	// ActionRegTicket returns NFT registration tickets.
	// Command `tickets get <txid>`.
	ActionRegTicket(ctx context.Context, regTxid string) (ActionRegTicket, error)

	// RegTickets returns all NFT registration tickets.
	// Command `tickets list art`.
	RegTickets(ctx context.Context) (RegTickets, error)

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
	GetTransaction(ctx context.Context, txID string) (*GetTransactionResult, error)

	// GetRawTransactionVerbose1 returns details of an raw transaction
	// Command `getrawtransaction  <txid> 1`
	GetRawTransactionVerbose1(ctx context.Context, txID string) (*GetRawTransactionVerbose1Result, error)

	// GetNetworkFeePerMB return network storage fee
	// Command `storagefee  getnetworkfee`
	GetNetworkFeePerMB(ctx context.Context) (int64, error)

	// GetNFTTicketFeePerKB return network ticket fee
	// Command `storagefee  getNFTticketfee`
	GetNFTTicketFeePerKB(ctx context.Context) (int64, error)

	// GetRegisterNFTFee return fee of ticket
	// refer https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration - step 12
	// Command `gettotalstoragefee ...`
	GetRegisterNFTFee(ctx context.Context, request GetRegisterNFTFeeRequest) (int64, error)

	// GetActionFee return fee of action
	// Command: `getactionfee getactionfees <ImgSizeInMb>`
	GetActionFee(ctx context.Context, ImgSizeInMb int64) (*GetActionFeesResult, error)

	// RegisterNFTTicket register an NFT ticket
	// Refer https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration - step 18
	// Command `tickets register NFT ...`
	// Return txid of transaction
	RegisterNFTTicket(ctx context.Context, request RegisterNFTRequest) (string, error)

	// RegisterActionTicket register an action ticket
	// Refer :
	// -https://pastel.wiki/en/Architecture/OpenAPI/Sense
	// - https://pastel.wiki/en/Architecture/Components/PastelOpenAPITicketStructures#sense-api-ticket
	// Command `tickets register action ...`
	// Return txid of transaction
	RegisterActionTicket(ctx context.Context, request RegisterActionRequest) (string, error)

	// ActivateActionTicket activate an action ticket
	// Refer: https://pastel.wiki/en/Architecture/Components/PastelOpenAPITicketStructures
	// Command `tickets activate action ...`
	// Return txid of transaction
	ActivateActionTicket(ctx context.Context, request ActivateActionRequest) (string, error)

	// RegisterActTicket activates an registered NFT ticket
	// Command `tickets register act "reg-ticket-tnxid" "artist-height" "fee" "PastelID" "passphrase"`
	// Return txid of NFTActivateTicket
	RegisterActTicket(ctx context.Context, regTicketTxid string, artistHeight int, fee int64, pastelID string, passphrase string) (string, error)

	// FindTicketByID returns the register ticket of pastelid
	// Command `tickets find id <pastelid>`
	FindTicketByID(ctx context.Context, pastelid string) (*IDTicket, error)

	// ActionTickets returns action tickets similar to RegTickets, but for action tickets
	// Command `tickets list action`
	ActionTickets(ctx context.Context) (ActionTicketDatas, error)

	// FindActionActByActionRegTxid returns the action activation ticket by ActionReg ticket txid
	// Command `tickets find action-act <action-reg-ticket-txid>`
	FindActionActByActionRegTxid(ctx context.Context, actionRegTxid string) (*IDTicket, error)

	// GetBalance returns the amount of PSL stored at address
	// Command `z_getbalance address`
	GetBalance(ctx context.Context, address string) (float64, error)

	// BurnAddress ...
	BurnAddress() string
}
