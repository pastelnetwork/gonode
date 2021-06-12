package probe

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
)

// Fingerprints is multiple Fingerprint
type Fingerprints []Fingerprint

// Single combines into single Fingerprint type.
func (fgs Fingerprints) Single() Fingerprint {
	var fingerprint Fingerprint

	for _, fg := range fgs {
		fingerprint = append(fingerprint, fg...)
	}
	return fingerprint
}

// Fingerprint represents an image fingerprint.
type Fingerprint []float32

// Bytes converts content to bytes.
func (fg Fingerprint) Bytes() []byte {
	output := new(bytes.Buffer)
	_ = binary.Write(output, binary.LittleEndian, fg)
	return output.Bytes()
}

func fingerprintBaseTypeSize() int {
	return int(reflect.TypeOf(float32(0)).Size())
}

// LSBTruncatedBytes converts content to bytes and truncates least significant byte (LSB) from each float32 value in the array
func (fg Fingerprint) LSBTruncatedBytes() []byte {
	typeSize := fingerprintBaseTypeSize()
	bytes := fg.Bytes()
	var output []byte

	for i, value := range bytes {
		if (i+typeSize)%typeSize != 0 {
			output = append(output, value)
		}
	}
	return output
}

// FingerprintFromTruncatedLSB returns a new Fingerprint instance by converting the given least significant byte (LSB) truncated bytes to its content.
func FingerprintFromTruncatedLSB(data []byte) Fingerprint {
	truncatedValueSize := fingerprintBaseTypeSize() - 1
	output := make([]float32, len(data)/truncatedValueSize)
	zero := []byte{0}

	for i := range output {
		bits := binary.LittleEndian.Uint32(append(zero[:], data[i*truncatedValueSize:(i+1)*truncatedValueSize]...))
		output[i] = math.Float32frombits(bits)
	}
	return output
}

// FingerprintFromBytes returns a new Fingerprint instance by converting the given bytes to its content.
func FingerprintFromBytes(data []byte) Fingerprint {
	typeSize := fingerprintBaseTypeSize()
	output := make([]float32, len(data)/typeSize)

	for i := range output {
		bits := binary.LittleEndian.Uint32(data[i*typeSize : (i+1)*typeSize])
		output[i] = math.Float32frombits(bits)
	}
	return output
}
