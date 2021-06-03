package aes256gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/conn/transport"
	"strconv"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenAES256GCMGCM = 10
	gcmTagSize              = 16
)

// AES256GCM struct to work with encode and decode
type AES256GCM struct {
	configuration []byte
	block         cipher.Block
	aead          cipher.AEAD
	inCounter     Counter
	outCounter    Counter
}

// KeySizeError - error when passed wrong size of AES256GCM key
func KeySizeError(size int) error {
	return errors.New("crypto/aes: invalid key size " + strconv.Itoa(size) + " expected 32 bytes.")
}

// New - Creates AES256GCMCrypto based on passed key that it's going to be used in encryption process
func New(key []byte) (transport.Crypto, error) {
	if len(key) != 32 {
		return nil, KeySizeError(len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Errorf("error during initialising cipher %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Errorf("error during initialising gcm %w", err)
	}

	return &AES256GCM{
		configuration: key,
		block:         block,
		aead:          aead,
		inCounter:     NewCounter(overflowLenAES256GCMGCM),
		outCounter:    NewCounter(overflowLenAES256GCMGCM),
	}, nil
}

// GetConfiguration returns configuration that using for encrypting and decrypting data
func (crypto *AES256GCM) GetConfiguration() []byte {
	return crypto.configuration
}

// Configure - configure the AES256GCM encryption with new params
func (crypto *AES256GCM) Configure(rawData []byte) error {
	crypto.configuration = rawData
	block, err := aes.NewCipher(rawData)
	if err != nil {
		return errors.Errorf("error during initialising cipher %w", err)
	}
	crypto.block = block

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return errors.Errorf("error during initialising gcm %w", err)
	}
	crypto.aead = aead

	// reset counters
	crypto.inCounter = NewCounter(overflowLenAES256GCMGCM)
	crypto.outCounter = NewCounter(overflowLenAES256GCMGCM)
	return nil
}

// About - returns name of encryption
func (crypto *AES256GCM) About() string {
	return "AES256GMC"
}

// EncryptionOverhead - returns overhead of encryption
func (crypto *AES256GCM) EncryptionOverhead() int {
	return gcmTagSize
}

// Encrypt - encrypt the data and returns cipher text and error if it happens
func (crypto *AES256GCM) Encrypt(plainText []byte) ([]byte, error) {
	seq, err := crypto.outCounter.Value()
	if err != nil {
		return nil, errors.Errorf("invalid counter value %w", err)
	}
	cipherText := crypto.aead.Seal(nil, seq, plainText, nil)
	crypto.outCounter.Inc()
	return cipherText, nil
}

// Decrypt - decrypt the data and returns plaintext text and error if it happens
func (crypto *AES256GCM) Decrypt(cipherText []byte) ([]byte, error) {
	seq, err := crypto.inCounter.Value()
	if err != nil {
		return nil, errors.Errorf("invalid counter value %w", err)
	}

	plainText, err := crypto.aead.Open(nil, seq, cipherText, nil)
	if err != nil {
		return nil, errors.Errorf("error during decrypting data %w", err)
	}
	crypto.inCounter.Inc()
	return plainText, nil
}

// Clone - create a clone of current instance
func (crypto *AES256GCM) Clone() (transport.Crypto, error) {
	return New(crypto.configuration)
}
