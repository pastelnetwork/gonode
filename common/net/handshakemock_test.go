package net

import (
	"github.com/cloudflare/circl/dh/x448"
	"github.com/pastelnetwork/gonode/common/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockCipher struct {
	mock.Mock
}

func (mock *mockCipher) Encrypt(dst, plaintext []byte) ([]byte, error) {
	args := mock.Called(dst, plaintext)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *mockCipher) Decrypt(dst, plaintext []byte) ([]byte, error) {
	args := mock.Called(dst, plaintext)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *mockCipher) PubKey() x448.Key {
	args := mock.Called()
	return args.Get(0).(x448.Key)
}

func (mock *mockCipher) Shared(pubKey x448.Key) x448.Key {
	args := mock.Called(pubKey)
	return args.Get(0).(x448.Key)
}

func (mock *mockCipher) SetSharedAndConfigureAES(rawPubKey []byte) error {
	args := mock.Called(rawPubKey)
	return args.Error(1)
}

func (mock *mockCipher) MaxOverhead() int {
	args := mock.Called()
	return args.Get(0).(int)
}


func TestClientHello(t *testing.T) {
	cipher := crypto.NewCipher()

	hankShaker := handshakerMock{
		clientCipher: cipher,
		serverCipher: cipher,
		isClient: false,
	}
	artistId, err := hankShaker.ClientHello()
	assert.Nil(t, err)
	assert.Equal(t, []byte("some-unique-artist-id"), artistId)
}

func TestServerHello(t *testing.T) {
	hankShaker := handshakerMock{
		clientCipher: crypto.NewCipher(),
		serverCipher: crypto.NewCipher(),
		isClient: false,
	}
	srvResponse, err := hankShaker.ServerHello()
	assert.Nil(t, err)
	assert.Equal(t, []byte("ok"), srvResponse)
}

func TestClientKeyExchange(t *testing.T) {
	mockCipher := mockCipher{}
	cipher := crypto.NewCipher()
	hankShaker := handshakerMock{
		clientCipher: cipher,
		serverCipher: cipher,
		isClient: false,
	}
	pub := cipher.PubKey()
	mockCipher.On("PubKey").Return(cipher.PubKey(), nil)
	pubKey, err := hankShaker.ClientKeyExchange()
	assert.Nil(t, err)
	assert.Equal(t, pub[:], pubKey)
}

func TestServerKeyExchange(t *testing.T) {
	mockCipher := mockCipher{}
	cipher := crypto.NewCipher()
	hankShaker := handshakerMock{
		clientCipher: cipher,
		serverCipher: cipher,
		isClient: false,
	}
	pub := cipher.PubKey()

	mockCipher.On("PubKey"). Return(cipher.PubKey(), nil)
	pubKey, err := hankShaker.ServerKeyExchange()
	assert.Nil(t, err)
	assert.Equal(t, pub[:], pubKey)
}

func TestClientKeyVerify(t *testing.T) {
	mockCipher := mockCipher{}
	hankShaker := handshakerMock{
		clientCipher: crypto.NewCipher(),
		serverCipher: crypto.NewCipher(),
		isClient: false,
	}

	var pub, prv x448.Key
	x448.KeyGen(&pub, &prv)

	serverPublicKey := make([]byte, x448.Size)
	mockCipher.On("SetSharedAndConfigureAES", serverPublicKey). Return(pub, nil)
	_, err := hankShaker.ClientKeyVerify(serverPublicKey)
	assert.Nil(t, err)
}