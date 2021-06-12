package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/p2p/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// GetMethod represent Get name method
	GetMethod = "Get"

	// StoreMethod represent Store name method
	StoreMethod = "Store"
)

// Client implementing pastel.Client for testing purpose
type Client struct {
	t *testing.T
	*mocks.Client
}

// NewMockClient new Client instance
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t:      t,
		Client: &mocks.Client{},
	}
}

// ListenOnGet listening MasterNodesTop and returns data, found flag, and error from args
func (client *Client) ListenOnGet(data []byte, found bool, err error) *Client {
	client.On(GetMethod, mock.Anything, mock.Anything).Return(data, found, err)
	return client
}

//  ListenOnStore listening  Store and returns id and error from args
func (client *Client) ListenOnStore(id string, err error) *Client {
	client.On(StoreMethod, mock.Anything, mock.Anything).Return(id, err)
	return client
}
