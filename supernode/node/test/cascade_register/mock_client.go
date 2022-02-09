package test

import (
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/pastelnetwork/gonode/supernode/node/mocks"
)

const (
	// ConnectMethod represent Connect name method
	ConnectMethod = "Connect"

	// RegisterCascadeMethod represent RegisterSenseInterface name method
	RegisterCascadeMethod = "RegisterCascade"

	// SendCascadeTicketSignatureMethod represent SendCascadeTicketSignature name method
	SendCascadeTicketSignatureMethod = "SendCascadeTicketSignature"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.RegisterCascadeInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                        t,
		ClientInterface:          &mocks.ClientInterface{},
		ConnectionInterface:      &mocks.ConnectionInterface{},
		RegisterCascadeInterface: &mocks.RegisterCascadeInterface{},
	}
}

// ListenOnRegisterCascade listening RegisterCascadeInterface call
func (client *Client) ListenOnRegisterCascade() *Client {
	client.ConnectionInterface.On(RegisterCascadeMethod).Return(client.RegisterCascadeInterface)
	return client
}

// AssertRegisterCascadeCall assertion AssertRegisterCascade call
func (client *Client) AssertRegisterCascadeCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ConnectionInterface.AssertCalled(client.t, RegisterCascadeMethod, arguments...)
	}
	client.ConnectionInterface.AssertNumberOfCalls(client.t, RegisterCascadeMethod, expectedCalls)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(addr string, returnErr error) *Client {
	if addr == "" {
		client.ClientInterface.On(ConnectMethod, mock.Anything, mock.IsType(string(""))).Return(client.ConnectionInterface, returnErr)
	} else {
		client.ClientInterface.On(ConnectMethod, mock.Anything, addr).Return(client.ConnectionInterface, returnErr)
	}

	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.ClientInterface.AssertCalled(client.t, ConnectMethod, arguments...)
	}
	client.ClientInterface.AssertNumberOfCalls(client.t, ConnectMethod, expectedCalls)
	return client
}

// ListenOnSendCascadeTicketSignature listens on send Sense ticket signature
func (client *Client) ListenOnSendCascadeTicketSignature(returnErr error) *Client {
	client.RegisterCascadeInterface.On(SendCascadeTicketSignatureMethod, mock.Anything, mock.Anything, mock.Anything).Return(returnErr)
	return client
}
