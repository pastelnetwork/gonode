package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/p2p/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// RetrieveMethod represent Get name method
	RetrieveMethod = "Retrieve"

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

// ListenOnRetrieve listening Retrieve and returns data, and error from args
func (client *Client) ListenOnRetrieve(data []byte, err error) *Client {
	client.On(RetrieveMethod, mock.Anything, mock.Anything).Return(data, err)
	return client
}

//  ListenOnStore listening  Store and returns id and error from args
func (client *Client) ListenOnStore(id string, err error) *Client {
	client.On(StoreMethod, mock.Anything, mock.Anything).Return(id, err)
	return client
}
