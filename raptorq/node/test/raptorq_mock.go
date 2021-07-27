package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/raptorq/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// EncodeInfoMethod represent EncodeInfo name method
	EncodeInfoMethod = "EncodeInfo"
)

// Client implementing pastel.Client for testing purpose
type RaptorQ struct {
	t *testing.T
	*mocks.RaptorQ
}

// NewMockRaptorQ new Client instance
func NewMockRaptorQ(t *testing.T) *RaptorQ {
	return &RaptorQ{
		t:       t,
		RaptorQ: &mocks.RaptorQ{},
	}
}

// ListenOnEncodeInfo listening EncodeInfo and returns Mn's and error from args
func (client *RaptorQ) ListenOnEncodeInfo(encodeInfo *node.EncodeInfo, err error) *RaptorQ {
	client.On(EncodeInfoMethod, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(encodeInfo, err)

	return client
}
