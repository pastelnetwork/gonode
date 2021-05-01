package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/random"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Client struct {
}

func (client *Client) Connect(ctx context.Context, address string) (node.Connection, error) {
	id, _ := random.String(8, random.Base62Chars)
	ctx = context.WithValue(ctx, log.PrefixKey, fmt.Sprintf("%s-%s", logPrefix, id))

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	grpcConn, err := grpc.DialContext(ctx, address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Errorf("fail to dial: %v", err).WithField("address", address)
	}
	log.WithContext(ctx).Debugf("Connected to %s", address)

	conn := NewConnection(id, grpcConn)
	go conn.waitForStateChange()
	go func() {
		<-conn.Done()
		log.WithContext(ctx).Debugf("Disconnected %s", conn.Target())
	}()
	return conn, nil
}

func NewClient() node.Client {
	return &Client{}
}

type Connection struct {
	*grpc.ClientConn
	pb.SuperNodeClient

	id     string
	doneCh chan struct{}
}

func (conn *Connection) waitForStateChange() {
	defer close(conn.doneCh)

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

func (conn *Connection) Done() <-chan struct{} {
	return conn.doneCh
}

func (conn *Connection) RegisterArtowrk(ctx context.Context) (node.RegisterArtowrk, error) {
	stream, err := conn.SuperNodeClient.RegisterArtowrk(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	return NewRegisterArtowrk(conn, stream), nil
}

func NewConnection(id string, conn *grpc.ClientConn) *Connection {
	return &Connection{
		ClientConn:      conn,
		SuperNodeClient: pb.NewSuperNodeClient(conn),
		id:              id,
		doneCh:          make(chan struct{}),
	}
}
