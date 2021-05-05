package encryption

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

type X448Cipher struct {
	prv, pub   x448.Key
	shared     x448.Key
	vi         []byte
	aead       cipher.AEAD
	inCounter  Counter
	outCounter Counter
}

func NewX448Cipher() *X448Cipher {
	var pub, prv x448.Key
	_, _ = io.ReadFull(rand.Reader, prv[:])
	x448.KeyGen(&pub, &prv)
	return &X448Cipher{
		pub: pub,
		prv: prv,
	}
}

// PubKey method returns x448 public key
func (c *X448Cipher) PubKey() x448.Key {
	return c.pub
}

// Shared method generates a shared x448 key based on own private
// and shared public key
func (c *X448Cipher) Shared(pubKey x448.Key) (shared x448.Key) {
	x448.Shared(&shared, &(c.prv), &pubKey)
	return
}

// SetShared method installs shared key from public
func (c *X448Cipher) SetSharedAndConfigureAES(rawPubKey, vi []byte) error {
	if len(rawPubKey) != x448.Size {
		return fmt.Errorf("not correct size of public key")
	}
	if len(vi) < aes.BlockSize {
		return fmt.Errorf("vi shouldn't be less than %d", aes.BlockSize)
	}
	var pubKey x448.Key
	copy(pubKey[:], rawPubKey[:x448.Size])
	c.shared = c.Shared(pubKey)
	c.vi = vi
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

func (c *X448Cipher) Encrypt(dst, plaintext []byte) ([]byte, error) {
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
func (c *X448Cipher) MaxOverhead() int {
	return GcmTagSize
}

// Decrypt method decrypts message using  private key
func (c *X448Cipher) Decrypt(dst, ciphertext []byte) ([]byte, error) {
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
