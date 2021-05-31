package ed448

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
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
	typeED448Msg
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
	arrayLen, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	array := make([]byte, arrayLen)
	if _, err := reader.Read(array); err != nil && err != io.EOF {
		return nil, err
	}

	return &array, nil
}

func readString(reader *bytes.Reader) (*string, error) {
	strLen, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	byteString := make([]byte, strLen)
	if _, err := reader.Read(byteString); err != nil && err != io.EOF {
		return nil, err
	}

	str := string(byteString)

	return &str, nil
}

func writeByteArray(writer *bytes.Buffer, array *[]byte) {
	arrayLen := len(*array)
	writeInt(writer, arrayLen)
	writer.Write(*array)
}

func writeString(writer *bytes.Buffer, str *string) {
	strLen := len(*str)
	writeInt(writer, strLen)
	writer.WriteString(*str)
}

func writeInt(writer *bytes.Buffer, val int) {
	var buff = make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(val))

	writer.Write(buff)
}

func readInt(reader *bytes.Reader) (int, error) {
	var buff = make([]byte, 4)
	if _, err := reader.Read(buff); err != nil && err != io.EOF {
		return 0, err
	}

	return int(binary.LittleEndian.Uint32(buff)), nil
}
