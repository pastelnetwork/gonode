package qrsignature

import (
	"image"
	"math"
)

const (
	textYPadding = 20
	rowYPadding  = 50
)

type Signature []*Data

func (sig Signature) images() []image.Image {
	var imgs []image.Image
	for _, data := range sig {
		imgs = append(imgs, data.images...)
	}
	return imgs
}

func (sig Signature) maxQRCodeSize() (int, int) {
	var maxWidth, maxHeight int

	for _, img := range sig.images() {
		x := img.Bounds().Size().X
		y := img.Bounds().Size().Y
		if maxWidth < x {
			maxWidth = x
		}
		if maxHeight < y {
			maxHeight = y
		}
	}

	return maxWidth, maxHeight + textYPadding + rowYPadding
}

func (sig Signature) artImageOverflow(artW, artH, qrW, qrH, qrTotal int) bool {
	capacityX := artW / qrW
	capacityY := artH / qrH

	if capacityX < 1 || capacityY < 1 {
		return true
	}

	// Substract 1 to reserve it for position vector
	totalCapacity := capacityX*capacityY - 1
	return totalCapacity < qrTotal
}

func (sig Signature) minArtImageDimension(artW, artH int) (int, int) {
	qrW, qrH := sig.maxQRCodeSize()
	qrTotal := len(sig.images())

	ratioW := float64(artH) / float64(artW)
	ratioH := float64(artW) / float64(artH)

	for {
		if !sig.artImageOverflow(artW, artH, qrW, qrH, qrTotal) {
			break
		}

		if artW > artH {
			artW++
			artH = int(math.Max(1.0, math.Floor(float64(artW)*ratioW+0.5)))
		} else {
			artH++
			artW = int(math.Max(1.0, math.Floor(float64(artH)*ratioH+0.5)))
		}
	}

	return artW, artH
}
