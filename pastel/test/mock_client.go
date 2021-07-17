package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// MasterNodesTopMethod represent MasterNodesTop name method
	MasterNodesTopMethod = "MasterNodesTop"

	// StorageNetWorkFeeMethod represent StorageNetworkFee name method
	StorageNetWorkFeeMethod = "StorageNetworkFee"

	// SignMethod represent Sign name method
	SignMethod = "Sign"

	// ActTicketsMethod represent ActTickets name method
	ActTicketsMethod = "ActTickets"

	// RegTicketMethod represent RegTicket name method
	RegTicketMethod = "RegTicket"

	// GetBlockVerbose1Method represent GetBlockVerbose1 method
	GetBlockVerbose1Method = "GetBlockVerbose1"

	// ListenOnGetBlockCountMethod represent  GetBlockCount method
	GetBlockCountMethod = "GetBlockCount"
)

// Client implementing pastel.Client for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
}

// NewMockClient new Client instance
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:      t,
		Client: &mocks.Client{},
	}
}

// ListenOnMasterNodesTop listening MasterNodesTop and returns Mn's and error from args
func (client *Client) ListenOnMasterNodesTop(nodes pastel.MasterNodes, err error) *Client {
	client.On(MasterNodesTopMethod, mock.Anything).Return(nodes, err)
	return client
}

// ListenOnStorageNetworkFee listening StorageNetworkFee call and returns pastel.StorageNetworkFee, error form args
func (client *Client) ListenOnStorageNetworkFee(fee float64, returnErr error) *Client {
	client.On(StorageNetWorkFeeMethod, mock.Anything).Return(fee, returnErr)
	return client
}

// ListenOnSign listening Sign call aand returns values from args
func (client *Client) ListenOnSign(signature []byte, returnErr error) *Client {
	client.On(SignMethod, mock.Anything, mock.IsType([]byte{}), mock.IsType(string("")), mock.IsType(string(""))).Return(signature, returnErr)
	return client
}

// AssertMasterNodesTopCall MasterNodesTop call assertion
func (client *Client) AssertMasterNodesTopCall(expectedCalls int, arguments ...interface{}) *Client {
	//don't check AssertCalled when expectedCall is 0. it become always fail
	if expectedCalls > 0 {
		client.AssertCalled(client.t, MasterNodesTopMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, MasterNodesTopMethod, expectedCalls)
	return client
}

// AssertStorageNetworkFeeCall StorageNetworkFee call assertion
func (client *Client) AssertStorageNetworkFeeCall(expectedCalls int, arguments ...interface{}) *Client {
	//don't check AssertCalled when expectedCall is 0. it become always fail
	if expectedCalls > 0 {
		client.AssertCalled(client.t, StorageNetWorkFeeMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, StorageNetWorkFeeMethod, expectedCalls)
	return client
}

// AssertSignCall Sign call assertion
func (client *Client) AssertSignCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, SignMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, SignMethod, expectedCalls)
	return client
}

// ListenOnActTickets listening ActTickets and returns tickets and error from args
func (client *Client) ListenOnActTickets(tickets pastel.ActTickets, err error) *Client {
	client.On(ActTicketsMethod, mock.Anything, mock.Anything, mock.Anything).Return(tickets, err)
	return client
}

// ListenOnRegTicket listening RegTicket and returns ticket and error from args
func (client *Client) ListenOnRegTicket(id string, ticket pastel.RegTicket, err error) *Client {
	client.On(RegTicketMethod, mock.Anything, id).Return(ticket, err)
	return client
}

// ListenOnGetBlockCount listening GetBlockCount and returns blockNum and error from args
func (client *Client) ListenOnGetBlockCount(blockNum int32, err error) *Client {
	client.On(GetBlockCountMethod, mock.Anything).Return(blockNum, err)
	return client
}

// ListenOnGetBlockVerbose1 listening GetBlockVerbose1 and returns blockNum and error from args
func (client *Client) ListenOnGetBlockVerbose1(blockInfo *pastel.GetBlockVerbose1Result, err error) *Client {
	client.On(GetBlockVerbose1Method, mock.Anything, mock.Anything).Return(blockInfo, err)
	return client
}
