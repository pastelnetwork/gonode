package test

import (
	"image"
	"image/png"
	"io/ioutil"
)

//CreateBlankImage create blank png image for testing purpose only
func CreateBlankImage(filePath string, w, h int) (string, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	file, err := ioutil.TempFile(filePath, "*.png")
	if err != nil {
		return "", err
	}

	defer file.Close()

	png.Encode(file, img)
	return file.Name(), nil
}
