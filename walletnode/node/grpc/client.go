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
	logPrefix = "walletnode-grpc-secclient"
)

type client struct {
	secClient alts.SecClient
}

// Connect implements node.Client.Connect()
func (client *client) Connect(ctx context.Context, address string, secInfo *alts.SecInfo) (node.ConnectionInterface, error) {
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

// NewClient will wrap the input client in the SecClient interface, providing Signing and Verification
//	functionality.  By wrapping the return in ClientInterface, the resulting client will be able to call
//  artwork registration, userdata, and sense stream functions.
func NewClient(secClient alts.SecClient) node.ClientInterface {
	return &client{
		secClient: secClient,
	}
}
