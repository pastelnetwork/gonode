package transport

import (
	"net"
)

type Transport interface {
	ClientHandshake(net.Conn) (net.Conn, error)
	ServerHandshake(net.Conn) (net.Conn, error)
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
