package common

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	ErrEmptyDatahash       = errors.Errorf("empty data hash")
	ErrEmptyFingerprints   = errors.Errorf("empty fingerprints")
	ErrEmptyPreviewHash    = errors.Errorf("empty preview or thumbnail hashes")
	ErrEmptyRaptorQSymbols = errors.Errorf("empty RaptorQ symbols identifiers")
)

// GetPubKey gets ED448 or PQ public key from base58 encoded key
func GetPubKey(in string) (key []byte, err error) {
	dec := base58.Decode(in)
	if len(dec) < 2 {
		return nil, errors.New("GetPubKey: invalid id")
	}

	return dec[2:], nil
}

func InIntRange(val int, min *int, max *int) bool {
	if min != nil && val < *min {
		return false
	}
	if max != nil && val > *max {
		return false
	}

	return true
}

func InFloatRange(val float64, min *float64, max *float64) bool {
	if min != nil && val < *min {
		return false
	}
	if max != nil && val > *max {
		return false
	}

	return true
}

func FromBase64(encoded string, to interface{}) error {
	bytes, err := b64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("decode: %s", err)
	}

	if err := json.Unmarshal(bytes, to); err != nil {
		return fmt.Errorf("unmarshal: %s", err)
	}

	return nil
}
