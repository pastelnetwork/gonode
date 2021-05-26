package conn

import (
	"context"
	"net"
)

func (t *Ed448Transport) ServerHandshake(ctx context.Context, conn net.Conn) (net.Conn, error) {

	clientHelloMessage, err := t.readClientHello(conn)
	if err != nil {
		return nil, err
	}

	chosenEncryption, err := t.chooseEncryption(clientHelloMessage.supportedEncryptions)
	if err != nil {
		// probably before exit we have to send an answer to client that we can't make any connection
		return nil, err
	}

	if err = t.writeRecord(t.prepareServerHello(chosenEncryption), conn); err != nil {
		return nil, err
	}

	clientHandshakeMessage, err := t.readClientHandshake(conn)
	if err != nil {
		return nil, err
	}

	if !t.verifySignature(clientHandshakeMessage.ctx, clientHandshakeMessage.signedPastelID, clientHandshakeMessage.pastelID, clientHandshakeMessage.pubKey) {
		return nil, ErrWrongSignature
	}

	serverHandshakeMsg, err := t.prepareServerHandshakeMsg()
	if err != nil {
		return nil, err
	}

	if err = t.writeRecord(serverHandshakeMsg, conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func (t *Ed448Transport) readClientHello(conn net.Conn) (*ClientHelloMessage, error) {
	clientHello, err := t.readRecord(conn)
	if err != nil {
		return nil, err
	}

	return clientHello.(*ClientHelloMessage), nil
}

func (t *Ed448Transport) chooseEncryption(clientEncryptions []string) (string, error) {
	for _, encryption := range clientEncryptions {
		if _, ok := t.cryptos[encryption]; ok == true {
			return encryption, nil
		}
	}
	return "", ErrUnsupportedEncryption
}

func (t *Ed448Transport) prepareServerHello(chosenEncryption string) *ServerHelloMessage {
	return &ServerHelloMessage{
		chosenEncryption: chosenEncryption,
		chosenSignature:  ED448,
	}
}

func (t *Ed448Transport) readClientHandshake(conn net.Conn) (*ClientHandshakeMessage, error) {
	clientHandshakeMessage, err := t.readRecord(conn)

	if err != nil {
		return nil, err
	}

	return clientHandshakeMessage.(*ClientHandshakeMessage), nil
}

func (t *Ed448Transport) prepareServerHandshakeMsg() (*ServerHandshakeMessage, error) {
	// ToDo replace bypass below with retrieving pastelID and same pastelID but signed, within pubKey and context from cNode
	signedPastelID := t.getSignedPastelId()

	crypto, ok := t.cryptos[t.chosenEncryptionScheme]
	if ok == false {
		return nil, ErrUnsupportedEncryption
	}

	return &ServerHandshakeMessage{
		pastelID:         signedPastelID.pastelID,
		signedPastelID:   signedPastelID.signedPastelID,
		pubKey:           signedPastelID.pubKey,
		ctx:              signedPastelID.ctx,
		encryptionParams: crypto.GenerateEncryptionParams(),
	}, nil
}
