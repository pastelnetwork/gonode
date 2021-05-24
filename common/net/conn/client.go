package conn

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"github.com/cloudflare/circl/sign/ed448"
)

var ErrServerContextDoesNotMatch = errors.New("Server context does not match")

// Client struct that implements handshake process
type Client struct {
	conn                   Conn
	secret                 ed448.PrivateKey
	public                 ed448.PublicKey
	serverPublicKey        ed448.PublicKey
	encryptionScheme       EncryptionScheme
	ctx                    string
	isHandshakeEstablished bool
}

// NewClient creates a new Client
func NewClient(conn Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) prepareHelloMessage() *ClientHelloMessage {
	return &ClientHelloMessage{
		supportedEncryptions:         []EncryptionScheme{AES256},
		supportedSignatureAlgorithms: []SignScheme{ED448},
	}
}

func (c *Client) readServerHello() (*ServerHelloMessage, error) {
	msg, err := c.conn.readRecord()
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHelloMessage), nil
}

func (c *Client) readServerHandshake() (*ServerHandshakeMessage, error) {
	msg, err := c.conn.readRecord()
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHandshakeMessage), nil
}

// ToDo: update with appropriate implementation
func (c *Client) retrievePastelID() uint32 {
	return 100
}

func (c *Client) prepareClientHandshake() (*ClientHandshakeMessage, error) {
	public, secret, err := ed448.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	c.secret = secret
	c.public = public
	c.ctx = randomString(512)

	pastelID := c.retrievePastelID()
	msg := make([]byte, 4)
	binary.LittleEndian.PutUint32(msg, pastelID)

	signedPastelID := ed448.Sign(c.secret, msg, c.ctx)

	return &ClientHandshakeMessage{
		signedPastelID: signedPastelID,
		pubKey:         c.public,
		ctx:            []byte(c.ctx),
	}, nil
}

// ToDo: update with appropriate implementation
func (c *Client) verifyPastelID(_ []byte) bool {
	return true
}

// Handshake - does handshake process with a server
func (c *Client) Handshake() error {
	clientHello := c.prepareHelloMessage()
	if _, err := c.conn.write(clientHello.marshall()); err != nil {
		return err
	}

	serverHelloMessage, err := c.readServerHello()
	if err != nil {
		return err
	}

	c.encryptionScheme = serverHelloMessage.chosenEncryption
	c.serverPublicKey = serverHelloMessage.publicKey

	clientHandshake, err := c.prepareClientHandshake()
	if err != nil {
		return err
	}
	if _, err := c.conn.write(clientHandshake.marshall()); err != nil {
		return err
	}

	serverHandshakeMessage, err := c.readServerHandshake()
	if err != nil {
		return nil
	}

	if bytes.Compare([]byte(c.ctx), serverHandshakeMessage.ctx) != 0 {
		return ErrServerContextDoesNotMatch
	}

	if ed448.Verify(serverHandshakeMessage.pubKey, serverHandshakeMessage.pastelID, serverHandshakeMessage.signedPastelID, c.ctx) {
		return ErrWrongSignature
	}

	// save public key from server
	c.serverPublicKey = serverHandshakeMessage.pubKey

	if !c.verifyPastelID(serverHandshakeMessage.pastelID) {
		return ErrIncorrectPastelID
	}

	c.isHandshakeEstablished = true
	return nil
}

// IsHandshakeEstablished - returns flag that says is handshake process finished
func (c *Client) IsHandshakeEstablished() bool {
	return c.isHandshakeEstablished
}

// GetCryptoProtocol - returns current encryption scheme that supports by server
func (c *Client) GetCryptoProtocol() EncryptionScheme {
	return c.encryptionScheme
}
