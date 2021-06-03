package ed448

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
)

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

const (
	msgTypeLen    = 1
	msgHeaderSize = 5
	msgMaxSize    = 1024 * 1024 // 1 MiB
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

func parseFramedMsg(buf []byte, maxLen uint32) ([]byte, []byte, error) {
	// if header isn't complete we return same buffer
	if len(buf) < msgHeaderSize {
		return nil, buf, nil
	}

	if buf[0] != typeEncryptedMsg {
		return nil, nil, errors.New(ErrWrongFormat)
	}

	msgLenField := buf[msgTypeLen:msgHeaderSize]
	length := binary.LittleEndian.Uint32(msgLenField)
	if length > maxLen {
		return nil, nil, fmt.Errorf("received the frame length %d larger than the limit %d", length, maxLen)
	}
	// ensure that frame is complete, either return same buffer again
	if len(buf) < int(length)+msgHeaderSize {
		// Frame is not complete yet.
		return nil, buf, nil
	}
	// split buf to frames
	return buf[:msgHeaderSize+length], buf[msgHeaderSize+length:], nil
}

func readByteArray(reader *bytes.Reader) (*[]byte, error) {
	arrayLen, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	array := make([]byte, arrayLen)
	if _, err := reader.Read(array); err != nil {
		return nil, errors.Errorf("can not read array %w", err)
	}

	return &array, nil
}

func readString(reader *bytes.Reader) (*string, error) {
	strLen, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	byteString := make([]byte, strLen)
	if _, err := reader.Read(byteString); err != nil {
		return nil, errors.Errorf("can not read string %w", err)
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
	if _, err := reader.Read(buff); err != nil {
		return 0, errors.Errorf("can not read integer value %w", err)
	}

	return int(binary.LittleEndian.Uint32(buff)), nil
}
