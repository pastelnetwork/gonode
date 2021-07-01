//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RaptorQ

package node

import (
	"context"
	//"github.com/pastelnetwork/gonode/common/service/artwork"
)

type EncoderParameters struct {
	TransferLength  uint64
	SymbolSize      uint32
	NumSourceBlocks uint32
	NumSubBlocks    uint32
	SymbolAlignment uint32
}

type EncodeInfo struct {
	SymbolIds    []string
	EncoderParam EncoderParameters
}

// Client represents a base connection interface.
type Client interface {
	// Connect connects to the server at the given address.
	Connect(ctx context.Context, address string) (Connection, error)
}

// Connection represents a client connection
type Connection interface {
	// Close closes connection.
	Close() error
	// Done returns a channel that's closed when connection is shutdown.
	Done() <-chan struct{}
	// RaptorQ returns a new RaptorQ stream.
	RaptorQ() RaptorQ
}

// RaptorQ contains methods for request services from RaptorQ service.
type RaptorQ interface {
	// Get list of encoded symbols
	Encode(ctx context.Context, data []byte) ([][]byte, error)

	// Get encode info(include symbol ids + encode parameters)
	// symbol id is hash of symbol
	EncodeInfo(ctx context.Context, data []byte) (*EncodeInfo, error)
}
