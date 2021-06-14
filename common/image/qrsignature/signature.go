package qrsignature

import (
	"image"
	"image/color"
	"sort"

	"github.com/fogleman/gg"
	"github.com/pastelnetwork/gonode/common/log"
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

// Encode encodes the data to QR codes and groups them on a new image with maintaining aspect ratios from the given image.
func (sig Signature) Encode(img image.Image) (image.Image, error) {
	// Generate QR codes from the payload data.
	for _, payload := range sig.payloads {
		if err := payload.Encode(); err != nil {
			return nil, err
		}
	}

	// Sort QR codes by size to be more densely placed on the canvas.
	qrCodes := sig.qrCodes()

	imgW := img.Bounds().Dx()
	imgH := img.Bounds().Dy()

	log.Debug("start")

	canvas := NewCanvas(imgW, imgH)
	for _, qrCode := range qrCodes {
		// Set coordinates for each QR code.
		canvas.FillPos(qrCode)
	}
	// Obtain the real size of the canvas that can be bigger or less than the original values but while maintaining aspect ratios.
	imgW, imgH = canvas.Size()

	log.Debug("end")

	// Generate (metadata) QR code that contains coordinates of all payload QR codes on the canvas.
	if err := sig.metadata.Encode(sig.payloads); err != nil {
		return nil, err
	}

	// Draw all QR codes in on single image.
	dc := gg.NewContext(imgW, imgH)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	dc.SetColor(color.White)

	for _, qrCode := range qrCodes {
		dc.DrawImageAnchored(qrCode, qrCode.X, qrCode.Y, 0, 0)
	}
	return dc.Image(), nil
}

// func (sig Signature) Dencode(img image.Image) (image.Image, error) {

// }

// New returns a new Signature instance.
func New(payload ...*Payload) *Signature {
	return &Signature{
		metadata: NewMetadata(),
		payloads: payload,
	}
}
