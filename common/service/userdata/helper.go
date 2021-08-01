package userdata

import (
	"bytes"
	"io"

	"golang.org/x/crypto/sha3"
)

// Sha3256hash generate hash value for the data
func Sha3256hash(msg []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(msg)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}
