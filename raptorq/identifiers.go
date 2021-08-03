package raptorq

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
)

// SymbolIDFiles is a list of raptorq symbolid file
type SymbolIDFiles []*SymbolIDFile

// SymbolIDFile represents structure of a raptorq symbolid file
type SymbolIDFile struct {
	ID                string   `json:"id,omitempty"`
	BlockHash         string   `json:"block_hash,omitempty"`
	PastelID          string   `json:"pastel_id,omitempty"`
	SymbolIdentifiers []string `json:"symbol_identifiers,omitempty"`
	Signature         []byte   `json:"signature,omitempty"`
	FileIdentifer     string   `json:"-"`
}

// ToMap translate SymbolIDFile struct to generic map[string]
func (files SymbolIDFiles) ToMap() (map[string][]byte, error) {
	m := make(map[string][]byte, len(files))
	for _, file := range files {
		js, err := json.Marshal(file)
		if err != nil {
			return nil, err
		}
		m[file.FileIdentifer] = js
	}
	return m, nil
}

// FileIdentifers return name of files in FileIdentifers
func (files SymbolIDFiles) FileIdentifers() ([]string, error) {
	var identifiers []string
	for _, file := range files {
		if file.FileIdentifer == "" {
			return nil, errors.Errorf("empty file identifier")
		}
		identifiers = append(identifiers, file.FileIdentifer)
	}
	return identifiers, nil
}
