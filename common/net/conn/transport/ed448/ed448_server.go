package ed448

import (
	"github.com/pastelnetwork/gonode/common/errors"
	"net"
)

// ServerHandshake does handshake on server side
func (trans *Ed448) ServerHandshake(conn net.Conn) (net.Conn, error) {
	if !acquire() {
		return nil, errors.New(ErrMaxHandshakes)
	}

	defer release()

	clientHelloMessage, err := trans.readClientHello(conn)
	if err != nil {
		return nil, err
	}

	chosenEncryption, err := trans.chooseEncryption(clientHelloMessage.supportedEncryptions)
	if err != nil {
		// probably before exit we have to send an answer to client that we can'trans make any connection
		return nil, err
	}

	if err = trans.writeRecord(trans.prepareServerHello(chosenEncryption), conn); err != nil {
		return nil, err
	}

	clientHandshakeMessage, err := trans.readClientHandshake(conn)
	if err != nil {
		return nil, err
	}

	if !trans.verifySignature(clientHandshakeMessage.ctx, clientHandshakeMessage.signedPastelID, clientHandshakeMessage.pastelID, clientHandshakeMessage.pubKey) {
		return nil, errors.New(ErrWrongSignature)
	}

	serverHandshakeMsg, err := trans.prepareServerHandshakeMsg()
	if err != nil {
		return nil, err
	}

	if err = trans.writeRecord(serverHandshakeMsg, conn); err != nil {
		return nil, err
	}

	return trans.initEncryptedConnection(conn, chosenEncryption, serverHandshakeMsg.encryptionParams)
}

func (trans *Ed448) readClientHello(conn net.Conn) (*ClientHelloMessage, error) {
	clientHello, err := trans.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return clientHello.(*ClientHelloMessage), nil
}

func (trans *Ed448) chooseEncryption(clientEncryptions []string) (string, error) {
	for _, encryption := range clientEncryptions {
		if _, ok := trans.cryptos[encryption]; ok {
			trans.chosenEncryptionScheme = encryption
			return encryption, nil
		}
	}
	return "", ErrUnsupportedEncryption
}

func (trans *Ed448) prepareServerHello(chosenEncryption string) *ServerHelloMessage {
	return &ServerHelloMessage{
		chosenEncryption: chosenEncryption,
		chosenSignature:  ED448,
	}
}

func (trans *Ed448) readClientHandshake(conn net.Conn) (*ClientHandshakeMessage, error) {
	clientHandshakeMessage, err := trans.readRecord(conn)

	if err != nil {
		return nil, err
	}

	return clientHandshakeMessage.(*ClientHandshakeMessage), nil
}

func (trans *Ed448) prepareServerHandshakeMsg() (*ServerHandshakeMessage, error) {
	// ToDo replace bypass below with retrieving pastelID and same pastelID but signed, within pubKey and context from cNode
	signedPastelID := trans.getSignedPastelID()

	crypto, ok := trans.cryptos[trans.chosenEncryptionScheme]
	if !ok {
		return nil, errors.New(ErrUnsupportedEncryption)
	}

	return &ServerHandshakeMessage{
		pastelID:         signedPastelID.pastelID,
		signedPastelID:   signedPastelID.signedPastelID,
		pubKey:           signedPastelID.pubKey,
		ctx:              signedPastelID.ctx,
		encryptionParams: crypto.GetConfiguration(),
	}, nil
}
