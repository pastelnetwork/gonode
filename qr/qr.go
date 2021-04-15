// Package qr generates image sequences of QR-codes from the input messages of any length
package qr

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/pastelnetwork/go-commons/errors"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	MaxMsgLength = 2200
)

type Image struct {
	raw   []byte
	image image.Image
	title string
}

// Encode splits input msg into chunks to fit max supported length of QR code message and generates an array of QR codes images
func Encode(msg string, outputDir string, outputFileTitle string, outputFileNamePattern string, outputFileNameSuffix string) ([]Image, error) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0770); err != nil {
			return nil, errors.New(err)
		}
	}
	pngs, err := toPngs(msg, MaxMsgLength)
	if err != nil {
		return nil, errors.New(err)
	}
	var images []Image
	for i, imageBytes := range pngs {
		filePathPartNumber := ""
		titlePartNumber := ""
		if len(pngs) > 1 {
			filePathPartNumber = fmt.Sprintf("__part_%v_of_%v", i+1, len(pngs))
			titlePartNumber = fmt.Sprintf(" part %v of %v", i+1, len(pngs))
		}
		partTitle := fmt.Sprintf("%v%v", outputFileTitle, titlePartNumber)
		filePath := filepath.Join(outputDir, fmt.Sprintf("%v%v%v.png", outputFileNamePattern, filePathPartNumber, outputFileNameSuffix))
		err := os.WriteFile(filePath, imageBytes, 0644)
		if err != nil {
			return nil, errors.New(err)
		}
		img, err := gg.LoadImage(filePath)
		if err != nil {
			return nil, errors.New(err)
		}
		images = append(images, Image{
			raw:   imageBytes,
			title: partTitle,
			image: img,
		})
	}
	return images, nil
}

func LoadImages(pattern string) ([]Image, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, errors.New(err)
	}
	var images []Image
	for _, filePath := range matches {
		img, err := gg.LoadImage(filePath)
		if err != nil {
			return nil, errors.New(err)
		}
		fileInfo, err := os.Lstat(filePath)
		if err != nil {
			return nil, errors.New(err)
		}
		imageBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, errors.New(err)
		}

		images = append(images, Image{
			raw:   imageBytes,
			title: fileInfo.Name(),
			image: img,
		})
	}
	return images, nil
}

// MapImages maps input images into output image of specified size
func MapImages(images []Image, outputSize image.Point, outputFilePath string) error {
	dc := gg.NewContext(outputSize.X, outputSize.Y)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	dc.SetColor(color.White)
	padding_pixels := 2
	currentX := 0
	currentY := 10
	textYPadding := 20
	rowYPadding := 50
	for _, image := range images {
		size := image.image.Bounds().Size()
		captionX := size.X / 2

		if currentX+size.X > outputSize.X {
			currentX = 0
			currentY += size.Y + rowYPadding
		}

		dc.DrawStringAnchored(image.title, float64(currentX+captionX), float64(currentY), 0.5, 0.5)
		dc.DrawImageAnchored(image.image, currentX+captionX, currentY+textYPadding+size.Y/2, 0.5, 0.5)
		currentX += size.X + padding_pixels
	}

	err := dc.SavePNG(outputFilePath)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func breakStringIntoChunks(str string, size int) []string {
	if len(str) <= size {
		return []string{str}
	}
	var result []string
	chunk := make([]rune, size)
	pos := 0
	for _, rune := range str {
		chunk[pos] = rune
		pos++
		if pos == size {
			pos = 0
			result = append(result, string(chunk))
		}
	}
	if pos > 0 {
		result = append(result, string(chunk[:pos]))
	}
	return result
}

func toPngs(s string, dataSize int) ([][]byte, error) {
	chunks := breakStringIntoChunks(s, dataSize)
	var qrs [][]byte
	for _, chunk := range chunks {
		png, err := qrcode.Encode(chunk, qrcode.Medium, 250)
		if err != nil {
			return nil, errors.New(err)
		}
		qrs = append(qrs, png)
	}
	return qrs, nil
}
