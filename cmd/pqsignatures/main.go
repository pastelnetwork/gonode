package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/pastelnetwork/pqsignatures/qr"
	"golang.org/x/crypto/sha3"
)

const (
	PastelIdSignatureFilesFolder = "pastel_id_signature_files"
)

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
	keyFilePath := "pastel_id_key_files"
	pkPngs, err := qr.Encode(pk, "pk", keyFilePath, "Pastel Public Key", "pastel_id_legroast_public_key_qr_code", "")
	if err != nil {
		return nil, errors.New(err)
	}
	_, err = qr.Encode(sk, "sk", keyFilePath, "", "pastel_id_legroast_private_key_qr_code", "")
	if err != nil {
		return nil, errors.New(err)
	}
	return pkPngs, nil
}

func sign(imagePath string) error {
	if _, err := os.Stat(OTPSecretFile); os.IsNotExist(err) {
		if err := setupGoogleAuthenticatorForPrivateKeyEncryption(); err != nil {
			return errors.New(err)
		}
	}

	if _, err := os.Stat(BoxKeyFilePath); os.IsNotExist(err) {
		if err := generateAndStoreKeyForNacl(); err != nil {
			return errors.New(err)
		}
	}

	fmt.Printf("\nApplying signature to file %v", imagePath)
	sha256HashOfImageToSign, err := getImageHashFromImageFilePath(imagePath)
	if err != nil {
		return errors.New(err)
	}
	fmt.Printf("\nSHA256 Hash of Image File: %v", sha256HashOfImageToSign)

	pkBase64, skBase64, err := pastelKeys()
	if err != nil {
		return errors.New(err)
	}

	pastelIdSignatureBase64, err := signAndVerify(sha256HashOfImageToSign, skBase64, pkBase64)
	if err != nil {
		return errors.New(err)
	}

	err = demonstrateSignatureQRCodeSteganography(pkBase64, skBase64, pastelIdSignatureBase64, imagePath)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func main() {
	imagePathPtr := flag.String("image", "sample_image2.png", "an image file path")
	flag.Parse()

	sampleImageFilePath := *imagePathPtr
	if err := sign(sampleImageFilePath); err != nil {
		fmt.Println(err.(*errors.Error).ErrorStack())
		panic(err)
	}
}
