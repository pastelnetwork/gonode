package test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/pastelnetwork/gonode/hermes/service/node/mocks"
)

const (
	// ConnectMethod represent Connect name method
	ConnectMethod = "Connect"

	// HermesP2PMethod represent HermesP2PInterface name method
	HermesP2PMethod = "HermesP2P"

	// DeleteMethod represent delete method
	DeleteMethod = "Delete"
	// RetrieveMethod represent Retrieve  method
	RetrieveMethod = "Retrieve"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.SNClientInterface
	*mocks.ConnectionInterface
	*mocks.HermesP2PInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                   t,
		SNClientInterface:   &mocks.SNClientInterface{},
		ConnectionInterface: &mocks.ConnectionInterface{},
		HermesP2PInterface:  &mocks.HermesP2PInterface{},
	}
}

// ListenOnHermesP2P listening ListenOnHermesP2P call
func (client *Client) ListenOnHermesP2P() *Client {
	client.ConnectionInterface.On(HermesP2PMethod).Return(client.HermesP2PInterface)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		client.SNClientInterface.On(ConnectMethod, mock.Anything, mock.IsType(string(""))).Return(client.ConnectionInterface, returnErr)
	} else {
		client.SNClientInterface.On(ConnectMethod, mock.Anything, addr).Return(client.ConnectionInterface, returnErr)
	}

	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.SNClientInterface.AssertCalled(client.t, ConnectMethod, arguments...)
	}
	client.SNClientInterface.AssertNumberOfCalls(client.t, ConnectMethod, expectedCalls)
	return client
}

// ListenOnRetrieve listens on p2p retrieve
func (client *Client) ListenOnRetrieve(data []byte, returnErr error) *Client {
	client.HermesP2PInterface.On(RetrieveMethod, mock.Anything, mock.Anything).Return(data, returnErr)
	return client
}

// ListenOnDelete listens on p2p delete
func (client *Client) ListenOnDelete(returnErr error) *Client {
	client.HermesP2PInterface.On(DeleteMethod, mock.Anything, mock.Anything).Return(returnErr)
	return client
}
