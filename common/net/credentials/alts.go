package credentials

import (
	"context"
	"net"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts/handshaker"
	"google.golang.org/grpc/credentials"
)

const (
	// defaultTimeout specifies the server handshake timeout.
	defaultTimeout = 30 * time.Second
)

// ClientOptions contains the client-side options of an ALTS channel. These
// options will be passed to the underlying ALTS handshaker
type ClientOptions struct {
}

// DefaultClientOptions creates a new ClientOptions object with the default
// values.
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{}
}

// ServerOptions contains the server-side options of an ALTS channel. These
// options will be passed to the underlying ALTS handshaker.
type ServerOptions struct {
}

// DefaultServerOptions creates a new ServerOptions object with the default
// values.
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{}
}

// altsTC is the credentials required for authenticating a connection using ALTS.
// It implements credentials.TransportCredentials interface.
type altsTC struct {
	info     *credentials.ProtocolInfo
	side     alts.Side
	signInfo *alts.SignInfo
	sign     alts.Sign
	verify   alts.Verify
}

// NewClientCreds constructs a client-side ALTS TransportCredentials object.
func NewClientCreds(sign alts.Sign, verify alts.Verify, info *alts.SignInfo) credentials.TransportCredentials {
	return newALTS(alts.ClientSide, sign, verify, info)
}

// NewServerCreds constructs a server-side ALTS TransportCredentials object.
func NewServerCreds(sign alts.Sign, verify alts.Verify, info *alts.SignInfo) credentials.TransportCredentials {
	return newALTS(alts.ServerSide, sign, verify, info)
}

func newALTS(side alts.Side, sign alts.Sign, verify alts.Verify, info *alts.SignInfo) credentials.TransportCredentials {
	return &altsTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "alts",
			SecurityVersion:  "0.1",
		},
		side:     side,
		sign:     sign,
		signInfo: info,
		verify:   verify,
	}
}

// ClientHandshake implements the client side handshake protocol.
func (g *altsTC) ClientHandshake(ctx context.Context, _ string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	opts := handshaker.DefaultClientHandshakerOptions()
	chs, err := handshaker.NewClientHandshaker(ctx, rawConn, opts)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			chs.Close()
		}
	}()

	secureConn, authInfo, err := chs.ClientHandshake(ctx, g.sign, g.verify, g.signInfo)
	if err != nil {
		return nil, nil, err
	}

	return secureConn, authInfo, nil
}

// ServerHandshake implements the server side ALTS handshaker.
func (g *altsTC) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	opts := handshaker.DefaultServerHandshakerOptions()
	shs, err := handshaker.NewServerHandshaker(ctx, rawConn, opts)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			shs.Close()
		}
	}()

	secureConn, authInfo, err := shs.ServerHandshake(ctx, g.sign, g.verify, g.signInfo)
	if err != nil {
		return nil, nil, err
	}

	return secureConn, authInfo, nil
}

func (g *altsTC) Info() credentials.ProtocolInfo {
	return *g.info
}

func (g *altsTC) Clone() credentials.TransportCredentials {
	info := *g.info
	return &altsTC{
		info: &info,
		side: g.side,
	}
}

func (g *altsTC) OverrideServerName(serverNameOverride string) error {
	g.info.ServerName = serverNameOverride
	return nil
}
