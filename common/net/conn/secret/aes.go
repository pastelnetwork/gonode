package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func aesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	// create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// create a nonce from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// encrypt the data using aesGCM.Seal, first nonce argument in Seal is the prefix
	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

func aesDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	// create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// the nonce size
	nonceSize := aesGCM.NonceSize()

	// extract the nonce from the encrypted data
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// decrypt the data
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
