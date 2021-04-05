package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kevinburke/nacl"
)

const (
	BoxKeyFilePath = "box_key.bin"
)

func generateAndStoreKeyForNacl() error {
	box_key := nacl.NewKey()
	box_key_base64 := base64.StdEncoding.EncodeToString(box_key[:])
	err := os.WriteFile(BoxKeyFilePath, []byte(box_key_base64), 0644)
	if err != nil {
		return fmt.Errorf("generateAndStoreKeyForNacl: %w", err)
	}

	fmt.Printf("\nThis is the key for encrypting the pastel ID private key (using NACL box) in Base64: %v", box_key)
	fmt.Printf("\nThe key has been saved as a file in the working directory. You should also write this key down as a backup.")
	return nil
}

func naclBoxKeyFromFile(boxKeyFilePath string) ([]byte, error) {
	box_key_base64, err := ioutil.ReadFile(boxKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("naclBoxKeyFromFile: %w", err)
	}

	box_key, err := base64.StdEncoding.DecodeString(string(box_key_base64))
	if err != nil {
		return nil, fmt.Errorf("naclBoxKeyFromFile: %w", err)
	}

	if len(box_key) != nacl.KeySize {
		return nil, errors.New("naclBoxKeyFromFile: nacl secret box key doesn't match its expected size")
	}

	return box_key, nil
}
