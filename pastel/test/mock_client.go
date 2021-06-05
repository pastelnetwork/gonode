package test

import (
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

// Client implementing pastel.Client for testing purpose
type Client struct {
	ClientMock *mocks.Client
}

// NewMockClient new Client instance
func NewMockClient() *Client {
	return &Client{
		ClientMock: &mocks.Client{},
	}
}

// ListenOnMasterNodesTop listening MasterNodesTop and returns Mn's and error from args
func (c *Client) ListenOnMasterNodesTop(nodes pastel.MasterNodes, err error) *Client {
	c.ClientMock.On("MasterNodesTop", mock.Anything).Return(nodes, err)
	return c
}

// ListenOnStorageFee listening StorageFee call and returns pastel.StorageFee, error form args
func (c *Client) ListenOnStorageFee(fee *pastel.StorageFee, returnErr error) *Client {
	c.ClientMock.On("StorageFee", mock.Anything).Return(fee, returnErr)
	return c
}
