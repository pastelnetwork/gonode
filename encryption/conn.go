package encryption

import (
	"io"
	"net"
	"time"
)

type conn struct {
	cypher          Cipher
	c               net.Conn
	encryptedWriter io.Writer
	encryptedReader io.Reader
}

func (c *conn) Read(b []byte) (n int, err error) {
	return c.encryptedReader.Read(b)
}

func (c *conn) Write(b []byte) (n int, err error) {
	return c.encryptedWriter.Write(b)
}

func (c *conn) Close() error {
	return c.c.Close()
}

func (c *conn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c *conn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

func NewConn(rawConn net.Conn, cypher Cipher) (net.Conn, error) {
	reader, err := cypher.WrapReader(rawConn)
	if err != nil{
		return nil, err
	}
	writer, err := cypher.WrapWriter(rawConn)
	if err != nil{
		return nil, err
	}
	return &conn{
		cypher:         cypher,
		c:              rawConn,
		encryptedReader: reader,
		encryptedWriter: writer,
	}, nil
}