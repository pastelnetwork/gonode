package conn

import (
	"io"
	"net"
)

const MaxBufferSize = 32768    // size which can store about 1024 frames of AES256
const MinBufferFrameSet = 1024 // can contains 32 frames of AES256

// Conn represents a secured connection.
// which implements the net.Conn interface.
type Conn struct {
	conn   net.Conn
	crypto Crypto
	buffer []byte
}

func NewClientConn(conn net.Conn) *Conn {
	var clientConn = Conn{
		conn:   conn,
		buffer: make([]byte, MaxBufferSize),
	}
	clientConn.crypto = NewClientCrypto(clientConn)
	return &clientConn
}

func (c *Conn) write(data []byte) (int, error) {
	n, err := c.conn.Write(data)
	return n, err
}

func (c *Conn) readRecord() (interface{}, error) {
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small buffer

	for {
		n, err := c.conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		buf = append(buf, tmp[:n]...)

	}
	if _, err := c.conn.Read(buf); err == nil {
		return nil, err
	}

	// trying to decrypt message
	switch buf[0] {
	case typeClientHello:
		return DecodeClientMsg(buf)
	case typeServerHello:
		return DecodeServerMsg(buf)
	case typeClientHandshakeMsg:
		return DecodeClientHandshakeMessage(buf)
	case typeServerHandshakeMsg:
		return DecodeServerHandshakeMsg(buf)
	default:
		return nil, WrongFormatErr
	}
}

func (c *Conn) Write(data []byte) (int, error) {
	if err := c.crypto.Handshake(); err != nil {
		return 0, err
	}

	encryptedData, err := c.crypto.Encrypt(data)
	if err != nil {
		return 0, err
	}

	return c.conn.Write(encryptedData)
}

func (c *Conn) Read(b []byte) (int, error) {
	if err := c.crypto.Handshake(); err != nil {
		return 0, err
	}

	expectedSize := cap(b)
	readed := 0
	var tmpBuffer = make([]byte, MinBufferFrameSet)

	for {
		n, err := c.conn.Read(tmpBuffer)
		if err != nil {
			if err != io.EOF {
				c.buffer = c.buffer[:0] // reset buffer in case of error
				return 0, err
			}
			break
		}
		readed += n
		c.buffer = append(c.buffer, tmpBuffer[:n]...)
	}

	// ToDo: how should we handle such case ?
	// readed != expectedSize

	decryptedData, err := c.crypto.Decrypt(c.buffer[:readed])
	if err != nil {
		return 0, err
	}

	if readed >= expectedSize {
		copy(decryptedData[:expectedSize], b)
	}

	return readed, nil
}
