package transport

import (
	"net"
)

// Transport defines the common interface for all the handshake protocols
type Transport interface {
	ClientHandshake(net.Conn) (net.Conn, error)
	ServerHandshake(net.Conn) (net.Conn, error)
	GetEncryptionInfo() string
	Clone() Transport
}

// Crypto defines the common interface for encryption
type Crypto interface {
	Configure(rawData []byte) error
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
	GetConfiguration() []byte
	EncryptionOverhead() int
	Clone() (Crypto, error)
	About() string
}
