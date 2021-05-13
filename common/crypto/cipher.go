package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/cloudflare/circl/dh/x448"
	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/crypto/sha3"
)

// GcmTagSize is the GCM tag size is the difference in length between
// plaintext and ciphertext. From crypto/cipher/gcm.go in Go crypto library.
const GcmTagSize = 16

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in each direction).
	overflowLenAES256 = 5
)

type Cipher struct {
	prv, pub x448.Key
	shared   x448.Key

	aead       cipher.AEAD
	inCounter  Counter
	outCounter Counter
}

func NewCipher() *Cipher {
	var pub, prv x448.Key
	_, _ = io.ReadFull(rand.Reader, prv[:])
	x448.KeyGen(&pub, &prv)
	return &Cipher{
		pub:        pub,
		prv:        prv,
		outCounter: NewCounter(overflowLenAES256),
		inCounter:  NewCounter(overflowLenAES256),
	}
}

// PubKey method returns x448 public key
func (c *Cipher) PubKey() x448.Key {
	return c.pub
}

// Shared method generates a shared x448 key based on own private
// and shared public key
func (c *Cipher) Shared(pubKey x448.Key) (shared x448.Key) {
	x448.Shared(&shared, &(c.prv), &pubKey)
	return
}

// SetSharedAndConfigureAES method installs shared key from public
func (c *Cipher) SetSharedAndConfigureAES(rawPubKey []byte) error {
	if len(rawPubKey) != x448.Size {
		return errors.Errorf("not correct size of public key")
	}
	var pubKey x448.Key
	copy(pubKey[:], rawPubKey[:x448.Size])
	c.shared = c.Shared(pubKey)
	h := sha3.New256()
	if _, err := h.Write(c.shared[:]); err != nil {
		return errors.Errorf("unable to write bytes %w", err)
	}
	hash := h.Sum(nil)
	block, err := aes.NewCipher(hash)
	if err != nil {
		return errors.Errorf("unable to create a new cipher block, %w", err)
	}
	a, err := cipher.NewGCM(block)
	if err != nil {
		return errors.Errorf("could not create New GCM, %w", err)
	}
	c.aead = a
	return nil
}

func (c *Cipher) Encrypt(dst, plaintext []byte) ([]byte, error) {
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
func (c *Cipher) MaxOverhead() int {
	return GcmTagSize
}

// Decrypt method decrypts message using  private key
func (c *Cipher) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	seq, err := c.inCounter.Value()
	if err != nil {
		return nil, err
	}
	// If dst is equal to ciphertext[:0], ciphertext storage is reused.
	plaintext, err := c.aead.Open(dst, seq, ciphertext, nil)
	if err != nil {
		return nil, errors.Errorf("could not decrypt the ciphertext, %w", err)
	}
	c.inCounter.Inc()
	return plaintext, nil
}