package kademlia

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const defaultCapacity = 128

type ConnectionItem struct {
	lastAccess time.Time
	conn       net.Conn
	rawConn    net.Conn
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

func (pool *ConnPool) Add(addr string, conn net.Conn, rawConn net.Conn) {
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
		rawConn:    rawConn,
	}
}

func (pool *ConnPool) Get(addr string) (net.Conn, net.Conn, error) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()
	item, ok := pool.conns[addr]
	if !ok {
		return nil, nil, fmt.Errorf("not found")
	}

	item.lastAccess = time.Now()
	return item.conn, item.rawConn, nil
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
		item.rawConn.Close()
		item.conn.Close()
		delete(pool.conns, addr)
	}
}
