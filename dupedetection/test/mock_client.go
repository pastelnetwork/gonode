package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/dupedetection"
	"github.com/pastelnetwork/gonode/dupedetection/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// GenerateMethod represent generate method
	GenerateMethod = "Generate"
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

// ListenOnGenerate listening generate
func (client *Client) ListenOnGenerate(dd *dupedetection.DupeDetection, err error) *Client {
	client.On(GenerateMethod, mock.Anything, mock.Anything, mock.Anything).Return(dd, err)
	return client
}
