package main

import (
	"context"
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
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/jsonrpc"
	pqtime "github.com/pastelnetwork/gonode/pqsignatures/internal/time"
	"github.com/pastelnetwork/gonode/pqsignatures/pkg/qr"
	"github.com/pastelnetwork/gonode/pqsignatures/pkg/steganography"
	"golang.org/x/crypto/sha3"
)

const (
	pastelIDSignatureFilesFolder = "pastel_id_signature_files"
	otpSecretFile                = "otp_secret.txt"
	userEmail                    = "user@user.com"
	otpQRCodeFilePath            = "Google_Authenticator_QR_Code.png"
	pastelKeysDirectoryPath      = "pastel_id_key_files"
)

var invalidSignature = errors.Errorf("signature is not valid")
var decodedPublicKeyNotMatch = errors.Errorf("decoded base64 public key doesn't match")
var decodedSignatureNotMatch = errors.Errorf("decoded base64 pastel id signature doesn't match")
var decodedFingerprintNotMatch = errors.Errorf("decoded base64 image fingerprint doesn't match")

var (
	defaultPath             = configurer.DefaultPath()
	defaultPastelConfigFile = filepath.Join(defaultPath, "pastel.conf")
)

const (
	// SignAlgorithmED448 is ED448 signature algorithm
	SignAlgorithmED448 = "ed448"
	// SignAlgorithmLegRoast is Efficient post-quantum signatures algorithm
	SignAlgorithmLegRoast = "legroast"
)

type client struct {
	jsonrpc.RPCClient
}

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

func demonstrateSignatureQRCodeSteganography(pkBase64 string, skBase64 string, pastelIDSignatureBase64 string, inputImagePath string) error {
	defer pqtime.Measure(time.Now())
	timestamp := time.Now().Format("Jan_02_2006_15_04_05")

	keypairImgs, err := generateKeypairQRs(pkBase64, skBase64)
	if err != nil {
		return err
	}

	signatureImgs, err := qr.Encode(pastelIDSignatureBase64, "sig", pastelIDSignatureFilesFolder, "Pastel Signature", "pastel_id_legroast_signature_qr_code", timestamp)
	if err != nil {
		return err
	}

	imgsToMap := append(keypairImgs, signatureImgs...)

	imgFingerprintBase64, err := loadImageFingerprint("fingerprint")
	if err != nil {
		return err
	}
	fingerprintImgs, err := qr.Encode(imgFingerprintBase64, "fin", pastelIDSignatureFilesFolder, "Fingerprint", "fingerprint_qr_code", timestamp)
	if err != nil {
		return err
	}

	imgsToMap = append(imgsToMap, fingerprintImgs...)

	signatureLayerImageOutputFilepath := filepath.Join(pastelIDSignatureFilesFolder, fmt.Sprintf("Complete_Signature_Image_Layer__%v.png", timestamp))
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
		return errors.New(decodedPublicKeyNotMatch)
	}
	if pastelIDSignatureBase64 != decodedSignatureBase64 {
		return errors.New(decodedSignatureNotMatch)
	}
	if imgFingerprintBase64 != decodedFingerprintBase64 {
		return errors.New(decodedFingerprintNotMatch)
	}

	fmt.Printf("\n\nBase64 public key and pastel id signature decoded from QR codes images are valid!\n")
	return nil
}

func sign(imagePath string, pastelClient pastel.Client, pastelID string, passphrase string) error {
	ctx := context.Background()
	defer pqtime.Measure(time.Now())

	fmt.Printf("\nApplying signature to file %v", imagePath)
	sha256HashOfImageToSign, err := getImageHashFromImageFilePath(imagePath)
	if err != nil {
		return err
	}
	fmt.Printf("\nSHA256 Hash of Image File: %v", sha256HashOfImageToSign)
	sha256HashOfImageToSignBytes := []byte(sha256HashOfImageToSign)

	fmt.Printf("\nGenerating Ed448 signature now...")
	//pastelIDSignatureBase64, err := pq.Sign(sha256HashOfImageToSign, skBase64, pkBase64)
	signature_ed448, err := pastelClient.Sign(ctx, sha256HashOfImageToSignBytes, pastelID, passphrase, "ed448")
	if err != nil {
		return err
	}
	fmt.Printf("\nGenerating LegRoast signature now...")
	signature_legroast, err := pastelClient.Sign(ctx, sha256HashOfImageToSignBytes, pastelID, passphrase, "legroast")
	if err != nil {
		return err
	}

	fmt.Printf("\nVerifying Ed448 signature now...")
	verified_ed448, err := pastelClient.Verify(ctx, sha256HashOfImageToSignBytes, string(signature_ed448), pastelID, "ed448")
	if err != nil {
		return err
	}

	fmt.Printf("\nVerifying LegRoast signature now...")
	verified_legroast, err := pastelClient.Verify(ctx, sha256HashOfImageToSignBytes, string(signature_legroast), pastelID, "legroast")
	if err != nil {
		return err
	}

	if verified_ed448 {
		fmt.Printf("\nEd448 Signature is valid!")
	} else {
		fmt.Printf("\nWarning! Ed448 Signature was NOT valid!")
		return errors.New(invalidSignature)
	}

	if verified_legroast {
		fmt.Printf("\nLegRoast Signature is valid!")
	} else {
		fmt.Printf("\nWarning! LegRoast Signature was NOT valid!")
		return errors.New(invalidSignature)
	}

	//err = demonstrateSignatureQRCodeSteganography(string(signature_ed448), string(signature_legroast), imagePath)
	err = demonstrateSignatureQRCodeSteganography(pastelID, string(signature_ed448), string(signature_legroast), imagePath)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var pastelConfig pastel.Config
	imagePathPtr := flag.String("image", "sample_image2.png", "an image file path")
	pastelConfigFile := flag.String("pastel-config-file", defaultPastelConfigFile, "Set `path` to the pastel config file")
	pastelID := flag.String("pastelID", "", "pastelID")
	passphrase := flag.String("passphrase", "", "passphrase")

	flag.Parse()
	if *pastelID == "" || *passphrase == "" {
		panic(errors.New("invalid pastelID or passphrase"))
	}
	if *pastelConfigFile != "" {
		if err := configurer.ParseFile(*pastelConfigFile, &pastelConfig); err != nil {
			panic(errors.Errorf("error parsing pastel config file: %v", err))
		}
	} else {
		panic(errors.New("empty -pastel-config-file"))
	}

	pastelClient := pastel.NewClient(&pastelConfig)

	sampleImageFilePath := *imagePathPtr
	if err := sign(sampleImageFilePath, pastelClient, *pastelID, *passphrase); err != nil {
		if err, isCommonError := err.(*errors.Error); isCommonError {
			fmt.Println(errors.ErrorStack(err))
		}
		panic(err)
	}
}
