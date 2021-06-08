package test

import (
	"context"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
	"github.com/stretchr/testify/mock"
)

// Client implementing node.Client mock for testing purpose
type Client struct {
	ClientMock     *mocks.Client
	ConnectionMock *mocks.Connection
	RegArtWorkMock *mocks.RegisterArtwork
}

// NewMockClient create new client mock
func NewMockClient() *Client {
	return &Client{
		ClientMock:     &mocks.Client{},
		ConnectionMock: &mocks.Connection{},
		RegArtWorkMock: &mocks.RegisterArtwork{},
	}
}

// ListenOnRegisterArtwork listening RegisterArtwork call
func (c *Client) ListenOnRegisterArtwork() *Client {
	c.ConnectionMock.On("RegisterArtwork").Return(c.RegArtWorkMock)
	return c
}

// ListenOnConnect listening Connect call and returns error from args
func (c *Client) ListenOnConnect(returnErr error) *Client {
	c.ClientMock.On("Connect", mock.Anything, mock.IsType(string(""))).Return(c.ConnectionMock, returnErr)
	return c
}

// ListenOnClose listening Close call and returns error from args
func (c *Client) ListenOnClose(returnErr error) *Client {
	c.ConnectionMock.On("Close").Return(returnErr)
	return c
}

// ListenOnDone listening Done call and returns channel from args
func (c *Client) ListenOnDone() *Client {
	c.ConnectionMock.On("Done").Return(nil)
	return c
}

// ListenOnProbeImage listening ProbeImage call and returns error from args
func (c *Client) ListenOnProbeImage(fingerprint []byte, returnErr error) *Client {
	c.RegArtWorkMock.On("ProbeImage", mock.Anything, mock.IsType(&artwork.File{})).Return(fingerprint, returnErr)
	return c
}

// ListenOnSession listening Session call and returns error from args
func (c *Client) ListenOnSession(returnErr error) *Client {
	c.RegArtWorkMock.On("Session", mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return c
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returns pastelIDs and error from args.
func (c *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(ctx context.Context) []string {
		//need block operation until context is done
		<-ctx.Done()
		return pastelIDs
	}

	c.RegArtWorkMock.On("AcceptedNodes", mock.Anything).Return(handleFunc, returnErr)
	return c
}

// ListenOnConnectTo listening ConnectTo call and returns error from args
func (c *Client) ListenOnConnectTo(returnErr error) *Client {
	c.RegArtWorkMock.On("ConnectTo", mock.Anything, mock.IsType(string("")), mock.IsType(string(""))).Return(returnErr)
	return c
}

// ListenOnSessID listening SessID call and returns sessID from args
func (c *Client) ListenOnSessID(sessID string) *Client {
	c.RegArtWorkMock.On("SessID").Return(sessID)
	return c
}
