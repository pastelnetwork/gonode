package conn

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/GoKillers/libsodium-go/crypto/aead/xchacha20poly1305ietf"
)

// rekeyAEAD holds the necessary information for an AEAD based on
// AES-GCM that performs nonce-based key derivation and XORs the
// nonce with a random mask.
type rekeyAEAD struct {
	kdfKey     []byte
	kdfCounter []byte
	nonceMask  []byte
	nonceBuf   []byte
	gcmAEAD    cipher.AEAD
}

// KeySizeError signals that the given key does not have the correct size.
type KeySizeError int

func (k KeySizeError) Error() string {
	return "alts/conn: invalid key size " + strconv.Itoa(int(k))
}

// newRekeyAEAD creates a new instance of aes128gcm with rekeying.
// The key argument should be 44 bytes, the first 32 bytes are used as a key
// for HKDF-expand and the remainining 12 bytes are used as a random mask for
// the counter.
func newRekeyAEAD(key []byte) (*rekeyAEAD, error) {
	k := len(key)
	if k != kdfKeyLen+nonceLen {
		return nil, KeySizeError(k)
	}
	return &rekeyAEAD{
		kdfKey:     key[:kdfKeyLen],
		kdfCounter: make([]byte, kdfCounterLen),
		nonceMask:  key[kdfKeyLen:],
		nonceBuf:   make([]byte, nonceLen),
		gcmAEAD:    nil,
	}, nil
}

// Seal rekeys if nonce[2:8] is different than in the last call, masks the nonce,
// and calls Seal for aes128gcm.
func (s *rekeyAEAD) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	if err := s.rekeyIfRequired(nonce); err != nil {
		panic(fmt.Sprintf("Rekeying failed with: %s", err.Error()))
	}
	maskNonce(s.nonceBuf, nonce, s.nonceMask)
	return s.gcmAEAD.Seal(dst, s.nonceBuf, plaintext, additionalData)
}

// Open rekeys if nonce[2:8] is different than in the last call, masks the nonce,
// and calls Open for aes128gcm.
func (s *rekeyAEAD) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if err := s.rekeyIfRequired(nonce); err != nil {
		return nil, err
	}
	maskNonce(s.nonceBuf, nonce, s.nonceMask)
	return s.gcmAEAD.Open(dst, s.nonceBuf, ciphertext, additionalData)
}

// rekeyIfRequired creates a new aes128gcm AEAD if the existing AEAD is nil
// or cannot be used with given nonce.
func (s *rekeyAEAD) rekeyIfRequired(nonce []byte) error {
	newKdfCounter := nonce[kdfCounterOffset : kdfCounterOffset+kdfCounterLen]
	if s.gcmAEAD != nil && bytes.Equal(newKdfCounter, s.kdfCounter) {
		return nil
	}
	copy(s.kdfCounter, newKdfCounter)
	a, err := aes.NewCipher(hkdfExpand(s.kdfKey, s.kdfCounter, aeadKeyLen))
	if err != nil {
		return err
	}
	s.gcmAEAD, err = cipher.NewGCM(a)
	return err
}

// maskNonce XORs the given nonce with the mask and stores the result in dst.
func maskNonce(dst, nonce, mask []byte) {
	nonce1 := binary.LittleEndian.Uint64(nonce[:sizeUint64])
	nonce2 := binary.LittleEndian.Uint32(nonce[sizeUint64:])
	mask1 := binary.LittleEndian.Uint64(mask[:sizeUint64])
	mask2 := binary.LittleEndian.Uint32(mask[sizeUint64:])
	binary.LittleEndian.PutUint64(dst[:sizeUint64], nonce1^mask1)
	binary.LittleEndian.PutUint32(dst[sizeUint64:], nonce2^mask2)
}

// NonceSize returns the required nonce size.
func (s *rekeyAEAD) NonceSize() int {
	return s.gcmAEAD.NonceSize()
}

// Overhead returns the ciphertext overhead.
func (s *rekeyAEAD) Overhead() int {
	return s.gcmAEAD.Overhead()
}

// hkdfExpand computes the first 16 bytes of the HKDF-expand function
// defined in RFC5869.
func hkdfExpand(key, info []byte, keyLen int) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(info)
	mac.Write([]byte{0x01}[:])
	return mac.Sum(nil)[:keyLen]
}

type rekeyAEADChaChaPoly struct {
	kdfKey     [xchacha20poly1305ietf.KeyBytes]byte
	nonceMask  [xchacha20poly1305ietf.NonceBytes]byte
	nonceBuf   [xchacha20poly1305ietf.NonceBytes]byte
	kdfCounter []byte
}

func newRekeyAEADChaChaPoly(key []byte) (*rekeyAEADChaChaPoly, error) {
	k := len(key)
	if k != xchacha20poly1305ietf.KeyBytes+xchacha20poly1305ietf.NonceBytes {
		return nil, KeySizeError(k)
	}

	c := &rekeyAEADChaChaPoly{
		kdfCounter: make([]byte, kdfCounterLen),
	}
	copy(c.kdfKey[:], key[:xchacha20poly1305ietf.KeyBytes])
	copy(c.nonceMask[:], key[xchacha20poly1305ietf.KeyBytes:])
	return c, nil
}

func (s *rekeyAEADChaChaPoly) Encrypt(plaintext, nonce []byte) ([]byte, error) {
	s.rekeyIfRequired(nonce)
	maskNonce(s.nonceBuf[:], nonce, s.nonceMask[:])
	ec := xchacha20poly1305ietf.Encrypt(plaintext, nil, &s.nonceBuf, &s.kdfKey)
	return ec, nil
}

func (s *rekeyAEADChaChaPoly) Decrypt(ciphertext, nonce []byte) ([]byte, error) {
	s.rekeyIfRequired(nonce)
	maskNonce(s.nonceBuf[:], nonce, s.nonceMask[:])
	p, err := xchacha20poly1305ietf.Decrypt(ciphertext, nil, &s.nonceBuf, &s.kdfKey)
	if err != nil {
		return nil, fmt.Errorf("decryption unexpectedly succeeded for %s", err)
	}
	return p, nil
}

func (s *rekeyAEADChaChaPoly) rekeyIfRequired(nonce []byte) {
	newKdfCounter := nonce[kdfCounterOffset : kdfCounterOffset+kdfCounterLen]
	if bytes.Equal(newKdfCounter, s.kdfCounter) {
		return
	}
	copy(s.kdfCounter, newKdfCounter)
	copy(s.kdfKey[:], hkdfExpand(s.kdfKey[:], s.kdfCounter, xchacha20poly1305ietf.KeyBytes))
}
