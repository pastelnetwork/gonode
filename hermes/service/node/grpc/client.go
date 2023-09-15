package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"google.golang.org/grpc"
)

const (
	logPrefix             = "hermes-grpc"
	defaultConnectTimeout = 45 * time.Second
)

type client struct{}

// Connect implements node.Client.Connect()
func (client *client) Connect(ctx context.Context, address string) (node.ConnectionInterface, error) {
	// Limits the dial timeout, prevent got stuck too long
	dialCtx, cancel := context.WithTimeout(ctx, defaultConnectTimeout)
	defer cancel()

	id, _ := random.String(8, random.Base62Chars)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, id))

	grpcConn, err := grpc.DialContext(dialCtx, address,
		//lint:ignore SA1019 we want to ignore this for now
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(35000000), grpc.MaxCallSendMsgSize(35000000)),
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
func NewClient() node.ClientInterface {
	return &client{}
}
