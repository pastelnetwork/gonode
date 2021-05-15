package crypto

import (
	"testing"

	"github.com/cloudflare/circl/dh/x448"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type cryptoTest struct {
	suite.Suite
	Cipher *Cipher
	Counter Counter
}

func TestCipher(t *testing.T) {
	suite.Run(t, new(cryptoTest))
}

func (crypto *cryptoTest) SetupSuite() {
	crypto.Cipher = NewCipher()
}

func (crypto *cryptoTest) TestSetSharedAndConfigureAES() {
	//Testing with valid public key
	validClientPublicKey := make([]byte, x448.Size)
	err := crypto.Cipher.SetSharedAndConfigureAES(validClientPublicKey)
	assert.Equal(crypto.T(), nil, err, "Err should be nil as public key is valid")
	assert.NotNil(crypto.T(), crypto.Cipher.aead, "Expected aead not to be nil as it is set with public key")

	//Testing with invalid public key
	invalidClientPublicKey := make([]byte, 12)
	err = crypto.Cipher.SetSharedAndConfigureAES(invalidClientPublicKey)
	assert.EqualError(crypto.T(), err, "not correct size of public key", "Expect invalid public key to have incorrect size")

}

func (crypto *cryptoTest) TestEncrypt() {
	rawPubKey := make([]byte, x448.Size)
	err := crypto.Cipher.SetSharedAndConfigureAES(rawPubKey)

	text := []byte("Hello World")
	_, err = crypto.Cipher.Encrypt(text, rawPubKey)
	assert.Nil(crypto.T(), err)
}
