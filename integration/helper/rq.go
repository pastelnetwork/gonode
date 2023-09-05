package helper

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
)

type RawSymbolIDFile struct {
	ID                string   `json:"id"`
	BlockHash         string   `json:"block_hash"`
	PastelID          string   `json:"pastel_id"`
	SymbolIdentifiers []string `json:"symbol_identifiers"`
}

func GetP2PID(in []byte) (string, error) {
	hash, err := utils.Sha3256hash(in)
	if err != nil {
		return "", errors.Errorf("sha3256hash: %w", err)
	}

	return base58.Encode(hash), nil
}

func (file *RawSymbolIDFile) GetIDFile() ([]byte, error) {
	dataJSON, err := json.Marshal(file)
	if err != nil {
		return nil, errors.Errorf("failed to marshal dd-data: %w", err)
	}

	encoded := utils.B64Encode(dataJSON)

	var buffer bytes.Buffer
	buffer.Write(encoded)
	buffer.WriteByte(46)
	buffer.WriteString("test-signature")
	buffer.WriteByte(46)
	buffer.WriteString(strconv.Itoa(55))

	compressedData, err := utils.Compress(buffer.Bytes(), 4)
	if err != nil {
		return nil, errors.Errorf("compress identifiers file: %w", err)
	}

	return compressedData, nil
}
