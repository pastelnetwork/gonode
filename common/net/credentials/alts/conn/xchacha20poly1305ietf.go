package conn

import (
	"crypto/cipher"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pkg/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenxChacha20Poly1305IETF = 8
)

type xchacha20poly1305ietfReKey struct {
	inCounter  Counter
	outCounter Counter
	inAEAD     cipher.AEAD
	outAEAD    cipher.AEAD
	nonceMask  [chacha20poly1305.NonceSizeX]byte
}

// NewxChaCha20Poly1305IETFReKey creates an instance that uses chacha20poly1305 with rekeying
// for ALTS record. The key argument should be 56 bytes, the first 32 bytes are used as
// a secret key and the remainining 24 bytes are used as public nonce
// Refer to: https://doc.libsodium.org/secret-key_cryptography/aead/chacha20-poly1305/xchacha20-poly1305_construction
func NewxChaCha20Poly1305IETFReKey(side alts.Side, key []byte) (ALTSRecordCrypto, error) {
	inCounter := NewInCounter(side, overflowLenxChacha20Poly1305IETF)
	outCounter := NewOutCounter(side, overflowLenxChacha20Poly1305IETF)

	if len(key) < chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX {
		return nil, errors.New("invalid keylength")
	}

	inAEAD, err := chacha20poly1305.NewX(key[:chacha20poly1305.KeySize])
	if err != nil {
		return nil, errors.Wrap(err, "new rekey aead in")
	}
	outAEAD, err := chacha20poly1305.NewX(key[:chacha20poly1305.KeySize])
	if err != nil {
		return nil, errors.Wrap(err, "new rekey aead out")
	}

	c := &xchacha20poly1305ietfReKey{
		inCounter:  inCounter,
		outCounter: outCounter,
		inAEAD:     inAEAD,
		outAEAD:    outAEAD,
	}
	copy(c.nonceMask[:], key[chacha20poly1305.KeySize:chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX])
	return c, nil
}

func (s *xchacha20poly1305ietfReKey) Encrypt(dst, plaintext []byte) ([]byte, error) {
	// If we need to allocate an output buffer, we want to include space for
	// GCM tag to avoid forcing ALTS record to reallocate as well.
	dlen := len(dst)
	dst, out := SliceForAppend(dst, len(plaintext)+GcmTagSize)
	seq, err := s.outCounter.Value()
	if err != nil {
		return nil, errors.Wrap(err, "get seq")
	}

	var nonceBuf [chacha20poly1305.NonceSizeX]byte
	maskNonce(nonceBuf[:], seq, s.nonceMask[:])

	data := out[:len(plaintext)]
	copy(data, plaintext) // data may alias plaintext

	// Seal appends the ciphertext and the tag to its first argument and
	// returns the updated slice. However, SliceForAppend above ensures that
	// dst has enough capacity to avoid a reallocation and copy due to the
	// append.
	dst = s.outAEAD.Seal(dst[:dlen], nonceBuf[:], data, nil)
	s.outCounter.Inc()
	return dst, nil
}

// EncryptionOverhead returns tag size
func (s *xchacha20poly1305ietfReKey) EncryptionOverhead() int {
	return chacha20poly1305.Overhead
}

func (s *xchacha20poly1305ietfReKey) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	seq, err := s.inCounter.Value()
	if err != nil {
		return nil, errors.Wrap(err, "get seq")
	}
	var nonceBuf [chacha20poly1305.NonceSizeX]byte
	maskNonce(nonceBuf[:], seq, s.nonceMask[:])

	plaintext, err := s.inAEAD.Open(dst, nonceBuf[:], ciphertext, nil)
	if err != nil {
		return nil, ErrAuthorize
	}
	s.inCounter.Inc()
	return plaintext, nil
}
