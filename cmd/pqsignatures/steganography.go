package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
	"github.com/pastelnetwork/go-commons/errors"
	pqtime "github.com/pastelnetwork/pqsignatures/internal/time"
	"github.com/pastelnetwork/pqsignatures/qr"
	"github.com/pastelnetwork/steganography"
)

func hideSignatureImageInInputImage(sample_image_file_path string, signature_layer_image_output_filepath string, signed_image_output_path string) error {
	img, err := gg.LoadImage(sample_image_file_path)
	if err != nil {
		return errors.New(err)
	}

	signature_layer_image_data, err := ioutil.ReadFile(signature_layer_image_output_filepath)
	if err != nil {
		return errors.New(err)
	}

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, signature_layer_image_data)
	if err != nil {
		return errors.New(err)
	}

	outFile, err := os.Create(signed_image_output_path)
	if err != nil {
		return errors.New(err)
	}

	w.WriteTo(outFile)
	outFile.Close()
	return nil
}

func extractSignatureImageInSampleImage(signed_image_output_path string, extracted_signature_layer_image_output_filepath string) error {
	img, err := gg.LoadImage(signed_image_output_path)
	if err != nil {
		return errors.New(err)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)

	decodedData := steganography.Decode(sizeOfMessage, img)
	err = os.WriteFile(extracted_signature_layer_image_output_filepath, decodedData, 0644)
	if err != nil {
		errors.New(err)
	}
	return nil
}

func demonstrateSignatureQRCodeSteganography(pkBase64 string, skBase64 string, pastelIdSignatureBase64 string, inputImagePath string) error {
	defer pqtime.Measure(time.Now())
	timestamp := time.Now().Format("Jan_02_2006_15_04_05")

	keypairImgs, err := generateKeypairQRs(pkBase64, skBase64)
	if err != nil {
		return errors.New(err)
	}

	signatureImags, err := qr.Encode(pastelIdSignatureBase64, "sig", PastelIdSignatureFilesFolder, "Pastel Signature", "pastel_id_legroast_signature_qr_code", timestamp)
	if err != nil {
		return errors.New(err)
	}

	imgsToMap := append(keypairImgs, signatureImags...)

	signatureLayerImageOutputFilepath := filepath.Join(PastelIdSignatureFilesFolder, fmt.Sprintf("Complete_Signature_Image_Layer__%v.png", timestamp))
	inputImage, err := gg.LoadImage(inputImagePath)
	if err != nil {
		return errors.New(err)
	}
	err = qr.MapImages(imgsToMap, inputImage.Bounds().Size(), signatureLayerImageOutputFilepath)
	if err != nil {
		return errors.New(err)
	}

	signedImageOutputPath := "final_watermarked_image.png"
	err = hideSignatureImageInInputImage(inputImagePath, signatureLayerImageOutputFilepath, signedImageOutputPath)
	if err != nil {
		return errors.New(err)
	}

	extractedSignatureLayerImageOutputFilepath := "extracted_signature_image.png"
	err = extractSignatureImageInSampleImage(signedImageOutputPath, extractedSignatureLayerImageOutputFilepath)
	if err != nil {
		return errors.New(err)
	}

	decodedMessages, err := qr.Decode(extractedSignatureLayerImageOutputFilepath)
	if err != nil {
		return errors.New(err)
	}

	for _, message := range decodedMessages {
		fmt.Printf("\nDecoded message with alias:%v and content:%v", message.Alias, message.Content)
	}

	return nil
}
