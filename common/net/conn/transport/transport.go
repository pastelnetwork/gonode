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
	Configure(rawData []byte) error
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
	GetConfiguration() []byte
	FrameSize() int
	About() string
}
