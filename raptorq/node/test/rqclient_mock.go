package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/raptorq/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// ConnectMethod represent MasterNodesTop name method
	ConnectMethod = "Connect"
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

// ListenOnMasterNodesTop listening MasterNodesTop and returns Mn's and error from args
func (client *Client) ListenOnConnect(conn node.Connection, err error) *Client {
	client.On(ConnectMethod, mock.Anything, mock.Anything).Return(conn, err)
	return client
}
