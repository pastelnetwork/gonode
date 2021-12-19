package b85

import (
	"bytes"
	"fmt"

	a85 "encoding/ascii85"
)

// Encode takes a byte array and returns a string of encoded data
func Encode(inData []byte) string {
	out := make([]byte, a85.MaxEncodedLen(len(inData)))
	encodedBytes := a85.Encode(out, inData)

	return string(out[:encodedBytes])
}

// Decode takes in a string of encoded data and returns a byte array of decoded data and an error
// code. The data is considered valid only if the error code is nil. Whitespace is ignored during
// decoding.
func Decode(inData string) ([]byte, error) {
	in := []byte(inData)
	decodedBytes := make([]byte, len(in))
	nDecodedBytes, _, err := a85.Decode(decodedBytes, in, true)
	if err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}
	decodedBytes = decodedBytes[:nDecodedBytes]

	//ascii85 adds /x00 null bytes at the end
	return bytes.Trim(decodedBytes, "\x00"), nil
}
