package conn

import (
	"fmt"

	"github.com/GoKillers/libsodium-go/crypto/aead/xchacha20poly1305ietf"
)

type _xchacha20poly1305ietf struct {
	key   [xchacha20poly1305ietf.KeyBytes]byte
	nonce [xchacha20poly1305ietf.NonceBytes]byte
}

// NewxChaCha20Poly1305IETF creates an instance that uses chacha20poly1305 with rekeying
// for ALTS record. The key argument should be 56 bytes, the first 32 bytes are used as
// a secret key and the remainining 24 bytes are used as public nonce
// Refer to: https://doc.libsodium.org/secret-key_cryptography/aead/chacha20-poly1305/xchacha20-poly1305_construction
func NewxChaCha20Poly1305IETF(key []byte) (ALTSRecordCrypto, error) {
	k := len(key)
	if k != xchacha20poly1305ietf.KeyBytes+xchacha20poly1305ietf.NonceBytes {
		return nil, KeySizeError(k)
	}

	c := &_xchacha20poly1305ietf{}
	copy(c.key[:], key[:xchacha20poly1305ietf.KeyBytes])
	copy(c.nonce[:], key[xchacha20poly1305ietf.KeyBytes:])
	return c, nil
}

func (s *_xchacha20poly1305ietf) Encrypt(dst, plaintext []byte) ([]byte, error) {
	ec := xchacha20poly1305ietf.Encrypt(plaintext, nil, &s.nonce, &s.key)
	dst = append(dst, ec...)
	return dst, nil
}

// EncryptionOverhead returns tag size
func (s *_xchacha20poly1305ietf) EncryptionOverhead() int {
	return xchacha20poly1305ietf.ABytes
}

func (s *_xchacha20poly1305ietf) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	p, err := xchacha20poly1305ietf.Decrypt(ciphertext, nil, &s.nonce, &s.key)
	if err != nil {
		return nil, fmt.Errorf("decryption unexpectedly succeeded for %s", err)
	}
	return p, nil
}
