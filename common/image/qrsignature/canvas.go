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
	return canvas.yx.maxPos(0, canvas.height), canvas.xy.maxPos(0, canvas.width)
}

// FillPos finds free space in the canvas and fills position points for the given qrcode.
func (canvas *Canvas) FillPos(qrCode *QRCode) {
	size := qrCode.Bounds().Size()

	if canvas.width == 0 {
		canvas.scale(size.X)
	}

	var posX, posY int
	for {
		posX, posY = canvas.width, canvas.height

		for x := 0; x <= canvas.width-size.X; x++ {
			var y int

			if max := canvas.xy.maxPos(x, size.X); y < max {
				if canvas.height < max+size.Y {
					continue
				}
				y = max
			}

			if max := canvas.yx.maxPos(y, size.Y); x < max {
				continue
			}
			if posY > y {
				posY = y
				posX = x
			}
		}

		if posX == canvas.width {
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
func NewCanvas(artW, artH int) *Canvas {
	return &Canvas{
		ratioW: float64(artH) / float64(artW),
		ratioH: float64(artW) / float64(artH),
	}
}
