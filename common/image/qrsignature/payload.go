package qrsignature

import (
	"context"
	"fmt"
	"math"

	"github.com/pastelnetwork/gonode/common/b85"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
)

const (
	payloadQRMinSize  = 185
	payloadQRCapacity = qrCodeCapacityBinary
	payloadQRTitleFmt = "%s part %d of %d"
)

// Payloads represents multiple Payload.
type Payloads []*Payload

// Raw returns raw data by payload name.
func (payloads *Payloads) Raw(name PayloadName) []byte {
	for _, payload := range *payloads {
		if payload.name == name {
			return payload.raw
		}
	}
	return nil
}

// NewPayloads returns a new Payloads instance.
func NewPayloads(payloads ...*Payload) *Payloads {
	p := Payloads(payloads)
	return &p
}

// Payload represents an independent piece of information to be saved as a sequence of qr codes.
type Payload struct {
	raw  []byte
	name PayloadName

	qrCodes []*QRCode
}

// Decode decodes raw data from QR code representation.
func (payload *Payload) Decode() error {
	var data string

	for _, qrCode := range payload.qrCodes {
		text, err := qrCode.Decode()
		if err != nil {
			return fmt.Errorf("qrDecode: %v", err)
		}
		data += text
	}

	raw, err := b85.Decode(data)
	if err != nil {
		return fmt.Errorf("b85Decode: %v", err)
	}

	if !(payload.name == PayloadPostQuantumPubKey || payload.name == PayloadEd448PubKey) {
		var err error
		raw, err = utils.Decompress(raw)
		if err != nil {
			return errors.Errorf("decompress: %w", err)
		}
	}
	payload.raw = raw

	return nil
}

// Encode splits raw data into chunks and encodes them to QR code representation.
func (payload *Payload) Encode() error {
	raw := payload.raw
	if !(payload.name == PayloadPostQuantumPubKey || payload.name == PayloadEd448PubKey) {
		var err error
		raw, err = utils.HighCompress(context.Background(), payload.raw)
		if err != nil {
			return errors.Errorf("compress: %w", err)
		}
	}

	data := b85.Encode(raw)

	size := int(payloadQRCapacity)
	total := int(math.Ceil(float64(len(data)) / float64(size)))

	var qrCodes []*QRCode
	for i := 0; i < total; i++ {
		data := data[i*size:]
		if i < total-1 {
			data = data[:size]
		}

		title := payload.name.String()
		if total > 1 {
			title = fmt.Sprintf(payloadQRTitleFmt, payload.name, i+1, total)
		}
		qrCode := NewQRCode(payloadQRMinSize)
		if err := qrCode.Encode(title, data); err != nil {
			return fmt.Errorf("%q: %w", payload.name, err)
		}
		qrCodes = append(qrCodes, qrCode)
	}

	payload.qrCodes = qrCodes
	return nil
}

// NewPayload returns a new Payload instance.
func NewPayload(raw []byte, name PayloadName) *Payload {
	return &Payload{
		raw:  raw,
		name: name,
	}
}
