package qrsignature

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	// QRCodeCapacityBinary represents size for Binary
	// ISO 8859-1
	qrCodeCapacityBinary QRCodeCapacity = 2953

	qrCodeBorder      = 1
	qrCodeTitleHeight = 40
)

// QRCodeCapacity represents maximum character storage dataCapacity that can be encoded with a single QRCodeWriter-code.
// https://stackoverflow.com/questions/12764334/qr-code-max-char-length
type QRCodeCapacity int

// Validate validates a length of the data.
func (cap QRCodeCapacity) Validate(data string) error {
	if int(cap) < len(data) {
		return errors.Errorf("exceeds the maximum allowed amount of data by %d characters", len(data)-int(cap))
	}
	return nil
}

// QRCode represents QR code image and its coordinates on the canvas.
type QRCode struct {
	image.Image
	imageSize int
	X         int
	Y         int
	writer    *qrcode.QRCodeWriter
	reader    gozxing.Reader
}

// Decode reads data from the QR code.
func (qr *QRCode) Decode() (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(qr.Image)
	if err != nil {
		return "", errors.Errorf("get bitmap from the QR code image: %w", err)
	}
	result, err := qr.reader.Decode(bmp, nil)
	if err != nil {
		return "", errors.Errorf("decode a QR code: %w", err)
	}

	return result.GetText(), nil
}

// Encode generates a new QR code image by the given data.
func (qr *QRCode) Encode(title, data string) error {
	img, err := qr.writer.Encode(data, gozxing.BarcodeFormat_QR_CODE, qr.imageSize, qr.imageSize, nil)
	if err != nil {
		return errors.Errorf("encode data: %w", err)
	}

	size := img.Bounds().Size()

	dc := gg.NewContext(size.X+qrCodeBorder*2, size.Y+qrCodeTitleHeight)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	dc.SetColor(color.White)

	dc.DrawStringAnchored(title, float64(size.X/2), float64(qrCodeTitleHeight/2), 0.5, 0.7)
	dc.DrawImageAnchored(img, qrCodeBorder, qrCodeTitleHeight, 0, 0)

	qr.Image = dc.Image()
	return nil
}

// Load crops QR code from the given image.
func (qr *QRCode) Load(img image.Image) error {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	imager, ok := img.(subImager)
	if !ok {
		return errors.New("image does not support cropping")
	}

	rect := image.Rectangle{image.Point{qr.X + qrCodeBorder, qr.Y + qrCodeTitleHeight}, image.Point{qr.X + qr.imageSize + qrCodeBorder, qr.Y + qr.imageSize + qrCodeTitleHeight}}
	qr.Image = imager.SubImage(rect)
	return nil
}

// NewQRCode returns a new QRCode instance.
func NewQRCode(imageSize int) *QRCode {
	return &QRCode{
		imageSize: imageSize,
		writer:    qrcode.NewQRCodeWriter(),
		reader:    qrcode.NewQRCodeReader(),
		Image: &image.NRGBA{
			Rect: image.Rectangle{
				Max: image.Point{
					X: imageSize + qrCodeBorder*2,
					Y: imageSize + qrCodeTitleHeight,
				},
			},
		},
	}
}
