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
	positionQRCapacity = qrCodeCapacityBinary
	positionQRTitle    = "QR Code Position Data"
)

// Metadata represents service information about QR codes such as their names, coordinates on the canvas.
type Metadata struct {
	qrCode *QRCode
}

// Decode decodes metadata from QR code representation.
func (metadata *Metadata) Decode(payloads *Payloads) error {
	data, err := metadata.qrCode.Decode()
	if err != nil {
		return errors.New(err)
	}

	payloadsData := strings.Split(strings.Trim(data, ";"), ";")
	if len(payloadsData) == 0 {
		return errors.Errorf("could not find payload data: %q", data)
	}

	for _, payloadData := range payloadsData {
		d := strings.Split(payloadData, ":")
		if len(d) != 2 {
			return errors.Errorf("could not parse payload data: %q", payloadData)
		}

		pos := strings.Split(d[1], ",")
		if len(pos) == 0 || math.Mod(float64(len(pos)), 3) != 0 {
			return errors.Errorf("could not parse positions of QR codes: %q", payloadData)
		}

		name := payloadName(d[0])
		if name.String() == "" {
			return errors.Errorf("empty data name: %q", payloadData)
		}
		payload := &Payload{
			name: name,
		}

		for i := 0; i < len(pos); i += 3 {
			x, err := strconv.Atoi(pos[i])
			if err != nil {
				return errors.Errorf(`could not parse "x" coordinate: %w`, err)
			}

			y, err := strconv.Atoi(pos[i+1])
			if err != nil {
				return errors.Errorf(`could not parse "y" coordinate: %w`, err)
			}

			imageSize, err := strconv.Atoi(pos[i+2])
			if err != nil {
				return errors.Errorf(`could not parse image "size": %w`, err)
			}

			qrCode := NewQRCode(imageSize)
			qrCode.X = x
			qrCode.Y = y

			payload.qrCodes = append(payload.qrCodes, qrCode)
		}
		*payloads = append(*payloads, payload)
	}
	return nil
}

// Encode encodes metadata to qrCode representation.
func (metadata *Metadata) Encode(payloads Payloads) error {
	payloadsData := make([]string, len(payloads))

	for i, payload := range payloads {
		qrCodesPos := make([]string, len(payload.qrCodes)*3)
		for t, qrCode := range payload.qrCodes {
			t := t * 3
			qrCodesPos[t] = strconv.Itoa(qrCode.X)
			qrCodesPos[t+1] = strconv.Itoa(qrCode.Y)
			qrCodesPos[t+2] = strconv.Itoa(qrCode.Bounds().Size().X)
		}
		payloadsData[i] = fmt.Sprintf("%v:%v", payload.name.String(), strings.Join(qrCodesPos, ","))
	}
	data := ";" + strings.Join(payloadsData, ";")
	if err := positionQRCapacity.Validate(data); err != nil {
		return err
	}
	return metadata.qrCode.Encode(positionQRTitle, data)
}

// NewMetadata returns a new Metadata instance.
func NewMetadata() *Metadata {
	return &Metadata{
		qrCode: NewQRCode(positionQRSize),
	}
}
