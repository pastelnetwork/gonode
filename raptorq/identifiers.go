package raptorq

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
)

type SymbolIdFiles []*SymbolIdFile

type SymbolIdFile struct {
	Id                string   `json:"id;omitempty"`
	BlockHash         string   `json:"block_hash;omitempty"`
	PastelId          string   `json:"pastel_id;omitempty"`
	SymbolIdentifiers []string `json:"symbol_identifiers;omitempty"`
	Signature         []byte   `json:"signature;omitempty"`
	FileIdentifer     string   `json:"-"`
}

func (files SymbolIdFiles) ToMap() (map[string][]byte, error) {
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

func (files SymbolIdFiles) FileIdentifers() ([]string, error) {
	var identifiers []string
	for _, file := range files {
		if file.FileIdentifer == "" {
			return nil, errors.Errorf("empty file identifier")
		}
		identifiers = append(identifiers, file.FileIdentifer)
	}
	return identifiers, nil
}
