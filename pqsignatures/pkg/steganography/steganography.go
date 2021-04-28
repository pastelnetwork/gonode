// Package steganography ciphers images with LSB steganography.
package steganography

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/fogleman/gg"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/steganography"
)

// Encode applies LSB steganography on imageFilePath to encode imageToHideFilePath and writes to encodedFilePath.
func Encode(imageFilePath string, imageToHideFilePath string, encodedFilePath string) error {
	img, err := gg.LoadImage(imageFilePath)
	if err != nil {
		return errors.New(err)
	}

	signature_layer_image_data, err := ioutil.ReadFile(imageToHideFilePath)
	if err != nil {
		return errors.New(err)
	}

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, signature_layer_image_data)
	if err != nil {
		return errors.New(err)
	}

	outFile, err := os.Create(encodedFilePath)
	if err != nil {
		return errors.New(err)
	}

	w.WriteTo(outFile)
	outFile.Close()
	return nil
}

// Decode runs reversive LSB steganography on imageFilePath and writes to decodedFilePath.
func Decode(imageFilePath string, decodedFilePath string) error {
	img, err := gg.LoadImage(imageFilePath)
	if err != nil {
		return errors.New(err)
	}
	sizeOfMessage := steganography.GetMessageSizeFromImage(img)
	decodedData := steganography.Decode(sizeOfMessage, img)
	err = os.WriteFile(decodedFilePath, decodedData, 0644)
	if err != nil {
		errors.New(err)
	}
	return nil
}
