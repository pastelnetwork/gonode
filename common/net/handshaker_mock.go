package net

import "github.com/pastelnetwork/gonode/common/crypto"

type handshakerMock struct {
	clientCipher *crypto.Cipher
	serverCipher *crypto.Cipher
	isClient     bool
}

func (h handshakerMock) ClientHello() ([]byte, error) {
	// TODO: implement later real returns artistID
	return []byte("some-unique-artist-id"), nil
}

func (h handshakerMock) ServerHello() ([]byte, error) {
	// returns artistID
	return []byte("ok"), nil
}

func (h *handshakerMock) ClientKeyExchange() ([]byte, error) {
	pubKey := h.clientCipher.PubKey()
	return pubKey[:], nil
}

func (h *handshakerMock) ServerKeyVerify(bytes []byte) (Cipher, error) {
	if err := h.serverCipher.SetSharedAndConfigureAES(bytes[:]); err != nil {
		return nil, err
	}
	return h.serverCipher, nil
}

func (h *handshakerMock) ServerKeyExchange() ([]byte, error) {
	// returns signature+publicKey
	pubKey := h.serverCipher.PubKey()
	return pubKey[:], nil
}

func (h handshakerMock) ClientKeyVerify(bytes []byte) (Cipher, error) {
	if err := h.clientCipher.SetSharedAndConfigureAES(bytes[:]); err != nil {
		return nil, err
	}
	return h.clientCipher, nil
}

func NewHandshakerMock(clientCipher *crypto.Cipher, serverCipher *crypto.Cipher, isClient bool) Handshaker {
	return &handshakerMock{
		clientCipher: clientCipher,
		serverCipher: serverCipher,
		isClient:     isClient,
	}
}
