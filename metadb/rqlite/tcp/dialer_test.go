package tcp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/metadb/rqlite/testdata/x509"
	"golang.org/x/sync/errgroup"
)

func Test_NewDialer(t *testing.T) {
	d := NewDialer(1, false, false)
	if d == nil {
		t.Fatal("failed to create a dialer")
	}
}

func Test_DialerNoConnect(t *testing.T) {
	d := NewDialer(87, false, false)
	_, err := d.Dial("127.0.0.1:0", 5*time.Second)
	if err == nil {
		t.Fatalf("no error connecting to bad address")
	}
}

func Test_DialerHeader(t *testing.T) {
	s := mustNewEchoServer()
	defer s.Close()
	go s.Start()

	d := NewDialer(64, false, false)
	conn, err := d.Dial(s.Addr(), 10*time.Second)
	if err != nil {
		t.Fatalf("failed to dial echo server: %s", err.Error())
	}

	buf := make([]byte, 1)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("failed to read from echo server: %s", err.Error())
	}
	if exp, got := buf[0], byte(64); exp != got {
		t.Fatalf("got wrong response from echo server, exp %d, got %d", exp, got)
	}
}

func Test_DialerHeaderTLS(t *testing.T) {
	s, cert, key := mustNewEchoServerTLS()
	defer s.Close()
	defer os.Remove(cert)
	defer os.Remove(key)
	go s.Start()

	d := NewDialer(23, true, true)
	conn, err := d.Dial(s.Addr(), 5*time.Second)
	if err != nil {
		t.Fatalf("failed to dial TLS echo server: %s", err.Error())
	}

	buf := make([]byte, 1)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("failed to read from TLS echo server: %s", err.Error())
	}
	if exp, got := buf[0], byte(23); exp != got {
		t.Fatalf("got wrong response from TLS echo server, exp %d, got %d", exp, got)
	}
}

func Test_DialerHeaderTLSBadConnect(t *testing.T) {
	s, cert, key := mustNewEchoServerTLS()
	defer s.Close()
	defer os.Remove(cert)
	defer os.Remove(key)
	go s.Start()

	// Connect to a TLS server with a unencrypted client, to make sure
	// code can handle that misconfig.
	d := NewDialer(56, false, false)
	conn, err := d.Dial(s.Addr(), 5*time.Second)
	if err != nil {
		t.Fatalf("failed to dial TLS echo server: %s", err.Error())
	}

	buf := make([]byte, 1)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Read(buf)
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("unexpected error from TLS echo server: %s", err.Error())
	}
}

type echoServer struct {
	ln net.Listener
}

// Addr returns the address of the echo server.
func (e *echoServer) Addr() string {
	return e.ln.Addr().String()
}

// Start starts the echo server.
func (e *echoServer) Start() {
	for {
		conn, err := e.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// Successful testing can cause this.
				return
			}

			panic(fmt.Sprintf("failed to accept a connection: %s", err.Error()))
		}
		group, _ := errgroup.WithContext(context.Background())
		group.Go(func() error {
			c := conn
			buf := make([]byte, 1)
			_, err := c.Read(buf)
			if err != nil {
				return nil
			}
			_, err = c.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to write byte: %s", err.Error())
			}
			c.Close()

			return nil
		})

		if err := group.Wait(); err != nil {
			panic(err)
		}
	}
}

// Close closes the echo server.
func (e *echoServer) Close() {
	e.ln.Close()
}

func mustNewEchoServer() *echoServer {
	return &echoServer{
		ln: mustTCPListener("127.0.0.1:0"),
	}
}

func mustNewEchoServerTLS() (*echoServer, string, string) {
	ln := mustTCPListener("127.0.0.1:0")
	cert := x509.CertFile("")
	key := x509.KeyFile("")

	tlsConfig, err := createTLSConfig(cert, key, "")
	if err != nil {
		panic("failed to create TLS config")
	}

	return &echoServer{
		ln: tls.NewListener(ln, tlsConfig),
	}, cert, key
}
