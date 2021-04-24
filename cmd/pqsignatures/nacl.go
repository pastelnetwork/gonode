package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pastelnetwork/go-commons/errors"

	"github.com/kevinburke/nacl"
)

const (
	BoxKeyFilePath = "box_key.bin"
)

var SecretBoxKeyWrongSize = errors.Errorf("nacl secret box key doesn't match its expected size")

func generateAndStoreKeyForNacl() error {
	box_key := nacl.NewKey()
	box_key_base64 := base64.StdEncoding.EncodeToString(box_key[:])
	err := os.WriteFile(BoxKeyFilePath, []byte(box_key_base64), 0644)
	if err != nil {
		return errors.New(err)
	}

	fmt.Printf("\nThis is the key for encrypting the pastel ID private key (using NACL box) in Base64: %v", box_key_base64)
	fmt.Printf("\nThe key has been saved as a file in the working directory. You should also write this key down as a backup.")
	return nil
}

func naclBoxKeyFromFile(boxKeyFilePath string) ([]byte, error) {
	box_key_base64, err := ioutil.ReadFile(boxKeyFilePath)
	if err != nil {
		return nil, errors.New(err)
	}

	box_key, err := base64.StdEncoding.DecodeString(string(box_key_base64))
	if err != nil {
		return nil, errors.New(err)
	}

	if len(box_key) != nacl.KeySize {
		return nil, SecretBoxKeyWrongSize
	}

	return box_key, nil
}
