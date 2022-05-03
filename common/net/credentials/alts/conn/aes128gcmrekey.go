package conn

import (
	"crypto/cipher"
	"fmt"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
)

const (
	// Overflow length n in bytes, never encrypt more than 2^(n*8) frames (in
	// each direction).
	overflowLenAES128GCMRekey = 8
	nonceLen                  = 12
	aeadKeyLen                = 16
	kdfKeyLen                 = 32
	kdfCounterOffset          = 2
	kdfCounterLen             = 6
	sizeUint64                = 8
)

// aes128gcmRekey is the struct that holds necessary information for ALTS record.
// The counter value is NOT included in the payload during the encryption and
// decryption operations.
type aes128gcmRekey struct {
	// inCounter is used in ALTS record to check that incoming counters are
	// as expected, since ALTS record guarantees that messages are unwrapped
	// in the same order that the peer wrapped them.
	inCounter  Counter
	outCounter Counter
	inAEAD     cipher.AEAD
	outAEAD    cipher.AEAD
}

// NewAES128GCMRekey creates an instance that uses aes128gcm with rekeying
// for ALTS record. The key argument should be 44 bytes, the first 32 bytes
// are used as a key for HKDF-expand and the remainining 12 bytes are used
// as a random mask for the counter.
func NewAES128GCMRekey(side alts.Side, key []byte) (ALTSRecordCrypto, error) {
	inCounter := NewInCounter(side, overflowLenAES128GCMRekey)
	outCounter := NewOutCounter(side, overflowLenAES128GCMRekey)
	inAEAD, err := newRekeyAEAD(key)
	if err != nil {
		return nil, err
	}
	outAEAD, err := newRekeyAEAD(key)
	if err != nil {
		return nil, err
	}
	return &aes128gcmRekey{
		inCounter,
		outCounter,
		inAEAD,
		outAEAD,
	}, nil
}

// Encrypt is the encryption function. dst can contain bytes at the beginning of
// the ciphertext that will not be encrypted but will be authenticated. If dst
// has enough capacity to hold these bytes, the ciphertext and the tag, no
// allocation and copy operations will be performed. dst and plaintext do not
// overlap.
func (s *aes128gcmRekey) Encrypt(dst, plaintext []byte) ([]byte, error) {
	// If we need to allocate an output buffer, we want to include space for
	// GCM tag to avoid forcing ALTS record to reallocate as well.
	dlen := len(dst)
	dst, out := SliceForAppend(dst, len(plaintext)+GcmTagSize)
	seq, err := s.outCounter.Value()
	if err != nil {
		return nil, err
	}
	data := out[:len(plaintext)]
	copy(data, plaintext) // data may alias plaintext

	// Seal appends the ciphertext and the tag to its first argument and
	// returns the updated slice. However, SliceForAppend above ensures that
	// dst has enough capacity to avoid a reallocation and copy due to the
	// append.
	if dst = s.outAEAD.Seal(dst[:dlen], seq, data, nil); dst == nil {
		return nil, fmt.Errorf("seal returned empty data")
	}
	s.outCounter.Inc()
	return dst, nil
}

func (s *aes128gcmRekey) EncryptionOverhead() int {
	return GcmTagSize
}

func (s *aes128gcmRekey) Decrypt(dst, ciphertext []byte) ([]byte, error) {
	seq, err := s.inCounter.Value()
	if err != nil {
		return nil, err
	}
	plaintext, err := s.inAEAD.Open(dst, seq, ciphertext, nil)
	if err != nil {
		return nil, ErrAuthorize
	}
	s.inCounter.Inc()
	return plaintext, nil
}
