package test

import (
	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
	"github.com/stretchr/testify/mock"
)

type Client struct {
	ClientMock     *mocks.Client
	ConnectionMock *mocks.Connection
	RegArtWorkMock *mocks.RegisterArtwork
}

func NewMockClient() *Client {
	return &Client{
		ClientMock:     &mocks.Client{},
		ConnectionMock: &mocks.Connection{},
		RegArtWorkMock: &mocks.RegisterArtwork{},
	}
}

func (c *Client) ListenOnRegisterArtwork() *Client {
	c.ConnectionMock.On("RegisterArtwork").Return(c.RegArtWorkMock)
	return c
}

func (c *Client) ListenOnConnect(returnErr error) *Client {
	c.ClientMock.On("Connect", mock.Anything, mock.AnythingOfType("string")).Return(c.ConnectionMock, returnErr)
	return c
}

func (c *Client) ListenOnClose(returnErr error) *Client {
	c.ConnectionMock.On("Close").Return(returnErr)
	return c
}

func (c *Client) ListenOnUploadImage(returnErr error) *Client {
	c.RegArtWorkMock.On("UploadImage", mock.Anything, mock.AnythingOfType("*artwork.File")).Return(returnErr)
	return c
}
