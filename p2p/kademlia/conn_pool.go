package kademlia

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/anacrolix/utp"
	"github.com/pastelnetwork/gonode/common/errors"
	"google.golang.org/grpc/credentials"
)

const defaultCapacity = 128

type connectionItem struct {
	lastAccess time.Time
	conn       net.Conn
}

// ConnPool is a manager of connection pool
type ConnPool struct {
	capacity int
	conns    map[string]*connectionItem
	mtx      sync.Mutex
}

// NewConnPool return a connection pool
func NewConnPool() *ConnPool {
	return &ConnPool{
		capacity: defaultCapacity,
		conns:    map[string]*connectionItem{},
	}
}

// Add a connection to pool
func (pool *ConnPool) Add(addr string, conn net.Conn) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	// if connection not in pool
	if _, ok := pool.conns[addr]; !ok {
		// if pool is full
		if len(pool.conns) >= pool.capacity {
			oldestAccess := time.Now()
			oldestAccessAddr := ""

			for addr, item := range pool.conns {
				if item.lastAccess.Before(oldestAccess) {
					oldestAccessAddr = addr
					oldestAccess = item.lastAccess
				}
			}

			delete(pool.conns, oldestAccessAddr)
		}
	}

	pool.conns[addr] = &connectionItem{
		lastAccess: time.Now(),
		conn:       conn,
	}
}

// Get return a connection from pool
func (pool *ConnPool) Get(addr string) (net.Conn, error) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()
	item, ok := pool.conns[addr]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	item.lastAccess = time.Now()
	return item.conn, nil
}

// Del remove a connection from pool
func (pool *ConnPool) Del(addr string) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()
	delete(pool.conns, addr)
}

// Release all connections in pool - used when exits
func (pool *ConnPool) Release() {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	for addr, item := range pool.conns {
		item.conn.Close()
		delete(pool.conns, addr)
	}
}

// connWrapper implements wrapper of secure connection
type connWrapper struct {
	secureConn net.Conn
	rawConn    net.Conn
	mtx        sync.Mutex
}

// NewSecureClientConn do client handshake and return a secure connection
func NewSecureClientConn(ctx context.Context, tpCredentials credentials.TransportCredentials, remoteAddr string) (net.Conn, error) {
	// dial the remote address with udp network
	rawConn, err := utp.DialContext(ctx, remoteAddr)
	if err != nil {
		return nil, errors.Errorf("dial %q: %w", remoteAddr, err)
	}

	// set the deadline for read and write
	rawConn.SetDeadline(time.Now().Add(defaultConnDeadline))

	conn, _, err := tpCredentials.ClientHandshake(ctx, "", rawConn)
	if err != nil {
		rawConn.Close()
		return nil, errors.Errorf("client secure establish %q: %w", remoteAddr, err)
	}

	return &connWrapper{
		secureConn: conn,
		rawConn:    rawConn,
	}, nil
}

// NewSecureServerConn do server handshake and create a secure connection
func NewSecureServerConn(_ context.Context, tpCredentials credentials.TransportCredentials, rawConn net.Conn) (net.Conn, error) {
	conn, _, err := tpCredentials.ServerHandshake(rawConn)
	if err != nil {
		return nil, errors.Errorf("server secure establish failed")
	}

	return &connWrapper{
		secureConn: conn,
		rawConn:    rawConn,
	}, nil
}

// Read implements net.Conn's Read interface
func (conn *connWrapper) Read(b []byte) (n int, err error) {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.secureConn.Read(b)
}

// Write implements net.Conn's Write interface
func (conn *connWrapper) Write(b []byte) (n int, err error) {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.secureConn.Write(b)
}

// Close implements net.Conn's Close interface
func (conn *connWrapper) Close() error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	conn.secureConn.Close()
	return conn.rawConn.Close()
}

// LocalAddr implements net.Conn's LocalAddr interface
func (conn *connWrapper) LocalAddr() net.Addr {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.LocalAddr()
}

// RemoteAddr implements net.Conn's RemoteAddr interface
func (conn *connWrapper) RemoteAddr() net.Addr {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.RemoteAddr()
}

// SetDeadline implements net.Conn's SetDeadline interface
func (conn *connWrapper) SetDeadline(t time.Time) error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.SetDeadline(t)
}

// SetReadDeadline implements net.Conn's SetReadDeadline interface
func (conn *connWrapper) SetReadDeadline(t time.Time) error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.SetReadDeadline(t)
}

// SetWriteDeadline implements net.Conn's SetWriteDeadline interface
func (conn *connWrapper) SetWriteDeadline(t time.Time) error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.SetWriteDeadline(t)
}
