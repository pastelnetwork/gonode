package pastel

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel/jsonrpc"
)

const (
	// SignAlgorithmED448 is ED448 signature algorithm
	SignAlgorithmED448 = "ed448"
	// SignAlgorithmLegRoast is Efficient post-quantum signatures algorithm
	SignAlgorithmLegRoast = "legroast"
	// TicketTypeInactive is inactive filter for tickets
	TicketTypeInactive RegTicketsFilter = "inactive"
	// TicketTypeAll is all filter for tickets
	TicketTypeAll RegTicketsFilter = "all"

	// CascadeMultiVolumeMetadata is the only contract type currently available
	CascadeMultiVolumeMetadata ContractType = "cascade_multi_volume_metadata"
)

type ContractType string

// RegTicketsFilter is filter for retrieving action & nft registration tickets
type RegTicketsFilter string

type client struct {
	jsonrpc.RPCClient
	burnAddress string
}

// ActionActivationTicketsFromBlockHeight returns action activation tickets from block height
func (client *client) ActionActivationTicketsFromBlockHeight(ctx context.Context, blockheight uint64) (ActTickets, error) {
	tickets := ActTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "action-act", "all", blockheight); err != nil {
		return nil, errors.Errorf("failed to get action-act tickets: %w", err)
	}

	return tickets, nil
}

// MasterNodeConfig implements pastel.Client.MasterNodeConfig
func (client *client) MasterNodeConfig(ctx context.Context) (*MasterNodeConfig, error) {
	listConf := make(map[string]MasterNodeConfig)

	if err := client.callFor(ctx, &listConf, "masternode", "list-conf"); err != nil {
		return nil, errors.Errorf("failed to get masternode configuration: %w", err)
	}

	if masterNodeConfig, ok := listConf["masternode"]; ok {
		return &masterNodeConfig, nil
	}
	return nil, errors.New("not found masternode configuration")
}

// MasterNodesTop implements pastel.Client.MasterNodesTop
func (client *client) MasterNodesTop(ctx context.Context) (MasterNodes, error) {
	blocknumMNs := make(map[string]MasterNodes)

	if err := client.callFor(ctx, &blocknumMNs, "masternode", "top"); err != nil {
		return nil, errors.Errorf("failed to get top masternodes: %w", err)
	}

	for _, masterNodes := range blocknumMNs {
		return masterNodes, nil
	}
	return nil, nil
}

// MasterNodesTopN implements pastel.Client.MasterNodesTop
func (client *client) MasterNodesTopN(ctx context.Context, n int) (MasterNodes, error) {
	blocknumMNs := make(map[string]MasterNodes)

	if err := client.callFor(ctx, &blocknumMNs, "masternode", "top", n); err != nil {
		return nil, errors.Errorf("failed to get top masternodes: %w", err)
	}

	for _, masterNodes := range blocknumMNs {
		return masterNodes, nil
	}

	return nil, nil
}

// MasterNodesList implements pastel.Client.MasterNodesList
func (client *client) MasterNodesList(ctx context.Context) (MasterNodes, error) {
	blocknumMNs := make(map[string]MasterNodes)
	if err := client.callFor(ctx, &blocknumMNs, "masternode", "list", "full"); err != nil {
		return nil, errors.Errorf("failed to get list of masternodes: %w", err)
	}
	for _, masterNodes := range blocknumMNs {
		return masterNodes, nil
	}
	return nil, nil
}

// MasterNodesTop implements pastel.Client.MasterNodeStatus
func (client *client) MasterNodeStatus(ctx context.Context) (*MasterNodeStatus, error) {
	var status MasterNodeStatus

	if err := client.callFor(ctx, &status, "masternode", "status"); err != nil {
		return nil, errors.Errorf("failed to get masternode status: %w", err)
	}
	return &status, nil
}

// StorageFee implements pastel.Client.StorageFee
func (client *client) StorageNetworkFee(ctx context.Context) (networkfee float64, err error) {
	var storagefee struct {
		NetworkFee float64 `json:"networkfee"`
	}

	if err := client.callFor(ctx, &storagefee, "storagefee", "getnetworkfee"); err != nil {
		return 0, errors.Errorf("failed to get storage fee: %w", err)
	}
	return storagefee.NetworkFee, nil
}

// IDTickets implements pastel.Client.IDTickets
func (client *client) IDTickets(ctx context.Context, idType IDTicketType) (IDTickets, error) {
	tickets := IDTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "id", string(idType)); err != nil {
		return nil, errors.Errorf("failed to get id tickets: %w", err)
	}
	return tickets, nil
}

func (client *client) FindTicketByID(ctx context.Context, pastelID string) (*IDTicket, error) {
	ticket := IDTicket{}

	if err := client.callFor(ctx, &ticket, "tickets", "find", "id", pastelID); err != nil {
		return nil, errors.Errorf("failed to get id tickets: %w", err)
	}
	return &ticket, nil
}

func (client *client) IncrementPoseBanScore(ctx context.Context, txid string, index int) error {
	res := make(map[string]interface{})
	if err := client.callFor(ctx, &res, "masternode", "pose-ban-score", "increment", txid, index); err != nil {
		return errors.Errorf("failed to increment pose-ban-score: %w", err)
	}

	return nil
}

// TicketOwnership implements pastel.Client.TicketOwnership
func (client *client) TicketOwnership(ctx context.Context, txID, pastelID, passphrase string) (string, error) {
	var ownership struct {
		NFT   string `json:"NFT"`   // txid from the request
		Trade string `json:"trade"` // txid from trade ticket
	}

	if err := client.callFor(ctx, &ownership, "tickets", "tools", "validateownership", txID, pastelID, passphrase); err != nil {
		return "", errors.Errorf("failed to get ticket ownership: %w", err)
	}
	return ownership.Trade, nil
}

// ListAvailableTradeTickets implements pastel.Client.ListAvailableTradeTickets
func (client *client) ListAvailableTradeTickets(ctx context.Context) ([]TradeTicket, error) {
	tradeTicket := []TradeTicket{}
	if err := client.callFor(ctx, &tradeTicket, "tickets", "list", "trade", "available"); err != nil {
		return nil, errors.Errorf("failed to get available trade tickets: %w", err)
	}
	return tradeTicket, nil
}

// Sign implements pastel.Client.Sign
func (client *client) Sign(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error) {
	var sign struct {
		Signature string `json:"signature"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	switch algorithm {
	case SignAlgorithmED448, SignAlgorithmLegRoast:
		if err = client.callFor(ctx, &sign, "pastelid", "sign", text, pastelID, passphrase, algorithm); err != nil {
			return nil, errors.Errorf("failed to sign data: %w", err)
		}
	default:
		return nil, errors.Errorf("unsupported algorithm %s", algorithm)
	}
	return []byte(sign.Signature), nil
}

// Verify implements pastel.Client.Verify
func (client *client) Verify(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error) {
	var verify struct {
		Verification string `json:"verification"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	switch algorithm {
	case SignAlgorithmED448, SignAlgorithmLegRoast:
		if err = client.callFor(ctx, &verify, "pastelid", "verify", text, signature, pastelID, algorithm); err != nil {
			return false, errors.Errorf("failed to verify data: %w", err)
		}
	default:
		return false, errors.Errorf("unsupported algorithm %s", algorithm)
	}

	return verify.Verification == "OK", nil
}

// StorageFee implements pastel.Client.StorageFee
func (client *client) SendToAddress(ctx context.Context, burnAddress string, amount int64) (txID string, error error) {
	res, err := client.CallWithContext(ctx, "sendtoaddress", burnAddress, fmt.Sprint(amount))
	if err != nil {
		return "", errors.Errorf("failed to call sendtoaddress: %w", err)
	}

	if res.Error != nil {
		return "", errors.Errorf("failed to send to address %s: %w", burnAddress, res.Error)
	}

	return res.GetString()
}

func (client *client) SendFromAddress(ctx context.Context, fromAddr string, toAddr string, amount float64) (txID string, error error) {
	amounts := []Amount{{toAddr, amount}}

	res, err := client.CallWithContext(ctx, "z_sendmanywithchangetosender", fromAddr, amounts)
	if err != nil {
		return "", errors.Errorf("failed to call z_sendmany: %w", err)
	}

	if res.Error != nil {
		return "", errors.Errorf("failed to sendmany: %w", res.Error)
	}

	opid, err := res.GetString()
	if err != nil {
		return "", errors.Errorf("failed to get operationid: %w", err)
	}

	opstatus := []GetOperationStatusResult{}
	for i := 0; i < 10; i++ {
		if err := client.callFor(ctx, &opstatus, "z_getoperationstatus", []string{opid}); err != nil {
			return "", errors.Errorf("failed to call z_getoperationstatus: %w", err)
		}

		if len(opstatus) == 0 {
			return "", errors.New("operationstatus is empty")
		}

		if opstatus[0].Error.Code != 0 {
			return "", errors.Errorf("operation failed code: %d, msg: %s", opstatus[0].Error.Code, opstatus[0].Error.Msg)
		}

		if opstatus[0].Status == "executing" {
			log.WithContext(ctx).Debugf("operation z_getoperationstatus() is executing - wait: %d", i)
			time.Sleep(5 * time.Second)
		}
	}

	if opstatus[0].Result.Txid == "" {
		return "", errors.New("empty txid")
	}

	return opstatus[0].Result.Txid, nil
}

// ActTickets implements pastel.Client.ActTickets
func (client *client) ActTickets(ctx context.Context, actType ActTicketType, minHeight int) (ActTickets, error) {
	tickets := ActTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "act", actType, minHeight); err != nil {
		return nil, errors.Errorf("failed to get act tickets: %w", err)
	}

	return tickets, nil
}

// RegTicket implements pastel.Client.RegTicket
func (client *client) RegTicket(ctx context.Context, regTxid string) (RegTicket, error) {
	ticket := RegTicket{}

	if err := client.callFor(ctx, &ticket, "tickets", "get", regTxid); err != nil {
		return ticket, errors.Errorf("failed to get reg ticket %s: %w", regTxid, err)
	}

	return ticket, nil
}

func (client *client) GetContractTicket(ctx context.Context, txid string) (Contract, error) {
	ticket := Contract{}

	if err := client.callFor(ctx, &ticket, "tickets", "get", txid); err != nil {
		return ticket, errors.Errorf("failed to get contract ticket %s: %w", txid, err)
	}

	return ticket, nil
}

// ActionRegTicket implements pastel.Client.RegTicket
func (client *client) ActionRegTicket(ctx context.Context, regTxid string) (ActionRegTicket, error) {
	ticket := ActionRegTicket{}

	if err := client.callFor(ctx, &ticket, "tickets", "get", regTxid); err != nil {
		return ticket, errors.Errorf("failed to get reg ticket %s: %w", regTxid, err)
	}

	return ticket, nil
}

// RegTickets implements pastel.Client.RegTickets
func (client *client) RegTickets(ctx context.Context) (RegTickets, error) {
	tickets := RegTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "nft"); err != nil {
		return nil, errors.Errorf("failed to get registration tickets: %w", err)
	}

	return tickets, nil
}

// RegTicketsFromBlockHeight implements pastel.Client.RegTicketsFromBlockHeight
func (client *client) RegTicketsFromBlockHeight(ctx context.Context, filter RegTicketsFilter, blockheight uint64) (RegTickets, error) {
	tickets := RegTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "nft", string(filter), blockheight); err != nil {
		return nil, errors.Errorf("failed to get registration tickets with block height: %w", err)
	}

	return tickets, nil
}

func (client *client) GetBlockVerbose1(ctx context.Context, blkHeight int32) (*GetBlockVerbose1Result, error) {
	result := &GetBlockVerbose1Result{}

	if err := client.callFor(ctx, result, "getblock", fmt.Sprint(blkHeight), 1); err != nil {
		return result, errors.Errorf("failed to get block: %w", err)
	}

	return result, nil
}

func (client *client) GetBlockCount(ctx context.Context) (int32, error) {
	res, err := client.CallWithContext(ctx, "getblockcount")
	if err != nil {
		return 0, errors.Errorf("failed to call getblockcount: %w", err)
	}

	if res.Error != nil {
		return 0, errors.Errorf("failed to get block count: %w", res.Error)
	}

	cnt, err := res.GetInt()

	return int32(cnt), err
}

func (client *client) GetBlockHash(ctx context.Context, blkIndex int32) (string, error) {
	res, err := client.CallWithContext(ctx, "getblockhash", blkIndex)
	if err != nil {
		return "", errors.Errorf("failed to call getblockhash: %w", err)
	}

	if res.Error != nil {
		return "", errors.Errorf("failed to get block hash: %w", res.Error)
	}

	return res.GetString()
}

func (client *client) GetInfo(ctx context.Context) (*GetInfoResult, error) {
	result := &GetInfoResult{}

	if err := client.callFor(ctx, result, "getinfo", ""); err != nil {
		return result, errors.Errorf("failed to get info: %w", err)
	}

	return result, nil
}

func (client *client) GetTransaction(ctx context.Context, txID string) (*GetTransactionResult, error) {
	result := &GetTransactionResult{}

	if err := client.callFor(ctx, result, "gettransaction", txID); err != nil {
		return result, errors.Errorf("failed to get transaction: %w", err)
	}

	return result, nil
}

func (client *client) GetRawTransactionVerbose1(ctx context.Context, txID string) (*GetRawTransactionVerbose1Result, error) {
	result := &GetRawTransactionVerbose1Result{}

	if err := client.callFor(ctx, result, "getrawtransaction", txID, 1); err != nil {
		return result, errors.Errorf("failed to get transaction: %w", err)
	}

	return result, nil
}

func (client *client) GetNetworkFeePerMB(ctx context.Context) (int64, error) {
	var networkFee struct {
		NetworkFee int64 `json:"networkfee"`
	}

	if err := client.callFor(ctx, &networkFee, "storagefee", "getnetworkfee"); err != nil {
		return 0, errors.Errorf("failed to call storagefee: %w", err)
	}
	return networkFee.NetworkFee, nil
}

func (client *client) GetNFTTicketFeePerKB(ctx context.Context) (int64, error) {
	var NFTticketFee struct {
		NFTticketFee int64 `json:"nftticketfee"`
	}

	if err := client.callFor(ctx, &NFTticketFee, "storagefee", "getnftticketfee"); err != nil {
		return 0, errors.Errorf("failed to call storagefee: %w", err)
	}
	return NFTticketFee.NFTticketFee, nil
}

func (client *client) GetRegisterNFTFee(ctx context.Context, request GetRegisterNFTFeeRequest) (int64, error) {
	var totalStorageFee struct {
		TotalStorageFee int64 `json:"totalstoragefee"`
	}

	// command : tickets tools gettotalstoragefee "ticket" "{signatures}" "pastelid" "passphrase" "label" "fee" "imagesize"
	ticket, err := EncodeNFTTicket(request.Ticket)
	if err != nil {
		return 0, errors.Errorf("failed to encode ticket: %w", err)
	}

	ticketBlob := base64.StdEncoding.EncodeToString(ticket)

	signatures, err := EncodeRegSignatures(*request.Signatures)
	if err != nil {
		return 0, errors.Errorf("failed to encode signatures: %w", err)
	}

	unitFee, err := client.GetNetworkFeePerMB(ctx)
	if err != nil {
		return 0, fmt.Errorf("get network fee failure: %w", err)
	}

	params := []interface{}{}
	params = append(params, "tools")
	params = append(params, "gettotalstoragefee")
	params = append(params, string(ticketBlob))
	params = append(params, string(signatures))
	params = append(params, request.Mn1PastelID)
	params = append(params, request.Passphrase)
	params = append(params, request.Label)
	params = append(params, unitFee*request.ImgSizeInMb)
	params = append(params, request.ImgSizeInMb)

	if err := client.callFor(ctx, &totalStorageFee, "tickets", params...); err != nil {
		return 0, errors.Errorf("failed to call gettotalstoragefee: %w", err)
	}
	return totalStorageFee.TotalStorageFee, nil
}

func (client *client) GetActionFee(ctx context.Context, ImgSizeInMb int64) (*GetActionFeesResult, error) {
	actionFees := &GetActionFeesResult{}

	params := []interface{}{}
	params = append(params, "getactionfees")
	params = append(params, ImgSizeInMb)

	if err := client.callFor(ctx, &actionFees, "storagefee", params...); err != nil {
		return nil, errors.Errorf("failed to call storagefee getactionfees: %w", err)
	}

	return actionFees, nil
}

func (client *client) RegisterNFTTicket(ctx context.Context, request RegisterNFTRequest) (string, error) {
	var txID struct {
		TxID string `json:"txid"`
	}

	ticket, err := EncodeNFTTicket(request.Ticket)
	if err != nil {
		return "", errors.Errorf("failed to encode ticket: %w", err)
	}
	ticketBlob := base64.StdEncoding.EncodeToString(ticket)

	signatures, err := EncodeRegSignatures(*request.Signatures)
	if err != nil {
		return "", errors.Errorf("failed to encode signatures: %w", err)
	}

	params := []interface{}{}
	params = append(params, "register")
	params = append(params, "nft")
	params = append(params, string(ticketBlob))
	params = append(params, string(signatures))
	params = append(params, request.Mn1PastelID)
	params = append(params, request.Passphrase)
	params = append(params, request.Label)
	params = append(params, fmt.Sprint(request.Fee))

	// command : tickets register NFT "ticket" "{signatures}" "pastelid" "passphrase" "label" "fee"
	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register NFT ticket: %w", err)
	}
	return txID.TxID, nil
}

func (client *client) RegisterActionTicket(ctx context.Context, request RegisterActionRequest) (string, error) {
	var txID struct {
		TxID string `json:"txid"`
	}

	ticket, err := EncodeActionTicket(request.Ticket)
	if err != nil {
		return "", errors.Errorf("failed to encode ticket: %w", err)
	}
	ticketBlob := base64.StdEncoding.EncodeToString(ticket)

	signatures, err := EncodeActionSignatures(*request.Signatures)
	if err != nil {
		return "", errors.Errorf("failed to encode signatures: %w", err)
	}

	params := []interface{}{}
	params = append(params, "register")
	params = append(params, "action")
	params = append(params, string(ticketBlob))
	params = append(params, string(signatures))
	params = append(params, request.Mn1PastelID)
	params = append(params, request.Passphrase)
	params = append(params, request.Label)
	params = append(params, fmt.Sprint(request.Fee))
	log.WithContext(ctx).WithField("ticket", ticketBlob).WithField("signatures", string(signatures)).WithField("pastelid", request.Mn1PastelID).WithField("label", request.Label).WithField("fee", request.Fee).Info("RegisterActionTicket Request")
	// command : tickets register action "ticket" "{signatures}" "pastelid" "passphrase" "label" "fee"
	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register action ticket: %w", err)
	}
	return txID.TxID, nil
}

func (client *client) RegisterCascadeMultiVolumeTicket(ctx context.Context, ticket CascadeMultiVolumeTicket) (string, error) {
	return client.registerContract(ctx, ticket, CascadeMultiVolumeMetadata)
}

func (client *client) registerContract(ctx context.Context, data interface{}, contractType ContractType) (string, error) {
	ticketJSON, err := json.Marshal(data)
	if err != nil {
		return "", errors.Errorf("failed to call register action ticket: %w", err)
	}
	ticketBlob := base64.StdEncoding.EncodeToString(ticketJSON)

	hash, _ := utils.Sha3256hash(ticketJSON)
	// Assuming some additional data or parameters are needed, similar to the RegisterNFTRequest example
	params := []interface{}{
		"register",
		"contract",
		ticketBlob,
		contractType,
		hash,
	}

	type ContractResponse struct {
		TxID string `json:"txid"`
		Key  string `json:"key"`
	}
	resp := ContractResponse{}

	if err := client.callFor(ctx, &resp, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register contract: %w", err)
	}
	if resp.TxID == "" {
		return "", errors.Errorf("failed to call register contract (no txid rcvd.): %w", err)
	}

	return resp.TxID, nil
}

func (client *client) ActivateActionTicket(ctx context.Context, request ActivateActionRequest) (string, error) {
	var txID struct {
		TxID string `json:"txid"`
	}

	params := []interface{}{}
	params = append(params, "activate")
	params = append(params, "action")
	params = append(params, request.RegTxID)
	params = append(params, fmt.Sprint(request.BlockNum))
	params = append(params, fmt.Sprint(request.Fee))
	params = append(params, request.PastelID)
	params = append(params, request.Passphrase)
	params = append(params, request.SpendableAddress)

	// command : tickets activate action "txid-of-action-reg-ticket" called_at_height-from_action-reg-ticket fee "PastelID-of-the-caller" "passphrase"
	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call activate action ticket: %w", err)
	}
	return txID.TxID, nil
}

func (client *client) FindActionActByActionRegTxid(ctx context.Context, actionRegTxid string) (*IDTicket, error) {
	ticket := IDTicket{}

	params := []interface{}{}
	params = append(params, "find")
	params = append(params, "action-act")
	params = append(params, actionRegTxid)

	if err := client.callFor(ctx, &ticket, "tickets", params...); err != nil {
		return nil, errors.Errorf("failed to call find action-act <actionRegTxid> : %w", err)
	}

	return &ticket, nil
}

func (client *client) FindActByRegTxid(ctx context.Context, regTxid string) (*IDTicket, error) {
	ticket := IDTicket{}

	params := []interface{}{}
	params = append(params, "find")
	params = append(params, "act")
	params = append(params, regTxid)

	if err := client.callFor(ctx, &ticket, "tickets", params...); err != nil {
		return nil, errors.Errorf("failed to call find action-act <actionRegTxid> : %w", err)
	}

	return &ticket, nil
}

func (client *client) ActivateNftTicket(ctx context.Context, regTicketTxid string, artistHeight int, fee int64,
	pastelID string, passphrase string,
	spendableAddress string) (string, error) {
	var txID struct {
		TxID string `json:"txid"`
	}

	params := []interface{}{}
	params = append(params, "activate")
	params = append(params, "nft")
	params = append(params, regTicketTxid)
	params = append(params, fmt.Sprint(artistHeight))
	params = append(params, fee)
	params = append(params, pastelID)
	params = append(params, passphrase)
	if spendableAddress != "" {
		params = append(params, spendableAddress)
	}

	// Command `tickets register act "reg-ticket-tnxid" "artist-height" "fee" "PastelID" "passphrase"`
	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register act ticket: %w", err)
	}

	return txID.TxID, nil
}

// ActionTickets implements pastel.Client.ActionTickets
func (client *client) ActionTickets(ctx context.Context) (ActionTicketDatas, error) {
	tickets := ActionTicketDatas{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "action"); err != nil {
		return nil, errors.Errorf("failed to get action tickets: %w", err)
	}

	return tickets, nil
}

// ActionTicketsFromBlockHeight implements pastel.Client.ActionTicketsFromBlockHeight
func (client *client) ActionTicketsFromBlockHeight(ctx context.Context, filter RegTicketsFilter, blockheight uint64) (ActionTicketDatas, error) {
	tickets := ActionTicketDatas{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "action", string(filter), blockheight); err != nil {
		return nil, errors.Errorf("failed to get action tickets: %w", err)
	}

	return tickets, nil
}

// CollectionActivationTicketsFromBlockHeight implements pastel.Client.NFTCollectionActivationTicketsFromBlockHeight
func (client *client) CollectionActivationTicketsFromBlockHeight(ctx context.Context, blockheight int) (ActTickets, error) {
	tickets := ActTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "collection-act", "all", blockheight); err != nil {
		return nil, errors.Errorf("failed to get collection-act tickets: %w", err)
	}

	return tickets, nil
}

// CollectionRegTicket implements pastel.Client.CollectionRegTicket
func (client *client) CollectionRegTicket(ctx context.Context, regTxid string) (CollectionRegTicket, error) {
	ticket := CollectionRegTicket{}

	if err := client.callFor(ctx, &ticket, "tickets", "get", regTxid); err != nil {
		return ticket, errors.Errorf("failed to get reg ticket %s: %w", regTxid, err)
	}

	return ticket, nil
}

// CollectionActTicket implements pastel.Client.CollectionActTicket
func (client *client) CollectionActTicket(ctx context.Context, actTxid string) (CollectionActTicket, error) {
	ticket := CollectionActTicket{}

	if err := client.callFor(ctx, &ticket, "tickets", "get", actTxid); err != nil {
		return ticket, errors.Errorf("failed to get reg ticket %s: %w", actTxid, err)
	}

	return ticket, nil
}

// RegisterCollectionTicket registers the collection ticket
func (client *client) RegisterCollectionTicket(ctx context.Context, req RegisterCollectionRequest) (txID string, err error) {
	var res struct {
		TxID string `json:"txid"`
	}

	ticket, err := EncodeCollectionTicket(req.Ticket)
	if err != nil {
		return "", errors.Errorf("failed to encode ticket: %w", err)
	}
	ticketBlob := base64.StdEncoding.EncodeToString(ticket)

	signatures, err := EncodeCollectionTicketSignatures(*req.Signatures)
	if err != nil {
		return "", errors.Errorf("failed to encode signatures: %w", err)
	}

	params := []interface{}{}
	params = append(params, "register")
	params = append(params, "collection")
	params = append(params, ticketBlob)
	params = append(params, string(signatures))
	params = append(params, req.Mn1PastelID)
	params = append(params, req.Passphrase)
	params = append(params, req.Label)
	params = append(params, fmt.Sprint(req.Fee))

	log.WithContext(ctx).WithField("ticket", ticketBlob).WithField("signatures", string(signatures)).WithField("pastelid", req.Mn1PastelID).WithField("label", req.Label).WithField("fee", req.Fee).Info("RegisterCollectionTicket Request")
	// command : tickets register action "ticket" "{signatures}" "pastelid" "passphrase" "label" "fee"
	if err := client.callFor(ctx, &res, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register collection ticket: %w", err)
	}

	return res.TxID, nil
}

func (client *client) ActivateCollectionTicket(ctx context.Context, request ActivateCollectionRequest) (string, error) {
	var txID struct {
		TxID string `json:"txid"`
	}

	params := []interface{}{}
	params = append(params, "activate")
	params = append(params, "collection")
	params = append(params, request.RegTxID)
	params = append(params, fmt.Sprint(request.BlockNum))
	params = append(params, fmt.Sprint(request.Fee))
	params = append(params, request.PastelID)
	params = append(params, request.Passphrase)

	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call activate collection ticket: %w", err)
	}
	return txID.TxID, nil
}

// SignCollectionTicket implements pastel.Client.SignCollectionTicket
func (client *client) SignCollectionTicket(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error) {
	var sign struct {
		Signature string `json:"signature"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	switch algorithm {
	case SignAlgorithmED448, SignAlgorithmLegRoast:
		if err = client.callFor(ctx, &sign, "pastelid", "sign-base64-encoded", text, pastelID, passphrase, algorithm); err != nil {
			return nil, errors.Errorf("failed to sign data: %w", err)
		}
	default:
		return nil, errors.Errorf("unsupported algorithm %s", algorithm)
	}
	return []byte(sign.Signature), nil
}

// VerifyCollectionTicket implements pastel.Client.VerifyCollectionTicket
func (client *client) VerifyCollectionTicket(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error) {
	var verify struct {
		Verification string `json:"verification"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	switch algorithm {
	case SignAlgorithmED448, SignAlgorithmLegRoast:
		if err = client.callFor(ctx, &verify, "pastelid", "verify-base64-encoded", text, signature, pastelID, algorithm); err != nil {
			return false, errors.Errorf("failed to verify data: %w", err)
		}
	default:
		return false, errors.Errorf("unsupported algorithm %s", algorithm)
	}

	return verify.Verification == "OK", nil
}

// FindNFTRegTicketsByLabel returns all NFT registration tickets with matching labels.
// Command `tickets findbylabel nft <label>`.
func (client *client) FindNFTRegTicketsByLabel(ctx context.Context, label string) (RegTickets, error) {
	tickets := RegTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "findbylabel", "nft", label); err != nil {
		return nil, errors.Errorf("failed to get registration tickets with block height: %w", err)
	}

	return tickets, nil
}

// FindActionRegTicketsByLabel returns all Action registration tickets with matching labels.
// Command `tickets findbylabel action <label>`.
func (client *client) FindActionRegTicketsByLabel(ctx context.Context, label string) (ActionTicketDatas, error) {
	tickets := ActionTicketDatas{}

	if err := client.callFor(ctx, &tickets, "tickets", "findbylabel", "action", label); err != nil {
		return nil, errors.Errorf("failed to find action tickets by label: %w", err)
	}

	return tickets, nil
}

func (client *client) GetBalance(ctx context.Context, _ string) (float64, error) {
	var balance float64
	if err := client.callFor(ctx, &balance, "getbalance"); err != nil {
		return 0.0, errors.Errorf("failed to call z_getbalance: %w", err)
	}
	return balance, nil
}

// MasterNodesExtra implements pastel.Client.MasterNodesExtra
func (client *client) MasterNodesExtra(ctx context.Context) (MasterNodes, error) {
	blocknumMNs := make(map[string]MasterNode)
	masterNodes := MasterNodes{}

	if err := client.callFor(ctx, &blocknumMNs, "masternode", "list", "extra"); err != nil {
		return nil, errors.Errorf("failed to get top masternodes: %w", err)
	}

	for _, masterNode := range blocknumMNs {
		masterNodes = append(masterNodes, masterNode)
	}

	return masterNodes, nil
}

// ParsedGetTotalBalanceResponse is the parsed response from the `z_gettotalbalance` command
type ParsedGetTotalBalanceResponse struct {
	Transparent string `json:"transparent"`
	Private     string `json:"private"`
	Total       string `json:"total"`
}

// Parse parses the response into a GetTotalBalanceResponse
func (resp *ParsedGetTotalBalanceResponse) Parse() (*GetTotalBalanceResponse, error) {
	transparent, err := strconv.ParseFloat(resp.Transparent, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing transparent balance: %w", err)
	}

	private, err := strconv.ParseFloat(resp.Private, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing private balance: %w", err)
	}

	total, err := strconv.ParseFloat(resp.Total, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing total balance: %w", err)
	}

	return &GetTotalBalanceResponse{
		Transparent: transparent,
		Private:     private,
		Total:       total,
	}, nil
}

// ZGetTotalBalance returns total balance
// Command `z_gettotalbalance`
func (client *client) ZGetTotalBalance(ctx context.Context) (*GetTotalBalanceResponse, error) {
	resp := &ParsedGetTotalBalanceResponse{}

	if err := client.callFor(ctx, &resp, "z_gettotalbalance"); err != nil {
		return nil, fmt.Errorf("failed to get total balance: %w", err)
	}

	return resp.Parse()
}

// NFTStorageFee returns the fee of NFT storage
// Command `tickets tools estimatenftstoragefee <sizeInMB>`
func (client *client) NFTStorageFee(ctx context.Context, sizeInMB int) (*NFTStorageFeeEstimate, error) {
	resp := &NFTStorageFeeEstimate{}

	if err := client.callFor(ctx, &resp, "tickets", "tools", "estimatenftstoragefee", sizeInMB); err != nil {
		return nil, errors.Errorf("failed to get estimated nft storage: %w", err)
	}

	return resp, nil
}

// GetRawMempool returns the list of in-progress transaction ids
func (client *client) GetRawMempool(ctx context.Context) ([]string, error) {
	var transactionIDs []string
	if err := client.callFor(ctx, &transactionIDs, "getrawmempool", false); err != nil {
		return nil, errors.Errorf("failed to call getrawmempool: %w", err)
	}

	return transactionIDs, nil
}

// GetInactiveActionTickets retrieve the inactive action tickets
func (client *client) GetInactiveActionTickets(ctx context.Context) (ActTickets, error) {
	tickets := ActTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "act", "inactive"); err != nil {
		return nil, errors.Errorf("failed to get act tickets: %w", err)
	}

	return tickets, nil
}

// GetInactiveNFTTickets retrieve the inactive NFT tickets
func (client *client) GetInactiveNFTTickets(ctx context.Context) (RegTickets, error) {
	tickets := RegTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "nft", "inactive"); err != nil {
		return nil, errors.Errorf("failed to get act tickets: %w", err)
	}

	return tickets, nil
}

func (client *client) callFor(ctx context.Context, object interface{}, method string, params ...interface{}) error {
	return client.CallForWithContext(ctx, object, method, params)
}

func (client *client) BurnAddress() string {
	return client.burnAddress
}

// NewClient returns a new Client instance.
// This client interface connects to the Core Pastel (cNode) RPC Server providing access to:
//
//	the blockchain DB, Masternodes DB, Tickets DB, and PastelID DB.
//
// Via testnet, this will connect over 19932, and over mainnet 9932
func NewClient(config *Config, burnAddress string) Client {
	//Configure network addressing
	endpoint := net.JoinHostPort(config.Hostname, strconv.Itoa(config.port()))
	if !strings.Contains(endpoint, "//") {
		endpoint = "http://" + endpoint
	}

	//Parse and configure RPC authorization headers
	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Username+":"+config.Password)),
		},
	}

	//Return a Client interface with the proper RPCClient configurations and burn address
	return &client{
		RPCClient:   jsonrpc.NewClientWithOpts(endpoint, opts),
		burnAddress: burnAddress,
	}
}
