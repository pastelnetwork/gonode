package conn

import (
	"errors"
	"math/rand"
)

var ErrWrongFormat = errors.New("Unknown message format")
var ErrWrongSignature = errors.New("Wrong signature")
var ErrIncorrectPastelID = errors.New("Incorrect Pastel Id")

// Handshake message types.
const (
	typeClientHello        byte = 1
	typeServerHello        byte = 2
	typeClientHandshakeMsg byte = 3
	typeServerHandshakeMsg byte = 4
)

// EncryptionScheme type defines all supported encryption
type EncryptionScheme byte

// EncryptionScheme types
const (
	AES128 EncryptionScheme = iota
	AES192
	AES256
)

// SignScheme type defines all supported signature methods
type SignScheme byte

// SignScheme methods
const (
	ED448 SignScheme = iota
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
