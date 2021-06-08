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

// ListenOnStorageNetworkFee listening StorageNetworkFee call and returns pastel.StorageNetworkFee, error form args
func (c *Client) ListenOnStorageNetworkFee(fee float64, returnErr error) *Client {
	c.ClientMock.On("StorageNetworkFee", mock.Anything).Return(fee, returnErr)
	return c
}
