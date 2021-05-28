package conn

import (
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"net"
)

func NewClient(conn net.Conn, transport transport.Transport) (net.Conn, error) {
	return transport.ClientHandshake(conn)
}
