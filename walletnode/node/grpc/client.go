package grpc

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc"
)

const (
	logPrefix = "client"
)

type client struct {
	secClient alts.SecClient
}

// Connect implements node.Client.Connect()
func (client *client) Connect(ctx context.Context, address string, secInfo *alts.SecInfo) (node.Connection, error) {
	id, _ := random.String(8, random.Base62Chars)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, id))

	altsTCClient := credentials.NewClientCreds(client.secClient, secInfo)
	grpcConn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(altsTCClient),
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
func NewClient(secClient alts.SecClient) node.Client {
	return &client{
		secClient: secClient,
	}
}
