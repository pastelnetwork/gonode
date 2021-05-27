package ed448

import (
	"bytes"
	"encoding/binary"
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

func readByteArray(reader *bytes.Reader) (*[]byte, error) {
	var arrayLen int
	if err := binary.Read(reader, binary.LittleEndian, &arrayLen); err != nil {
		return nil, err
	}

	array := make([]byte, arrayLen)
	if _, err := reader.Read(array); err != nil {
		return nil, err
	}

	return &array, nil
}

func readString(reader *bytes.Reader) (*string, error) {
	var strLen int
	if err := binary.Read(reader, binary.LittleEndian, &strLen); err != nil {
		return nil, err
	}

	byteString := make([]byte, strLen)
	if _, err := reader.Read(byteString); err != nil {
		return nil, err
	}

	str := string(byteString)

	return &str, nil
}

func writeByteArray(writer *bytes.Buffer, array *[]byte) error {
	arrayLen := len(*array)
	if err := binary.Write(writer, binary.LittleEndian, arrayLen); err != nil {
		return err
	}
	writer.Write(*array)
	return nil
}

func writeString(writer *bytes.Buffer, str *string) error {
	strLen := len(*str)
	if err := binary.Write(writer, binary.LittleEndian, strLen); err != nil {
		return err
	}
	writer.WriteString(*str)
	return nil
}
