package test

import (
	"image"
	"image/png"
	"os"
)

//CreateBlankImage create blank png image for testing purpose only
func CreateBlankImage(filePath string, w, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	png.Encode(file, img)
	return nil
}
