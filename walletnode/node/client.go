package node

import (
	"context"
)

type Client interface {
	Connect(ctx context.Context, address string) (Connection, error)
}

type Connection interface {
	Close() error
	RegisterArtowrk(ctx context.Context) (RegisterArtowrk, error)
}

type RegisterArtowrk interface {
	Handshake(connID string, IsPrimary bool) error
	PrimaryAcceptSecondary() (SuperNodes, error)
	SecondaryConnectToPrimary(nodeKey string) error
}
