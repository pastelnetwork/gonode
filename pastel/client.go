package pastel

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
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

// Sign implements pastel.Client.Sign
func (client *client) Sign(ctx context.Context, data []byte, pastelID, passphrase string) (signature []byte, err error) {
	var sign struct {
		Signature string `json:"signature"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	if err = client.callFor(ctx, &sign, "pastelid", "sign", text, pastelID, passphrase); err != nil {
		return nil, errors.Errorf("failed to sign data: %w", err)
	}
	return []byte(sign.Signature), nil
}

// Verify implements pastel.Client.Verify
func (client *client) Verify(ctx context.Context, data []byte, signature, pastelID string) (ok bool, err error) {
	var verify struct {
		Verification bool `json:"verification"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	if err = client.callFor(ctx, &verify, "pastelid", "verify", text, signature, pastelID); err != nil {
		return false, errors.Errorf("failed to verify data: %w", err)
	}
	return verify.Verification, nil
}

// StorageFee implements pastel.Client.StorageFee
func (client *client) SendToAddress(ctx context.Context, pastelID string, amount int64) (txID TxIDType, error error) {
	var transactionID struct {
		TransactionID TxIDType `json:"transactionid"`
	}

	if err := client.callFor(ctx, &transactionID, "sendtoaddress", pastelID, fmt.Sprint(amount)); err != nil {
		return "", errors.Errorf("failed to send to address: %w", err)
	}
	return transactionID.TransactionID, nil
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

func (client *client) GetBlockCount(ctx context.Context) (int64, error) {
	res, err := client.CallWithContext(ctx, "getblockcount", "")
	if err != nil {
		return 0, errors.Errorf("failed to call getblockcount: %w", err)
	}

	if res.Error != nil {
		return 0, errors.Errorf("failed to get block count: %w", err)
	}

	return res.GetInt()
}

func (client *client) GetBlockHash(ctx context.Context, blkIndex int32) (string, error) {
	res, err := client.CallWithContext(ctx, "getblockhash", fmt.Sprint(blkIndex))
	if err != nil {
		return "", errors.Errorf("failed to call getblockhash: %w", err)
	}

	if res.Error != nil {
		return "", errors.Errorf("failed to get block hash: %w", err)
	}

	return res.GetString()
}

func (client *client) GetInfo(ctx context.Context) (*GetInfoResult, error) {
	info := &GetInfoResult{}

	if err := client.callFor(ctx, info, "getinfo", ""); err != nil {
		return info, errors.Errorf("failed to get info: %w", err)
	}

	return info, nil
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
