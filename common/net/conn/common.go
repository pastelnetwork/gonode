package conn

import (
	"errors"
	"math/rand"
)

var WrongFormatErr = errors.New("Unknown message format")
var WrongSignatureErr = errors.New("Wrong signature")
var IncorrectPastelIdErr = errors.New("Incorrect Pastel Id")

// Handshake message types.
const (
	typeClientHello        byte = 1
	typeServerHello        byte = 2
	typeClientHandshakeMsg byte = 3
	typeServerHandshakeMsg byte = 4
)

type EncryptionScheme byte

const (
	AES128 EncryptionScheme = iota
	AES192
	AES256
)

type SignScheme byte

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
