package net

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/core"

	"net"
	"time"

	"google.golang.org/grpc/credentials"
)

type X446Info struct{}

// AuthType returns the type of TLSInfo as a string.
func (t X446Info) AuthType() string {
	return "x66"
}

type TransportCredentials struct {
	h Handshaker
}

func (t TransportCredentials) ClientHandshake(ctx context.Context, authority string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, err := NewConn(conn, core.SideClient, t.h)
	fmt.Printf("cleint handshake %s\n", err)
	if err != nil {
		return nil, nil, err
	}
	return conn, X446Info{}, nil
}

func (t TransportCredentials) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, err := NewConn(conn, core.SideServer, t.h)
	fmt.Printf("server handshake, %s\n", err)

	if err != nil {
		return nil, nil, err
	}
	time.Sleep(1 * time.Second)
	return conn, X446Info{}, nil
}

func (t TransportCredentials) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "x443",
		SecurityVersion:  "1.0",
		ServerName:       "",
	}
}

func (t TransportCredentials) Clone() credentials.TransportCredentials {
	return NewTransportCredentials(t.h)
}

func (t TransportCredentials) OverrideServerName(serverName string) error {
	// TODO: maybe implement works with servernames ?
	return nil
}

func NewTransportCredentials(h Handshaker) credentials.TransportCredentials {
	return TransportCredentials{
		h: h,
	}
}
