// Package pqsignatures helps to manage keys.
package pqsignatures

import (
	"context"
	"encoding/base64"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel/jsonrpc"
)

const (
	// SignAlgorithmED448 is ED448 signature algorithm
	SignAlgorithmED448 = "ed448"
	// SignAlgorithmLegRoast is Efficient post-quantum signatures algorithm
	SignAlgorithmLegRoast = "legroast"
)

type client struct {
	jsonrpc.RPCClient
}

func (client *client) Sign(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error) {
	var sign struct {
		Signature string `json:"signature"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	switch algorithm {
	case SignAlgorithmED448, SignAlgorithmLegRoast:
		if err = client.callFor(ctx, &sign, "pastelid", "sign", text, pastelID, passphrase, algorithm); err != nil {
			return nil, errors.Errorf("failed to sign data: %w", err)
		}
	default:
		return nil, errors.Errorf("unsupported algorithm %s", algorithm)
	}
	return []byte(sign.Signature), nil
}

func (client *client) Verify(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error) {
	var verify struct {
		Verification string `json:"verification"`
	}
	text := base64.StdEncoding.EncodeToString(data)

	switch algorithm {
	case SignAlgorithmED448, SignAlgorithmLegRoast:
		if err = client.callFor(ctx, &verify, "pastelid", "verify", text, signature, pastelID, algorithm); err != nil {
			return false, errors.Errorf("failed to verify data: %w", err)
		}
	default:
		return false, errors.Errorf("unsupported algorithm %s", algorithm)
	}

	return verify.Verification == "OK", nil
}

/*
// Sign signs data with provided pair of keys.
func Sign(data string, skBase64 string, pkBase64 string) (string, error) {
	defer pqtime.Measure(time.Now())

	sk, err := base64.StdEncoding.DecodeString(skBase64)
	if err != nil {
		return "", errors.New(err)
	}
	pk, err := base64.StdEncoding.DecodeString(pkBase64)
	if err != nil {
		return "", errors.New(err)
	}
	pqtime.Sleep()
	pastelIDSignature := legroast.Sign(pk, sk, ([]byte)(data[:]))
	pastelIDSignatureBase64 := base64.StdEncoding.EncodeToString(pastelIDSignature)
	pqtime.Sleep()
	return pastelIDSignatureBase64, nil
}

// Verify validates previously signed data.
func Verify(data string, signedData string, pkBase64 string) (int, error) {
	defer pqtime.Measure(time.Now())

	pastelIDSignature, err := base64.StdEncoding.DecodeString(signedData)
	if err != nil {
		return 0, errors.New(err)
	}
	pk, err := base64.StdEncoding.DecodeString(pkBase64)
	if err != nil {
		return 0, errors.New(err)
	}
	pqtime.Sleep()
	verified := legroast.Verify(pk, ([]byte)(data[:]), pastelIDSignature)
	pqtime.Sleep()
	return verified, nil
}
*/
