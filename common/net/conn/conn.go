package conn

import (
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"net"
)

type Side = uint8

const (
	Server Side = iota
	Client
)

func New(conn net.Conn, side Side, transport transport.Transport) (net.Conn, error) {
	if side == Client {
		return transport.ClientHandshake(conn)
	}
	return transport.ServerHandshake(conn)
}
