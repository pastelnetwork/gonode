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
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

const (
	logPrefix = "wn-grpc"
)

type client struct {
	secClient alts.SecClient
}

// Connect implements node.Client.Connect()
func (client *client) Connect(ctx context.Context, address string, secInfo *alts.SecInfo, taskID string) (node.ConnectionInterface, error) {
	grpclog.SetLoggerV2(log.NewLoggerWithErrorLevel())
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, taskID))

	// Define the keep-alive parameters
	ka := keepalive.ClientParameters{
		Time:                30 * time.Minute, // Send pings every 5 minutes  if there is no activity
		Timeout:             30 * time.Minute, // Wait 1 minute for ping ack before considering the connection dead
		PermitWithoutStream: true,             // Allow pings to be sent without a stream
	}

	altsTCClient := credentials.NewClientCreds(client.secClient, secInfo)
	var grpcConn *grpc.ClientConn
	var err error
	if os.Getenv("INTEGRATION_TEST_ENV") == "true" {
		grpcConn, err = grpc.DialContext(ctx, address,
			//lint:ignore SA1019 we want to ignore this for now
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100000000), grpc.MaxCallSendMsgSize(100000000)),
			grpc.WithKeepaliveParams(ka),
		)

	} else {
		grpcConn, err = grpc.DialContext(ctx, address,
			grpc.WithTransportCredentials(altsTCClient),
			grpc.WithBlock(),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100000000), grpc.MaxCallSendMsgSize(100000000)),
			grpc.WithKeepaliveParams(ka),
		)
	}

	if err != nil {
		return nil, errors.Errorf("fail to dial: %w", err).WithField("address", address)
	}
	log.WithContext(ctx).Debugf("Connected to %s", address)

	conn := newClientConn(taskID, grpcConn)
	go func() {
		<-conn.Done()
		log.WithContext(ctx).Debugf("Disconnected %s", grpcConn.Target())
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
