package ddclient

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	defaultConnectTimeout = 1200 * time.Second
)

type client struct{}

// Connect implements node.Client.Connect()
func (cl *client) Connect(ctx context.Context, address string) (*clientConn, error) {
	// Limits the dial timeout, prevent got stuck too long
	dialCtx, cancel := context.WithTimeout(ctx, defaultConnectTimeout)
	defer cancel()

	id, _ := random.String(8, random.Base62Chars)

	grpcConn, err := grpc.DialContext(dialCtx, address,
		//lint:ignore SA1019 we want to ignore this for now
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(35000000),
			grpc.MaxCallSendMsgSize(35000000),
		))
	if err != nil {
		return nil, errors.Errorf("fail to dial: %w", err).WithField("address", address)
	}

	log.DD().WithContext(ctx).Infof("Connected to %s with max send & recv size 35 MB", address)

	conn := newClientConn(id, grpcConn)
	go func() {
		<-conn.Done()
		log.DD().WithContext(ctx).Debugf("Disconnected %s", grpcConn.Target())
	}()
	return conn, nil
}

// NewClient returns a new client instance.
func NewClient() *client {
	return &client{}
}
