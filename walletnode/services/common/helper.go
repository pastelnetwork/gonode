package common

import (
	b64 "encoding/base64"
	"fmt"
	"time"

	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	//ErrEmptyDatahash error when empty data hash
	ErrEmptyDatahash = errors.Errorf("empty data hash")
	// ErrEmptyFingerprints error when empty fingerprints
	ErrEmptyFingerprints = errors.Errorf("empty fingerprints")
	// ErrEmptyPreviewHash error when empty preview
	ErrEmptyPreviewHash = errors.Errorf("empty preview or thumbnail hashes")
	// ErrEmptyRaptorQSymbols error when empty RaptorQ symbols identifiers
	ErrEmptyRaptorQSymbols = errors.Errorf("empty RaptorQ symbols identifiers")
)

const (
	//ProbeImageSNCallTimeout sets the timeout for probeImage call to SN from WN
	ProbeImageSNCallTimeout = 30 * time.Minute
)

// InIntRange checks value in the range
func InIntRange(val int, min *int, max *int) bool {
	if min != nil && val < *min {
		return false
	}
	if max != nil && val > *max {
		return false
	}

	return true
}

// InFloatRange checks value in the range
func InFloatRange(val float64, min *float64, max *float64) bool {
	if min != nil && val < *min {
		return false
	}
	if max != nil && val > *max {
		return false
	}

	return true
}

// FromBase64 checks value in the range
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
