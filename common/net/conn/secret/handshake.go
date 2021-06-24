package secret

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/otrv4/ed448"
)

const (
	// ED448PubKeySize - public key size of ed448
	ED448PubKeySize = 56
	// ED448PriKeySize - private key size of ed448
	ED448PriKeySize = 144
	// AESSecretKeySize - secret key size of aes256
	AESSecretKeySize = 32
	// HandshakeTimeout - the default timeout for handshake read and write
	HandshakeTimeout = 2 * time.Second
)

// key exchange request
type kexECDHRequest struct {
	PubKey []byte // the client's public key
}

// key exchange response
type kexECDHResponse struct {
	PubKey []byte // the server's public key
}

// read and decode the handshake record
func (s *Conn) readHandshake(typ *recordType, value interface{}) error {
	// read handshake record from connection
	record, err := s.readRecord(typ)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return io.EOF
		}
		return fmt.Errorf("read handshake record: %w", err)
	}

	// new a gob decoder for the record payload
	decoder := gob.NewDecoder(bytes.NewBuffer(record))
	// decode the record payload
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode record payload: %w", err)
	}
	return nil
}

// encode and write the handshake record
func (s *Conn) writeHandshake(typ recordType, value interface{}) error {
	var buf bytes.Buffer
	// new a gob encoder for the key exchange record
	encoder := gob.NewEncoder(&buf)

	// encode the key exchange record
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode kex record: %w", err)
	}

	// init the key exchange record
	d := s.initRecord(typ, buf.Bytes())

	// write the key exchange record
	if _, err := s.conn.Write(d); err != nil {
		return fmt.Errorf("write record: %w", err)
	}
	return nil
}

// serverHandshake performs a secret handshake as a server
func (s *Conn) serverHandshake() error {
	// set deadline time for connection, and disable after handshake
	s.conn.SetDeadline(time.Now().Add(HandshakeTimeout))
	defer func() {
		s.conn.SetDeadline(time.Time{})
	}()

	var typ recordType
	// read and decode the client's handshake record
	var request kexECDHRequest
	if err := s.readHandshake(&typ, &request); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}
	if len(request.PubKey) != ED448PubKeySize {
		return errors.New("client's public key is invalid")
	}
	if typ != clientHandshake {
		return fmt.Errorf("record type %d is invalid", typ)
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
	sharedSecret := curve.ComputeSecret(priv, clientPubKey)
	// truncate the secret to 256 bit
	copy(s.secret, sharedSecret[:AESSecretKeySize])

	response := kexECDHResponse{PubKey: pub[:]}
	// encode and write the server's handshake record
	if err := s.writeHandshake(serverHandshake, &response); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	// update the handshake status
	atomic.StoreUint32(&s.handshakeStatus, 1)

	return nil
}

// clientHandshake performs a secret handshake as a client
func (s *Conn) clientHandshake() error {
	// set deadline time for connection, and disable after handshake
	s.conn.SetDeadline(time.Now().Add(HandshakeTimeout))
	defer func() {
		s.conn.SetDeadline(time.Time{})
	}()

	curve := ed448.NewCurve()
	// generate the private and public key for client
	priv, pub, ok := curve.GenerateKeys()
	if !ok {
		return errors.New("failed to generate keys")
	}

	request := kexECDHRequest{PubKey: pub[:]}
	// encode and write the client's handshake record
	if err := s.writeHandshake(clientHandshake, &request); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}

	// read and decode the server's handshake record
	var typ recordType
	// read and decode the client's handshake record
	var response kexECDHResponse
	if err := s.readHandshake(&typ, &response); err != nil {
		return fmt.Errorf("read handshake: %w", err)
	}
	if typ != serverHandshake {
		return fmt.Errorf("record type %d is invalid", typ)
	}
	if len(response.PubKey) != ED448PubKeySize {
		return errors.New("server's public key is invalid")
	}

	var serverPubKey [ED448PubKeySize]byte
	copy(serverPubKey[:], response.PubKey)
	// compute the secret key for connection
	sharedSecret := curve.ComputeSecret(priv, serverPubKey)
	// truncate the secret to 256 bit
	copy(s.secret, sharedSecret[:AESSecretKeySize])

	// update the handshake status
	atomic.StoreUint32(&s.handshakeStatus, 1)

	return nil
}
