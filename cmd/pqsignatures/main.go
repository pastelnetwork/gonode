package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"

	"fmt"
	"time"

	pqtime "github.com/PastelNetwork/pqsignatures/internal/time"
	"github.com/PastelNetwork/pqsignatures/legroast"
	"github.com/PastelNetwork/pqsignatures/qr"

	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/sha3"

	"github.com/auyer/steganography"
)

const (
	PastelIdSignatureFilesFolder = "pastel_id_signature_files"
)

func getImageHashFromImageFilePath(sampleImageFilePath string) (string, error) {
	f, err := os.Open(sampleImageFilePath)
	if err != nil {
		return "", fmt.Errorf("getImageHashFromImageFilePath: %w", err)
	}

	defer f.Close()
	hash := sha3.New256()
	if _, err := io.Copy(hash, f); err != nil {
		return "", fmt.Errorf("getImageHashFromImageFilePath: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func generateKeypairQRs(pk string, sk string) ([]qr.Image, error) {
	key_file_path := "pastel_id_key_files"
	pkPngs, err := qr.Encode(pk, key_file_path, "Pastel Public Key", "pastel_id_legroast_public_key_qr_code", "")
	if err != nil {
		return nil, err
	}
	_, err = qr.Encode(sk, key_file_path, "", "pastel_id_legroast_private_key_qr_code", "")
	if err != nil {
		return nil, err
	}
	return pkPngs, nil
}

func pastel_id_write_signature_on_data_func(input_data_or_string string, pastel_id_private_key_b16_encoded string, pastel_id_public_key_b16_encoded string) string {
	fmt.Printf("\nGenerating LegRoast signature now...")
	defer pqtime.Measure(time.Now())

	pastel_id_private_key, _ := base64.StdEncoding.DecodeString(pastel_id_private_key_b16_encoded)
	pastel_id_public_key, _ := base64.StdEncoding.DecodeString(pastel_id_public_key_b16_encoded)
	pqtime.Sleep()
	pastel_id_signature := legroast.Sign(pastel_id_public_key, pastel_id_private_key, ([]byte)(input_data_or_string[:]))
	pastel_id_signature_b16_encoded := base64.StdEncoding.EncodeToString(pastel_id_signature)
	pqtime.Sleep()
	return pastel_id_signature_b16_encoded
}

func pastel_id_verify_signature_with_public_key_func(input_data_or_string string, pastel_id_signature_b16_encoded string, pastel_id_public_key_b16_encoded string) int {
	fmt.Printf("\nVerifying LegRoast signature now...")
	defer pqtime.Measure(time.Now())

	pastel_id_signature, _ := base64.StdEncoding.DecodeString(pastel_id_signature_b16_encoded)
	pastel_id_public_key, _ := base64.StdEncoding.DecodeString(pastel_id_public_key_b16_encoded)
	pqtime.Sleep()
	verified := legroast.Verify(pastel_id_public_key, ([]byte)(input_data_or_string[:]), pastel_id_signature)
	pqtime.Sleep()
	if verified > 0 {
		fmt.Printf("\nSignature is valid!")
	} else {
		fmt.Printf("\nWarning! Signature was NOT valid!")
	}
	return verified
}

func hide_signature_image_in_sample_image_func(sample_image_file_path string, signature_layer_image_output_filepath string, signed_image_output_path string) {
	img, err := gg.LoadImage(sample_image_file_path)
	if err != nil {
		log.Fatal(err)
	}

	signature_layer_image_data, err := ioutil.ReadFile(signature_layer_image_output_filepath)
	if err != nil {
		panic(err)
	}

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, signature_layer_image_data)
	if err != nil {
		panic(err)
	}
	outFile, _ := os.Create(signed_image_output_path)
	w.WriteTo(outFile)
	outFile.Close()
}

func extract_signature_image_in_sample_image_func(signed_image_output_path string, extracted_signature_layer_image_output_filepath string) {
	img, err := gg.LoadImage(signed_image_output_path)
	if err != nil {
		panic(err)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)

	decodedData := steganography.Decode(sizeOfMessage, img)
	err = os.WriteFile(extracted_signature_layer_image_output_filepath, decodedData, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {

	imagePathPtr := flag.String("image", "sample_image2.png", "an image file path")
	flag.Parse()

	if _, err := os.Stat(OTPSecretFile); os.IsNotExist(err) {
		if err := setupGoogleAuthenticatorForPrivateKeyEncryption(); err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(BoxKeyFilePath); os.IsNotExist(err) {
		if err := generateAndStoreKeyForNacl(); err != nil {
			panic(err)
		}
	}

	sampleImageFilePath := *imagePathPtr

	fmt.Printf("\nApplying signature to file %v", sampleImageFilePath)
	sha256HashOfImageToSign, err := getImageHashFromImageFilePath(sampleImageFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nSHA256 Hash of Image File: %v", sha256HashOfImageToSign)

	pkBase64, skBase64, err := pastelKeys()
	if err != nil {
		panic(err)
	}

	keypairImgs, err := generateKeypairQRs(pkBase64, skBase64)
	if err != nil {
		panic(err)
	}

	pastel_id_signature_b16_encoded := pastel_id_write_signature_on_data_func(sha256HashOfImageToSign, skBase64, pkBase64)
	_ = pastel_id_verify_signature_with_public_key_func(sha256HashOfImageToSign, pastel_id_signature_b16_encoded, pkBase64)

	timestamp := time.Now().Format("Jan_02_2006_15_04_05")
	signatureImags, err := qr.Encode(pastel_id_signature_b16_encoded, PastelIdSignatureFilesFolder, "Pastel Signature", "pastel_id_legroast_signature_qr_code", timestamp)
	if err != nil {
		panic(err)
	}

	imgsToMap := append(keypairImgs, signatureImags...)

	signatureLayerImageOutputFilepath := filepath.Join(PastelIdSignatureFilesFolder, fmt.Sprintf("Complete_Signature_Image_Layer__%v.png", timestamp))
	sample_image, err := gg.LoadImage(sampleImageFilePath)
	if err != nil {
		panic(err)
	}
	err = qr.MapImages(imgsToMap, sample_image.Bounds().Size(), signatureLayerImageOutputFilepath)
	if err != nil {
		panic(err)
	}

	signed_image_output_path := "final_watermarked_image.png"
	hide_signature_image_in_sample_image_func(sampleImageFilePath, signatureLayerImageOutputFilepath, signed_image_output_path)

	extracted_signature_layer_image_output_filepath := "extracted_signature_image.png"
	extract_signature_image_in_sample_image_func(signed_image_output_path, extracted_signature_layer_image_output_filepath)
}
