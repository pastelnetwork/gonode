package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

type AES256Crypto struct {
	block cipher.Block
}

func KeySizeError(size int) error {
	return errors.New("crypto/aes: invalid key size " + strconv.Itoa(size) + " expected 32 bytes.")
}

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
