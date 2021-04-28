package pqsignatures

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/kevinburke/nacl"
	"github.com/pastelnetwork/gonode/common/errors"
)

var SecretBoxKeyWrongSize = errors.Errorf("nacl secret box key doesn't match its expected size")

// SetupNaclKey setups nacl secret box with generated key store as keyFilePath.
func SetupNaclKey(keyFilePath string) (string, error) {
	boxKey := nacl.NewKey()
	boxKeyBase64 := base64.StdEncoding.EncodeToString(boxKey[:])
	err := os.WriteFile(keyFilePath, []byte(boxKeyBase64), 0644)
	if err != nil {
		return "", errors.New(err)
	}
	return boxKeyBase64, nil
}

func naclBoxKeyFromFile(boxKeyFilePath string) ([]byte, error) {
	boxKeyBase64, err := ioutil.ReadFile(boxKeyFilePath)
	if err != nil {
		return nil, errors.New(err)
	}
	boxKey, err := base64.StdEncoding.DecodeString(string(boxKeyBase64))
	if err != nil {
		return nil, errors.New(err)
	}
	if len(boxKey) != nacl.KeySize {
		return nil, SecretBoxKeyWrongSize
	}
	return boxKey, nil
}
