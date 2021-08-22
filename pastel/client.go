package pastel

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel/jsonrpc"
)

type client struct {
	jsonrpc.RPCClient
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
	case "ed448", "legroast":
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
	case "ed448", "legroast":
		if err = client.callFor(ctx, &verify, "pastelid", "verify", text, signature, pastelID, algorithm); err != nil {
			return false, errors.Errorf("failed to verify data: %w", err)
		}
	default:
		return false, errors.Errorf("unsupported algoritm %s", algorithm)
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
			return "", errors.Errorf("operationstatus is empty")
		}

		if opstatus[0].Error.Code != 0 {
			return "", errors.Errorf("operation failed code: %d, msg: %s", opstatus[0].Error.Code, opstatus[0].Error.Msg)
		}

		if opstatus[0].Status == "executing" {
			log.WithContext(ctx).Debugf("operation is executing - wait: %d", i)
			time.Sleep(5 * time.Second)
		}
	}

	if opstatus[0].Result.Txid == "" {
		return "", errors.Errorf("empty txid")
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
		return ticket, errors.Errorf("failed to get reg ticket: %w", err)
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
	res, err := client.CallWithContext(ctx, "getblockhash", fmt.Sprint(blkIndex))
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

	// command : tickets tools gettotalstoragefee "ticket" "{signatures}" "pastelid" "passphrase" "key1" "key2" "fee" "imagesize"
	ticket, err := EncodeNFTTicket(request.Ticket)
	if err != nil {
		return 0, errors.Errorf("failed to encode ticket: %w", err)
	}

	ticketBlob := base64.StdEncoding.EncodeToString(ticket)

	signatures, err := EncodeSignatures(*request.Signatures)
	if err != nil {
		return 0, errors.Errorf("failed to encode signatures: %w", err)
	}

	params := []interface{}{}
	params = append(params, "tools")
	params = append(params, "gettotalstoragefee")
	params = append(params, string(ticketBlob))
	params = append(params, string(signatures))
	params = append(params, request.Mn1PastelID)
	params = append(params, request.Passphrase)
	params = append(params, request.Key1)
	params = append(params, request.Key2)
	params = append(params, request.Fee)
	params = append(params, request.ImgSizeInMb)

	if err := client.callFor(ctx, &totalStorageFee, "tickets", params...); err != nil {
		return 0, errors.Errorf("failed to call gettotalstoragefee: %w", err)
	}
	return totalStorageFee.TotalStorageFee, nil
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

	signatures, err := EncodeSignatures(*request.Signatures)
	if err != nil {
		return "", errors.Errorf("failed to encode signatures: %w", err)
	}

	params := []interface{}{}
	params = append(params, "register")
	params = append(params, "nft")
	params = append(params, string(ticketBlob))
	params = append(params, string(signatures))
	params = append(params, request.Mn1PastelID)
	params = append(params, request.Pasphase)
	params = append(params, request.Key1)
	params = append(params, request.Key2)
	params = append(params, fmt.Sprint(request.Fee))

	// command : tickets register NFT "ticket" "{signatures}" "pastelid" "passphrase" "key1" "key2" "fee"
	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register NFT ticket: %w", err)
	}
	return txID.TxID, nil
}

func (client *client) RegisterActTicket(ctx context.Context, regTicketTxid string, artistHeight int, fee int64, pastelID string, passphrase string) (string, error) {
	var txID struct {
		TxID string `json:"txid"`
	}

	params := []interface{}{}
	params = append(params, "register")
	params = append(params, "act")
	params = append(params, regTicketTxid)
	params = append(params, fmt.Sprint(artistHeight))
	params = append(params, fee)
	params = append(params, pastelID)
	params = append(params, passphrase)

	// Command `tickets register act "reg-ticket-tnxid" "artist-height" "fee" "PastelID" "passphrase"`
	if err := client.callFor(ctx, &txID, "tickets", params...); err != nil {
		return "", errors.Errorf("failed to call register act ticket: %w", err)
	}

	return txID.TxID, nil
}

func (client *client) GetBalance(ctx context.Context, address string) (float64, error) {
	var balance float64
	if err := client.callFor(ctx, &balance, "z_getbalance", address); err != nil {
		return 0.0, errors.Errorf("failed to call z_getbalance: %w", err)
	}
	return balance, nil
}

func (client *client) GenBlock(ctx context.Context, amount int) ([]string, error) {
	var blockHashes []string
	if err := client.callFor(ctx, &blockHashes, "generate", amount); err != nil {
		return nil, errors.Errorf("failed to call generate: %w", err)
	}
	return blockHashes, nil
}

func (client *client) FindActTicketByCreatorHeight(ctx context.Context, creatorHeight int32) ([]ActTicket, error) {
	tickets := []ActTicket{}
	if err := client.callFor(ctx, &tickets, "tickets", "find", "act", strconv.FormatInt(int64(creatorHeight), 10)); err != nil {
		return nil, errors.Errorf("failed to call tickets get act %d: %w", creatorHeight, err)
	}
	return tickets, nil
}

func (client *client) callFor(ctx context.Context, object interface{}, method string, params ...interface{}) error {
	return client.CallForWithContext(ctx, object, method, params)
}

// NewClient returns a new Client instance.
func NewClient(config *Config) Client {
	endpoint := net.JoinHostPort(config.hostname(), strconv.Itoa(config.port()))
	if !strings.Contains(endpoint, "//") {
		endpoint = "http://" + endpoint
	}

	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.username()+":"+config.password())),
		},
	}

	return &client{
		RPCClient: jsonrpc.NewClientWithOpts(endpoint, opts),
	}
}
