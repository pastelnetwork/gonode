package transport

import (
	"context"
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"github.com/pastelnetwork/gonode/common/net/conn/transport/ed448"
	"google.golang.org/grpc/credentials"
	"net"
)

// Ed448TransportCredentials structure that implements TransportCredentials
// protocol by using ed448 handshake protocol
type Ed448TransportCredentials struct {
	transport transport.Transport
	protocol  credentials.ProtocolInfo
}

// New creates Ed448TransportCredentials
func New(crypto ...transport.Crypto) credentials.TransportCredentials {
	return &Ed448TransportCredentials{
		transport: ed448.New(crypto...),
		protocol: credentials.ProtocolInfo{
			ProtocolVersion: "ed448",
		},
	}
}

// ClientHandshake does the authentication handshake specified by the
// corresponding authentication protocol on rawConn for clients. It returns
// the authenticated connection and the corresponding auth information about the connection
func (transportCredentials *Ed448TransportCredentials) ClientHandshake(_ context.Context, _ string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, err := transportCredentials.transport.ClientHandshake(rawConn)
	transportCredentials.protocol.SecurityProtocol = transportCredentials.transport.GetEncryptionInfo()
	if err != nil {
		return nil, nil, err
	}

	return conn, newAuth(), nil
}

// ServerHandshake does the authentication handshake for servers. It returns
// the authenticated connection and the corresponding auth information about
// the connection. The auth information should embed CommonAuthInfo to return additional information
// about the credentials
func (transportCredentials *Ed448TransportCredentials) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, err := transportCredentials.transport.ServerHandshake(rawConn)
	transportCredentials.protocol.SecurityProtocol = transportCredentials.transport.GetEncryptionInfo()
	if err != nil {
		return nil, nil, err
	}

	return conn, newAuth(), nil
}

// Info provides the ProtocolInfo of this TransportCredentials.
func (transportCredentials *Ed448TransportCredentials) Info() credentials.ProtocolInfo {
	return transportCredentials.protocol
}

// Clone makes a copy of this TransportCredentials.
func (transportCredentials *Ed448TransportCredentials) Clone() credentials.TransportCredentials {
	return &Ed448TransportCredentials{
		transport: transportCredentials.transport.Clone(),
	}
}

// OverrideServerName overrides the server name used to verify the hostname on the returned certificates from the server.
// gRPC internals also use it to override the virtual hosting name if it is set.
// It must be called before dialing. Currently, this is only used by grpclb.
func (transportCredentials *Ed448TransportCredentials) OverrideServerName(serverName string) error {
	transportCredentials.protocol.ServerName = serverName
	return nil
}
