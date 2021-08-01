package conn

import (
	"testing"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
)

// getChaCha20Poly1305RekeyCryptoPair outputs a client/server pair on xchacha20poly1305ietfReKey.
func getChaCha20Poly1305RekeyCryptoPair(key []byte, t *testing.T) (ALTSRecordCrypto, ALTSRecordCrypto) {
	client, err := NewxChaCha20Poly1305IETFReKey(alts.ClientSide, key)
	if err != nil {
		t.Fatalf("NewAES128GCMRekey(ClientSide, key) = %v", err)
	}
	server, err := NewxChaCha20Poly1305IETFReKey(alts.ServerSide, key)
	if err != nil {
		t.Fatalf("NewAES128GCMRekey(ServerSide, key) = %v", err)
	}
	return client, server
}

// Test encrypt and decrypt on roundtrip messages for xchacha20poly1305ietfReKey.
func TestChaCha20Poly1305RekeyEncryptRoundtrip(t *testing.T) {
	// Test for xchacha20poly1305ietfReKey.
	key := make([]byte, 56)
	client, server := getChaCha20Poly1305RekeyCryptoPair(key, t)
	testRekeyEncryptRoundtrip(client, server, t)
}
