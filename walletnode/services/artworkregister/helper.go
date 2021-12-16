package artworkregister

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	errEmptyFingerprints        = errors.Errorf("empty fingerprints")
	errEmptyDatahash            = errors.Errorf("empty data hash")
	errEmptyPreviewHash         = errors.Errorf("empty preview hash")
	errEmptyMediumThumbnailHash = errors.Errorf("empty medium thumbnail hash")
	errEmptySmallThumbnailHash  = errors.Errorf("empty small thumbnail hash")
	errEmptyRaptorQSymbols      = errors.Errorf("empty RaptorQ symbols identifiers")
)

// getPubKey gets ED448 or PQ public key from base58 encoded key
func getPubKey(in string) (key []byte, err error) {
	dec := base58.Decode(in)
	if len(dec) < 2 {
		return nil, errors.New("getPubKey: invalid id")
	}

	return dec[2:], nil
}
