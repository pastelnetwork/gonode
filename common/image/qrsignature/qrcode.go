package qrsignature

import "image"

// QRCode represents QR code image and its coordinates on the canvas.
type QRCode struct {
	image.Image
	X int
	Y int
}

// NewQRCode returns a new QRCode instance.
func NewQRCode(img image.Image) *QRCode {
	return &QRCode{
		Image: img,
		X:     -1,
		Y:     -1,
	}
}

// NewEmptyQRCode returns a new QRCode instance with empty image.
func NewEmptyQRCode(x, y int) *QRCode {
	return NewQRCode(&image.NRGBA{
		Rect: image.Rectangle{
			Max: image.Point{
				X: x,
				Y: y,
			},
		},
	})
}
