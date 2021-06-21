package qrsignature

import (
	"math"
)

// CanvasAxis keeps the farthest occupied pixels for the axis.
type CanvasAxis []int

func (points *CanvasAxis) maxPos(start, length int) int {
	var max int
	for i := start; i < start+length; i++ {
		if max < (*points)[i] {
			max = (*points)[i]
		}
	}
	return max
}

func (points *CanvasAxis) setPos(start, length, pos int) {
	for i := start; i < start+length; i++ {
		(*points)[i] = pos
	}
}

func (points *CanvasAxis) setSize(cap int) {
	p := *points
	(*points) = make(CanvasAxis, cap)
	copy((*points), p)
}

// Canvas represents a matrix to keep allocated points to qrcodes for each axis.
type Canvas struct {
	ratioW float64
	ratioH float64

	width  int
	height int

	xy CanvasAxis
	yx CanvasAxis
}

// Size returns the canvas size with clipping black edges.
func (canvas *Canvas) Size() (int, int) {
	return canvas.width, canvas.height
}

// FillPos finds free space in the canvas and fills position points for the given qrcode.
func (canvas *Canvas) FillPos(qrCodes []*QRCode) {
	for i := 0; i < 2; i++ {
		for _, qrCode := range qrCodes {
			canvas.fillPos(qrCode)
		}
		// reset axis
		canvas.xy = make(CanvasAxis, canvas.width)
		canvas.yx = make(CanvasAxis, canvas.height)
	}
}

func (canvas *Canvas) fillPos(qrCode *QRCode) {
	size := qrCode.Bounds().Size()

	if canvas.width == 0 {
		canvas.scale(size.X)
	}

	var posX, posY int
	for {
		posX, posY = canvas.width, canvas.height

		for x, y := 0, 0; x <= canvas.width-size.X; x++ {
			if y = canvas.xy.maxPos(x, size.X); y+size.Y > canvas.height {
				continue
			}
			if max := canvas.yx.maxPos(y, size.Y); x < max {
				continue
			}
			if posY > y {
				posY = y
				posX = x
			}
		}

		if posX+size.X > canvas.width || posY+size.Y > canvas.height {
			canvas.scale(1)
			continue
		}
		break
	}

	canvas.xy.setPos(posX, size.X, size.Y+posY)
	canvas.yx.setPos(posY, size.Y, size.X+posX)

	qrCode.X = posX
	qrCode.Y = posY
}

func (canvas *Canvas) scale(pixels int) {
	if canvas.ratioH == canvas.ratioW {
		canvas.width += pixels
		canvas.height += pixels
	} else if canvas.ratioH < canvas.ratioW {
		canvas.width += pixels
		canvas.height = int(math.Max(1.0, math.Floor(float64(canvas.width)*canvas.ratioW+0.5)))
	} else {
		canvas.height += pixels
		canvas.width = int(math.Max(1.0, math.Floor(float64(canvas.height)*canvas.ratioH+0.5)))
	}

	canvas.xy.setSize(canvas.width)
	canvas.yx.setSize(canvas.height)
}

// NewCanvas returns a new Canvas instance.
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		ratioW: float64(height) / float64(width),
		ratioH: float64(width) / float64(height),
	}
}
