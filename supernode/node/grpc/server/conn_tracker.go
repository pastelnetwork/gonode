package server

import (
	"net"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

type connTrackListener struct {
	net.Listener
	connTrack *ConnectionTracker
}

// Accept wraps the Listener.Accept method and adds the connection reference to the map.
func (l *connTrackListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	conn := &TrackedConn{
		Conn:      c,
		connTrack: l.connTrack,
		OpenedAt:  time.Now().UTC(),
	}
	log.WithField("con", c.RemoteAddr()).Info("new connection request received")

	l.connTrack.add(c.RemoteAddr().String(), conn)

	return conn, nil
}

// TrackedConn represents an open conn
type TrackedConn struct {
	net.Conn  `json:"-"`
	connTrack *ConnectionTracker `json:"-"`
	OpenedAt  time.Time          `json:"opened_at"`
}

// Close wraps the Conn.Close method and removes the connection reference from the map.
func (c *TrackedConn) Close() error {
	err := c.Conn.Close()
	if err != nil {
		log.Warnf("Error closing connection: %v, from: %v\n", err, c.RemoteAddr())
	}
	c.connTrack.remove(c.LocalAddr().String(), c)
	return err
}

// ConnectionTracker keeps track of open connection
type ConnectionTracker struct {
	mu    sync.Mutex
	conns map[string][]*TrackedConn
}

// NewConnectionTracker initializes the connection tracker map.
func NewConnectionTracker() *ConnectionTracker {
	return &ConnectionTracker{
		conns: make(map[string][]*TrackedConn),
	}
}

// add inserts a new connection reference into the in-memory map.
func (ct *ConnectionTracker) add(serverAddr string, c *TrackedConn) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	ct.conns[serverAddr] = append(ct.conns[serverAddr], c)
	log.Infof("New connection: %v, total: %d for server: %s\n", c.RemoteAddr(), len(ct.conns[serverAddr]), serverAddr)
}

// remove deletes a connection reference from the map.
func (ct *ConnectionTracker) remove(serverAddr string, c *TrackedConn) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	connsForServer, ok := ct.conns[serverAddr]
	if !ok {
		return
	}

	for i, conn := range connsForServer {
		if conn == c {
			ct.conns[serverAddr] = append(connsForServer[:i], connsForServer[i+1:]...)
			break
		}
	}

	if len(ct.conns[serverAddr]) == 0 {
		delete(ct.conns, serverAddr)
	}

	log.Infof("Closed connection: %v, total: %d for server: %s\n", c.RemoteAddr(), len(ct.conns[serverAddr]), serverAddr)
}
