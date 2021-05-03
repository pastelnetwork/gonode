package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// ClientConn represents custom grpc client connection.
type ClientConn struct {
	*grpc.ClientConn
	closedCh chan struct{}
}

// Done returns a channel that's closed when connection is shutdown.
func (conn *ClientConn) Done() <-chan struct{} {
	return conn.closedCh
}

func (conn *ClientConn) watchConnStatus() {
	defer close(conn.closedCh)

	for {
		conn.ClientConn.WaitForStateChange(context.Background(), conn.GetState())
		if conn.GetState() == connectivity.TransientFailure {
			conn.Close()
		}
		if conn.GetState() == connectivity.Shutdown {
			return
		}
	}
}

// NewClientConn returns a new ClientConn instance.
func NewClientConn(grpcConn *grpc.ClientConn) *ClientConn {
	conn := &ClientConn{
		ClientConn: grpcConn,
		closedCh:   make(chan struct{}),
	}
	go conn.watchConnStatus()
	return conn
}
