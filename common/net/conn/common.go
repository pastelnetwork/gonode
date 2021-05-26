package conn

import (
	"errors"
)

// ErrWrongFormat - error when structure can't be parsed
var ErrWrongFormat = errors.New("unknown message format")

// ErrWrongSignature - error when signature is incorrect
var ErrWrongSignature = errors.New("wrong signature")

// ErrIncorrectPastelID - error when pastel ID is wrong
var ErrIncorrectPastelID = errors.New("incorrect Pastel Id")

var ErrUnsupportedEncryption = errors.New("unsupported encryption")

// Handshake message types.
const (
	typeClientHello = iota + 1
	typeServerHello
	typeClientHandshakeMsg
	typeServerHandshakeMsg
	typeEncryptedMsg
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

type message interface {
	marshall() ([]byte, error)
}

type signedPastelID struct {
	pastelID       []byte
	signedPastelID []byte
	pubKey         []byte
	ctx            []byte
}
