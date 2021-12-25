package client

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc"
)

const (
	logPrefix = "client"
)

type client struct {
	secClient alts.SecClient
	secInfo   *alts.SecInfo
}

// Connect implements node.Client.Connect()
func (client *client) Connect(ctx context.Context, address string) (node.Connection, error) {
	id, _ := random.String(8, random.Base62Chars)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, id))

	if client.secClient == nil || client.secInfo == nil {
		return nil, errors.Errorf("secClient or secInfo don't initialize")
	}
	altsTCClient := credentials.NewClientCreds(client.secClient, client.secInfo)
	grpcConn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(altsTCClient),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Errorf("dial address %s: %w", address, err)
	}
	log.WithContext(ctx).Debugf("Connected to %s", address)

	conn := newClientConn(id, grpcConn)
	go func() {
		<-conn.Done()
		log.WithContext(ctx).Debugf("Disconnected %s", grpcConn.Target())
	}()
	return conn, nil
}

// New returns a new client instance.
func New(secClient alts.SecClient, secInfo *alts.SecInfo) node.Client {
	return &client{
		secClient: secClient,
		secInfo:   secInfo,
	}
}
