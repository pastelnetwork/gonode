package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/DataDog/zstd"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/pastelnetwork/gonode/common/errors"
	pqtime "github.com/pastelnetwork/pqsignatures/internal/time"
	pq "github.com/pastelnetwork/pqsignatures/pkg/pqsignatures"
	"github.com/pastelnetwork/pqsignatures/pkg/qr"
	"github.com/pastelnetwork/pqsignatures/pkg/steganography"
	"golang.org/x/crypto/sha3"
)

const (
	PastelIdSignatureFilesFolder = "pastel_id_signature_files"
	otpSecretFile                = "otp_secret.txt"
	userEmail                    = "user@user.com"
	otpQRCodeFilePath            = "Google_Authenticator_QR_Code.png"
	naclBoxKeyFilePath           = "box_key.bin"
	pastelKeysDirectoryPath      = "pastel_id_key_files"
)

var InvalidSignature = errors.Errorf("signature is not valid")
var DecodedPublicKeyNotMatch = errors.Errorf("decoded base64 public key doesn't match")
var DecodedSignatureNotMatch = errors.Errorf("decoded base64 pastel id signature doesn't match")
var DecodedFingerprintNotMatch = errors.Errorf("decoded base64 image fingerprint doesn't match")

func getImageHashFromImageFilePath(sampleImageFilePath string) (string, error) {
	f, err := os.Open(sampleImageFilePath)
	if err != nil {
		return "", errors.New(err)
	}

	defer f.Close()
	hash := sha3.New256()
	if _, err := io.Copy(hash, f); err != nil {
		return "", errors.New(err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func generateKeypairQRs(pk string, sk string) ([]qr.Image, error) {
	pkPngs, err := qr.Encode(pk, "pk", pastelKeysDirectoryPath, "Pastel Public Key", "pastel_id_legroast_public_key_qr_code", "")
	if err != nil {
		return nil, err
	}
	_, err = qr.Encode(sk, "sk", pastelKeysDirectoryPath, "", "pastel_id_legroast_private_key_qr_code", "")
	if err != nil {
		return nil, err
	}
	return pkPngs, nil
}

func loadImageFingerprint(fingerprintFilePath string) (string, error) {
	fingerprintData, err := ioutil.ReadFile(fingerprintFilePath)
	if err != nil {
		return "", errors.New(err)
	}

	output, err := zstd.CompressLevel(nil, fingerprintData, 22)
	if err != nil {
		return "", errors.New(err)
	}
	return base64.StdEncoding.EncodeToString(output), nil
}

func demonstrateSignatureQRCodeSteganography(pkBase64 string, skBase64 string, pastelIdSignatureBase64 string, inputImagePath string) error {
	defer pqtime.Measure(time.Now())
	timestamp := time.Now().Format("Jan_02_2006_15_04_05")

	keypairImgs, err := generateKeypairQRs(pkBase64, skBase64)
	if err != nil {
		return err
	}

	signatureImgs, err := qr.Encode(pastelIdSignatureBase64, "sig", PastelIdSignatureFilesFolder, "Pastel Signature", "pastel_id_legroast_signature_qr_code", timestamp)
	if err != nil {
		return err
	}

	imgsToMap := append(keypairImgs, signatureImgs...)

	imgFingerprintBase64, err := loadImageFingerprint("fingerprint")
	if err != nil {
		return err
	}
	fingerprintImgs, err := qr.Encode(imgFingerprintBase64, "fin", PastelIdSignatureFilesFolder, "Fingerprint", "fingerprint_qr_code", timestamp)
	if err != nil {
		return err
	}

	imgsToMap = append(imgsToMap, fingerprintImgs...)

	signatureLayerImageOutputFilepath := filepath.Join(PastelIdSignatureFilesFolder, fmt.Sprintf("Complete_Signature_Image_Layer__%v.png", timestamp))
	inputImage, err := gg.LoadImage(inputImagePath)
	if err != nil {
		return errors.New(err)
	}

	inputImgSize := inputImage.Bounds().Size()
	err = qr.ImagesFitOutputSize(imgsToMap, inputImgSize)
	if err != nil {
		if err == qr.OutputSizeTooSmall {
			inputImage = resize.Resize(1700, 0, inputImage, resize.Lanczos3)
			f, err := os.Create(inputImagePath)
			if err != nil {
				return errors.New(err)
			}
			defer f.Close()
			png.Encode(f, inputImage)
		} else {
			return err
		}
	}

	err = qr.MapImages(imgsToMap, inputImage.Bounds().Size(), signatureLayerImageOutputFilepath)
	if err != nil {
		return err
	}

	signedImageOutputPath := "final_watermarked_image.png"
	err = steganography.Encode(inputImagePath, signatureLayerImageOutputFilepath, signedImageOutputPath)
	if err != nil {
		return err
	}

	extractedSignatureLayerImageOutputFilepath := "extracted_signature_image.png"
	err = steganography.Decode(signedImageOutputPath, extractedSignatureLayerImageOutputFilepath)
	if err != nil {
		return err
	}

	decodedMessages, err := qr.Decode(extractedSignatureLayerImageOutputFilepath)
	if err != nil {
		return err
	}

	var decodedPKBase64 string
	var decodedSignatureBase64 string
	var decodedFingerprintBase64 string
	for _, message := range decodedMessages {
		fmt.Printf("\nDecoded message with alias:%v and content:%v", message.Alias, message.Content)
		if message.Alias == "pk" {
			decodedPKBase64 = message.Content
		} else if message.Alias == "sig" {
			decodedSignatureBase64 = message.Content
		} else if message.Alias == "fin" {
			decodedFingerprintBase64 = message.Content
		}
	}
	if pkBase64 != decodedPKBase64 {
		return errors.New(DecodedPublicKeyNotMatch)
	}
	if pastelIdSignatureBase64 != decodedSignatureBase64 {
		return errors.New(DecodedSignatureNotMatch)
	}
	if imgFingerprintBase64 != decodedFingerprintBase64 {
		return errors.New(DecodedFingerprintNotMatch)
	}

	fmt.Printf("\n\nBase64 public key and pastel id signature decoded from QR codes images are valid!\n")
	return nil
}

func sign(imagePath string) error {
	defer pqtime.Measure(time.Now())

	if _, err := os.Stat(otpSecretFile); os.IsNotExist(err) {
		if err := pq.SetupOTPAuthenticator(userEmail, otpSecretFile, otpQRCodeFilePath); err != nil {
			return err
		}
	}

	if _, err := os.Stat(naclBoxKeyFilePath); os.IsNotExist(err) {
		if key, err := pq.SetupNaclKey(naclBoxKeyFilePath); err != nil {
			return err
		} else {
			fmt.Printf("\nThis is the key for encrypting the pastel ID private key (using NACL box) in Base64: %v", key)
			fmt.Printf("\nThe key has been saved as a file in the working directory. You should also write this key down as a backup.")
		}
	}

	fmt.Printf("\nApplying signature to file %v", imagePath)
	sha256HashOfImageToSign, err := getImageHashFromImageFilePath(imagePath)
	if err != nil {
		return err
	}
	fmt.Printf("\nSHA256 Hash of Image File: %v", sha256HashOfImageToSign)

	pkBase64, skBase64, err := pq.ImportPastelKeys(pastelKeysDirectoryPath, naclBoxKeyFilePath, otpSecretFile)
	if err == pq.KeyNotFound {
		pkBase64, skBase64, err = pq.GeneratePastelKeys(pastelKeysDirectoryPath, naclBoxKeyFilePath)
	}
	if err != nil {
		return err
	}

	fmt.Printf("\nGenerating LegRoast signature now...")
	pastelIdSignatureBase64, err := pq.Sign(sha256HashOfImageToSign, skBase64, pkBase64)
	if err != nil {
		return err
	}
	fmt.Printf("\nVerifying LegRoast signature now...")
	verified, err := pq.Verify(sha256HashOfImageToSign, pastelIdSignatureBase64, pkBase64)
	if err != nil {
		return err
	}
	if verified > 0 {
		fmt.Printf("\nSignature is valid!")
	} else {
		fmt.Printf("\nWarning! Signature was NOT valid!")
		return errors.New(InvalidSignature)
	}

	err = demonstrateSignatureQRCodeSteganography(pkBase64, skBase64, pastelIdSignatureBase64, imagePath)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	imagePathPtr := flag.String("image", "sample_image2.png", "an image file path")
	flag.Parse()

	sampleImageFilePath := *imagePathPtr
	if err := sign(sampleImageFilePath); err != nil {
		if err, isCommonError := err.(*errors.Error); isCommonError {
			fmt.Println(err.ErrorStack())
		}
		panic(err)
	}
}
