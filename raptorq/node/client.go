//go:generate mockery --name=Client
//go:generate mockery --name=Connection
//go:generate mockery --name=RaptorQ

package node

import (
	"context"
)

type Config struct {
	RqFilesDir string
}

type RawSymbolIdFile struct {
	Id                string
	BlockHash         string
	PastelID          string
	SymbolIdentifiers []string
}

type EncoderParameters struct {
	Oti []byte
}

type EncodeInfo struct {
	SymbolIdFiles map[string]RawSymbolIdFile
	EncoderParam  EncoderParameters
}

type Encode struct {
	Symbols      map[string][]byte
	EncoderParam EncoderParameters
}

type Decode struct {
	File []byte
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
	RaptorQ(config *Config) RaptorQ
}

// RaptorQ contains methods for request services from RaptorQ service.
type RaptorQ interface {
	// Get map of symbols
	Encode(ctx context.Context, data []byte) (*Encode, error)

	// Get encode info(include encode parameters + symbol id files)
	EncodeInfo(ctx context.Context, data []byte, copies uint32, blockHash string, pastelID string) (*EncodeInfo, error)

	// Decode returns a path to restored file.
	Decode(ctx context.Context, encodeInfo *Encode) (*Decode, error)
}
