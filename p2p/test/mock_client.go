package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/p2p/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// GetMethod represent Get method name
	GetMethod = "Get"

	// StoreMethod represent Store method name
	StoreMethod = "Store"
)

// Client implementing p2p.Client mock for testing purpose
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

// ListenOnGet listening Get call and returns values from args
func (client *Client) ListenOnGet(data []byte, found bool, err error) *Client {
	client.On(GetMethod, mock.Anything, mock.IsType(string(""))).Return(data, found, err)
	return client
}

// ListenOnStore listening STore call and returns values from args
func (client *Client) ListenOnStore(id string, err error) *Client {
	client.On(StoreMethod, mock.Anything, mock.IsType([]byte{})).Return(id, err)
	return client
}

// AssertGetCall assertion Get call
func (client *Client) AssertGetCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, GetMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, GetMethod, expectedCalls)
	return client
}

// AssertStoreCall assertion Store call
func (client *Client) AssertStoreCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, StoreMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, StoreMethod, expectedCalls)
	return client
}
