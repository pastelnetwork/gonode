package test

import (
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

// ListenOnConnect listening Connect call and returning error from args
func (c *Client) ListenOnConnect(returnErr error) *Client {
	c.ClientMock.On("Connect", mock.Anything, mock.IsType(string(""))).Return(c.ConnectionMock, returnErr)
	return c
}

// ListenOnClose listening Close call and returning error from args
func (c *Client) ListenOnClose(returnErr error) *Client {
	c.ConnectionMock.On("Close").Return(returnErr)
	return c
}

// ListenOnProbeImage listening ProbeImage call and returning error from args
func (c *Client) ListenOnProbeImage(fingerprint []byte, returnErr error) *Client {
	c.RegArtWorkMock.On("ProbeImage", mock.Anything, mock.IsType(&artwork.File{})).Return(fingerprint, returnErr)
	return c
}
