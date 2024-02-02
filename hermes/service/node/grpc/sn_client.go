package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"google.golang.org/grpc"
)

type snClient struct {
	secClient alts.SecClient
	secInfo   *alts.SecInfo
}

// Connect implements node.Client.Connect()
func (client *snClient) Connect(ctx context.Context, address string) (node.ConnectionInterface, error) {
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
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100000000), grpc.MaxCallSendMsgSize(100000000)),
		)
	} else {
		grpcConn, err = grpc.DialContext(ctx, address,
			grpc.WithTransportCredentials(altsTCClient),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100000000), grpc.MaxCallSendMsgSize(100000000)),
		)
	}

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

// NewSNClient returns a new client instance.
func NewSNClient(secClient alts.SecClient, secInfo *alts.SecInfo) node.SNClientInterface {
	return &snClient{
		secClient: secClient,
		secInfo:   secInfo,
	}
}
