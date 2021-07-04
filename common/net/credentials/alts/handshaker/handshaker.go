package handshaker

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts/authinfo"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts/conn"
	"google.golang.org/grpc/credentials"
)

const (
	// maxPendingHandshakes represents the maximum number of concurrent handshakes.
	maxPendingHandshakes = 100
	// handshake header length
	handshakeHeader = 2
	// ED448PubKeySize - public key size of ed448
	ED448PubKeySize = 56
	// ED448PriKeySize - private key size of ed448
	ED448PriKeySize = 144
	// record protocol name
	rekeyRecordProtocol = "ALTSRP_GCM_AES128_REKEY"
)

var (
	altsRecordFuncs = map[string]conn.ALTSRecordFunc{
		// ALTS handshaker protocols.
		rekeyRecordProtocol: func(s alts.Side, key []byte) (conn.ALTSRecordCrypto, error) {
			return conn.NewAES128GCMRekey(s, key)
		},
	}
	recordKeyLen = map[string]int{
		rekeyRecordProtocol: 44,
	}
	// control number of concurrent created (but not closed) handshakers.
	mu                   sync.Mutex
	concurrentHandshakes = int64(0)
	// errDropped occurs when maxPendingHandshakes is reached.
	errDropped = errors.New("maximum number of concurrent ALTS handshakes is reached")
)

func init() {
	for protocol, f := range altsRecordFuncs {
		if err := conn.RegisterProtocol(protocol, f); err != nil {
			panic(err)
		}
	}
}

func acquire() bool {
	mu.Lock()
	// If we need n to be configurable, we can pass it as an argument.
	n := int64(1)
	success := maxPendingHandshakes-concurrentHandshakes >= n
	if success {
		concurrentHandshakes += n
	}
	mu.Unlock()
	return success
}

func release() {
	mu.Lock()
	// If we need n to be configurable, we can pass it as an argument.
	n := int64(1)
	concurrentHandshakes -= n
	if concurrentHandshakes < 0 {
		mu.Unlock()
		panic("bad release")
	}
	mu.Unlock()
}

// key exchange request
type kexECDHRequest struct {
	PubKey []byte // the client's public key
}

// key exchange response
type kexECDHResponse struct {
	PubKey []byte // the server's public key
}

// ClientHandshakerOptions contains the client handshaker options that can
// provided by the caller.
type ClientHandshakerOptions struct {
}

// ServerHandshakerOptions contains the server handshaker options that can
// provided by the caller.
type ServerHandshakerOptions struct {
}

// DefaultClientHandshakerOptions returns the default client handshaker options.
func DefaultClientHandshakerOptions() *ClientHandshakerOptions {
	return &ClientHandshakerOptions{}
}

// DefaultServerHandshakerOptions returns the default client handshaker options.
func DefaultServerHandshakerOptions() *ServerHandshakerOptions {
	return &ServerHandshakerOptions{}
}

// altsHandshaker is used to complete a ALTS handshaking between client and server.
type altsHandshaker struct {
	// the record protocol
	protocol string
	// the record key
	key []byte
	// the connection to the peer.
	conn net.Conn
	// client handshake options.
	clientOpts *ClientHandshakerOptions
	// server handshake options.
	serverOpts *ServerHandshakerOptions
	// defines the side doing the handshake, client or server.
	side alts.Side
}

// NewClientHandshaker creates a ALTS handshaker
func NewClientHandshaker(_ context.Context, conn net.Conn, opts *ClientHandshakerOptions) (alts.Handshaker, error) {
	return &altsHandshaker{
		protocol:   rekeyRecordProtocol,
		conn:       conn,
		clientOpts: opts,
		side:       alts.ClientSide,
	}, nil
}

// NewServerHandshaker creates a ALTS handshaker
func NewServerHandshaker(_ context.Context, conn net.Conn, opts *ServerHandshakerOptions) (alts.Handshaker, error) {
	return &altsHandshaker{
		protocol:   rekeyRecordProtocol,
		conn:       conn,
		serverOpts: opts,
		side:       alts.ServerSide,
	}, nil
}

// read and decode the handshake record
func (s *altsHandshaker) readHandshake(value interface{}) error {
	hdr := make([]byte, handshakeHeader)
	// read handshake header from connection
	if _, err := s.conn.Read(hdr); err != nil {
		return fmt.Errorf("read handshake header: %w", err)
	}

	// the handshake body length
	n := int(hdr[0])<<8 | int(hdr[1])

	body := make([]byte, n)
	// read handshake body from connection
	if _, err := s.conn.Read(body); err != nil {
		return fmt.Errorf("read handshake body: %w", err)
	}

	// new a gob decoder for the handshake record
	decoder := gob.NewDecoder(bytes.NewBuffer(body))
	// decode the handshake record
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode handshake record: %w", err)
	}
	return nil
}

// encode and write the handshake record
func (s *altsHandshaker) writeHandshake(value interface{}) error {
	var buf bytes.Buffer
	// new a gob encoder for the handshake record
	encoder := gob.NewEncoder(&buf)

	// encode the handshake record
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode handshake record: %w", err)
	}

	m := buf.Len()
	d := make([]byte, handshakeHeader+m)
	d[0] = byte(m >> 8)
	d[1] = byte(m)
	copy(d[handshakeHeader:], buf.Bytes())

	// write the handshake record
	if _, err := s.conn.Write(d); err != nil {
		return fmt.Errorf("write handshake record: %w", err)
	}
	return nil
}

func (s *altsHandshaker) doClientHandshake() error {
	curve := ed448.NewCurve()
	// generate the private and public key for client
	priv, pub, ok := curve.GenerateKeys()
	if !ok {
		return errors.New("failed to generate keys")
	}

	request := kexECDHRequest{PubKey: pub[:]}
	// encode and write the client's handshake record
	if err := s.writeHandshake(&request); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	// read and decode the client's handshake record
	var response kexECDHResponse
	if err := s.readHandshake(&response); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	var serverPubKey [ED448PubKeySize]byte
	copy(serverPubKey[:], response.PubKey)
	// compute the secret key for connection
	secret := curve.ComputeSecret(priv, serverPubKey)

	// update the record key
	s.key = secret[:]

	return nil
}

// ClientHandshake starts and completes a client ALTS handshaking. Once
// done, ClientHandshake returns a secure connection.
func (s *altsHandshaker) ClientHandshake(_ context.Context) (net.Conn, credentials.AuthInfo, error) {
	if !acquire() {
		return nil, nil, errDropped
	}
	defer release()

	if s.side != alts.ClientSide {
		return nil, nil, errors.New("only handshakers created using NewClientHandshaker can perform a client handshaker")
	}

	// do the client handshake
	if err := s.doClientHandshake(); err != nil {
		return nil, nil, fmt.Errorf("do client handshake: %w", err)
	}
	log.Print("client handshake is complete")

	// The handshaker returns a 128 bytes key. It should be truncated based
	// on the returned record protocol.
	len, ok := recordKeyLen[s.protocol]
	if !ok {
		return nil, nil, fmt.Errorf("unknown resulted record protocol: %v", s.protocol)
	}
	// new a secure connection
	sc, err := conn.NewConn(s.side, s.conn, s.protocol, s.key[:len])
	if err != nil {
		return nil, nil, fmt.Errorf("new secure conn: %w", err)
	}

	return sc, authinfo.New(), nil
}

func (s *altsHandshaker) doServerHandshake() error {
	// read and decode the client's handshake record
	var request kexECDHRequest
	if err := s.readHandshake(&request); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	curve := ed448.NewCurve()
	// generate the private key and public key for server
	priv, pub, ok := curve.GenerateKeys()
	if !ok {
		return errors.New("failed to generate keys")
	}

	var clientPubKey [ED448PubKeySize]byte
	copy(clientPubKey[:], request.PubKey)
	// compute the secret key for connection
	secret := curve.ComputeSecret(priv, clientPubKey)
	// update the record key
	s.key = secret[:]

	response := kexECDHResponse{PubKey: pub[:]}
	// encode and write the server's handshake record
	if err := s.writeHandshake(&response); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	return nil
}

// ServerHandshake starts and completes a server ALTS handshaking for GCP. Once
// done, ServerHandshake returns a secure connection.
func (s *altsHandshaker) ServerHandshake(_ context.Context) (net.Conn, credentials.AuthInfo, error) {
	if !acquire() {
		return nil, nil, errDropped
	}
	defer release()

	if s.side != alts.ServerSide {
		return nil, nil, errors.New("only handshakers created using NewServerHandshaker can perform a server handshaker")
	}

	// do the server handshake
	if err := s.doServerHandshake(); err != nil {
		return nil, nil, fmt.Errorf("do server handshake: %w", err)
	}
	log.Print("server handshake is complete")

	// The handshaker returns a 128 bytes key. It should be truncated based
	// on the returned record protocol.
	len, ok := recordKeyLen[s.protocol]
	if !ok {
		return nil, nil, fmt.Errorf("unknown resulted record protocol: %v", s.protocol)
	}
	// new a secure connection
	sc, err := conn.NewConn(s.side, s.conn, s.protocol, s.key[:len])
	if err != nil {
		return nil, nil, fmt.Errorf("new secure conn: %w", err)
	}

	return sc, authinfo.New(), nil
}

// Close terminates the Handshaker. It should be called when the caller obtains
// the secure connection.
func (s *altsHandshaker) Close() {}
