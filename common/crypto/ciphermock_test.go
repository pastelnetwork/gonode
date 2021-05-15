package crypto

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/cloudflare/circl/dh/x448"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCounter struct {
	mock.Mock
}

func (mock *mockCounter) Value() ([]byte, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (mock *mockCounter) Inc() {
	mock.Called()
}


func TestEncrypt(t *testing.T) {
	var value [counterLen]byte
	value[counterLen-1] = 0x80

	plainText := []byte("SomeRandomStringToTest")
	rawPubKey := make([]byte, x448.Size)
	mockCounter := mockCounter{}

	var pub, prv x448.Key
	_, _ = io.ReadFull(rand.Reader, prv[:])
	x448.KeyGen(&pub, &prv)

	testService := Cipher{
		pub: pub,
		prv: prv,
		inCounter: &mockCounter,
		outCounter:&mockCounter,
	}

	mockCounter.On("Value").Return(value[:], nil)
	mockCounter.On("Inc")
	_ = testService.SetSharedAndConfigureAES(rawPubKey)
	_, err := testService.Encrypt(plainText, rawPubKey)
	assert.Nil(t, nil, err)
}
