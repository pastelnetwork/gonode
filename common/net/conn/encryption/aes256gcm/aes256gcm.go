package aes256gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenAES256GCMGCM = 10
	AES256GCMGCMFrameSize   = 32
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

// New - Creates AES256GCMCrypto based on passed key that it's going to be used in ecryption process
func New(key []byte) (*AES256GCM, error) {
	if len(key) != 32 {
		return nil, KeySizeError(len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AES256GCM{
		configuration: key,
		block:         block,
		aead:          aead,
		inCounter:     NewCounter(overflowLenAES256GCMGCM),
		outCounter:    NewCounter(overflowLenAES256GCMGCM),
	}, nil
}

func (crypto *AES256GCM) GetConfiguration() []byte {
	return crypto.configuration
}

func (crypto *AES256GCM) Configure(rawData []byte) error {
	crypto.configuration = rawData
	block, err := aes.NewCipher(rawData)
	if err != nil {
		return err
	}
	crypto.block = block

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	crypto.aead = aead

	// reset counters
	crypto.inCounter = NewCounter(overflowLenAES256GCMGCM)
	crypto.outCounter = NewCounter(overflowLenAES256GCMGCM)
	return nil
}

func (crypto *AES256GCM) FrameSize() int {
	return AES256GCMGCMFrameSize
}

func (crypto *AES256GCM) About() string {
	return "AES256GMC"
}

func (crypto *AES256GCM) Encrypt(plainText []byte) ([]byte, error) {
	seq, err := crypto.outCounter.Value()
	if err != nil {
		return nil, err
	}
	cipherText := crypto.aead.Seal(nil, seq, plainText, nil)
	crypto.outCounter.Inc()
	return cipherText, nil
}

func (crypto *AES256GCM) Decrypt(cipherText []byte) ([]byte, error) {
	seq, err := crypto.inCounter.Value()
	if err != nil {
		return nil, err
	}
	// If dst is equal to ciphertext[:0], ciphertext storage is reused.
	plainText, err := crypto.aead.Open(nil, seq, cipherText, nil)
	if err != nil {
		return nil, err
	}
	crypto.inCounter.Inc()
	return plainText, nil
}
