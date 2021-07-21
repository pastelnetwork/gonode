

package userdata

import (
	"bytes"
	"io"

	"golang.org/x/crypto/sha3"
)

func sha3256hash(msg []byte) ([]byte, error) {
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, bytes.NewReader(msg)); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func safeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}