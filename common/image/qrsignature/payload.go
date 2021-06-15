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
			return errors.New(err)
		}
		data += text
	}

	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return errors.Errorf("failed to decode text to binary data: %w", err)
	}

	raw, err = zstd.Decompress(nil, raw)
	if err != nil {
		return errors.Errorf("failed to decompress binary data: %w", err)
	}
	payload.raw = raw
	return nil
}

// Encode splits raw data into chunks and encodes them to QR code representation.
func (payload *Payload) Encode() error {
	raw, err := zstd.CompressLevel(nil, payload.raw, 22)
	if err != nil {
		return errors.Errorf("failed to compress binary data: %w", err)
	}
	data := base64.StdEncoding.EncodeToString(raw)
	size := int(payloadQRCapacity)

	var qrCodes []*QRCode
	for i, last := 0, len(data); i < last; i += size {
		if i+size > last {
			size = last - i
		}

		qrCode := NewQRCode(payloadQRMinSize)
		if err := qrCode.Encode(data[i : i+size]); err != nil {
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
