package test

import (
	"image"
	"image/png"
	"os"
)

//CreateBlankImage create blank png image for testing purpose only
func CreateBlankImage(filePath string, w, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	png.Encode(f, img)
	return nil
}
