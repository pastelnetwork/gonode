package ed448

import (
	"net"
)

func (transport *Ed448) ServerHandshake(conn net.Conn) (net.Conn, error) {

	clientHelloMessage, err := transport.readClientHello(conn)
	if err != nil {
		return nil, err
	}

	chosenEncryption, err := transport.chooseEncryption(clientHelloMessage.supportedEncryptions)
	if err != nil {
		// probably before exit we have to send an answer to client that we can'transport make any connection
		return nil, err
	}

	if err = transport.writeRecord(transport.prepareServerHello(chosenEncryption), conn); err != nil {
		return nil, err
	}

	clientHandshakeMessage, err := transport.readClientHandshake(conn)
	if err != nil {
		return nil, err
	}

	if !transport.verifySignature(clientHandshakeMessage.ctx, clientHandshakeMessage.signedPastelID, clientHandshakeMessage.pastelID, clientHandshakeMessage.pubKey) {
		return nil, ErrWrongSignature
	}

	serverHandshakeMsg, err := transport.prepareServerHandshakeMsg()
	if err != nil {
		return nil, err
	}

	if err = transport.writeRecord(serverHandshakeMsg, conn); err != nil {
		return nil, err
	}

	return transport.initEncryptedConnection(conn, chosenEncryption, serverHandshakeMsg.encryptionParams)
}

func (transport *Ed448) readClientHello(conn net.Conn) (*ClientHelloMessage, error) {
	clientHello, err := transport.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return clientHello.(*ClientHelloMessage), nil
}

func (transport *Ed448) chooseEncryption(clientEncryptions []string) (string, error) {
	for _, encryption := range clientEncryptions {
		if _, ok := transport.cryptos[encryption]; ok == true {
			transport.chosenEncryptionScheme = encryption
			return encryption, nil
		}
	}
	return "", ErrUnsupportedEncryption
}

func (transport *Ed448) prepareServerHello(chosenEncryption string) *ServerHelloMessage {
	return &ServerHelloMessage{
		chosenEncryption: chosenEncryption,
		chosenSignature:  ED448,
	}
}

func (transport *Ed448) readClientHandshake(conn net.Conn) (*ClientHandshakeMessage, error) {
	clientHandshakeMessage, err := transport.readRecord(conn)

	if err != nil {
		return nil, err
	}

	return clientHandshakeMessage.(*ClientHandshakeMessage), nil
}

func (transport *Ed448) prepareServerHandshakeMsg() (*ServerHandshakeMessage, error) {
	// ToDo replace bypass below with retrieving pastelID and same pastelID but signed, within pubKey and context from cNode
	signedPastelID := transport.getSignedPastelId()

	crypto, ok := transport.cryptos[transport.chosenEncryptionScheme]
	if ok == false {
		return nil, ErrUnsupportedEncryption
	}

	return &ServerHandshakeMessage{
		pastelID:         signedPastelID.pastelID,
		signedPastelID:   signedPastelID.signedPastelID,
		pubKey:           signedPastelID.pubKey,
		ctx:              signedPastelID.ctx,
		encryptionParams: crypto.GetConfiguration(),
	}, nil
}
