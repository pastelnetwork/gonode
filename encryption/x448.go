package encryption

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/cloudflare/circl/dh/x448"
)

type X448Cipher struct {
	prv, pub x448.Key
	shared x448.Key
	vi []byte
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
func (c *X448Cipher) SetShared(rawPubKey, vi []byte) error {
	if len(rawPubKey) != x448.Size{
		return fmt.Errorf("not correct size of public key")
	}
	if len(vi) < aes.BlockSize{
		return fmt.Errorf("vi shouldn't be less than %d", aes.BlockSize)
	}
	var pubKey x448.Key
	copy(pubKey[:], rawPubKey[:x448.Size])
	c.shared = c.Shared(pubKey)
	c.vi = vi
	return nil
}

// WrapWriter methods wraps to encrypt stream
func (c *X448Cipher) WrapWriter(writer io.Writer) (io.Writer, error){
	if len(c.vi) < aes.BlockSize{
		return nil, fmt.Errorf("vi is should be lower than %d", aes.BlockSize)
	}
	return NewStreamEncrypter(c.shared[:aes.BlockSize], c.vi, writer)
}

// WrapReader methods wraps reader to decrypt stream
func (c *X448Cipher) WrapReader(reader io.Reader) (io.Reader, error){
	if len(c.vi) < aes.BlockSize{
		return nil, fmt.Errorf("vi shouldn't be less than %d", aes.BlockSize)
	}
	return NewStreamDecrypter(c.shared[:aes.BlockSize], c.vi, reader)
}
