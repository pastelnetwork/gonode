package secret

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAndDecrypt(t *testing.T) {
	// generate a random 32 byte key for AES-256
	key := make([]byte, 32)
	rand.Read(key)

	source := make([]byte, 256*256)
	rand.Read(source)

	encrypted, err := aesEncrypt(source, key)
	if err != nil {
		t.Fatalf("aes encrypt: %v", err)
	}
	decrypted, err := aesDecrypt(encrypted, key)
	if err != nil {
		t.Fatalf("aes decrypt: %v", err)
	}
	assert.Equal(t, source, decrypted)
}
