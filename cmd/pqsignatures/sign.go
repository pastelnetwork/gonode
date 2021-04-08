package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	pqtime "github.com/pastelnetwork/pqsignatures/internal/time"
	"github.com/pastelnetwork/pqsignatures/legroast"
)

func pastelIdWriteSignatureOnData(inputData string, skBase64 string, pkBase64 string) (string, error) {
	fmt.Printf("\nGenerating LegRoast signature now...")
	defer pqtime.Measure(time.Now())

	sk, err := base64.StdEncoding.DecodeString(skBase64)
	if err != nil {
		return "", fmt.Errorf("pastelIdWriteSignatureOnData: %w", err)
	}
	pk, err := base64.StdEncoding.DecodeString(pkBase64)
	if err != nil {
		return "", fmt.Errorf("pastelIdWriteSignatureOnData: %w", err)
	}
	pqtime.Sleep()
	pastelIdSignature := legroast.Sign(pk, sk, ([]byte)(inputData[:]))
	pastelIdSignatureBase64 := base64.StdEncoding.EncodeToString(pastelIdSignature)
	pqtime.Sleep()
	return pastelIdSignatureBase64, nil
}

func pastelIdVerifySignatureWithPublicKey(inputData string, pastelIdSignatureBase64 string, pkBase64 string) (int, error) {
	fmt.Printf("\nVerifying LegRoast signature now...")
	defer pqtime.Measure(time.Now())

	pastelIdSignature, err := base64.StdEncoding.DecodeString(pastelIdSignatureBase64)
	if err != nil {
		return 0, fmt.Errorf("pastelIdVerifySignatureWithPublicKey: %w", err)
	}
	pk, err := base64.StdEncoding.DecodeString(pkBase64)
	if err != nil {
		return 0, fmt.Errorf("pastelIdVerifySignatureWithPublicKey: %w", err)
	}
	pqtime.Sleep()
	verified := legroast.Verify(pk, ([]byte)(inputData[:]), pastelIdSignature)
	pqtime.Sleep()
	return verified, nil
}

func signAndVerify(inputData string, skBase64 string, pkBase64 string) (string, error) {
	pastelIdSignatureBase64, err := pastelIdWriteSignatureOnData(inputData, skBase64, pkBase64)
	if err != nil {
		return "", fmt.Errorf("signAndVerify: %w", err)
	}
	verified, err := pastelIdVerifySignatureWithPublicKey(inputData, pastelIdSignatureBase64, pkBase64)
	if err != nil {
		return "", fmt.Errorf("signAndVerify: %w", err)
	}
	if verified > 0 {
		fmt.Printf("\nSignature is valid!")
	} else {
		fmt.Printf("\nWarning! Signature was NOT valid!")
		return "", fmt.Errorf("signAndVerify: %w", errors.New("signature is not valid"))
	}
	return pastelIdSignatureBase64, nil
}
