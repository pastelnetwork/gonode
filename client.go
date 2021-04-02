package pastel

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/go-pastel/models"

	"github.com/ybbus/jsonrpc/v2"
)

type Client struct {
	jsonrpc.RPCClient
}

func (client *Client) Getblockchaininfo() (*models.BlockchainInfo, error) {
	info := &models.BlockchainInfo{}
	err := client.callFor(info, "getblockchaininfo")
	return info, err
}

func (client *Client) ListIDTickets(idType string) (*[]models.IdTicket, error) {
	tickets := &[]models.IdTicket{}
	err := client.callFor(tickets, "tickets", "list", "id", idType)
	return tickets, err
}

func (client *Client) FindIDTicket(search string) (*models.IdTicket, error) {
	tickets := &models.IdTicket{}
	err := client.callFor(tickets, "tickets", "find", "id", search)
	return tickets, err
}

func (client *Client) FindIDTickets(search string) (*[]models.IdTicket, error) {
	tickets := &[]models.IdTicket{}
	err := client.callFor(tickets, "tickets", "find", "id", search)
	return tickets, err
}

func (client *Client) ListPastelIDs() (*[]models.PastelID, error) {
	pastelids := &[]models.PastelID{}
	err := client.callFor(pastelids, "pastelid", "list")
	return pastelids, err
}

func (client *Client) GetMNRegFee() (int, error) {
	r, err := client.call("storagefee", "getnetworkfee")
	if err != nil {
		return -1, err
	}
	return r.Result.(map[string]interface{})["networkfee"].(int), nil
}

func (client *Client) callFor(object interface{}, method string, params ...interface{}) error {
	err := client.CallFor(&object, method, params)
	if err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			if err.Value == "string" {
				var errMsg string
				if err := client.CallFor(&errMsg, method, params); err != nil {
					return errors.New(errMsg)
				}
			}
		}
		return errors.Errorf("could not call method %q, %s", method, err)
	}
	if object == nil {
		return errors.New("nothing found")
	}
	return nil
}

func (client *Client) call(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	response, err := client.Call(method, params)
	if err != nil {
		return nil, errors.Errorf("could not call %q, %s", method, err)
	}
	if response == nil {
		return nil, errors.Errorf("empty response on call %q", method)
	}
	if response.Error != nil {
		return nil, errors.Errorf("call %q returns error: %s", method, response.Error.Message)
	}
	if response.Result == nil {
		return nil, errors.Errorf("call %q returns empty result", method)
	}
	return response, nil
}

func NewClient(endpoint, username, password string) *Client {
	return &Client{
		RPCClient: jsonrpc.NewClientWithOpts(endpoint,
			&jsonrpc.RPCClientOpts{
				CustomHeaders: map[string]string{
					"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
				},
			}),
	}
}
