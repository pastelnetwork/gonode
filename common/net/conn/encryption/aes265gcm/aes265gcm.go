package aes265gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenAES265GCMGCM = 10
	AES265GCMGCMFrameSize   = 32
)

// AES265GCM struct to work with encode and decode
type AES265GCM struct {
	configuration []byte
	block         cipher.Block
	aead          cipher.AEAD
	inCounter     Counter
	outCounter    Counter
}

// KeySizeError - error when passed wrong size of AES265GCM key
func KeySizeError(size int) error {
	return errors.New("crypto/aes: invalid key size " + strconv.Itoa(size) + " expected 32 bytes.")
}

// New - Creates AES265GCMCrypto based on passed key that it's going to be used in ecryption process
func New(key []byte) (*AES265GCM, error) {
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

	return &AES265GCM{
		configuration: key,
		block:         block,
		aead:          aead,
		inCounter:     NewCounter(overflowLenAES265GCMGCM),
		outCounter:    NewCounter(overflowLenAES265GCMGCM),
	}, nil
}

func (crypto *AES265GCM) GetConfiguration() []byte {
	return crypto.configuration
}

func (crypto *AES265GCM) Configure(rawData []byte) error {
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
	crypto.inCounter = NewCounter(overflowLenAES265GCMGCM)
	crypto.outCounter = NewCounter(overflowLenAES265GCMGCM)
	return nil
}

func (crypto *AES265GCM) FrameSize() int {
	return AES265GCMGCMFrameSize
}

func (crypto *AES265GCM) About() string {
	return "AES265GCMGCM"
}

func (crypto *AES265GCM) Encrypt(plainText []byte) ([]byte, error) {
	seq, err := crypto.outCounter.Value()
	if err != nil {
		return nil, err
	}
	cipherText := crypto.aead.Seal(nil, seq, plainText, nil)
	crypto.outCounter.Inc()
	return cipherText, nil
}

func (crypto *AES265GCM) Decrypt(cipherText []byte) ([]byte, error) {
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
