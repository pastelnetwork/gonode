package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/supernode/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// ConnectMethod represent Connect name method
	ConnectMethod = "Connect"

	// CloseMethod represent Close name method
	CloseMethod = "Close"

	// DoneMethod represent Done name method
	DoneMethod = "Done"

	// RegisterArtworkMethod represent RegisterArtwork name method
	RegisterArtworkMethod = "RegisterArtwork"

	// SessIDMethod represent SessID name method
	SessIDMethod = "SessID"

	// SessionMethod represent Session name method
	SessionMethod = "Session"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
	*mocks.Connection
	*mocks.RegisterArtwork
}

// NewMockClient create new Client mock
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:               t,
		Client:          &mocks.Client{},
		Connection:      &mocks.Connection{},
		RegisterArtwork: &mocks.RegisterArtwork{},
	}
}

// ListenOnConnectCall listening Connect call and returns error from args
func (client *Client) ListenOnConnectCall(returnErr error) *Client {
	client.Client.On(ConnectMethod, mock.Anything, mock.IsType(string(""))).Return(client.Connection, returnErr)
	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(expectedCalls int, argument ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Client.AssertCalled(client.t, ConnectMethod, argument...)
	}
	client.Client.AssertNumberOfCalls(client.t, ConnectMethod, expectedCalls)
	return client
}

// ListenOnCloseCall listening Close call and returns error from args
func (client *Client) ListenOnCloseCall(returnErr error) *Client {
	client.Connection.On(CloseMethod).Return(returnErr)
	return client
}

// AssertCloseCall assertion Close call
func (client *Client) AssertCloseCall(expectedCalls int, argument ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, CloseMethod, argument...)
	}
	client.Connection.AssertNumberOfCalls(client.t, CloseMethod, expectedCalls)
	return client
}

// ListenOnDoneCall listening Done call and returns c from args
func (client *Client) ListenOnDoneCall(c <-chan struct{}) *Client {
	client.Connection.On(DoneMethod).Return(c)
	return client
}

// AssertDoneCall assertion Done call
func (client *Client) AssertDoneCall(expectedCalls int, argument ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, DoneMethod, argument...)
	}
	client.Connection.AssertNumberOfCalls(client.t, DoneMethod, expectedCalls)
	return client
}

// ListenOnRegisterArtworkCall listening RegisterArtwork call
func (client *Client) ListenOnRegisterArtworkCall() *Client {
	client.Connection.On(RegisterArtworkMethod).Return(client.RegisterArtwork)
	return client
}

// AssertRegisterArtworkCall assertion RegisterArtwork call
func (client *Client) AssertRegisterArtworkCall(expectedCalls int, argument ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(client.t, RegisterArtworkMethod, argument...)
	}
	client.Connection.AssertNumberOfCalls(client.t, RegisterArtworkMethod, expectedCalls)
	return client
}

// ListenOnSessIDCall listening SessID call and return taskID from args
func (client *Client) ListenOnSessIDCall(taskID string) *Client {
	client.RegisterArtwork.On(SessIDMethod).Return(taskID)
	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(expectedCalls int, argument ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, SessIDMethod, argument...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, SessIDMethod, expectedCalls)
	return client
}

// ListenOnSessionCall listening Session call and return error from args
func (client *Client) ListenOnSessionCall(returnErr error) *Client {
	client.RegisterArtwork.On(SessionMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return client
}

// AssertSessionCall assertion Session call
func (client *Client) AssertSessionCall(expectedCalls int, argument ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(client.t, SessionMethod, argument...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(client.t, SessionMethod, expectedCalls)
	return client
}
