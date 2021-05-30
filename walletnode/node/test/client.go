package test

import (
	"context"
	"time"

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
	c.ClientMock.On("Connect", mock.Anything, mock.AnythingOfType("string")).Return(c.ConnectionMock, returnErr)
	return c
}

// ListenOnClose listening Close call and returning error from args
func (c *Client) ListenOnClose(returnErr error) *Client {
	c.ConnectionMock.On("Close").Return(returnErr)
	return c
}

// ListenOnUploadImage listening UploadImage call and returning error from args
func (c *Client) ListenOnUploadImage(returnErr error) *Client {
	c.RegArtWorkMock.On("UploadImage", mock.Anything, mock.AnythingOfType("*artwork.File")).Return(returnErr)
	return c
}

// ListenOnSession listening Session call and returning error from args
func (c *Client) ListenOnSession(returnErr error) *Client {
	c.RegArtWorkMock.On("Session", mock.Anything, mock.AnythingOfType("bool")).Return(returnErr)
	return c
}

// ListenOnAcceptedNodes listening AcceptedNodes call and returning pastelIDs and error from args.
func (c *Client) ListenOnAcceptedNodes(pastelIDs []string, returnErr error) *Client {
	handleFunc := func(cxt context.Context) []string {
		time.Sleep(time.Second * 2)
		return pastelIDs
	}

	c.RegArtWorkMock.On("AcceptedNodes", mock.Anything).Return(handleFunc, returnErr)
	return c
}

// ListenOnConnectTo listening ConnectTo call and returning error from args
func (c *Client) ListenOnConnectTo(returnErr error) *Client {
	c.RegArtWorkMock.On("ConnectTo", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(returnErr)
	return c
}

// ListenOnSessID listening SessID call and returning sessID from args
func (c *Client) ListenOnSessID(sessID string) *Client {
	c.RegArtWorkMock.On("SessID").Return(sessID)
	return c
}

type Clients []*Client

func NewClients() Clients {
	return make(Clients, 0)
}

func (c *Clients) Add(i *Client) {
	*c = append(*c, i)
}
