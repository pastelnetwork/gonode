package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/p2p/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// RetrieveMethod represent Retrieve method name
	RetrieveMethod = "Retrieve"

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

// ListenOnRetrieve listening Retrieve call and returns values from args
func (client *Client) ListenOnRetrieve(data []byte, err error) *Client {
	client.On(RetrieveMethod, mock.Anything, mock.IsType(string(""))).Return(data, err)
	return client
}

// ListenOnStore listening Store call and returns values from args
func (client *Client) ListenOnStore(id string, err error) *Client {
	client.On(StoreMethod, mock.Anything, mock.IsType([]byte{})).Return(id, err)
	return client
}

// AssertRetrieveCall assertion Retrieve call
func (client *Client) AssertRetrieveCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, RetrieveMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, RetrieveMethod, expectedCalls)
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
