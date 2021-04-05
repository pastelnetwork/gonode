package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PastelNetwork/pqsignatures/legroast"
	"github.com/kevinburke/nacl/secretbox"
)

func pastelIdKeypairGeneration() (string, string) {
	fmt.Println("\nGenerating LegRoast keypair now...")
	pk, sk := legroast.Keygen()
	skBase64 := base64.StdEncoding.EncodeToString(sk)
	pkBase64 := base64.StdEncoding.EncodeToString(pk)
	return skBase64, pkBase64
}

func importPastelPublicAndPrivateKeysFromPemFiles(boxKeyFilePath string) (string, string, error) {
	keyFilePath := "pastel_id_key_files"
	if info, err := os.Stat(keyFilePath); os.IsNotExist(err) || !info.IsDir() {
		fmt.Printf("\nCan't find key storage directory, trying to use current working directory instead!")
		keyFilePath = "."
	}
	pkPemFilePath := keyFilePath + "/pastel_id_legroast_public_key.pem"
	skPemFilePath := keyFilePath + "/pastel_id_legroast_private_key.pem"

	pkBase64 := ""
	skBase64 := ""
	pkExportFormat := ""
	skExportFormatEncrypted := ""
	otpCorrect := false

	infoPK, errPK := os.Stat(pkPemFilePath)
	infoSK, errSK := os.Stat(skPemFilePath)
	if !os.IsNotExist(errPK) && !infoPK.IsDir() && !os.IsNotExist(errSK) && !infoSK.IsDir() {
		pkExportData, errPK := ioutil.ReadFile(pkPemFilePath)
		if errPK != nil {
			return "", "", fmt.Errorf("importPastelPublicAndPrivateKeysFromPemFiles: %w", errPK)
		}
		skExportDataEncrypted, errSK := ioutil.ReadFile(skPemFilePath)
		if errSK != nil {
			return "", "", fmt.Errorf("importPastelPublicAndPrivateKeysFromPemFiles: %w", errSK)
		}

		pkExportFormat = string(pkExportData)
		skExportFormatEncrypted = string(skExportDataEncrypted)

		otp := generateCurrentOtpString()
		if otp == "" {
			otp = generateCurrentOtpStringFromUserInput()
		}
		fmt.Println("\n\nPlease Enter your pastel Google Authenticator Code:")
		fmt.Println(otp)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		otp_from_user_input := scanner.Text()

		if len(otp_from_user_input) != 6 {
			return "", "", fmt.Errorf("importPastelPublicAndPrivateKeysFromPemFiles: %w", errors.New("One time password must contain 6 digits"))
		}
		otpCorrect = (otp_from_user_input == otp)
	}

	if otpCorrect {
		boxKey, err := naclBoxKeyFromFile(boxKeyFilePath)
		if err != nil {
			return "", "", fmt.Errorf("importPastelPublicAndPrivateKeysFromPemFiles: %w", err)
		}
		var key [32]byte
		copy(key[:], boxKey)
		skExportFormat, err := secretbox.EasyOpen(([]byte)(skExportFormatEncrypted[:]), &key)
		if err != nil {
			return "", "", fmt.Errorf("importPastelPublicAndPrivateKeysFromPemFiles: %w", err)
		}
		sk := strings.ReplaceAll(string(skExportFormat), "-----BEGIN LEGROAST PRIVATE KEY-----\n", "")
		skBase64 = strings.ReplaceAll(string(sk), "\n-----END LEGROAST PRIVATE KEY-----", "")

		pkExportFormat = strings.ReplaceAll(pkExportFormat, "-----BEGIN LEGROAST PUBLIC KEY-----\n", "")
		pkBase64 = strings.ReplaceAll(pkExportFormat, "\n-----END LEGROAST PUBLIC KEY-----", "")
	}
	return pkBase64, skBase64, nil
}

func writePastelPublicAndPrivateKeyToFile(pkBase64 string, skBase64 string, boxKeyFilePath string) error {
	pkExportFormat := "-----BEGIN LEGROAST PUBLIC KEY-----\n" + pkBase64 + "\n-----END LEGROAST PUBLIC KEY-----"
	skExportFormat := "-----BEGIN LEGROAST PRIVATE KEY-----\n" + skBase64 + "\n-----END LEGROAST PRIVATE KEY-----"
	boxKey, err := naclBoxKeyFromFile(boxKeyFilePath)
	if err != nil {
		return fmt.Errorf("writePastelPublicAndPrivateKeyToFile: %w", err)
	}
	var key [32]byte
	copy(key[:], boxKey)
	encrypted := secretbox.EasySeal(([]byte)(skExportFormat[:]), &key)

	keyFilePath := "pastel_id_key_files"
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		if err = os.MkdirAll(keyFilePath, 0770); err != nil {
			return fmt.Errorf("writePastelPublicAndPrivateKeyToFile: %w", err)
		}
	}
	err = os.WriteFile(keyFilePath+"/pastel_id_legroast_public_key.pem", []byte(pkExportFormat), 0644)
	if err != nil {
		return fmt.Errorf("writePastelPublicAndPrivateKeyToFile: %w", err)
	}
	err = os.WriteFile(keyFilePath+"/pastel_id_legroast_private_key.pem", []byte(encrypted), 0644)
	if err != nil {
		return fmt.Errorf("writePastelPublicAndPrivateKeyToFile: %w", err)
	}
	return nil
}

func pastelKeys() (string, string, error) {
	pkBase64, skBase64, err := importPastelPublicAndPrivateKeysFromPemFiles(BoxKeyFilePath)
	if err != nil {
		return "", "", fmt.Errorf("pastelKeys: %w", err)
	}
	if pkBase64 == "" {
		skBase64, pkBase64 = pastelIdKeypairGeneration()
		if err = writePastelPublicAndPrivateKeyToFile(pkBase64, skBase64, BoxKeyFilePath); err != nil {
			return "", "", fmt.Errorf("pastelKeys: %w", err)
		}
	}
	return pkBase64, skBase64, nil
}
