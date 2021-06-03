package ed448

import (
	"github.com/pastelnetwork/gonode/common/errors"
	"net"
)

// ClientHandshake does handshake on client side
func (trans *Ed448) ClientHandshake(conn net.Conn) (net.Conn, error) {
	if !acquire() {
		return nil, errors.New(ErrMaxHandshakes)
	}

	defer release()

	clientHello := trans.prepareHelloMessage()
	if err := trans.writeRecord(clientHello, conn); err != nil {
		return nil, err
	}

	serverHelloMessage, err := trans.readServerHello(conn)
	if err != nil {
		return nil, err
	}

	if !trans.isClientSupportServerEncryption(serverHelloMessage.chosenEncryption) {
		return nil, ErrUnsupportedEncryption
	}

	clientHandshake, err := trans.prepareClientHandshake()
	if err != nil {
		return nil, err
	}
	if err := trans.writeRecord(clientHandshake, conn); err != nil {
		return nil, err
	}

	serverHandshakeMessage, err := trans.readServerHandshake(conn)
	if err != nil {
		return nil, err
	}

	if !trans.verifySignature(serverHandshakeMessage.ctx, serverHandshakeMessage.signedPastelID, serverHandshakeMessage.pastelID, serverHandshakeMessage.pubKey) {
		return nil, errors.New(ErrWrongSignature)
	}

	if !trans.verifyPastelID(serverHandshakeMessage.pastelID) {
		return nil, errors.New(ErrIncorrectPastelID)
	}

	trans.isHandshakeEstablished = true

	return trans.initEncryptedConnection(conn, serverHelloMessage.chosenEncryption, serverHandshakeMessage.encryptionParams)
}

func (trans *Ed448) isClientSupportServerEncryption(serverEncryption string) bool {
	_, ok := trans.cryptos[serverEncryption]
	return ok
}

func (trans *Ed448) prepareHelloMessage() *ClientHelloMessage {
	var supportedEncryptions = make([]string, len(trans.cryptos))
	for crypto := range trans.cryptos {
		supportedEncryptions = append(supportedEncryptions, crypto)
	}
	return &ClientHelloMessage{
		supportedEncryptions:         supportedEncryptions,
		supportedSignatureAlgorithms: []SignScheme{ED448},
	}
}

func (trans *Ed448) readServerHello(conn net.Conn) (*ServerHelloMessage, error) {
	msg, err := trans.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHelloMessage), nil
}

func (trans *Ed448) readServerHandshake(conn net.Conn) (*ServerHandshakeMessage, error) {
	msg, err := trans.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return msg.(*ServerHandshakeMessage), nil
}

func (trans *Ed448) prepareClientHandshake() (*ClientHandshakeMessage, error) {
	signedPastelID := trans.getSignedPastelID()

	return &ClientHandshakeMessage{
		pastelID:       signedPastelID.pastelID,
		signedPastelID: signedPastelID.signedPastelID,
		pubKey:         signedPastelID.ctx,
		ctx:            signedPastelID.ctx,
	}, nil
}
