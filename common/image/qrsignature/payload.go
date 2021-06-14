package qrsignature

import (
	"encoding/base64"
	"fmt"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
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
	output, err := zstd.CompressLevel(nil, payload.raw, 22)
	if err != nil {
		return errors.Errorf("failed to compress raw data: %w", err)
	}
	data := base64.StdEncoding.EncodeToString(output)
	size := int(payloadQRCapacity)

	var qrCodes []*QRCode
	for i, last := 0, len(data); i < last; i += size {
		if i+size > last {
			size = last - i
		}
		img, err := payload.writer.Encode(data[i : i+size])
		if err != nil {
			return fmt.Errorf("%q: %w", payload.name, err)
		}
		qrCodes = append(qrCodes, NewQRCode(img))
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
