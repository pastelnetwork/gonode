package b85

// This package implements the RFC 1924 Base 85 algorithm

import (
	"errors"
	"fmt"
	"strings"
)

var ErrDecodingB85 = errors.New("base85 decoding error")

const b85chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~"

var decodeMap map[byte]int

func init() {
	length := len(b85chars)
	decodeMap = make(map[byte]int)
	for i := 0; i < length; i++ {
		decodeMap[b85chars[i]] = i
	}
}

// Encode takes a byte array and returns a string of encoded data
func Encode(inData []byte) string {
	var outData strings.Builder

	length := len(inData)
	chunkCount := uint32(length / 4)
	var dataIndex uint32

	for i := uint32(0); i < chunkCount; i++ {
		var decnum, remainder uint32
		decnum = uint32(inData[dataIndex])<<24 | uint32(inData[dataIndex+1])<<16 |
			uint32(inData[dataIndex+2])<<8 | uint32(inData[dataIndex+3])
		outData.WriteByte(b85chars[decnum/52200625])
		remainder = decnum % 52200625
		outData.WriteByte(b85chars[remainder/614125])
		remainder %= 614125
		outData.WriteByte(b85chars[remainder/7225])
		remainder %= 7225
		outData.WriteByte(b85chars[remainder/85])
		outData.WriteByte(b85chars[remainder%85])
		dataIndex += 4
	}

	extraBytes := length % 4
	if extraBytes != 0 {
		lastChunk := uint32(0)
		for i := length - extraBytes; i < length; i++ {
			lastChunk <<= 8
			lastChunk |= uint32(inData[i])
		}

		// Pad extra bytes with zeroes
		for i := (4 - extraBytes); i > 0; i-- {
			lastChunk <<= 8
		}
		outData.WriteByte(b85chars[lastChunk/52200625])
		remainder := lastChunk % 52200625
		outData.WriteByte(b85chars[remainder/614125])
		if extraBytes > 1 {
			remainder %= 614125
			outData.WriteByte(b85chars[remainder/7225])
			if extraBytes > 2 {
				remainder %= 7225
				outData.WriteByte(b85chars[remainder/85])
			}
		}
	}
	return outData.String()
}

// Decode takes in a string of encoded data and returns a byte array of decoded data and an error
// code. The data is considered valid only if the error code is nil. Whitespace is ignored during
// decoding.
func Decode(inData string) ([]byte, error) {

	length := uint32(len(inData))
	outData := make([]byte, length)
	var accumulator, outCount uint32
	var inIndex uint32
	chunkCount := length / 5
	for chunk := uint32(0); chunk < chunkCount; chunk++ {
		accumulator = 0
		for i := 0; i < 5; i++ {
			switch inData[inIndex] {
			case 32, 10, 11, 13:
				// Ignore whitespace
				i--
				inIndex++
				continue
			}
			value, ok := decodeMap[inData[inIndex]]
			if !ok {
				return outData, fmt.Errorf("bad value %v in data", inData[inIndex])
			}
			accumulator = (accumulator * 85) + uint32(value)
			inIndex++
		}
		outData[outCount] = byte(accumulator >> 24)
		outData[outCount+1] = byte((accumulator >> 16) & 255)
		outData[outCount+2] = byte((accumulator >> 8) & 255)
		outData[outCount+3] = byte(accumulator & 255)
		outCount += 4
	}

	remainder := length % 5
	if remainder > 0 {
		accumulator = 0
		for i := uint32(0); i < 5; i++ {
			var value int
			if i < remainder {
				switch inData[inIndex] {
				case 32, 10, 11, 13:
					// Ignore whitespace
					i--
					inIndex++
					continue
				}
				var ok bool
				value, ok = decodeMap[inData[inIndex]]
				if !ok {
					return outData, fmt.Errorf("bad value %v in data", inData[inIndex])
				}
			} else {
				value = 126
			}
			accumulator = (accumulator * 85) + uint32(value)
			inIndex++
		}
		switch remainder {
		case 4:
			outData[outCount] = byte(accumulator >> 24)
			outData[outCount+1] = byte((accumulator >> 16) & 255)
			outData[outCount+2] = byte((accumulator >> 8) & 255)
			outCount += 3
		case 3:
			outData[outCount] = byte(accumulator >> 24)
			outData[outCount+1] = byte((accumulator >> 16) & 255)
			outCount += 2
		case 2:
			outData[outCount] = byte(accumulator >> 24)
			outData[outCount+1] = byte((accumulator >> 16) & 255)
			outCount++
		}
	}
	return outData[:outCount], nil
}
