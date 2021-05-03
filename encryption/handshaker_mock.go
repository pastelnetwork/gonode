package encryption

import (
	"crypto/aes"
)

type handshakerMock struct {
	clientCipher *X448Cipher
	serverCipher *X448Cipher
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
	// TODO: add signature
	// returns signature+publicKey
	pubKey := h.clientCipher.PubKey()
	return pubKey[:], nil
}

func (h *handshakerMock) ServerKeyVerify(bytes []byte) (Cipher, error) {
	// TODO: verify signature
	vi, err := h.ClientHello()
	if err != nil{
		return nil, err
	}
	vi = vi[:aes.BlockSize]
	if err := h.serverCipher.SetShared(bytes[:], vi); err != nil{
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
	// verify client's signature
	vi, err := h.ClientHello()
	if err != nil{
		return nil, err
	}
	vi = vi[:aes.BlockSize]
	if err := h.clientCipher.SetShared(bytes[:], vi); err != nil{
		return nil, err
	}
	return h.clientCipher, nil
}

func NewHandshakerMock(clientCipher *X448Cipher, serverCipher *X448Cipher) Handshaker {
	return &handshakerMock{
		clientCipher: clientCipher,
		serverCipher: serverCipher,
	}
}
