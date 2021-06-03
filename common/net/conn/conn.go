package conn

import (
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"net"
)

// Side identifies the  role: client or server.
type Side = uint8

const (
	// Server side role
	Server Side = iota
	// Client side role
	Client
)

// New creates an encrypted connection instance that uses transport during client or server handshake.
func New(conn net.Conn, side Side, transport transport.Transport) (net.Conn, error) {
	if side == Client {
		return transport.ClientHandshake(conn)
	}
	return transport.ServerHandshake(conn)
}
