package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/p2p/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// RetrieveDataMethod represent retrieve data name method
	RetrieveDataMethod = "RetrieveData"
	// RetrieveThumbnailsMethod represent retrieve thumbnails name method
	RetrieveThumbnailsMethod = "RetrieveThumbnails"
	// RetrieveFingerprintsMethod represent retrieve Fingerprints name method
	RetrieveFingerprintsMethod = "RetrieveFingerprints"

	// StoreDataMethod represent Store data name method
	StoreDataMethod = "StoreData"
	// StoreThumbnailsMethod represent Store Thumbnails name method
	StoreThumbnailsMethod = "StoreThumbnails"
	// StoreFingerprintsMethod represent Store Fingerprints name method
	StoreFingerprintsMethod = "StoreFingerprints"
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

// ListenOnRetrieveData listening Retrieve and returns data, and error from args
func (client *Client) ListenOnRetrieveData(data []byte, err error) *Client {
	client.On(RetrieveDataMethod, mock.Anything, mock.Anything).Return(data, err)
	return client
}

// ListenOnRetrieveThumbnails listening Retrieve and returns data, and error from args
func (client *Client) ListenOnRetrieveThumbnails(data []byte, err error) *Client {
	client.On(RetrieveThumbnailsMethod, mock.Anything, mock.Anything).Return(data, err)
	return client
}

// ListenOnRetrieveFingerprints listening Retrieve and returns data, and error from args
func (client *Client) ListenOnRetrieveFingerprints(data []byte, err error) *Client {
	client.On(RetrieveFingerprintsMethod, mock.Anything, mock.Anything).Return(data, err)
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

//  ListenOnStoreData listening  Store and returns id and error from args
func (client *Client) ListenOnStoreData(id string, err error) *Client {
	client.On(StoreDataMethod, mock.Anything, mock.Anything).Return(id, err)
	return client
}

//  ListenOnStoreThumbnails listening  Store and returns id and error from args
func (client *Client) ListenOnStoreThumbnails(id string, err error) *Client {
	client.On(StoreThumbnailsMethod, mock.Anything, mock.Anything).Return(id, err)
	return client
}

//  ListenOnStoreFingerprints listening  Store and returns id and error from args
func (client *Client) ListenOnStoreFingerprints(id string, err error) *Client {
	client.On(StoreFingerprintsMethod, mock.Anything, mock.Anything).Return(id, err)
	return client
}
