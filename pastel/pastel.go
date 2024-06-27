//go:generate mockery --name=Client

package pastel

import "context"

// Client represents pastel RPC API client.
type Client interface {
	// MasterNodesTop returns 10 top masternodes for the current or n-th block.
	// Command `masternode top`.
	MasterNodesTop(ctx context.Context) (MasterNodes, error)

	// MasterNodesTopN returns 10 top masternodes for n-th block.
	// Command `masternode top N`.
	MasterNodesTopN(ctx context.Context, n int) (MasterNodes, error)

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
	// Command `tickets list nft`.
	RegTickets(ctx context.Context) (RegTickets, error)

	// RegTicketsFromBlockHeight returns all NFT registration tickets higher than height <blockheight> or higher.
	// Command `tickets list nft <height>`.
	RegTicketsFromBlockHeight(ctx context.Context, filter RegTicketsFilter, blockheight uint64) (RegTickets, error)

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
	ActivateNftTicket(ctx context.Context, regTicketTxid string, artistHeight int, fee int64, pastelID string, passphrase string, spendableAddress string) (string, error)

	// FindTicketByID returns the register ticket of pastelid
	// Command `tickets find id <pastelid>`
	FindTicketByID(ctx context.Context, pastelid string) (*IDTicket, error)

	// ActionTickets returns action tickets similar to RegTickets, but for action tickets
	// Command `tickets list action`
	ActionTickets(ctx context.Context) (ActionTicketDatas, error)

	// ActionTicketsFromBlockHeight returns action tickets similar to RegTickets, but for action tickets
	// Command `tickets list action <blockheight>`
	ActionTicketsFromBlockHeight(ctx context.Context, filter RegTicketsFilter, blockheight uint64) (ActionTicketDatas, error)

	// FindActionActByActionRegTxid returns the action activation ticket by ActionReg ticket txid
	// Command `tickets find action-act <action-reg-ticket-txid>`
	FindActionActByActionRegTxid(ctx context.Context, actionRegTxid string) (*IDTicket, error)

	// GetBalance returns the amount of PSL stored at address
	// Command `z_getbalance address`
	GetBalance(ctx context.Context, address string) (float64, error)

	// FindNFTRegTicketsByLabel returns all NFT registration tickets with matching labels.
	// Command `tickets findbylabel nft <label>`.
	FindNFTRegTicketsByLabel(ctx context.Context, label string) (RegTickets, error)

	// FindActionRegTicketsByLabel returns all Action registration tickets with matching labels.
	// Command `tickets findbylabel action <label>`.
	FindActionRegTicketsByLabel(ctx context.Context, label string) (ActionTicketDatas, error)

	// IncrementPoseBanScore increments pose-ban score
	IncrementPoseBanScore(ctx context.Context, txid string, index int) error

	// ActionActivationTicketsFromBlockHeight returns action activation tickets
	// Command `tickets list action-act <block-height>`
	ActionActivationTicketsFromBlockHeight(ctx context.Context, blockheight uint64) (ActTickets, error)

	// CollectionActivationTicketsFromBlockHeight returns collection activation tickets
	// Command `tickets list nft-collection-act <block-height>`
	CollectionActivationTicketsFromBlockHeight(ctx context.Context, blockheight int) (ActTickets, error)

	// CollectionRegTicket returns collection registration ticket.
	// Command `tickets get <txid>`.
	CollectionRegTicket(ctx context.Context, regTxid string) (CollectionRegTicket, error)

	// CollectionActTicket returns collection activation ticket.
	// Command `tickets get <txid>`.
	CollectionActTicket(ctx context.Context, actTxid string) (CollectionActTicket, error)

	// RegisterCollectionTicket registers collection ticket and returns the TxID.
	// Command `tickets register collection "{collection-ticket}" "{signatures}" "pastelid" "passphrase" "label" "fee" ["address"] `.
	RegisterCollectionTicket(ctx context.Context, request RegisterCollectionRequest) (txID string, err error)

	// SignCollectionTicket signs data by the given pastelID and passphrase, if successful returns signature.
	// Command `pastelid sign-base64-encoded "base64-encoded-text" "PastelID" <"passphrase"> ("algorithm")"
	// Algorithm is ed448[default] or legroast
	SignCollectionTicket(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error)

	// VerifyCollectionTicket verifies signed data by the given its signature and pastelID, if successful returns true.
	// Command `pastelid verify-base64-encoded "base64-encoded-text" "PastelID" <"passphrase"> ("algorithm")"
	VerifyCollectionTicket(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error)

	// ActivateCollectionTicket activates the collection ticket, if successful returns act-txid.
	// Command `tickets activate collection "reg-ticket-txid" "called-at-height" "fee" "PastelID" "passphrase" ["address"]`
	ActivateCollectionTicket(ctx context.Context, request ActivateCollectionRequest) (string, error)
	// FindActByRegTxid returns the activation ticket by nReg ticket txid
	// Command `tickets find act <reg-ticket-txid>`
	FindActByRegTxid(ctx context.Context, actionRegTxid string) (*IDTicket, error)

	//GetRawMempool returns the list of in-progress transaction ids
	//Command `getrawmempool false`
	GetRawMempool(ctx context.Context) ([]string, error)

	//GetInactiveActionTickets returns inactive action tickets
	//Command `tickets list act inactive`
	GetInactiveActionTickets(ctx context.Context) (ActTickets, error)

	//GetInactiveNFTTickets returns inactive NFT tickets
	//Command `tickets list nft inactive`
	GetInactiveNFTTickets(ctx context.Context) (RegTickets, error)

	// BurnAddress ...
	BurnAddress() string

	//ZGetTotalBalance returns total balance
	//Command `z_gettotalbalance`
	ZGetTotalBalance(ctx context.Context) (*GetTotalBalanceResponse, error)

	//NFTStorageFee returns the fee of NFT storage
	//Command `tickets tools estimatenftstoragefee <sizeInMB>`
	NFTStorageFee(ctx context.Context, sizeInMB int) (*NFTStorageFeeEstimate, error)

	// RegisterCascadeMultiVolumeTicket registers a cascade multi-volume ticket
	// Command `tickets register contract <<ticket>>, <<sub-type>>, <<hash of the ticket data>>`
	RegisterCascadeMultiVolumeTicket(ctx context.Context, ticket CascadeMultiVolumeTicket) (string, error)

	// GetContractTicket returns contract ticket.
	// Command `tickets get <txid>`.
	GetContractTicket(ctx context.Context, txid string) (Contract, error)
}
