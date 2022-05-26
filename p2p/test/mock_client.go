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
	// DeleteMethod represent Store name method
	DeleteMethod = "Delete"
	//NClosestMethod mocks getting the n closest nodes to a given string
	NClosestMethod = "NClosestNodes"
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
	client.On(RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(data, err)
	return client
}

// ListenOnDelete listening Delete and returns error from args
func (client *Client) ListenOnDelete(err error) *Client {
	client.On(DeleteMethod, mock.Anything, mock.Anything).Return(err)
	return client
}

/*
// AssertRetrieveCall is Retrieve call assertion
func (client *Client) AssertRetrieveCall(expectedCalls int, arguments ...interface{}) *Client {
	if expectedCalls > 0 {
		client.AssertCalled(client.t, RetrieveMethod, arguments...)
	}
	client.AssertNumberOfCalls(client.t, RetrieveMethod, expectedCalls)
	return client
}
*/

//  ListenOnStore listening  Store and returns id and error from args
func (client *Client) ListenOnStore(id string, err error) *Client {
	client.On(StoreMethod, mock.Anything, mock.Anything).Return(id, err)
	return client
}

// ListenOnNClosestNodes returns retArr and error from args
func (client *Client) ListenOnNClosestNodes(retArr []string, err error) *Client {

	client.On(NClosestMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(retArr, err)
	return client
}
