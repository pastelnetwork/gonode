//go:generate mockery --name=Client

package p2p

import "context"

// Client exposes the interfaces for p2p service
type Client interface {
	// RetrieveData retrieve data from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the data
	RetrieveData(ctx context.Context, key string) ([]byte, error)

	// RetrieveThumbnails retrieve thumbnails from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the dat
	RetrieveThumbnails(ctx context.Context, key string) ([]byte, error)

	// RetrieveData retrieve fingerprints from the kademlia network by key, return nil if not found
	// - key is the base58 encoded identifier of the dat
	RetrieveFingerprints(ctx context.Context, key string) ([]byte, error)

	// StoreData store data to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	StoreData(ctx context.Context, data []byte) (string, error)

	// StoreThumbnails store thumbnail to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	StoreThumbnails(ctx context.Context, data []byte) (string, error)

	// StoreFingerprints store fingerprints to the network, which will trigger the iterative store message
	// - the base58 encoded identifier will be returned
	StoreFingerprints(ctx context.Context, data []byte) (string, error)

	// Stats return status of p2p
	Stats(ctx context.Context) (map[string]interface{}, error)
}
