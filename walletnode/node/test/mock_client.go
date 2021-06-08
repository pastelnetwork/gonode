package test

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
	"github.com/stretchr/testify/mock"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	*mocks.Client
	*mocks.Connection
	*mocks.RegisterArtwork
	acceptedNodesMethod   string
	closeMethod           string
	connectMethod         string
	connectToMethod       string
	doneMethod            string
	probeImageMethod      string
	registerArtWorkMethod string
	sessionMethod         string
	sessIDMethod          string
}

// NewMockClient create new client mock
func NewMockClient() *Client {
	return &Client{
		Client:                &mocks.Client{},
		Connection:            &mocks.Connection{},
		RegisterArtwork:       &mocks.RegisterArtwork{},
		acceptedNodesMethod:   "AcceptedNodes",
		closeMethod:           "Close",
		connectMethod:         "Connect",
		connectToMethod:       "ConnectTo",
		doneMethod:            "Done",
		probeImageMethod:      "ProbeImage",
		registerArtWorkMethod: "RegisterArtwork",
		sessionMethod:         "Session",
		sessIDMethod:          "SessID",
	}
}

// ListenOnRegisterArtwork listening RegisterArtwork call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.Connection.On(client.registerArtWorkMethod).Return(client.RegisterArtwork)
	return client
}

// AssertRegisterArtworkCall assertion RegisterArtwork call
func (client *Client) AssertRegisterArtworkCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(t, client.registerArtWorkMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(t, client.registerArtWorkMethod, expectedCalls)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(returnErr error) *Client {
	client.Client.On(client.connectMethod, mock.Anything, mock.IsType(string(""))).Return(client.Connection, returnErr)
	return client
}

// AssertConnectCall assertion Connect call
func (client *Client) AssertConnectCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Client.AssertCalled(t, client.connectMethod, arguments...)
	}
	client.Client.AssertNumberOfCalls(t, client.connectMethod, expectedCalls)
	return client
}

// ListenOnClose listening Close call and returns error from args
func (client *Client) ListenOnClose(returnErr error) *Client {
	client.Connection.On(client.closeMethod).Return(returnErr)
	return client
}

// AssertCloseCall assertion Close call
func (client *Client) AssertCloseCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(t, client.closeMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(t, client.closeMethod, expectedCalls)
	return client
}

// ListenOnDone listening Done call and returns channel from args
func (client *Client) ListenOnDone() *Client {
	client.Connection.On(client.doneMethod).Return(nil)
	return client
}

// AssertDoneCall assertion Done call
func (client *Client) AssertDoneCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.Connection.AssertCalled(t, client.doneMethod, arguments...)
	}
	client.Connection.AssertNumberOfCalls(t, client.doneMethod, expectedCalls)
	return client
}

// ListenOnProbeImage listening ProbeImage call and returns args value
func (client *Client) ListenOnProbeImage(arguments ...interface{}) *Client {
	client.RegisterArtwork.On(client.probeImageMethod, mock.Anything, mock.IsType(&artwork.File{})).Return(arguments...)
	return client
}

// AssertProbeImageCall assertion ProbeImage call
func (client *Client) AssertProbeImageCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(t, client.probeImageMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(t, client.probeImageMethod, expectedCalls)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterArtwork.On(client.sessionMethod, mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// AssertSessionCall assertion Session Call
func (client *Client) AssertSessionCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(t, client.sessionMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(t, client.sessionMethod, expectedCalls)
	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterArtwork.On(client.acceptedNodesMethod, mock.Anything).Return(handleFunc, returnErr)
	return client
}

// AssertAcceptedNodesCall assertion AcceptedNodes call
func (client *Client) AssertAcceptedNodesCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(t, client.acceptedNodesMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(t, client.acceptedNodesMethod, expectedCalls)
	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.RegisterArtwork.On(client.connectToMethod, mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return client
}

// AssertConnectToCall assertion ConnectTo call
func (client *Client) AssertConnectToCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(t, client.connectToMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(t, client.connectToMethod, expectedCalls)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.RegisterArtwork.On(client.sessIDMethod).Return(sessID)
	return client
}

// AssertSessIDCall assertion SessID call
func (client *Client) AssertSessIDCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.RegisterArtwork.AssertCalled(t, client.sessIDMethod, arguments...)
	}
	client.RegisterArtwork.AssertNumberOfCalls(t, client.sessIDMethod, expectedCalls)
	return client
}
