package secret

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	recordHeaderLen = 4                // record header length
	maxPayloadSize  = 16*1024*1024 - 1 // maximum payload size we support (16 MB)
)

// secret record types
type recordType uint8

const (
	clientHandshake recordType = 22
	serverHandshake recordType = 23
	applicationData recordType = 24
)

// Conn represents a secured connection, implements the net.Conn interface
type Conn struct {
	conn            net.Conn
	config          *Config
	handshakeFn     func() error // (*Conn).clientHandshake or serverHandshake
	handshakeStatus uint32       // handshakeStatus is 1 if the connection is currently transferring data
	handshakeMutex  sync.Mutex   // constant after handshake; protected by handshakeMutex
	handshakeErr    error        // error resulting from handshake

	inputMux sync.Mutex   // input mutex for concurrent reads
	input    bytes.Reader // application data waiting to be read, from rawInput.Next
	rawInput bytes.Buffer // raw input, starting with a record header

	// activeCall is an atomic int32; the low bit is whether Close has
	// been called. the rest of the bits are the number of goroutines
	// in Conn.Write.
	active int32

	// the secret key for connection's encryption
	secret []byte
}

// Implements net.Conn
// LocalAddr local address of connection
func (s *Conn) LocalAddr() net.Addr { return s.conn.LocalAddr() }

// RemoteAddr remote address of connection
func (s *Conn) RemoteAddr() net.Addr { return s.conn.RemoteAddr() }

// SetDeadline set deadline for connection
func (s *Conn) SetDeadline(t time.Time) error { return s.conn.SetDeadline(t) }

// SetReadDeadline set read deadline for connection
func (s *Conn) SetReadDeadline(t time.Time) error {
	return s.conn.SetReadDeadline(t)
}

// SetWriteDeadline set write deadline for connection
func (s *Conn) SetWriteDeadline(t time.Time) error {
	return s.conn.SetWriteDeadline(t)
}

// Write data to the connection
//
// As Write calls Handshake, in order to prevent indefinite blocking a deadline
// must be set for both Read and Write before Write is called when the handshake
// has not yet completed. See SetDeadline, SetReadDeadline, and SetWriteDeadline
func (s *Conn) Write(b []byte) (int, error) {
	var n int
	// interlock with Close below
	for {
		x := atomic.LoadInt32(&s.active)
		if x&1 != 0 {
			return n, net.ErrClosed
		}
		if atomic.CompareAndSwapInt32(&s.active, x, x+2) {
			break
		}
	}
	defer atomic.AddInt32(&s.active, -2)

	// do the handshake before write
	if err := s.Handshake(); err != nil {
		return n, err
	}
	if !s.handshakeComplete() {
		return n, errors.New("handshake not complete")
	}

	// write a record to the connection
	return s.writeRecord(b)
}

// Read data from the connection
//
// As Read calls Handshake, in order to prevent indefinite blocking a deadline
// must be set for both Read and Write before Read is called when the handshake
// has not yet completed. See SetDeadline, SetReadDeadline, and SetWriteDeadline.
func (s *Conn) Read(b []byte) (int, error) {
	var m int

	// do the handshake before read
	if err := s.Handshake(); err != nil {
		return m, err
	}

	if len(b) == 0 {
		// Put this after Handshake, in case calling Read(nil) for the side effect of the Handshake
		return 0, nil
	}

	s.inputMux.Lock()
	defer s.inputMux.Unlock()

	// input is empty, read and decrypt record by rawInput
	for s.input.Len() == 0 {
		var typ recordType

		// read the record and validate the record type
		payload, err := s.readRecord(&typ)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return m, io.EOF
			}
			return m, fmt.Errorf("read record: %w", err)
		}
		if typ != applicationData {
			return m, fmt.Errorf("record %d is invalid", typ)
		}

		// decrypt the record payload
		decrypted, err := aesDecrypt(payload, s.secret)
		if err != nil {
			return m, fmt.Errorf("decrypt record: %w", err)
		}

		// reset the input with decrypted data
		s.input.Reset(decrypted)
	}

	// read the data from input
	return s.input.Read(b)
}

// atLeastReader reads from R, stopping with EOF once at least N bytes have been
// read. It is different from an io.LimitedReader in that it doesn't cut short
// the last Read call, and in that it considers an early EOF an error.
type atLeastReader struct {
	R io.Reader
	N int64
}

func (r *atLeastReader) Read(p []byte) (int, error) {
	if r.N <= 0 {
		return 0, io.EOF
	}
	n, err := r.R.Read(p)

	// won't underflow unless len(p) >= n > 9223372036854775809
	r.N -= int64(n)
	if r.N > 0 && err == io.EOF {
		return n, io.ErrUnexpectedEOF
	}
	if r.N <= 0 && err == nil {
		return n, io.EOF
	}
	return n, err
}

// readFromUntil reads from r into c.rawInput until c.rawInput contains
// at least n bytes or else returns an error.
func (s *Conn) readFromUntil(r io.Reader, n int) error {
	if s.rawInput.Len() >= n {
		return nil
	}
	needs := n - s.rawInput.Len()

	// There might be extra input waiting on the wire. Make a best effort
	// attempt to fetch it so that it can be used in (*Conn).Read to
	// "predict" closeNotify alerts.
	s.rawInput.Grow(needs + bytes.MinRead)

	// read from connection
	if _, err := s.rawInput.ReadFrom(&atLeastReader{r, int64(needs)}); err != nil {
		return err
	}
	return nil
}

// read a record from connection
func (s *Conn) readRecord(typ *recordType) ([]byte, error) {
	// read the record header from connection
	if err := s.readFromUntil(s.conn, recordHeaderLen); err != nil {
		return nil, err
	}
	hdr := s.rawInput.Next(recordHeaderLen)

	// the record type
	*typ = recordType(hdr[0])
	// the record length
	n := int(hdr[1])<<16 | int(hdr[2])<<8 | int(hdr[3])

	// read the record body from connection
	if err := s.readFromUntil(s.conn, n); err != nil {
		return nil, err
	}
	record := s.rawInput.Next(n)

	return record, nil
}

// Handshake runs the client or server handshake protocol if it has not yet been run
//
// Most uses of this package need not call Handshake explicitly: the
// first Read or Write will call it automatically.
//
// For control over canceling or setting a timeout on a handshake, use
// the Dialer's DialContext method.
func (s *Conn) Handshake() error {
	s.handshakeMutex.Lock()
	defer s.handshakeMutex.Unlock()

	// if it happens error for handshake
	if err := s.handshakeErr; err != nil {
		return err
	}

	// the handshake is complete
	if s.handshakeComplete() {
		return nil
	}

	// do the handshake process
	return s.handshakeFn()
}

func (s *Conn) initRecord(typ recordType, d []byte) []byte {
	m := len(d)

	// init the record header
	out := make([]byte, recordHeaderLen+m)
	out[0] = byte(typ)
	out[1] = byte(m >> 16)
	out[2] = byte(m >> 8)
	out[3] = byte(m)

	// init the record body
	copy(out[recordHeaderLen:], d)

	return out
}

// write a record to the connection
func (s *Conn) writeRecord(record []byte) (int, error) {
	var n int

	// encrypt the payload with secret key
	encrypted, err := aesEncrypt(record, s.secret)
	if err != nil {
		return n, fmt.Errorf("encrypt: %w", err)
	}
	m := recordHeaderLen + len(encrypted)
	if m >= maxPayloadSize {
		return n, fmt.Errorf("message of length %d bytes exceeds maximum of %d bytes", m, maxPayloadSize)
	}

	// init the header and body of record
	d := s.initRecord(applicationData, encrypted)

	// write the key exchange record
	n, err = s.conn.Write(d)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (s *Conn) handshakeComplete() bool {
	return atomic.LoadUint32(&s.handshakeStatus) == 1
}

// Close the connection
func (s *Conn) Close() error {
	// Interlock with Conn.Write above.
	var x int32
	for {
		x = atomic.LoadInt32(&s.active)
		if x&1 != 0 {
			return net.ErrClosed
		}
		if atomic.CompareAndSwapInt32(&s.active, x, x|1) {
			break
		}
	}

	// close the connection
	return s.conn.Close()
}
