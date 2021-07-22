package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/raptorq/node/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// RaptorQMethod represent RaptorQ name method
	RaptorQMethod = "RaptorQ"
	// CloseMethod represent Close name method
	CloseMethod = "Close"
)

// Connection implementing node.Connection for testing purpose
type Connection struct {
	t *testing.T
	*mocks.Connection
}

// NewMockConnection new Connection instance
func NewMockConnection(t *testing.T) *Connection {
	return &Connection{
		t:          t,
		Connection: &mocks.Connection{},
	}
}

// ListenOnRaptorQ listening ListenOnRaptorQ and returns node and error from args
func (conn *Connection) ListenOnRaptorQ(raptorq node.RaptorQ, err error) *Connection {
	conn.On(RaptorQMethod, mock.Anything).Return(raptorq, err)

	return conn
}

// ListenOnClose listening ListenOnClose and returns  error from args
func (conn *Connection) ListenOnClose(err error) *Connection {
	conn.On(CloseMethod).Return(err)

	return conn
}
