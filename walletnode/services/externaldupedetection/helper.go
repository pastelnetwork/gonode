package externaldupedetetion

import (
	"bytes"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/crypto/sha3"
)

var (
	errEmptyFingerprints         = errors.Errorf("empty fingerprints")
	errEmptyFingerprintsHash     = errors.Errorf("empty fingerprints hash")
	errEmptyFingerprintSignature = errors.Errorf("empty fingerprint signature")
	errEmptyDatahash             = errors.Errorf("empty data hash")
	errEmptyRaptorQSymbols       = errors.Errorf("empty RaptorQ symbols identifiers")
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
