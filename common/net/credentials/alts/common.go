package alts

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)

const (
	// ClientSide identifies the client in this communication.
	ClientSide Side = iota
	// ServerSide identifies the server in this communication.
	ServerSide
)

// Side identifies the party's role: client or server.
type Side int

// Handshaker defines a ALTS handshaker interface.
type Handshaker interface {
	// ClientHandshake starts and completes a client-side handshaking and
	// returns a secure connection and corresponding auth information.
	ClientHandshake(ctx context.Context, secClient SecClient, secInfo *SecInfo) (net.Conn, credentials.AuthInfo, error)
	// ServerHandshake starts and completes a server-side handshaking and
	// returns a secure connection and corresponding auth information.
	ServerHandshake(ctx context.Context, secClient SecClient, secInfo *SecInfo) (net.Conn, credentials.AuthInfo, error)
}

// ProtocolInfo provides information regarding the gRPC wire protocol version,
// security protocol, security protocol version in use, server name, etc.
type ProtocolInfo struct {
	// ProtocolVersion is the gRPC wire protocol version.
	ProtocolVersion string
	// SecurityProtocol is the security protocol in use.
	SecurityProtocol string
	// SecurityVersion is the security protocol version.  It is a static version string from the
	// credentials, not a value that reflects per-connection protocol negotiation.  To retrieve
	// details about the credentials used for a connection, use the Peer's AuthInfo field instead.
	//
	// Deprecated: please use Peer.AuthInfo.
	SecurityVersion string
	// ServerName is the user-configured server name.
	ServerName string
}

// SecClient defines a secure client interface
type SecClient interface {
	// Sign signs data by the given pastelID and passphrase, if successful returns signature.
	Sign(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error)
	// Verify verifies signed data by the given its signature and pastelID, if successful returns true.
	Verify(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error)
}

// SecInfo is information for signing
type SecInfo struct {
	PastelID   string
	PassPhrase string
	Algorithm  string
}
