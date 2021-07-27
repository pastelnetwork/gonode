package conn

import (
	"fmt"

	"github.com/GoKillers/libsodium-go/crypto/aead/xchacha20poly1305ietf"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenxChacha20Poly1305IETF = 8
)

type _xchacha20poly1305ietf struct {
	key        [xchacha20poly1305ietf.KeyBytes]byte
	nonce      [xchacha20poly1305ietf.NonceBytes]byte
	inCounter  Counter
	outCounter Counter
}

// NewxChaCha20Poly1305IETF creates an instance that uses chacha20poly1305 with rekeying
// for ALTS record. The key argument should be 56 bytes, the first 32 bytes are used as
// a secret key and the remainining 24 bytes are used as public nonce
// Refer to: https://doc.libsodium.org/secret-key_cryptography/aead/chacha20-poly1305/xchacha20-poly1305_construction
func NewxChaCha20Poly1305IETF(side alts.Side, key []byte) (ALTSRecordCrypto, error) {
	k := len(key)
	if k != xchacha20poly1305ietf.KeyBytes+xchacha20poly1305ietf.NonceBytes {
		return nil, KeySizeError(k)
	}
	inCounter := NewInCounter(side, overflowLenxChacha20Poly1305IETF)
	outCounter := NewOutCounter(side, overflowLenxChacha20Poly1305IETF)

	c := &_xchacha20poly1305ietf{
		inCounter:  inCounter,
		outCounter: outCounter,
	}
	copy(c.key[:], key[:xchacha20poly1305ietf.KeyBytes])
	copy(c.nonce[:], key[xchacha20poly1305ietf.KeyBytes:])
	return c, nil
}

func (s *_xchacha20poly1305ietf) Encrypt(dst, plaintext []byte) ([]byte, error) {
	seq, err := s.outCounter.Value()
	if err != nil {
		return nil, err
	}

	ec := xchacha20poly1305ietf.Encrypt(plaintext, seq, &s.nonce, &s.key)
	dst = append(dst, ec...)
	s.outCounter.Inc()
	return dst, nil
}

// EncryptionOverhead returns tag size
func (s *_xchacha20poly1305ietf) EncryptionOverhead() int {
	return xchacha20poly1305ietf.ABytes
}

func (s *_xchacha20poly1305ietf) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	seq, err := s.inCounter.Value()
	if err != nil {
		return nil, err
	}

	p, err := xchacha20poly1305ietf.Decrypt(ciphertext, seq, &s.nonce, &s.key)
	if err != nil {
		return nil, fmt.Errorf("decryption unexpectedly succeeded for %s", err)
	}
	s.inCounter.Inc()
	return p, nil
}
