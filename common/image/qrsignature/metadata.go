package qrsignature

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	positionQRSize     = 185
	positionQRCapacity = QRCodeCapacityBinary
)

// Metadata represents service information about QR codes such as their names, coordinates on the canvas.
type Metadata struct {
	qrCode *QRCode
}

// Decode decodes metadata from QR code representation.
func (metadata *Metadata) Decode() (Payloads, error) {
	data, err := metadata.qrCode.Decode()
	if err != nil {
		return nil, errors.New(err)
	}

	var payloads Payloads

	for _, payloadData := range strings.Split(data, ";") {
		d := strings.Split(payloadData, ":")
		if len(d) != 2 {
			return nil, errors.Errorf("failed to parse metadata: %q", payloadData)
		}

		pos := strings.Split(d[1], ",")
		if len(pos) == 0 || math.Mod(float64(len(pos)), 3) != 0 {
			return nil, errors.Errorf("failed to parse metadata: %q", payloadData)
		}

		name := payloadName(d[0])
		if name.String() == "" {
			return nil, errors.Errorf("empty data name: %q", payloadData)
		}
		payload := &Payload{
			name: name,
		}

		for i := 0; i < len(pos); i += 3 {
			x, err := strconv.Atoi(pos[i])
			if err != nil {
				return nil, errors.Errorf("failed to parse QR code X: %w", err)
			}

			y, err := strconv.Atoi(pos[i+1])
			if err != nil {
				return nil, errors.Errorf("failed to parse QR code Y: %w", err)
			}

			imageSize, err := strconv.Atoi(pos[i+2])
			if err != nil {
				return nil, errors.Errorf("failed to parse QR code image size: %w", err)
			}

			qrCode := NewQRCode(imageSize)
			qrCode.X = x
			qrCode.Y = y

			payload.qrCodes = append(payload.qrCodes, qrCode)
		}
		payloads = append(payloads, payload)
	}

	return payloads, nil
}

// Encode encodes metadata to qrCode representation.
func (metadata *Metadata) Encode(payloads Payloads) error {
	var data string

	for _, payload := range payloads {
		data += fmt.Sprintf("%v:", payload.name.String())
		for _, qrCode := range payload.qrCodes {
			data += fmt.Sprintf("%v,%v,%v", qrCode.X, qrCode.Y, qrCode.Bounds().Size().X)
		}
		data += ";"
	}

	if err := positionQRCapacity.Validate(strings.Trim(data, ";")); err != nil {
		return err
	}
	return metadata.qrCode.Encode(data)
}

// NewMetadata returns a new Metadata instance.
func NewMetadata() *Metadata {
	return &Metadata{
		qrCode: NewQRCode(positionQRSize),
	}
}
