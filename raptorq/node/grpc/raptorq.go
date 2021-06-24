package grpc

import (
	pb "github.com/pastelnetwork/gonode/raptorq"
	"github.com/pastelnetwork/gonode/raptorq/node"
)

type raptorQ struct {
	conn   *clientConn
	client pb.RaptorQClient
}

func newRaptorq(conn *clientConn) node.RaptorQ {
	return &raptorQ{
		conn:   conn,
		client: pb.NewRaptorQClient(conn),
	}
}
