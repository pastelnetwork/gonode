package artworkregister

import (
	"github.com/pastelnetwork/gonode/common/errors"
)

var (
	errEmptyFingerprints         = errors.Errorf("empty fingerprints")
	errEmptyFingerprintsHash     = errors.Errorf("empty fingerprints hash")
	errEmptyFingerprintSignature = errors.Errorf("empty fingerprint signature")
	errEmptyDatahash             = errors.Errorf("empty data hash")
	errEmptyPreviewHash          = errors.Errorf("empty preview hash")
	errEmptyMediumThumbnailHash  = errors.Errorf("empty medium thumbnail hash")
	errEmptySmallThumbnailHash   = errors.Errorf("empty small thumbnail hash")
	errEmptyRaptorQSymbols       = errors.Errorf("empty RaptorQ symbols identifiers")
)
