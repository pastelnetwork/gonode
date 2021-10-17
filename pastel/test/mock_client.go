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

	// MasterNodeStatusMethod represent MasterNodesTop name method
	MasterNodeStatusMethod = "MasterNodeStatus"

	// StorageNetWorkFeeMethod represent StorageNetworkFee name method
	StorageNetWorkFeeMethod = "StorageNetworkFee"

	// SignMethod represent Sign name method
	SignMethod = "Sign"

	// ActTicketsMethod represent ActTickets name method
	ActTicketsMethod = "ActTickets"

	// RegTicketMethod represent RegTicket name method
	RegTicketMethod = "RegTicket"

	// RegTicketsMethod represent RegTickets name method
	RegTicketsMethod = "RegTickets"

	// GetBlockVerbose1Method represent GetBlockVerbose1 method
	GetBlockVerbose1Method = "GetBlockVerbose1"

	// GetBlockCountMethod represent  GetBlockCount method
	GetBlockCountMethod = "GetBlockCount"

	// FindTicketByIDMethod represent find ticket by ID method
	FindTicketByIDMethod = "FindTicketByID"

	// SendFromAddressMethod represent send from address method
	SendFromAddressMethod = "SendFromAddress"

	// GetRawTransactionVerbose1Method  represent GetRawTransactionVerbose1 method
	GetRawTransactionVerbose1Method = "GetRawTransactionVerbose1"

	// TicketOwnershipMethod represents TicketOwnership method name
	TicketOwnershipMethod = "TicketOwnership"

	// ListAvailableTradeTicketsMethod represents ListAvailableTradeTickets method name
	ListAvailableTradeTicketsMethod = "ListAvailableTradeTickets"

	// VerifyMethod represents Verify method name
	VerifyMethod = "Verify"

	// RegisterNFTTicketMethod represents RegisterNFTTicket method name
	RegisterNFTTicketMethod = "RegisterNFTTicket"

	// RegisterExDDTicketMethod represents RegisterExDDTicket method name
	RegisterExDDTicketMethod = "RegisterExDDTicket"

	// RegisterArtTicketMethod represents RegisterArtTicket method
	RegisterArtTicketMethod = "RegisterArtTicket"

	// GetRegisterNFTFeeMethod represents GetRegisterNFTFee method
	GetRegisterNFTFeeMethod = "GetRegisterNFTFee"

	// GetRegisterExDDFeeMethod represents GetRegisterExDDFee method
	GetRegisterExDDFeeMethod = "GetRegisterExDDFee"

	// MasterNodesExtraMethod represent MasterNodesExtra name method
	MasterNodesExtraMethod = "MasterNodesExtra"
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

// ListenOnMasterNodesExtra listening MasterNodesExtra and returns Mn's and error from args
func (client *Client) ListenOnMasterNodesExtra(nodes pastel.MasterNodes, err error) *Client {
	client.On(MasterNodesExtraMethod, mock.Anything).Return(nodes, err)
	return client
}

// ListenOnMasterNodeStatus listening MasterNodeStatus and returns master node status and error from args
func (client *Client) ListenOnMasterNodeStatus(status *pastel.MasterNodeStatus, err error) *Client {
	client.On(MasterNodeStatusMethod, mock.Anything).Return(status, err)
	return client
}

// ListenOnStorageNetworkFee listening StorageNetworkFee call and returns pastel.StorageNetworkFee, error form args
func (client *Client) ListenOnStorageNetworkFee(fee float64, returnErr error) *Client {
	client.On(StorageNetWorkFeeMethod, mock.Anything).Return(fee, returnErr)
	return client
}

// ListenOnSign listening Sign call aand returns values from args
func (client *Client) ListenOnSign(signature []byte, returnErr error) *Client {
	client.On(SignMethod, mock.Anything, mock.IsType([]byte{}), mock.IsType(string("")), mock.IsType(string("")), mock.Anything).Return(signature, returnErr)
	return client
}

// ListenOnSendFromAddress listening Send From Address Method & return txn id & err from args
func (client *Client) ListenOnSendFromAddress(burnTxnID string, returnErr error) *Client {
	client.On(SendFromAddressMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string("")), mock.Anything).Return(burnTxnID, returnErr)
	return client
}

// ListenOnGetRawTransactionVerbose1 listens on GetRawTransactionVerbose1 & return result & err
func (client *Client) ListenOnGetRawTransactionVerbose1(res *pastel.GetRawTransactionVerbose1Result, returnErr error) *Client {
	client.On(GetRawTransactionVerbose1Method, mock.Anything, mock.IsType(string(""))).Return(res, returnErr)
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

// ListenOnRegTickets listening RegTickets and returns ticket and error from args
func (client *Client) ListenOnRegTickets(ticket pastel.RegTickets, err error) *Client {
	client.On(RegTicketsMethod, mock.Anything).Return(ticket, err)
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

// ListenOnFindTicketByID listening FindTicketByID
func (client *Client) ListenOnFindTicketByID(idticket *pastel.IDTicket, err error) *Client {
	client.On(FindTicketByIDMethod, mock.Anything, mock.Anything).Return(idticket, err)
	return client
}

// ListenOnRegisterNFTTicket listening on RegisterNFTTicket
func (client *Client) ListenOnRegisterNFTTicket(txid string, err error) *Client {
	client.On(RegisterNFTTicketMethod, mock.Anything, mock.Anything).Return(txid, err)
	return client
}

// ListenOnRegisterExDDTicket listening on RegisterExDDTicket
func (client *Client) ListenOnRegisterExDDTicket(txid string, err error) *Client {
	client.On(RegisterExDDTicketMethod, mock.Anything, mock.Anything).Return(txid, err)
	return client
}

// ListenOnGetRegisterNFTFee listen on get register NFT Fee & return fee & err
func (client *Client) ListenOnGetRegisterNFTFee(retFee int64, retErr error) *Client {
	client.On(GetRegisterNFTFeeMethod, mock.Anything, mock.Anything).Return(retFee, retErr)
	return client
}

// ListenOnGetRegisterExDDFee listen on get register EDD Fee & return fee & err
func (client *Client) ListenOnGetRegisterExDDFee(retFee int64, retErr error) *Client {
	client.On(GetRegisterExDDFeeMethod, mock.Anything, mock.Anything).Return(retFee, retErr)
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
	client.On(VerifyMethod, mock.Anything, mock.IsType([]byte{}), mock.IsType(string("")), mock.IsType(string("")), mock.Anything).Return(isValid, returnErr)
	return client
}

// ListenOnRegisterArtTicket listenes register art ticket & return id & err
func (client *Client) ListenOnRegisterArtTicket(retID string, retErr error) *Client {
	client.On(RegisterArtTicketMethod, mock.Anything, mock.Anything).Return(retID, retErr)
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
