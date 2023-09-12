package grpc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/random"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

const (
	logPrefix = "walletnode-grpc-secclient"
)

type client struct {
	secClient alts.SecClient
}

// Connect implements node.Client.Connect()
func (client *client) Connect(ctx context.Context, address string, secInfo *alts.SecInfo) (node.ConnectionInterface, error) {
	grpclog.SetLoggerV2(log.DefaultLogger)
	id, _ := random.String(8, random.Base62Chars)
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, id))

	// Define the keep-alive parameters
	ka := keepalive.ClientParameters{
		Time:                120 * time.Second, // Send pings every 120 seconds if there is no activity
		Timeout:             15 * time.Second,  // Wait 5 second for ping ack before considering the connection dead
		PermitWithoutStream: true,              // Allow pings to be sent without a stream
	}

	altsTCClient := credentials.NewClientCreds(client.secClient, secInfo)
	var grpcConn *grpc.ClientConn
	var err error
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		grpcConn, err = grpc.DialContext(ctx, address,
			//lint:ignore SA1019 we want to ignore this for now
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(35000000), grpc.MaxCallSendMsgSize(35000000)),
			grpc.WithKeepaliveParams(ka),
		)

	} else {
		grpcConn, err = grpc.DialContext(ctx, address,
			grpc.WithTransportCredentials(altsTCClient),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(35000000), grpc.MaxCallSendMsgSize(35000000)),
			grpc.WithKeepaliveParams(ka),
		)
	}

	if err != nil {
		return nil, errors.Errorf("fail to dial: %w", err).WithField("address", address)
	}
	log.WithContext(ctx).Infof("Connected to %s", address)

	conn := newClientConn(id, grpcConn)
	go func() {
		<-conn.Done()
		log.WithContext(ctx).Infof("Disconnected %s", grpcConn.Target())
	}()
	return conn, nil
}

// NewClient will wrap the input client in the SecClient interface, providing Signing and Verification
//
//		functionality.  By wrapping the return in ClientInterface, the resulting client will be able to call
//	 nft registration, userdata, and sense stream functions.
func NewClient(secClient alts.SecClient) node.ClientInterface {
	return &client{
		secClient: secClient,
	}
}
