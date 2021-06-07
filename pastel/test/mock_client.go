package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

// Client implementing pastel.Client for testing purpose
type Client struct {
	*mocks.Client
	masterNodesTopMethod    string
	storageNetWorkFeeMethod string
}

// NewMockClient new Client instance
func NewMockClient() *Client {
	return &Client{
		Client:                  &mocks.Client{},
		masterNodesTopMethod:    "MasterNodesTop",
		storageNetWorkFeeMethod: "StorageNetworkFee",
	}
}

// ListenOnMasterNodesTop listening MasterNodesTop and returns Mn's and error from args
func (client *Client) ListenOnMasterNodesTop(nodes pastel.MasterNodes, err error) *Client {
	client.On(client.masterNodesTopMethod, mock.Anything).Return(nodes, err)
	return client
}

// ListenOnStorageNetworkFee listening StorageNetworkFee call and returns pastel.StorageNetworkFee, error form args
func (client *Client) ListenOnStorageNetworkFee(fee float64, returnErr error) *Client {
	client.On(client.storageNetWorkFeeMethod, mock.Anything).Return(fee, returnErr)
	return client
}

// AssertMasterNodesTopCall MasterNodesTop call assertion
func (client *Client) AssertMasterNodesTopCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	//don't check AssertCalled when expectedCall is 0. it become always fail
	if expectedCalls > 0 {
		client.AssertCalled(t, client.masterNodesTopMethod, arguments...)
	}
	client.AssertNumberOfCalls(t, client.masterNodesTopMethod, expectedCalls)
	return client
}

// AssertStorageNetworkFeeCall StorageNetworkFee call assertion
func (client *Client) AssertStorageNetworkFeeCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Client {
	//don't check AssertCalled when expectedCall is 0. it become always fail
	if expectedCalls > 0 {
		client.AssertCalled(t, client.storageNetWorkFeeMethod, arguments...)
	}
	client.AssertNumberOfCalls(t, client.storageNetWorkFeeMethod, expectedCalls)
	return client
}
