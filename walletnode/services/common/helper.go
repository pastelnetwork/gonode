package common

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	ErrEmptyFingerprints        = errors.Errorf("empty fingerprints")
	ErrEmptyDatahash            = errors.Errorf("empty data hash")
	ErrEmptyPreviewHash         = errors.Errorf("empty preview hash")
	ErrEmptyMediumThumbnailHash = errors.Errorf("empty medium thumbnail hash")
	ErrEmptySmallThumbnailHash  = errors.Errorf("empty small thumbnail hash")
	ErrEmptyRaptorQSymbols      = errors.Errorf("empty RaptorQ symbols identifiers")
)

// GetPubKey gets ED448 or PQ public key from base58 encoded key
func GetPubKey(in string) (key []byte, err error) {
	dec := base58.Decode(in)
	if len(dec) < 2 {
		return nil, errors.New("GetPubKey: invalid id")
	}

	return dec[2:], nil
}
