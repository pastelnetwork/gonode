package ed448

import (
	"encoding/binary"
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"net"
)

// conn represents a secured connection. It implements the net.Conn interface.
type conn struct {
	net.Conn
	crypto transport.Crypto
}

// NewConn creates a new secure channel instance
func NewConn(c net.Conn, crypto transport.Crypto) net.Conn {
	return &conn{
		Conn:   c,
		crypto: crypto,
	}
}

// Read reads and decrypts a frame from the underlying connection, and copies the
// decrypted payload into b
func (c *conn) Read(b []byte) (n int, err error) {
	var msgHeader = make([]byte, 5)
	if _, err := c.Conn.Read(msgHeader); err != nil {
		return 0, err
	}

	if msgHeader[0] != typeEncryptedMsg {
		return 0, ErrWrongFormat
	}

	expectedFrameSize := binary.LittleEndian.Uint32(msgHeader[1:])
	var frame = make([]byte, expectedFrameSize)
	var cryptoFrame = make([]byte, c.crypto.FrameSize())
	var readFrameSize uint32 = 0
	for readFrameSize != expectedFrameSize {
		read, err := c.Conn.Read(cryptoFrame)
		if err != nil {
			return 0, err
		}
		currentSize := readFrameSize + uint32(read)
		toRead := uint32(read)
		if currentSize > expectedFrameSize {
			toRead -= uint32(read) - (currentSize - expectedFrameSize)
		}
		copy(frame[readFrameSize:], cryptoFrame[:toRead])
		readFrameSize += uint32(read)
	}

	decryptedData, err := c.crypto.Decrypt(frame)
	if err != nil {
		return 0, err
	}
	n = copy(b, decryptedData)
	return n, nil
}

// Write encrypts, frames, and writes bytes from b to the underlying connection.
func (c *conn) Write(b []byte) (int, error) {
	size := len(b)
	// calculate msg size
	frameCount := size / c.crypto.FrameSize()
	remainder := size % c.crypto.FrameSize()
	if remainder > 0 {
		frameCount += 1
	}

	header := make([]byte, 5+frameCount*c.crypto.FrameSize())
	header[0] = typeEncryptedMsg
	binary.LittleEndian.PutUint32(header[1:], uint32(frameCount*c.crypto.FrameSize()))
	// encrypt data
	encrypted, err := c.crypto.Encrypt(b)
	if err != nil {
		return 0, err
	}
	if _, err := c.Conn.Write(header); err != nil {
		return 0, err
	}
	if _, err := c.Conn.Write(encrypted); err != nil {
		return 0, err
	}

	// return size of b
	return size, nil
}
