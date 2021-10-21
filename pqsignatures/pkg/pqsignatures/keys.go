// Package pqsignatures helps to manage keys.
package pqsignatures

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/darkwyrm/b85"
	"github.com/kevinburke/nacl/secretbox"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/legroast"
	pqtime "github.com/pastelnetwork/gonode/pqsignatures/internal/time"
)

const (
	publicKeyFileName  = "pastel_id_legroast_public_key.pem"
	privateKeyFileName = "pastel_id_legroast_private_key.pem"
)

// WrongOTPFormat describes errors message for wrong format of entered OTP password
var WrongOTPFormat = errors.Errorf("one time password must contain 6 digits")

// IncorrectEnteredOTP describes errors message if entered OTP doesn't match expected one
var IncorrectEnteredOTP = errors.Errorf("entered OTP is incorrect")

// KeyNotFound describes errors message if public or private keys are not found on their import attempt
var KeyNotFound = errors.Errorf("public or private key is not found")

// KeyReadError describes errors message if public or private keys files are not readable
var KeyReadError = errors.Errorf("public or private key file read error")

// ImportPastelKeys imports previously generated public and private keys.
func ImportPastelKeys(importDirectoryPath, naclBoxKeyFilePath, otpSecretFilePath string) (string, string, error) {
	pkPemFilePath := filepath.Join(importDirectoryPath, publicKeyFileName)
	skPemFilePath := filepath.Join(importDirectoryPath, privateKeyFileName)

	infoPK, errPK := os.Stat(pkPemFilePath)
	infoSK, errSK := os.Stat(skPemFilePath)
	if os.IsNotExist(errPK) || infoPK.IsDir() || os.IsNotExist(errSK) || infoSK.IsDir() {
		return "", "", KeyNotFound
	}

	pkExportData, errPK := ioutil.ReadFile(pkPemFilePath)
	if errPK != nil {
		return "", "", KeyReadError
	}
	skExportDataEncrypted, errSK := ioutil.ReadFile(skPemFilePath)
	if errSK != nil {
		return "", "", KeyReadError
	}

	pkExportFormat := string(pkExportData)
	skExportFormatEncrypted := string(skExportDataEncrypted)

	otp := generateCurrentOtpString(otpSecretFilePath)
	if otp == "" {
		otp = generateCurrentOtpStringFromUserInput()
	}
	fmt.Println("\n\nPlease Enter your pastel Google Authenticator Code:")
	fmt.Println(otp)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	otpFromUserInput := scanner.Text()

	if len(otpFromUserInput) != 6 {
		return "", "", errors.New(WrongOTPFormat)
	}
	if otpFromUserInput != otp {
		return "", "", IncorrectEnteredOTP
	}

	boxKey, err := naclBoxKeyFromFile(naclBoxKeyFilePath)
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
	skBase85 := strings.ReplaceAll(string(sk), "\n-----END LEGROAST PRIVATE KEY-----", "")

	pkExportFormat = strings.ReplaceAll(pkExportFormat, "-----BEGIN LEGROAST PUBLIC KEY-----\n", "")
	pkBase85 := strings.ReplaceAll(pkExportFormat, "\n-----END LEGROAST PUBLIC KEY-----", "")

	return pkBase85, skBase85, nil
}

// GeneratePastelKeys generates public and private keys, encodes private key with nacl secret box.
func GeneratePastelKeys(targetDirectoryPath, naclBoxKeyFilePath string) (string, string, error) {
	pk, sk := legroast.Keygen()
	skBase85 := b85.Encode(sk)
	pkBase85 := b85.Encode(pk)
	pkExportFormat := "-----BEGIN LEGROAST PUBLIC KEY-----\n" + pkBase85 + "\n-----END LEGROAST PUBLIC KEY-----"
	skExportFormat := "-----BEGIN LEGROAST PRIVATE KEY-----\n" + skBase85 + "\n-----END LEGROAST PRIVATE KEY-----"
	boxKey, err := naclBoxKeyFromFile(naclBoxKeyFilePath)
	if err != nil {
		return "", "", err
	}
	var key [32]byte
	copy(key[:], boxKey)
	encrypted := secretbox.EasySeal(([]byte)(skExportFormat[:]), &key)

	if _, err := os.Stat(targetDirectoryPath); os.IsNotExist(err) {
		if err = os.MkdirAll(targetDirectoryPath, 0770); err != nil {
			return "", "", errors.New(err)
		}
	}
	err = os.WriteFile(filepath.Join(targetDirectoryPath, publicKeyFileName), []byte(pkExportFormat), 0644)
	if err != nil {
		return "", "", errors.New(err)
	}
	err = os.WriteFile(filepath.Join(targetDirectoryPath, privateKeyFileName), []byte(encrypted), 0644)
	if err != nil {
		return "", "", errors.New(err)
	}

	return pkBase85, skBase85, nil
}

// Sign signs data with provided pair of keys.
func Sign(data string, skBase85 string, pkBase85 string) (string, error) {
	defer pqtime.Measure(time.Now())

	sk, err := b85.Decode(skBase85)
	if err != nil {
		return "", errors.New(err)
	}
	pk, err := b85.Decode(pkBase85)
	if err != nil {
		return "", errors.New(err)
	}
	pqtime.Sleep()
	pastelIDSignature := legroast.Sign(pk, sk, ([]byte)(data[:]))
	pastelIDSignatureBase85 := b85.Encode(pastelIDSignature)
	pqtime.Sleep()
	return pastelIDSignatureBase85, nil
}

// Verify validates previously signed data.
func Verify(data string, signedData string, pkBase85 string) (int, error) {
	defer pqtime.Measure(time.Now())

	pastelIDSignature, err := b85.Decode(signedData)
	if err != nil {
		return 0, errors.New(err)
	}
	pk, err := b85.Decode(pkBase85)
	if err != nil {
		return 0, errors.New(err)
	}
	pqtime.Sleep()
	verified := legroast.Verify(pk, ([]byte)(data[:]), pastelIDSignature)
	pqtime.Sleep()
	return verified, nil
}
