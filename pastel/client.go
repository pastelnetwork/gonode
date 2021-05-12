package pastel

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel/jsonrpc"
)

type client struct {
	jsonrpc.RPCClient
}

// MasterNodeConfig implements pastel.Client.MasterNodeConfig()
func (client *client) MasterNodeConfig(ctx context.Context) (*MasterNodeConfig, error) {
	listConf := make(map[string]MasterNodeConfig)
	err := client.callFor(ctx, &listConf, "masternode", "list-conf")
	if err != nil {
		return nil, errors.Errorf("failed to get masternode configuration: %w", err)
	}

	if masterNodeConfig, ok := listConf["masternode"]; ok {
		return &masterNodeConfig, nil
	}
	return nil, errors.New("not found masternode configuration")
}

// MasterNodesTop implements pastel.Client.MasterNodesTop()
func (client *client) MasterNodesTop(ctx context.Context) (MasterNodes, error) {
	blocknumMNs := make(map[string]MasterNodes)
	err := client.callFor(ctx, &blocknumMNs, "masternode", "top")
	if err != nil {
		return nil, errors.Errorf("failed to get top masternodes: %w", err)
	}
	for _, masterNodes := range blocknumMNs {
		return masterNodes, nil
	}
	return nil, nil
}

// MasterNodesTop implements pastel.Client.MasterNodeStatus()
func (client *client) MasterNodeStatus(ctx context.Context) (*MasterNodeStatus, error) {
	var status MasterNodeStatus
	err := client.callFor(ctx, &status, "masternode", "status")
	if err != nil {
		return nil, errors.Errorf("failed to get masternode status: %w", err)
	}
	return &status, nil
}

// StorageFee implements pastel.Client.StorageFee()
func (client *client) StorageFee(ctx context.Context) (*StorageFee, error) {
	var storagefee StorageFee
	err := client.callFor(ctx, &storagefee, "storagefee", "getnetworkfee")
	if err != nil {
		return nil, errors.Errorf("failed to get storage fee: %w", err)
	}
	return &storagefee, nil
}

// IDTickets implements pastel.Client.IDTickets()
func (client *client) IDTickets(ctx context.Context, idType IDTicketType) (IDTickets, error) {
	tickets := IDTickets{}
	err := client.callFor(ctx, &tickets, "tickets", "list", "id", string(idType))
	if err != nil {
		errors.Errorf("failed to get id tickets: %w", err)
	}
	return tickets, nil
}

func (client *client) callFor(ctx context.Context, object interface{}, method string, params ...interface{}) error {
	err := client.CallForWithContext(ctx, &object, method, params)
	if err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return errors.New(err)
		}
		return errors.Errorf("could not call method %q, %w", method, err)
	}
	if object == nil {
		return errors.New("nothing found")
	}
	return nil
}

// func (client *client) call(ctx context.Context, method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
// 	response, err := client.CallWithContext(ctx, method, params)
// 	if err != nil {
// 		return nil, errors.Errorf("could not call %q, %w", method, err)
// 	}
// 	if response == nil {
// 		return nil, errors.Errorf("empty response on call %q", method)
// 	}
// 	if response.Error != nil {
// 		return nil, errors.Errorf("call %q returns error: %s", method, response.Error.Message)
// 	}
// 	if response.Result == nil {
// 		return nil, errors.Errorf("call %q returns empty result", method)
// 	}
// 	return response, nil
// }

// NewClient returns a new Client instance.
func NewClient(config *Config) Client {
	endpoint := net.JoinHostPort(config.Hostname(), strconv.Itoa(config.Port()))
	if !strings.Contains(endpoint, "//") {
		endpoint = "http://" + endpoint
	}

	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Username()+":"+config.Password())),
		},
	}

	return &client{
		RPCClient: jsonrpc.NewClientWithOpts(endpoint, opts),
	}
}
