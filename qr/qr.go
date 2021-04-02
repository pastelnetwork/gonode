// Package qr generates image sequences of QR-codes from the input messages of any length
package qr

import (
	"fmt"
	"os"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
)

func Encode(msg string, outputDir string, outputFileNamePattern string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0770); err != nil {
			return err
		}
	}
	maxMsgLength := 2200
	pngs, err := ToPngs(msg, maxMsgLength)
	if err != nil {
		return err
	}
	for i, imageData := range pngs {
		partNumber := ""
		if len(pngs) > 1 {
			partNumber = fmt.Sprintf("__part_%v_of_%v.png", i+1, len(pngs))
		}
		err := os.WriteFile(filepath.Join(outputDir, fmt.Sprintf("%v%v.png", outputFileNamePattern, partNumber)), []byte(imageData), 0644)
		if err != nil {
			return err
		}
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

func ToPngs(s string, dataSize int) ([][]byte, error) {
	chunks := breakStringIntoChunks(s, dataSize)
	var qrs [][]byte
	for _, chunk := range chunks {
		png, err := qrcode.Encode(chunk, qrcode.Medium, 250)
		if err != nil {
			return nil, err
		}
		qrs = append(qrs, png)
	}
	return qrs, nil
}
