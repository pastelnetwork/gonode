package test

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
	"github.com/stretchr/testify/mock"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	*mocks.Client
	*mocks.Connection
	*mocks.RegisterArtwork
	connectMethod         string
	registerArtWorkMethod string
}

// NewMockClient create new client mock
func NewMockClient() *Client {
	return &Client{
		Client:                &mocks.Client{},
		Connection:            &mocks.Connection{},
		RegisterArtwork:       &mocks.RegisterArtwork{},
		connectMethod:         "Connect",
		registerArtWorkMethod: "RegisterArtwork",
	}
}

// ListenOnRegisterArtwork listening RegisterArtwork call
func (client *Client) ListenOnRegisterArtwork() *Client {
	client.Connection.On("RegisterArtwork").Return(client.RegisterArtwork)
	return client
}

// ListenOnConnect listening Connect call and returns error from args
func (client *Client) ListenOnConnect(returnErr error) *Client {
	client.Client.On("Connect", mock.Anything, mock.IsType(string(""))).Return(client.Connection, returnErr)
	return client
}

// ListenOnClose listening Close call and returns error from args
func (client *Client) ListenOnClose(returnErr error) *Client {
	client.Connection.On("Close").Return(returnErr)
	return client
}

// ListenOnDone listening Done call and returns channel from args
func (client *Client) ListenOnDone() *Client {
	client.Connection.On("Done").Return(nil)
	return client
}

// ListenOnProbeImage listening ProbeImage call and returns error from args
func (client *Client) ListenOnProbeImage(fingerprint []byte, returnErr error) *Client {
	client.RegisterArtwork.On("ProbeImage", mock.Anything, mock.IsType(&artwork.File{})).Return(fingerprint, returnErr)
	return client
}

// ListenOnSession listening Session call and returns error from args
func (client *Client) ListenOnSession(returnErr error) *Client {
	client.RegisterArtwork.On("Session", mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return client
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (client *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	client.RegisterArtwork.On("AcceptedNodes", mock.Anything).Return(handleFunc, returnErr)
	return client
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (client *Client) ListenOnConnectTo(returnErr error) *Client {
	client.RegisterArtwork.On("ConnectTo", mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return client
}

// ListenOnSessID listening SessID call and returns sessID from args
func (client *Client) ListenOnSessID(sessID string) *Client {
	client.RegisterArtwork.On("SessID").Return(sessID)
	return client
}
