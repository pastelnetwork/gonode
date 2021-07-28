package conn

import (
	"github.com/GoKillers/libsodium-go/crypto/aead/xchacha20poly1305ietf"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenxChacha20Poly1305IETF = 8
)

type xchacha20poly1305ietfReKey struct {
	inCounter  Counter
	outCounter Counter
	inAEAD     *rekeyAEADChaChaPoly
	outAEAD    *rekeyAEADChaChaPoly
}

// NewxChaCha20Poly1305IETFReKey creates an instance that uses chacha20poly1305 with rekeying
// for ALTS record. The key argument should be 56 bytes, the first 32 bytes are used as
// a secret key and the remainining 24 bytes are used as public nonce
// Refer to: https://doc.libsodium.org/secret-key_cryptography/aead/chacha20-poly1305/xchacha20-poly1305_construction
func NewxChaCha20Poly1305IETFReKey(side alts.Side, key []byte) (ALTSRecordCrypto, error) {
	inCounter := NewInCounter(side, overflowLenxChacha20Poly1305IETF)
	outCounter := NewOutCounter(side, overflowLenxChacha20Poly1305IETF)

	inAEAD, err := newRekeyAEADChaChaPoly(key)
	if err != nil {
		return nil, err
	}
	outAEAD, err := newRekeyAEADChaChaPoly(key)
	if err != nil {
		return nil, err
	}

	return &xchacha20poly1305ietfReKey{
		inCounter:  inCounter,
		outCounter: outCounter,
		inAEAD:     inAEAD,
		outAEAD:    outAEAD,
	}, nil
}

func (s *xchacha20poly1305ietfReKey) Encrypt(dst, plaintext []byte) ([]byte, error) {
	seq, err := s.outCounter.Value()
	if err != nil {
		return nil, err
	}

	ec, err := s.outAEAD.Encrypt(plaintext, seq)
	if err != nil {
		return nil, err
	}
	dst = append(dst, ec...)
	s.outCounter.Inc()
	return dst, nil
}

// EncryptionOverhead returns tag size
func (s *xchacha20poly1305ietfReKey) EncryptionOverhead() int {
	return xchacha20poly1305ietf.ABytes
}

func (s *xchacha20poly1305ietfReKey) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	seq, err := s.inCounter.Value()
	if err != nil {
		return nil, err
	}

	p, err := s.inAEAD.Decrypt(ciphertext, seq)
	if err != nil {
		return nil, err
	}
	s.inCounter.Inc()
	return p, nil
}
