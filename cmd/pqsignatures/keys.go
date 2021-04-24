package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pastelnetwork/go-commons/errors"

	"github.com/kevinburke/nacl/secretbox"
	legroast "github.com/pastelnetwork/go-legroast"
)

var WrongOTPFormat = errors.Errorf("one time password must contain 6 digits")

func pastelIdKeypairGeneration() (string, string) {
	fmt.Println("\nGenerating LegRoast keypair now...")
	pk, sk := legroast.Keygen()
	fmt.Printf("\npk length: %v ;sk: %v", len(pk), len(sk))
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
			return "", "", errors.New(errPK)
		}
		skExportDataEncrypted, errSK := ioutil.ReadFile(skPemFilePath)
		if errSK != nil {
			return "", "", errors.New(errSK)
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
			return "", "", errors.New(WrongOTPFormat)
		}
		otpCorrect = (otp_from_user_input == otp)
	}

	if otpCorrect {
		boxKey, err := naclBoxKeyFromFile(boxKeyFilePath)
		if err != nil {
			return "", "", err
		}
		var key [32]byte
		copy(key[:], boxKey)
		skExportFormat, err := secretbox.EasyOpen(([]byte)(skExportFormatEncrypted[:]), &key)
		if err != nil {
			return "", "", errors.New(err)
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
		return err
	}
	var key [32]byte
	copy(key[:], boxKey)
	encrypted := secretbox.EasySeal(([]byte)(skExportFormat[:]), &key)

	keyFilePath := "pastel_id_key_files"
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		if err = os.MkdirAll(keyFilePath, 0770); err != nil {
			return errors.New(err)
		}
	}
	err = os.WriteFile(keyFilePath+"/pastel_id_legroast_public_key.pem", []byte(pkExportFormat), 0644)
	if err != nil {
		return errors.New(err)
	}
	err = os.WriteFile(keyFilePath+"/pastel_id_legroast_private_key.pem", []byte(encrypted), 0644)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func pastelKeys() (string, string, error) {
	pkBase64, skBase64, err := importPastelPublicAndPrivateKeysFromPemFiles(BoxKeyFilePath)
	if err != nil {
		return "", "", err
	}
	if pkBase64 == "" {
		skBase64, pkBase64 = pastelIdKeypairGeneration()
		if err = writePastelPublicAndPrivateKeyToFile(pkBase64, skBase64, BoxKeyFilePath); err != nil {
			return "", "", err
		}
	}
	return pkBase64, skBase64, nil
}
