package common

import (
	"context"
	"encoding/base64"

	"github.com/otrv4/ed448"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
)

const (
	// pre-generate keys based on code :
	// Create a preconfig keys, used in common package

	// curve := ed448.NewCurve()
	// pri, pub, _ := curve.GenerateKeys()
	// fmt.Println(base64.StdEncoding.EncodeToString(pri[:]))
	// fmt.Println(base64.StdEncoding.EncodeToString(pub[:]))

	// ServerB64PriKey - preconfig private key
	ServerB64PriKey = "MSTJkuJI7WV1H5UZ09FqCNCx8UNor18e89lSHxm3A7M1fuWTC7ED0e8Ss1t/KGzNfG3lz7dX8Rawc4aLFRchii2gRq049U/V3XXnV4zRRUZAb9hBx2XB+jF/YFQyT+gKfzfGjKkCb1EO+ee393yGPm/NSSmzvMv1brX5PfNeuCPtRZJAdvvxHCEw2rysJq9s"
	// ServerB64PubKey - preconfig public key
	ServerB64PubKey = "sHOGixUXIYotoEatOPVP1d1151eM0UVGQG/YQcdlwfoxf2BUMk/oCn83xoypAm9RDvnnt/d8hj4="

	// ClientB64PriKey - preconfig private key
	ClientB64PriKey = "ACb1SPfXNjBqYLk4JIkCTFlYQjDjYyBGVbMOxcOeKPjhWob22+nGIaNJM5icElP/VXrQf28A3hL0TLznDNDZfonyNOFe9FGTAO1dXU3SgszjA56XC1Jvp9FNCFtfEZ0ItSsAXVVg17lAh8lxOUEjXlqbuQxwEB6l37yzw9VGOCNhxGoimI551WSJ/w1n32Ov"
	// ClientB64PubKey - preconfig public key
	ClientB64PubKey = "9Ey85wzQ2X6J8jThXvRRkwDtXV1N0oLM4wOelwtSb6fRTQhbXxGdCLUrAF1VYNe5QIfJcTlBI14="
)

// SecClient is implemented Sign/Verify functions
type SecClient struct {
	// Client is mock client
	*pastelMock.Client
	// Curve is a curve
	Curve ed448.Curve
	// ServerPri is Server private key
	ServerPri [144]byte
	// ServerPub is Server pub key
	ServerPub [56]byte
	// ClientPri is Client private key
	ClientPri [144]byte
	// ClientPub is Client pub key
	ClientPub [56]byte
}

// NewSecClient return a SecClient
func NewSecClient() *SecClient {
	sec := &SecClient{
		Client: nil,
		Curve:  ed448.NewCurve(),
	}

	Priv, err := base64.StdEncoding.DecodeString(ServerB64PriKey)
	if err != nil {
		panic("failed to decode Private key :" + err.Error())
	}
	copy(sec.ServerPri[:], Priv)

	Pub, err := base64.StdEncoding.DecodeString(ServerB64PubKey)
	if err != nil {
		panic("failed to decode Pub key :" + err.Error())
	}
	copy(sec.ServerPub[:], Pub)

	Priv, err = base64.StdEncoding.DecodeString(ClientB64PriKey)
	if err != nil {
		panic("failed to decode Private key :" + err.Error())
	}
	copy(sec.ClientPri[:], Priv)

	Pub, err = base64.StdEncoding.DecodeString(ClientB64PubKey)
	if err != nil {
		panic("failed to decode Pub key :" + err.Error())
	}
	copy(sec.ClientPub[:], Pub)

	return sec
}

// // Sign signs data by the given pastelID and passphrase, if successful returns signature.
// // Command `pastelid sign "text" "PastelID" "passphrase"` "algorithm"
// // Algorithm is ed448[default] or legroast
// Sign(ctx context.Context, data []byte, pastelID, passphrase string, algorithm string) (signature []byte, err error)

// // Verify verifies signed data by the given its signature and pastelID, if successful returns true.
// // Command `pastelid verify "text" "signature" "PastelID"`.
// Verify(ctx context.Context, data []byte, signature, pastelID string, algorithm string) (ok bool, err error)

// Sign a data
func (c *SecClient) Sign(_ context.Context, data []byte, pastelID string, _ string, _ string) ([]byte, error) {
	pri := c.ServerPri
	if pastelID == ClientB64PubKey {
		pri = c.ClientPri
	}
	signature, ok := c.Curve.Sign(pri, data)
	if !ok {
		return nil, errors.New("sign failed")
	}
	signatureStr := base64.StdEncoding.EncodeToString(signature[:])
	return []byte(signatureStr), nil
}

// Verify data + signature
func (c *SecClient) Verify(_ context.Context, data []byte, signature, pastelID string, _ string) (ok bool, err error) {
	pub := c.ServerPub
	if pastelID == ClientB64PubKey {
		pub = c.ClientPub
	}

	signatureData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.Errorf("decode failed %w", err)
	}
	var copiedSignature [112]byte
	copy(copiedSignature[:], signatureData)

	ok = c.Curve.Verify(copiedSignature, data, pub)
	return ok, nil
}

// MasterNodesExtra - return master node extra requests
func (c *SecClient) MasterNodesExtra(_ context.Context) (pastel.MasterNodes, error) {
	// // For test error
	// return pastel.MasterNodes{}, nil
	return pastel.MasterNodes{
		pastel.MasterNode{
			ExtAddress: "127.0.0.1:1000", // port is not important here
			ExtKey:     ServerB64PubKey,
		},
		pastel.MasterNode{
			ExtAddress: "127.0.0.1:anyport",
			ExtKey:     ClientB64PubKey, // port is not important here
		},
	}, nil
}
