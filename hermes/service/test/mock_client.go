package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/hermes/service/mocks"
)

const (
	// GetActionFeeMethod represents GetActionFee method
	GetActionFeeMethod = "GetActionFee"
)

// Client implementing service.Service for testing purpose
type Client struct {
	t *testing.T
	*mocks.SvcInterface
}

// NewMockClient new Client instance
func NewMockClient(t *testing.T) *Client {
	return &Client{
		t: t,
	}
}
