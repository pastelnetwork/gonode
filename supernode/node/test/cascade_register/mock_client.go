package test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/pastelnetwork/gonode/supernode/node/mocks"
)

const (
	// ConnectMethod represent Connect name method
	ConnectMethod = "Connect"

	// RegisterCascadeMethod represent RegisterSenseInterface name method
	RegisterCascadeMethod = "RegisterCascade"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.ClientInterface
	*mocks.ConnectionInterface
	*mocks.RegisterSenseInterface
}

// NewMockClient create new client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:                      t,
		ClientInterface:        &mocks.ClientInterface{},
		ConnectionInterface:    &mocks.ConnectionInterface{},
		RegisterSenseInterface: &mocks.RegisterSenseInterface{},
	}
}

// ListenOnRegisterSense listening RegisterSenseInterface call
func (client *Client) ListenOnRegisterSense() *Client {
	client.ConnectionInterface.On(RegisterCascadeMethod).Return(client.RegisterSenseInterface)
	return client
}

// AssertRegisterSenseCall assertion RegisterSenseInterface call
func (client *Client) AssertRegisterSenseCall(expectedCalls int, arguments ...interface{}) *Client {
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
