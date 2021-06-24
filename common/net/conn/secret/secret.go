package secret

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Config structure is used to configure a secret client or server
type Config struct {
}

// A listener implements a network listener (net.Listener) for secret connections.
type listener struct {
	net.Listener
	config *Config
}

// Accept returns the next incoming secret connection
func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return Server(c, l.config), nil
}

// Server returns a new secret server side connection which
// use conn as the underlying transport
func Server(conn net.Conn, config *Config) *Conn {
	s := &Conn{
		conn:   conn,
		config: config,
		secret: make([]byte, AESSecretKeySize),
	}
	s.handshakeFn = s.serverHandshake

	return s
}

// Client returns a new secret client side connection which
// use conn as the underlying transport
func Client(conn net.Conn, config *Config) *Conn {
	s := &Conn{
		conn:   conn,
		config: config,
		secret: make([]byte, AESSecretKeySize),
	}
	s.handshakeFn = s.clientHandshake

	return s
}

// NewListener creates a Listener with accepts connections from an inner Listener
// and wraps each connection with server
func NewListener(inner net.Listener, config *Config) net.Listener {
	return &listener{
		Listener: inner,
		config:   config,
	}
}

// Listen creates a secret listener accepting connections on the network address
func Listen(network, laddr string, config *Config) (net.Listener, error) {
	l, err := net.Listen(network, laddr)
	if err != nil {
		return nil, err
	}
	return NewListener(l, config), nil
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "secret: DialWithDialer timed out" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

// DialWithDialer connects to the given network address using dialer.Dial and
// then initiates a secret handshake, returning the resulting secret connection. Any
// timeout or deadline given in the dialer apply to connection and secret
// handshake as a whole.
//
// DialWithDialer interprets a nil configuration as equivalent to the zero
// configuration; see the documentation of Config for the defaults.
func DialWithDialer(dialer *net.Dialer, network, addr string, config *Config) (*Conn, error) {
	return dial(context.Background(), dialer, network, addr, config)
}

func dial(ctx context.Context, netDialer *net.Dialer, network, addr string, config *Config) (*Conn, error) {
	// We want the Timeout and Deadline values from dialer to cover the
	// whole process: TCP connection and secret handshake. This means that we
	// also need to start our own timers now.
	timeout := netDialer.Timeout

	if !netDialer.Deadline.IsZero() {
		deadlineTimeout := time.Until(netDialer.Deadline)
		if timeout == 0 || deadlineTimeout < timeout {
			timeout = deadlineTimeout
		}
	}

	// hsErrCh is non-nil if we might not wait for Handshake to complete
	var hsErrCh chan error
	if timeout != 0 || ctx.Done() != nil {
		hsErrCh = make(chan error, 2)
	}
	if timeout != 0 {
		timer := time.AfterFunc(timeout, func() {
			hsErrCh <- timeoutError{}
		})
		defer timer.Stop()
	}

	rawConn, err := netDialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("dial context: %w", err)
	}
	conn := Client(rawConn, config)

	if hsErrCh == nil {
		err = conn.Handshake()
	} else {
		go func() {
			hsErrCh <- conn.Handshake()
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-hsErrCh:
			if err != nil {
				// If the error was due to the context
				// closing, prefer the context's error, rather
				// than some random network teardown error.
				if e := ctx.Err(); e != nil {
					err = e
				}
			}
		}
	}
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	return conn, nil
}

// Dial connects to the given network address using net.Dial
// and then initiates a secret handshake, returning the resulting
// secret connection.
// Dial interprets a nil configuration as equivalent to
// the zero configuration; see the documentation of Config
// for the defaults.
func Dial(network, addr string, config *Config) (*Conn, error) {
	return DialWithDialer(new(net.Dialer), network, addr, config)
}

// Dialer dials secret connections given a configuration and a Dialer for the
// underlying connection.
type Dialer struct {
	// NetDialer is the optional dialer to use for the secret connections'
	// underlying TCP connections.
	// A nil NetDialer is equivalent to the net.Dialer zero value.
	NetDialer *net.Dialer

	// Config is the secret configuration to use for new connections.
	// A nil configuration is equivalent to the zero configuration
	// defaults.
	Config *Config
}

// Dial connects to the given network address and initiates a secret
// handshake, returning the resulting secret connection.
//
// The returned Conn, if any, will always be of type *Conn.
func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

func (d *Dialer) netDialer() *net.Dialer {
	if d.NetDialer != nil {
		return d.NetDialer
	}
	return new(net.Dialer)
}

// DialContext connects to the given network address and initiates a secret
// handshake, returning the resulting secret connection.
//
// The provided Context must be non-nil. If the context expires before
// the connection is complete, an error is returned. Once successfully
// connected, any expiration of the context will not affect the
// connection.
//
// The returned Conn, if any, will always be of type *Conn.
func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	c, err := dial(ctx, d.netDialer(), network, addr, d.Config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
