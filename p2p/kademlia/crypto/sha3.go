package crypto

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

// GetKey returns the key for data
func GetKey(data []byte) []byte {
	h := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(h, data)
	return []byte(fmt.Sprintf("%x", h))
}
