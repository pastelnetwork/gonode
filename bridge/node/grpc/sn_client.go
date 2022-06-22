package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/bridge/node"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/random"
	"google.golang.org/grpc"
)

type snclient struct {
	secClient alts.SecClient
}

// Connect implements node.Client.Connect()
func (client *snclient) Connect(ctx context.Context, address string, secInfo *alts.SecInfo) (node.ConnectionInterface, error) {
	id, _ := random.String(8, random.Base62Chars)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, id))

	altsTCClient := credentials.NewClientCreds(client.secClient, secInfo)
	var grpcConn *grpc.ClientConn
	var err error
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		grpcConn, err = grpc.DialContext(ctx, address,
			//lint:ignore SA1019 we want to ignore this for now
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)

	} else {
		grpcConn, err = grpc.DialContext(ctx, address,
			grpc.WithTransportCredentials(altsTCClient),
			grpc.WithBlock(),
		)
	}

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

// NewSNClient will wrap the input client in the SecClient interface, providing Signing and Verification
//	functionality.  By wrapping the return in ClientInterface, the resulting client will be able to call
//  nft registration, userdata, and sense stream functions.
func NewSNClient(secClient alts.SecClient) node.SNClientInterface {
	return &snclient{
		secClient: secClient,
	}
}
