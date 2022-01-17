package common

import (
	"github.com/pastelnetwork/gonode/pastel"
)

type FingerprintsHandler struct {
	FingerprintAndScores      *pastel.DDAndFingerprints
	FingerprintAndScoresBytes []byte // JSON bytes of FingerprintAndScores
	Signature                 []byte
}
