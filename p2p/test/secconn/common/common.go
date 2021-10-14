package common

import (
	"context"
	"encoding/base64"

	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	// B64PriKey - preconfig private key
	B64PriKey = "MSTJkuJI7WV1H5UZ09FqCNCx8UNor18e89lSHxm3A7M1fuWTC7ED0e8Ss1t/KGzNfG3lz7dX8Rawc4aLFRchii2gRq049U/V3XXnV4zRRUZAb9hBx2XB+jF/YFQyT+gKfzfGjKkCb1EO+ee393yGPm/NSSmzvMv1brX5PfNeuCPtRZJAdvvxHCEw2rysJq9s"
	// B64PubKey - preconfig public key
	B64PubKey = "sHOGixUXIYotoEatOPVP1d1151eM0UVGQG/YQcdlwfoxf2BUMk/oCn83xoypAm9RDvnnt/d8hj4="
)

// GetKeys return a couple of (private/public) keys to sign/verify handshake authentication
func GetKeys() (PriKey [144]byte, PubKey [56]byte) {
	Priv, err := base64.StdEncoding.DecodeString(B64PriKey)
	if err != nil {
		panic("failed to decode Private key :" + err.Error())
	}
	copy(PriKey[:], Priv)

	Pub, err := base64.StdEncoding.DecodeString(B64PubKey)
	if err != nil {
		panic("failed to decode Pub key :" + err.Error())
	}
	copy(PubKey[:], Pub)

	return
}

// SecClient is implemented Sign/Verify functions
type SecClient struct {
	Curve ed448.Curve
	Pri   [144]byte
	Pub   [56]byte
}

// Sign a data
func (c *SecClient) Sign(_ context.Context, data []byte, _, _ string, _ string) ([]byte, error) {
	signature, ok := c.Curve.Sign(c.Pri, data)
	if !ok {
		return nil, errors.New("sign failed")
	}
	signatureStr := base64.StdEncoding.EncodeToString(signature[:])
	return []byte(signatureStr), nil
}

// Verify data + signature
func (c *SecClient) Verify(_ context.Context, data []byte, signature, _ string, _ string) (ok bool, err error) {
	signatureData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.Errorf("decode failed %w", err)
	}
	var copiedSignature [112]byte
	copy(copiedSignature[:], signatureData)

	ok = c.Curve.Verify(copiedSignature, data, c.Pub)
	return ok, nil
}
