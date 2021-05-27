package transport

import (
	"context"
	"net"
)

type Transport interface {
	ClientHandshake(context.Context, net.Conn) (net.Conn, error)
	ServerHandshake(context.Context, net.Conn) (net.Conn, error)
	IsHandshakeEstablished() bool
}

type Crypto interface {
	Configure(string) error
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
	GetConfiguration() string
	FrameSize() int
	About() string
}
