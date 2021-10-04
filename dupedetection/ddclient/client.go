package ddclient

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"

	"google.golang.org/grpc"
)

const (
	defaultConnectTimeout = 5 * time.Second
)

type client struct{}

// Connect implements node.Client.Connect()
func (cl *client) Connect(ctx context.Context, address string) (*clientConn, error) {
	// Limits the dial timeout, prevent got stuck too long
	dialCtx, cancel := context.WithTimeout(ctx, defaultConnectTimeout)
	defer cancel()

	id, _ := random.String(8, random.Base62Chars)

	grpcConn, err := grpc.DialContext(dialCtx, address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Errorf("fail to dial: %w", err).WithField("address", address)
	}
	log.WithContext(ctx).Debugf("Connected to %s", address)

	conn := newClientConn(id, grpcConn)
	go func() {
		<-conn.Done()
		log.WithContext(ctx).Debugf("Disconnected %s", grpcConn.Target())
	}()
	return conn, nil
}

// NewClient returns a new client instance.
func NewClient() *client {
	return &client{}
}