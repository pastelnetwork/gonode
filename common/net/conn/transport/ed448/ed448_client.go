package ed448

import (
	"context"
	"net"
)

func (transport *Ed448) ClientHandshake(context context.Context, conn net.Conn) (net.Conn, error) {
	clientHello := transport.prepareHelloMessage()
	if err := transport.writeRecord(clientHello, conn); err != nil {
		return nil, err
	}

	serverHelloMessage, err := transport.readServerHello(conn)
	if err != nil {
		return nil, err
	}

	if !transport.isClientSupportServerEncryption(serverHelloMessage.chosenEncryption) {
		return nil, ErrUnsupportedEncryption
	}

	clientHandshake, err := transport.prepareClientHandshake()
	if err != nil {
		return nil, err
	}
	if err := transport.writeRecord(clientHandshake, conn); err != nil {
		return nil, err
	}

	serverHandshakeMessage, err := transport.readServerHandshake(conn)
	if err != nil {
		return nil, err
	}

	if transport.verifySignature(serverHandshakeMessage.ctx, serverHandshakeMessage.signedPastelID, serverHandshakeMessage.pastelID, serverHandshakeMessage.pubKey) {
		return nil, ErrWrongSignature
	}

	if !transport.verifyPastelID(serverHandshakeMessage.pastelID) {
		return nil, ErrIncorrectPastelID
	}

	transport.isHandshakeEstablished = true

	return transport.initEncryptedConnection(conn, serverHelloMessage.chosenEncryption, serverHandshakeMessage.encryptionParams)
}

func (transport *Ed448) isClientSupportServerEncryption(serverEncryption string) bool {
	_, ok := transport.cryptos[serverEncryption]
	return ok
}

func (transport *Ed448) prepareHelloMessage() *ClientHelloMessage {
	var supportedEncryptions = make([]string, len(transport.cryptos))
	for crypto := range transport.cryptos {
		supportedEncryptions = append(supportedEncryptions, crypto)
	}
	return &ClientHelloMessage{
		supportedEncryptions:         supportedEncryptions,
		supportedSignatureAlgorithms: []SignScheme{ED448},
	}
}

func (transport *Ed448) readServerHello(conn net.Conn) (*ServerHelloMessage, error) {
	msg, err := transport.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHelloMessage), nil
}

func (transport *Ed448) readServerHandshake(conn net.Conn) (*ServerHandshakeMessage, error) {
	msg, err := transport.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHandshakeMessage), nil
}

// ToDo: update with appropriate implementation
func (transport *Ed448) retrievePastelID() uint32 {
	return 100
}

func (transport *Ed448) prepareClientHandshake() (*ClientHandshakeMessage, error) {
	signedPastelId := transport.getSignedPastelId()

	return &ClientHandshakeMessage{
		pastelID:       signedPastelId.pastelID,
		signedPastelID: signedPastelId.signedPastelID,
		pubKey:         signedPastelId.ctx,
		ctx:            signedPastelId.ctx,
	}, nil
}
