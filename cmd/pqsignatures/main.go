package main

import (
	"flag"
	"io"
	"log"
	"os"

	"fmt"

	"github.com/PastelNetwork/pqsignatures/qr"

	"encoding/hex"

	"golang.org/x/crypto/sha3"
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

	pastelIdSignatureBase64, err := signAndVerify(sha256HashOfImageToSign, skBase64, pkBase64)
	if err != nil {
		panic(err)
	}

	err = demonstrateSignatureQRCodeSteganography(pkBase64, skBase64, pastelIdSignatureBase64, sampleImageFilePath)
	if err != nil {
		panic(err)
	}
}
