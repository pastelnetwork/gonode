package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

// ToDo: complete encryption

// AES256Crypto struct to work with encode and decode
type AES256Crypto struct {
	block cipher.Block
}

// KeySizeError - error when passed wrong size of AES256 key
func KeySizeError(size int) error {
	return errors.New("crypto/aes: invalid key size " + strconv.Itoa(size) + " expected 32 bytes.")
}

// New - Creates AES256Crypto based on passed key that it's going to be used in ecryption process
func New(key []byte) (*AES256Crypto, error) {
	if len(key) != 32 {
		return nil, KeySizeError(len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AES256Crypto{
		block: block,
	}, nil
}
