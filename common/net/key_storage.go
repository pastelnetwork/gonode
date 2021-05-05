package net

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/sha3"

	"github.com/cloudflare/circl/dh/x448"
)

const GcmTagSize = 16

type KeyStorage struct {
	prv, pub x448.Key
	shared   x448.Key

	aead       cipher.AEAD
	inCounter  Counter
	outCounter Counter
}

func NewKeyStorage() *KeyStorage {
	var pub, prv x448.Key
	_, _ = io.ReadFull(rand.Reader, prv[:])
	x448.KeyGen(&pub, &prv)
	return &KeyStorage{
		pub: pub,
		prv: prv,
	}
}

// PubKey method returns x448 public key
func (c *KeyStorage) PubKey() x448.Key {
	return c.pub
}

// Shared method generates a shared x448 key based on own private
// and shared public key
func (c *KeyStorage) Shared(pubKey x448.Key) (shared x448.Key) {
	x448.Shared(&shared, &(c.prv), &pubKey)
	return
}

// SetShared method installs shared key from public
func (c *KeyStorage) SetSharedAndConfigureAES(rawPubKey []byte) error {
	if len(rawPubKey) != x448.Size {
		return fmt.Errorf("not correct size of public key")
	}
	var pubKey x448.Key
	copy(pubKey[:], rawPubKey[:x448.Size])
	c.shared = c.Shared(pubKey)
	h := sha3.New256()
	h.Write(c.shared[:])
	hash := h.Sum(nil)
	block, err := aes.NewCipher(hash)
	if err != nil {
		return err
	}
	a, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	c.aead = a
	c.outCounter = NewOutCounter(GcmTagSize)
	c.inCounter = NewInCounter(GcmTagSize)
	return nil
}

func (c *KeyStorage) Encrypt(dst, plaintext []byte) ([]byte, error) {
	// If we need to allocate an output buffer, we want to include space for
	// GCM tag to avoid forcing ALTS record to reallocate as well.
	dlen := len(dst)
	dst, out := SliceForAppend(dst, len(plaintext)+GcmTagSize)
	seq, err := c.outCounter.Value()
	if err != nil {
		return nil, err
	}
	data := out[:len(plaintext)]
	copy(data, plaintext) // data may alias plaintext

	// Seal appends the ciphertext and the tag to its first argument and
	// returns the updated slice. However, SliceForAppend above ensures that
	// dst has enough capacity to avoid a reallocation and copy due to the
	// append.
	dst = c.aead.Seal(dst[:dlen], seq, data, nil)
	c.outCounter.Inc()
	return dst, nil
}

// MaxOverhead implements Cipher.
func (c *KeyStorage) MaxOverhead() int {
	return GcmTagSize
}

// Decrypt method decrypts message using  private key
func (c *KeyStorage) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	seq, err := c.inCounter.Value()
	if err != nil {
		return nil, err
	}
	// If dst is equal to ciphertext[:0], ciphertext storage is reused.
	plaintext, err := c.aead.Open(dst, seq, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	c.inCounter.Inc()
	return plaintext, nil
}
