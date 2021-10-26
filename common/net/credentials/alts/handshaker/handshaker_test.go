package handshaker

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"

	"github.com/cloudflare/circl/dh/x448"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
)

const (
	aeadSizeOverhead = chacha20poly1305.Overhead // overhead of poly 1305 authentication tag
	aeadKeySize      = chacha20poly1305.KeySize
	aeadNonceSize    = chacha20poly1305.NonceSizeX
)

var (
	secretConnKeySalt = []byte("PASTEL_SECRET_CONNECTION_KEY_SALT")
	secretConnKeyInfo = []byte("PASTEL_SECRET_CONNECTION_KEY_INFO")
)

func TestCirclX448(t *testing.T) {
	var secretClient x448.Key
	rand.Read(secretClient[:])
	var publicClient x448.Key
	x448.KeyGen(&publicClient, &secretClient)

	var secretServer x448.Key
	_, err := rand.Read(secretServer[:])
	if err != nil {
		t.Fatalf("err in rand.Read: %v", err)
	}

	var publicServer x448.Key
	x448.KeyGen(&publicServer, &secretServer)

	var sharedClient x448.Key
	if ok := x448.Shared(&sharedClient, &secretClient, &publicServer); !ok {
		t.Fatal("x448 shared key for client")
	}
	var sharedServer x448.Key
	if ok := x448.Shared(&sharedServer, &secretServer, &publicClient); !ok {
		t.Fatal("x448 shared key for server")
	}
	assert.Equal(t, sharedClient, sharedServer)

	hash := sha3.New512
	hkdf := hkdf.New(hash, sharedClient[:], secretConnKeySalt, secretConnKeyInfo)
	out := make([]byte, aeadKeySize)
	_, err = io.ReadFull(hkdf, out)
	if err != nil {
		t.Fatalf("io read full: %v", err)
	}
	t.Logf("key after hkdf: %v", hex.EncodeToString(out))

	secretAead, err := chacha20poly1305.NewX(out)
	if err != nil {
		t.Fatalf("new secret aead: %v", err)
	}
	t.Logf("overhead: %v", secretAead.Overhead())

	// data := hex.EncodeToString(out)
	// secretAead.Seal()
}
