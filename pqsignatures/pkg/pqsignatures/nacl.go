package pqsignatures

import (
	"io/ioutil"
	"os"

	"github.com/darkwyrm/b85"
	"github.com/kevinburke/nacl"
	"github.com/pastelnetwork/gonode/common/errors"
)

// SecretBoxKeyWrongSize describes errors message if secret key from nacl secret box doesn't match its expected size
var SecretBoxKeyWrongSize = errors.Errorf("nacl secret box key doesn't match its expected size")

// SetupNaclKey setups nacl secret box with generated key store as keyFilePath.
func SetupNaclKey(keyFilePath string) (string, error) {
	boxKey := nacl.NewKey()
	boxKeyBase85 := b85.Encode(boxKey[:])
	err := os.WriteFile(keyFilePath, []byte(boxKeyBase85), 0644)
	if err != nil {
		return "", errors.New(err)
	}
	return boxKeyBase85, nil
}

func naclBoxKeyFromFile(boxKeyFilePath string) ([]byte, error) {
	boxKeyBase85, err := ioutil.ReadFile(boxKeyFilePath)
	if err != nil {
		return nil, errors.New(err)
	}
	boxKey, err := b85.Decode(string(boxKeyBase85))
	if err != nil {
		return nil, errors.New(err)
	}
	if len(boxKey) != nacl.KeySize {
		return nil, SecretBoxKeyWrongSize
	}
	return boxKey, nil
}
