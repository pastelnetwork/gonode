package qrsignature

import (
	"fmt"

	"github.com/pastelnetwork/gonode/common/collection"
)

const (
	payloadQRMinSize  = 0
	payloadQRCapacity = QRCodeCapacityBinary
)

// Payload represents an independent piece of information to be saved as a sequence of qr codes.
type Payload struct {
	raw  []byte
	name PayloadName

	writer  *QRCodeWriter
	qrCodes []*QRCode
}

// Encode encodes raw data to qrCode representation.
func (payload *Payload) Encode() error {
	raws := collection.SplitBytesBySize(payload.raw, int(payloadQRCapacity))
	qrCodes := make([]*QRCode, len(raws))

	for i, raw := range raws {
		img, err := payload.writer.Encode(raw)
		if err != nil {
			return fmt.Errorf("%q: %w", payload.name, err)
		}
		qrCodes[i] = NewQRCode(img)
	}

	payload.qrCodes = qrCodes
	return nil
}

// NewPayload returns a new Payload instance.
func NewPayload(raw []byte, name PayloadName) *Payload {
	return &Payload{
		raw:  raw,
		name: name,

		writer: NewQRCodeWriter(payloadQRMinSize),
	}
}
