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

	// TicketOwnershipMethod represents TicketOwnership method name
	TicketOwnershipMethod = "TicketOwnership"

	// ListAvailableTradeTicketsMethod represents ListAvailableTradeTickets method name
	ListAvailableTradeTicketsMethod = "ListAvailableTradeTickets"

	// VerifyMethod represents Verify method name
	VerifyMethod = "Verify"
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

// AssertRegTicketCall RegTicket call assertion
func (client *Client) AssertRegTicketCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, RegTicketMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, RegTicketMethod, expectedCalls)
	return client
}

// ListenOnTicketOwnership listening TicketOwnership call and returns values from args
func (client *Client) ListenOnTicketOwnership(ttxID string, returnErr error) *Client {
	client.On(TicketOwnershipMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string("")), mock.IsType(string(""))).Return(ttxID, returnErr)
	return client
}

// AssertTicketOwnershipCall TicketOwnership call assertion
func (client *Client) AssertTicketOwnershipCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, TicketOwnershipMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, TicketOwnershipMethod, expectedCalls)
	return client
}

// ListenOnListAvailableTradeTickets listening ListAvailableTradeTickets call and returns values from args
func (client *Client) ListenOnListAvailableTradeTickets(tradeTickets []pastel.TradeTicket, returnErr error) *Client {
	client.On(ListAvailableTradeTicketsMethod, mock.Anything).Return(tradeTickets, returnErr)
	return client
}

// AssertListAvailableTradeTicketsCall ListAvailableTradeTickets call assertion
func (client *Client) AssertListAvailableTradeTicketsCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, ListAvailableTradeTicketsMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, ListAvailableTradeTicketsMethod, expectedCalls)
	return client
}

// ListenOnVerify listening Verify call and returns values from args
func (client *Client) ListenOnVerify(isValid bool, returnErr error) *Client {
	client.On(VerifyMethod, mock.Anything, mock.IsType([]byte{}), mock.IsType(string("")), mock.IsType(string(""))).Return(isValid, returnErr)
	return client
}

// AssertVerifyCall Verify call assertion
func (client *Client) AssertVerifyCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, VerifyMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, VerifyMethod, expectedCalls)
	return client
}
