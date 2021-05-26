package conn

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
	Init(...[]byte) error
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
	GenerateEncryptionParams() [][]byte
	FrameSize() int
	About() string
}
