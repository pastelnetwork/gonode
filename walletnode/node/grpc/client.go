package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	logPrefix = "grpc"
)

type Client struct {
}

func (client *Client) Connect(ctx context.Context, address string) (node.Connection, error) {
	id, _ := random.String(8, random.Base62Chars)
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, id))

	ctxConnect, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	grpcConn, err := grpc.DialContext(ctxConnect, address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("address", address).Errorf("Could not connect to node")
		return nil, errors.Errorf("fail to dial: %v", err).WithField("address", address)
	}
	log.WithContext(ctx).Debugf("Connected to %s", address)

	conn := NewConnection(id, grpcConn)
	go conn.waitForStateChange()
	go func() {
		<-conn.Closed()
		log.WithContext(ctx).Debugf("Disconnected %s", conn.Target())
	}()
	return conn, nil
}

func NewClient() node.Client {
	return &Client{}
}

type Connection struct {
	*grpc.ClientConn
	pb.WalletNodeClient

	id       string
	closedCh chan struct{}
}

func (conn *Connection) waitForStateChange() {
	defer close(conn.closedCh)

	for {
		conn.WaitForStateChange(context.Background(), conn.GetState())
		if conn.GetState() == connectivity.TransientFailure {
			conn.Close()
		}
		if conn.GetState() == connectivity.Shutdown {
			return
		}
	}
}

func (conn *Connection) Closed() <-chan struct{} {
	return conn.closedCh
}

func (conn *Connection) RegisterArtowrk(ctx context.Context) (node.RegisterArtowrk, error) {
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, conn.id))

	stream, err := conn.WalletNodeClient.RegisterArtowrk(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("address", conn.Target()).Errorf("Could not start stream")
		return nil, errors.New(err)
	}

	regArtwork := NewRegisterArtowrk(conn, stream)
	if err := regArtwork.Start(ctx); err != nil {
		return nil, err
	}

	return regArtwork, nil
}

func NewConnection(id string, conn *grpc.ClientConn) *Connection {
	return &Connection{
		ClientConn:       conn,
		WalletNodeClient: pb.NewWalletNodeClient(conn),
		id:               id,
		closedCh:         make(chan struct{}),
	}
}
