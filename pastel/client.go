package pastel

import (
	"context"
	"encoding/base64"
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

// ActTickets implements pastel.Client.ActTickets
func (client *client) ActTickets(ctx context.Context, actType ActivationTicketType, minHeight int) (ActivationTickets, error) {
	tickets := ActivationTickets{}

	if err := client.callFor(ctx, &tickets, "tickets", "list", "act", actType, minHeight); err != nil {
		return nil, errors.Errorf("failed to get id tickets: %w", err)
	}

	return tickets, nil
}

// GetRegTicket implements pastel.Client.GetRegTicket
func (client *client) GetRegTicket(ctx context.Context, regTxid string) (*Ticket, error) {
	ticket := &Ticket{}

	if err := client.callFor(ctx, &ticket, "tickets", "get", regTxid); err != nil {
		return nil, errors.Errorf("failed to get reg ticket: %w", err)
	}

	return ticket, nil
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
