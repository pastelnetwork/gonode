package auth

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
)

const (
	// clientSide identifies the client in this communication.
	clientSide Side = iota
	// serverSide identifies the server in this communication.
	serverSide
)

// side identifies the party's role: client or server.
type Side int

// Handshaker defines a ALTS handshaker interface.
type Handshaker interface {
	// ClientHandshake starts and completes a client-side handshaking and
	// returns a secure connection and corresponding auth information.
	ClientHandshake(ctx context.Context) (net.Conn, error)
	// ServerHandshake starts and completes a server-side handshaking and
	// returns a secure connection and corresponding auth information.
	ServerHandshake(ctx context.Context) (net.Conn, error)
}

// PeerAuthInfo is information for signing
type PeerAuthInfo struct {
	// Address of peer
	Address string
	// PastelID of peer
	PastelID string
	// Timestamp is signed
	Timestamp time.Time
	// Signature of timestamp
	Signature []byte
}

// PeerAuthHelper inteface defines functions for authentication
type PeerAuthHelper interface {
	// GenAuthInfo generatue auth info of peer to show his finger print
	GenAuthInfo(ctx context.Context) (*PeerAuthInfo, error)
	// VerifyPeer verify auth info of other side
	VerifyPeer(ctx context.Context, authInfo *PeerAuthInfo) error
}

type handshakeResultType int

const (
	handshakeResultSuccess handshakeResultType = 0
	handshakeResultFailed  handshakeResultType = 1
)

const (
	// handshake header length
	handshakeHeader = 2
)

// data exchange request
type kexExchangeRequest struct {
	PastelID  string
	Timestamp time.Time
	Signature []byte
}

// data exchange response
type kexExchangeResponse struct {
	Result handshakeResultType
	ErrStr string
}

// authHandshaker is used to complete a ALTS handshaking between client and server.
type authHandshaker struct {
	// peer functions for authentication
	peer PeerAuthHelper
	// the connection to the peer.
	conn net.Conn
	// defines the side doing the handshake, client or server.
	side Side
}

// NewClientHandshaker creates a ALTS handshaker
func NewClientHandshaker(_ context.Context, peer PeerAuthHelper, conn net.Conn) (Handshaker, error) {
	return &authHandshaker{
		peer: peer,
		conn: conn,
		side: clientSide,
	}, nil
}

// NewServerHandshaker creates a ALTS handshaker
func NewServerHandshaker(_ context.Context, peer PeerAuthHelper, conn net.Conn) (Handshaker, error) {
	return &authHandshaker{
		peer: peer,
		conn: conn,
		side: serverSide,
	}, nil
}

// read and decode the handshake record
func (s *authHandshaker) readHandshake(value interface{}) error {
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
func (s *authHandshaker) writeHandshake(value interface{}) error {
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
		return fmt.Errorf("flush: %s", err)
	}
	return nil
}

func (s *authHandshaker) doClientHandshake(ctx context.Context) error {
	authInfo, err := s.peer.GenAuthInfo(ctx)
	if err != nil {
		return fmt.Errorf("gen auth info: %w", err)
	}
	request := kexExchangeRequest{
		PastelID:  authInfo.PastelID,
		Timestamp: authInfo.Timestamp,
		Signature: authInfo.Signature,
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

	if response.Result != handshakeResultSuccess {
		return errors.New(response.ErrStr)
	}

	return nil
}

// ClientHandshake starts and completes a client ALTS handshaking. Once
// done, ClientHandshake returns a secure connection.
func (s *authHandshaker) ClientHandshake(ctx context.Context) (net.Conn, error) {

	if s.side != clientSide {
		return nil, errors.New("only handshakers created using NewClientHandshaker can perform a client handshaker")
	}

	// do the client handshake
	if err := s.doClientHandshake(ctx); err != nil {
		return nil, fmt.Errorf("do client handshake: %w", err)
	}
	log.WithContext(ctx).Debugf("client handshake is complete")

	return s.conn, nil
}

func (s *authHandshaker) doServerHandshake(ctx context.Context) error {
	// read and decode the client's handshake record
	var request kexExchangeRequest

	if err := s.readHandshake(&request); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}

	authInfo := &PeerAuthInfo{
		Address:   s.conn.RemoteAddr().String(),
		PastelID:  request.PastelID,
		Timestamp: request.Timestamp,
		Signature: request.Signature,
	}

	// verify request
	if err := s.peer.VerifyPeer(ctx, authInfo); err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to verify client public key")
		err = fmt.Errorf("verify client public key: %w", err)
		response := kexExchangeResponse{
			Result: handshakeResultFailed,
			ErrStr: err.Error(),
		}

		// encode and write the server's handshake record, don't care error
		s.writeHandshake(&response)

		return err
	}

	response := kexExchangeResponse{
		Result: handshakeResultSuccess,
		ErrStr: "",
	}

	// encode and write the server's handshake record
	if err := s.writeHandshake(&response); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	return nil
}

// ServerHandshake starts and completes a server ALTS handshaking for GCP. Once
// done, ServerHandshake returns a secure connection.
func (s *authHandshaker) ServerHandshake(ctx context.Context) (net.Conn, error) {

	if s.side != serverSide {
		return nil, errors.New("only handshakers created using NewServerHandshaker can perform a server handshaker")
	}

	// do the server handshake
	if err := s.doServerHandshake(ctx); err != nil {
		return nil, fmt.Errorf("do server handshake: %w", err)
	}
	log.WithContext(ctx).Debugf("server handshake is complete")

	return s.conn, nil
}
