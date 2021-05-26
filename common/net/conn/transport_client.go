package conn

import (
	"bytes"
	"context"
	"net"
)

func (t *Ed448Transport) ClientHandshake(context context.Context, conn net.Conn) (net.Conn, error) {
	clientHello := t.prepareHelloMessage()
	if err := t.writeRecord(clientHello, conn); err != nil {
		return nil, err
	}

	serverHelloMessage, err := t.readServerHello(conn)
	if err != nil {
		return nil, err
	}

	if !t.isClientSupportServerEncryption(serverHelloMessage.chosenEncryption) {
		return nil, ErrUnsupportedEncryption
	}

	clientHandshake, err := t.prepareClientHandshake()
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(clientHandshake.marshall()); err != nil {
		return nil, err
	}

	serverHandshakeMessage, err := t.readServerHandshake(conn)
	if err != nil {
		return nil, err
	}

	if bytes.Compare([]byte(t.ctx), serverHandshakeMessage.ctx) != 0 {
		return nil, ErrServerContextDoesNotMatch
	}

	if t.verifySignature(serverHandshakeMessage.ctx, serverHandshakeMessage.signedPastelID, serverHandshakeMessage.pastelID, serverHandshakeMessage.pubKey) {
		return nil, ErrWrongSignature
	}

	if !t.verifyPastelID(serverHandshakeMessage.pastelID) {
		return nil, ErrIncorrectPastelID
	}

	t.isHandshakeEstablished = true
	return conn, nil
}

func (t *Ed448Transport) isClientSupportServerEncryption(serverEncryption string) bool {
	_, ok := t.cryptos[serverEncryption]
	return ok
}

func (t *Ed448Transport) prepareHelloMessage() *ClientHelloMessage {
	var supportedEncryptions = make([]string, len(t.cryptos))
	for crypto := range t.cryptos {
		supportedEncryptions = append(supportedEncryptions, crypto)
	}
	return &ClientHelloMessage{
		supportedEncryptions:         supportedEncryptions,
		supportedSignatureAlgorithms: []SignScheme{ED448},
	}
}

func (t *Ed448Transport) readServerHello(conn net.Conn) (*ServerHelloMessage, error) {
	msg, err := t.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHelloMessage), nil
}

func (t *Ed448Transport) readServerHandshake(conn net.Conn) (*ServerHandshakeMessage, error) {
	msg, err := t.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHandshakeMessage), nil
}

// ToDo: update with appropriate implementation
func (t *Ed448Transport) retrievePastelID() uint32 {
	return 100
}

func (t *Ed448Transport) prepareClientHandshake() (*ClientHandshakeMessage, error) {
	signedPastelId := t.getSignedPastelId()

	return &ClientHandshakeMessage{
		pastelID:       signedPastelId.pastelID,
		signedPastelID: signedPastelId.signedPastelID,
		pubKey:         signedPastelId.ctx,
		ctx:            signedPastelId.ctx,
	}, nil
}
