//go:generate mockery --name=Client

package p2p

import "context"

// PeerInfo is every thing related to an p2p peer
type PeerInfo struct {
	// id is a 32 byte unique identifier
	ID []byte

	// ip address of the node
	IP string

	// port of the node
	Port int
}

// return current status of backend Database
type DbStatus struct {
	SizeInMb int // Current size of DB
	KeyCount int // number of keys in the table
	//.. others put here
}

// Client exposes the interfaces for p2p service
type Client interface {
	// Retrieve data from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the data
	Retrieve(ctx context.Context, key string) ([]byte, error)

	// Store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	Store(ctx context.Context, data []byte) (string, error)

	// Peers return info of peers
	Peers(ctx context.Context) ([]PeerInfo, error)

	// DatabaseStatus return status of backend database
	DatabaseStatus(ctx context.Context) (DbStatus, error)
}
