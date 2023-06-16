package client

import (
	"context"
	"fmt"
	"os"

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
func (client *client) Connect(ctx context.Context, address string) (node.ConnectionInterface, error) {
	id, _ := random.String(8, random.Base62Chars)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, id))

	if client.secClient == nil || client.secInfo == nil {
		return nil, errors.Errorf("secClient or secInfo don't initialize")
	}

	altsTCClient := credentials.NewClientCreds(client.secClient, client.secInfo)
	var grpcConn *grpc.ClientConn
	var err error
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		grpcConn, err = grpc.DialContext(ctx, address,
			//lint:ignore SA1019 we want to ignore this for now
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(35000000), grpc.MaxCallSendMsgSize(35000000)),
		)
	} else {
		grpcConn, err = grpc.DialContext(ctx, address,
			grpc.WithTransportCredentials(altsTCClient),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(35000000), grpc.MaxCallSendMsgSize(35000000)),
		)
	}

	if err != nil {
		log.WithContext(ctx).WithError(err).Error("DialContext")
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
func New(secClient alts.SecClient, secInfo *alts.SecInfo) node.ClientInterface {
	return &client{
		secClient: secClient,
		secInfo:   secInfo,
	}
}
