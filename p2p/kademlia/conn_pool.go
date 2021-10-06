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

type ConnectionItem struct {
	lastAccess time.Time
	conn       net.Conn
}

type ConnPool struct {
	capacity int
	conns    map[string]*ConnectionItem
	mtx      sync.Mutex
}

func NewConnPool() *ConnPool {
	return &ConnPool{
		capacity: defaultCapacity,
		conns:    map[string]*ConnectionItem{},
	}
}

func (pool *ConnPool) Add(addr string, conn net.Conn) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

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

	pool.conns[addr] = &ConnectionItem{
		lastAccess: time.Now(),
		conn:       conn,
	}
}

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

func (pool *ConnPool) Del(addr string) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()
	delete(pool.conns, addr)
}

func (pool *ConnPool) Release() {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	for addr, item := range pool.conns {
		item.conn.Close()
		delete(pool.conns, addr)
	}
}

type connWrapper struct {
	secureConn net.Conn
	rawConn    net.Conn
	mtx        sync.Mutex
}

func NewSecureClientConn(tpCredentials credentials.TransportCredentials, ctx context.Context, remoteAddr string) (net.Conn, error) {
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

func NewSecureServerConn(tpCredentials credentials.TransportCredentials, ctx context.Context, rawConn net.Conn) (net.Conn, error) {
	conn, _, err := tpCredentials.ServerHandshake(rawConn)
	if err != nil {
		return nil, errors.Errorf("server secure establish failed")
	}

	return &connWrapper{
		secureConn: conn,
		rawConn:    rawConn,
	}, nil
}

func (conn *connWrapper) Read(b []byte) (n int, err error) {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.secureConn.Read(b)
}

func (conn *connWrapper) Write(b []byte) (n int, err error) {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.secureConn.Write(b)
}

func (conn *connWrapper) Close() error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	conn.secureConn.Close()
	return conn.rawConn.Close()
}

func (conn *connWrapper) LocalAddr() net.Addr {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.LocalAddr()
}

func (conn *connWrapper) RemoteAddr() net.Addr {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.RemoteAddr()
}

func (conn *connWrapper) SetDeadline(t time.Time) error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.SetDeadline(t)
}

func (conn *connWrapper) SetReadDeadline(t time.Time) error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.SetReadDeadline(t)
}

func (conn *connWrapper) SetWriteDeadline(t time.Time) error {
	conn.mtx.Lock()
	defer conn.mtx.Unlock()
	return conn.rawConn.SetWriteDeadline(t)
}
