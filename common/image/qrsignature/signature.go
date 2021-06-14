package qrsignature

import (
	"image"
	"image/color"
	"sort"

	"github.com/fogleman/gg"
)

// Signature represents set of various data(payloads) that is represented in QR codes on the single canvas.
type Signature struct {
	metadata *Metadata
	payloads []*Payload
}

func (sig *Signature) qrCodes() []*QRCode {
	var qrCodes []*QRCode

	for _, payload := range sig.payloads {
		qrCodes = append(qrCodes, payload.qrCodes...)
	}

	sort.Slice(qrCodes, func(i, j int) bool {
		return qrCodes[i].Bounds().Size().X > qrCodes[j].Bounds().Size().X
	})

	return append([]*QRCode{sig.metadata.qrCode}, qrCodes...)
}

// Encode encodes the data to QR codes and groups them on a new image with maintaining aspect ratios from the given `width`, `height` params.
func (sig Signature) Encode(width, height int) (image.Image, error) {
	// Generate QR codes from the payload data.
	for _, payload := range sig.payloads {
		if err := payload.Encode(); err != nil {
			return nil, err
		}
	}

	// Sort QR codes by size to be more densely placed on the canvas.
	qrCodes := sig.qrCodes()

	// Set coordinates for each QR code.
	canvas := NewCanvas(width, height)
	for _, qrCode := range qrCodes {
		canvas.FillPos(qrCode)
	}

	// Generate (metadata) QR code that contains coordinates of all payload QR codes on the canvas.
	if err := sig.metadata.Encode(sig.payloads); err != nil {
		return nil, err
	}

	// Obtain the real size of the canvas that can be bigger or less than the original values.
	width, height = canvas.Size()

	// Create a new image from the QR codes.
	dc := gg.NewContext(width, height)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	dc.SetColor(color.White)

	for _, qrCode := range qrCodes {
		dc.DrawImageAnchored(qrCode, qrCode.X, qrCode.Y, 0, 0)
	}

	return dc.Image(), nil
}

// New returns a new Signature instance.
func New(payload ...*Payload) *Signature {
	return &Signature{
		metadata: NewMetadata(),
		payloads: payload,
	}
}
