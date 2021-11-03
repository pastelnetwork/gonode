package qrsignature

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/steganography"
)

// Signature represents set of various data(payloads) that is represented in QR codes on the single canvas.
type Signature struct {
	*Payloads
	metadata *Metadata
}

func (sig *Signature) qrCodes() []*QRCode {
	var qrCodes []*QRCode

	for _, payload := range *sig.Payloads {
		qrCodes = append(qrCodes, payload.qrCodes...)
	}

	return append([]*QRCode{sig.metadata.qrCode}, qrCodes...)
}

// Encode encodes the data to QR codes and groups them on a new image with maintaining aspect ratios from the given image.
func (sig Signature) Encode(img image.Image) (image.Image, error) {
	// Generate QR codes from the payload data.
	for _, payload := range *sig.Payloads {
		if err := payload.Encode(); err != nil {
			return nil, err
		}
	}

	// Sort QR codes by size to be more densely placed on the canvas.
	qrCodes := sig.qrCodes()

	imgW := img.Bounds().Dx()
	imgH := img.Bounds().Dy()

	canvas := NewCanvas(imgW, imgH)
	// Set coordinates for each QR code.
	canvas.FillPos(qrCodes)
	// Obtain the real size of the canvas that can be bigger or less than the original values but while maintaining aspect ratio.
	imgW, imgH = canvas.Size()

	// Generate (metadata) QR code that contains coordinates of all payload QR codes on the canvas.
	if err := sig.metadata.Encode(*sig.Payloads); err != nil {
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
	sigImg := dc.Image()

	sigData := new(bytes.Buffer)
	if err := png.Encode(sigData, sigImg); err != nil {
		return nil, errors.Errorf("encode signature: %w", err)
	}

	if sigW := sigImg.Bounds().Dx(); sigW > img.Bounds().Dx() {
		img = imaging.Resize(img, sigW, 0, imaging.Lanczos)
	}

	img, err := steganography.Encode(img, sigData.Bytes())
	if err != nil {
		return nil, errors.New(err)
	}
	return img, nil
}

// Decode decodes QR codes from the given image to raw data. The top-left QR code should contain the metadata of the rest.
func (sig Signature) Decode(img image.Image) error {
	msgLen := steganography.MaxEncodeSize(img)
	data := steganography.Decode(msgLen, img)

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return errors.Errorf("decode image: %w", err)
	}

	if err := sig.metadata.qrCode.Load(img); err != nil {
		return errors.Errorf("load: %w", err)
	}

	if err := sig.metadata.Decode(sig.Payloads); err != nil {
		return errors.Errorf("metadata decode: %w", err)
	}

	for _, payload := range *sig.Payloads {
		for _, qrCode := range payload.qrCodes {
			if err := qrCode.Load(img); err != nil {
				return errors.Errorf("qrcode load: %w", err)
			}
		}

		if err := payload.Decode(); err != nil {
			return errors.Errorf("decode payload: %w", err)
		}
	}
	return nil
}

// New returns a new Signature instance.
func New(payloads ...*Payload) *Signature {
	return &Signature{
		Payloads: NewPayloads(payloads...),
		metadata: NewMetadata(),
	}
}
