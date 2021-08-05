package handshaker

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/GoKillers/libsodium-go/cryptokdf"
	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/log"
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
	rekeyRecordAesProtocol    = "ALTSRP_GCM_AES128_REKEY"
	recordxChaChaPolyProtocol = "ALTSRP_XCHACHA20_POLY1305_IETF"
	kdfContext                = "__auth__"
	kdfSubKeyID               = 99
	challengeSize             = 16
)

var (
	altsRecordFuncs = map[string]conn.ALTSRecordFunc{
		// ALTS handshaker protocols.
		rekeyRecordAesProtocol: func(s alts.Side, key []byte) (conn.ALTSRecordCrypto, error) {
			return conn.NewAES128GCMRekey(s, key)
		},
		recordxChaChaPolyProtocol: func(s alts.Side, key []byte) (conn.ALTSRecordCrypto, error) {
			return conn.NewxChaCha20Poly1305IETFReKey(s, key)
		},
	}
	recordKeyLen = map[string]int{
		rekeyRecordAesProtocol:    44,
		recordxChaChaPolyProtocol: 56,
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

// data exchange request
type kexExchangeRequest struct {
	PubKey    [ED448PubKeySize]byte // the client's public key
	PastelID  string
	Signature []byte
}

// data exchange response
type kexExchangeResponse struct {
	PubKey    [ED448PubKeySize]byte // the server's public key
	PastelID  string
	Signature []byte
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
		protocol:   recordxChaChaPolyProtocol,
		conn:       conn,
		clientOpts: opts,
		side:       alts.ClientSide,
	}, nil
}

// NewServerHandshaker creates a ALTS handshaker
func NewServerHandshaker(_ context.Context, conn net.Conn, opts *ServerHandshakerOptions) (alts.Handshaker, error) {
	return &altsHandshaker{
		protocol:   recordxChaChaPolyProtocol,
		conn:       conn,
		serverOpts: opts,
		side:       alts.ServerSide,
	}, nil
}

// read and decode the handshake record
func (s *altsHandshaker) readHandshake(value interface{}) error {
	hdr := make([]byte, handshakeHeader)
	reader := bufio.NewReader(s.conn)
	// read handshake header from connection
	if _, err := reader.Read(hdr); err != nil {
		return fmt.Errorf("read handshake header: %w", err)
	}

	// the handshake body length
	n := binary.LittleEndian.Uint16(hdr[0:])

	body := make([]byte, n)
	// read handshake body from connection
	if _, err := reader.Read(body); err != nil {
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
	binary.BigEndian.PutUint16(d[0:], uint16(m))
	copy(d[handshakeHeader:], buf.Bytes())

	// write the handshake record
	writer := bufio.NewWriter(s.conn)
	if _, err := writer.Write(d); err != nil {
		return fmt.Errorf("write handshake record: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush, err: %s", err)
	}
	return nil
}

func (s *altsHandshaker) doClientHandshake(ctx context.Context, secClient alts.SecClient, secInfo *alts.SecInfo) error {
	challengeA := make([]byte, challengeSize)
	rand.Read(challengeA)

	if err := s.writeHandshake(&challengeA); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	var challengeB []byte
	if err := s.readHandshake(&challengeB); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	curve := ed448.NewCurve()
	// generate the private and public key for client
	priv, pub, ok := curve.GenerateKeys()
	if !ok {
		return errors.New("failed to generate keys")
	}

	dataSign := make([]byte, challengeSize+ED448PubKeySize)
	dataSign = append(dataSign, pub[:]...)
	dataSign = append(dataSign, challengeB...)
	signature, err := secClient.Sign(ctx, dataSign, secInfo.PastelID, secInfo.PassPhrase, secInfo.Algorithm)
	if err != nil {
		log.WithContext(ctx).Errorf("failed to generate signature: %v", err)
		return fmt.Errorf("failed to generate signature: %w", err)
	}

	request := kexExchangeRequest{
		PubKey:    pub,
		Signature: signature,
		PastelID:  secInfo.PastelID,
	}
	// encode and write the client's handshake record
	if err := s.writeHandshake(&request); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	// read and decode the client's handshake record
	var response kexExchangeResponse
	if err := s.readHandshake(&response); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	dataVerify := make([]byte, challengeSize+ED448PubKeySize)
	dataVerify = append(dataVerify, response.PubKey[:]...)
	dataVerify = append(dataVerify, challengeA...)
	if ok, err := secClient.Verify(ctx, dataVerify, string(response.Signature), response.PastelID, secInfo.Algorithm); err != nil || !ok {
		log.WithContext(ctx).Errorf("failed to verify server public key: %v", err)
		return fmt.Errorf("failed to verify server public key: %w", err)
	}

	// compute the secret key for connection
	secret := curve.ComputeSecret(priv, response.PubKey)
	if len(secret) < cryptokdf.CryptoKdfKeybytes() {
		return fmt.Errorf("secret key length is missmatch, len=%d", len(secret))
	}

	subkeyLen, ok := recordKeyLen[s.protocol]
	if !ok {
		return fmt.Errorf("unknown resulted record protocol: %v", s.protocol)
	}

	kdfValue, exit := cryptokdf.CryptoKdfDeriveFromKey(subkeyLen, kdfSubKeyID, kdfContext,
		secret[:cryptokdf.CryptoKdfKeybytes()])
	if exit != 0 {
		return fmt.Errorf("failed to compute kdf, exit=%d", exit)
	}
	// update the record key
	s.key = kdfValue

	return nil
}

// ClientHandshake starts and completes a client ALTS handshaking. Once
// done, ClientHandshake returns a secure connection.
func (s *altsHandshaker) ClientHandshake(ctx context.Context, secClient alts.SecClient, secInfo *alts.SecInfo) (net.Conn, credentials.AuthInfo, error) {
	if !acquire() {
		return nil, nil, errDropped
	}
	defer release()

	if s.side != alts.ClientSide {
		return nil, nil, errors.New("only handshakers created using NewClientHandshaker can perform a client handshaker")
	}

	// do the client handshake
	if err := s.doClientHandshake(ctx, secClient, secInfo); err != nil {
		return nil, nil, fmt.Errorf("do client handshake: %w", err)
	}
	log.WithContext(ctx).Debugf("client handshake is complete")

	// new a secure connection
	sc, err := conn.NewConn(s.side, s.conn, s.protocol, s.key)
	if err != nil {
		return nil, nil, fmt.Errorf("new secure conn: %w", err)
	}

	return sc, authinfo.New(), nil
}

func (s *altsHandshaker) doServerHandshake(ctx context.Context, secClient alts.SecClient, secInfo *alts.SecInfo) error {
	var challengeA []byte
	if err := s.readHandshake(&challengeA); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	challengeB := make([]byte, challengeSize)
	rand.Read(challengeB)
	if err := s.writeHandshake(&challengeB); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	// read and decode the client's handshake record
	var request kexExchangeRequest
	if err := s.readHandshake(&request); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	dataVerify := make([]byte, challengeSize+ED448PubKeySize)
	dataVerify = append(dataVerify, request.PubKey[:]...)
	dataVerify = append(dataVerify, challengeB...)
	if ok, err := secClient.Verify(ctx, dataVerify, string(request.Signature), request.PastelID, secInfo.Algorithm); err != nil || !ok {
		log.WithContext(ctx).Errorf("failed to verify client public key: %v", err)
		return fmt.Errorf("failed to verify client public key: %w", err)
	}

	curve := ed448.NewCurve()
	// generate the private key and public key for server
	priv, pub, ok := curve.GenerateKeys()
	if !ok {
		return errors.New("failed to generate keys")
	}

	dataSign := make([]byte, challengeSize+ED448PubKeySize)
	dataSign = append(dataSign, pub[:]...)
	dataSign = append(dataSign, challengeA...)
	signature, err := secClient.Sign(ctx, dataSign, secInfo.PastelID, secInfo.PassPhrase, secInfo.Algorithm)
	if err != nil {
		log.WithContext(ctx).Errorf("failed to generate signature: %v", err)
		return fmt.Errorf("failed to generate signature: %w", err)
	}

	// compute the secret key for connection
	secret := curve.ComputeSecret(priv, request.PubKey)
	if len(secret) < cryptokdf.CryptoKdfKeybytes() {
		return fmt.Errorf("secret key length is missmatch, len=%d", len(secret))
	}

	subkeyLen, ok := recordKeyLen[s.protocol]
	if !ok {
		return fmt.Errorf("unknown resulted record protocol: %v", s.protocol)
	}

	kdfValue, exit := cryptokdf.CryptoKdfDeriveFromKey(subkeyLen, kdfSubKeyID, kdfContext,
		secret[:cryptokdf.CryptoKdfKeybytes()])
	if exit != 0 {
		return fmt.Errorf("failed to compute kdf, exit=%d", exit)
	}
	// update the record key
	s.key = kdfValue
	response := kexExchangeResponse{
		PubKey:    pub,
		Signature: signature,
		PastelID:  secInfo.PastelID,
	}
	// encode and write the server's handshake record
	if err := s.writeHandshake(&response); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	return nil
}

// ServerHandshake starts and completes a server ALTS handshaking for GCP. Once
// done, ServerHandshake returns a secure connection.
func (s *altsHandshaker) ServerHandshake(ctx context.Context, secClient alts.SecClient, secInfo *alts.SecInfo) (net.Conn, credentials.AuthInfo, error) {
	if !acquire() {
		return nil, nil, errDropped
	}
	defer release()

	if s.side != alts.ServerSide {
		return nil, nil, errors.New("only handshakers created using NewServerHandshaker can perform a server handshaker")
	}

	// do the server handshake
	if err := s.doServerHandshake(ctx, secClient, secInfo); err != nil {
		return nil, nil, fmt.Errorf("do server handshake: %w", err)
	}
	log.WithContext(ctx).Debugf("server handshake is complete")

	// new a secure connection
	sc, err := conn.NewConn(s.side, s.conn, s.protocol, s.key)
	if err != nil {
		return nil, nil, fmt.Errorf("new secure conn: %w", err)
	}

	return sc, authinfo.New(), nil
}
