package pqsignature

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// Encode splits input msg into chunks to fit max supported length of QR code message and generates an array of QR codes images.
func Encode(msg string, alias string, outputDir string, outputFileTitle string, outputFileNamePattern string, outputFileNameSuffix string) ([]Image, error) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0770); err != nil {
			return nil, errors.New(err)
		}
	}
	pngs, err := toPngs(msg, MaxMsgLength, DataQRImageSize)
	if err != nil {
		return nil, err
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
			alias: alias,
		})
	}
	return images, nil
}

func toPngs(s string, dataSize int, qrImageSize int) ([][]byte, error) {
	chunks := breakStringIntoChunks(s, dataSize)
	var qrs [][]byte
	for _, chunk := range chunks {
		enc := qrcode.NewQRCodeWriter()
		image, err := enc.Encode(chunk, gozxing.BarcodeFormat_QR_CODE, qrImageSize, qrImageSize, nil)
		if err != nil {
			return nil, errors.New(err)
		}
		pngBuffer := new(bytes.Buffer)
		err = png.Encode(pngBuffer, image)
		if err != nil {
			return nil, errors.New(err)
		}
		qrs = append(qrs, pngBuffer.Bytes())
	}
	return qrs, nil
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
