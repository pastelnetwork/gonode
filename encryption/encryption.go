package encryption

import (
	"io"
)

type Handshaker interface {
	ClientHello() ([]byte, error)           // returns artistID
	ServerHello() ([]byte, error)           // returns OK
	ClientKeyExchange() ([]byte, error)     // returns signature+publicKey
	ServerKeyVerify([]byte) (Cipher, error) // verify client's signature
	ServerKeyExchange() ([]byte, error)     // returns signature+publicKey
	ClientKeyVerify([]byte) (Cipher, error) // verify serverMock's signature
}

type Cipher interface {
	WrapWriter(writer io.Writer) (io.Writer, error)
	WrapReader(reader io.Reader) (io.Reader, error)
}