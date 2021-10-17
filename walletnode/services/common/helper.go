package common

import (
	"bytes"
	"io"

	"golang.org/x/crypto/sha3"
)

// SafeString returns string values of string pointer, empty string if nil pointer
func SafeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

// Sha3256hash func
func Sha3256hash(msg []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(msg)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}
