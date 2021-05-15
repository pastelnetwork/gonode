package net

import (
	"fmt"
	"github.com/cloudflare/circl/dh/x448"
	"github.com/stretchr/testify/mock"
	"net"
	"testing"
)

type mockHandShaker struct {
	mock.Mock
}

func (mock *mockHandShaker) ClientHello() ([]byte, error) {
	args := mock.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *mockHandShaker) ServerHello() ([]byte, error) {
	args := mock.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *mockHandShaker) ClientKeyExchange() ([]byte, error) {
	args := mock.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *mockHandShaker) ServerKeyVerify(publicKey []byte) (Cipher, error) {
	args := mock.Called(publicKey)
	return args.Get(0).(Cipher), args.Error(1)
}

func (mock *mockHandShaker) ServerKeyExchange() ([]byte, error) {
	args := mock.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *mockHandShaker) ClientKeyVerify(publicKey []byte) (Cipher, error) {
	args := mock.Called(publicKey)
	return args.Get(0).(Cipher), args.Error(1)
}

func TestClientHandshake(t *testing.T) {
	var pub, prv x448.Key
	x448.KeyGen(&pub, &prv)
	conn, _ := net.Dial("tcp", ":80")

	serverPublicKey := make([]byte, x448.Size)

	mockHandShake := mockHandShaker{}
	mockCipher := new(mockCipher)
	handshakeService := handShake{
		Conn: conn,
		HandshakeService: &mockHandShake,
	}

	mockHandShake.On("ClientHello").Return([]byte("some-unique-artist-id"), nil)
	mockHandShake.On("ClientKeyExchange").Return(pub[:], nil)
	mockHandShake.On("ClientKeyVerify", serverPublicKey).Return(mockCipher, nil)


	cipher, err := handshakeService.ClientHandshake()
	fmt.Println(cipher, err)
}