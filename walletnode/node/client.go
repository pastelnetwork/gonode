package node

import (
	"context"
)

type Client interface {
	Connect(ctx context.Context, address string) (Connection, error)
}

type Connection interface {
	Close() error
	Done() <-chan struct{}
	RegisterArtowrk(ctx context.Context) (RegisterArtowrk, error)
}

type RegisterArtowrk interface {
	Handshake(ctx context.Context, connID string, IsPrimary bool) error
	PrimaryAcceptSecondary(ctx context.Context) (SuperNodes, error)
	SecondaryConnectToPrimary(ctx context.Context, nodeKey string) error
}
