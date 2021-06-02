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

// FingerprintFromBytes returns a new Fingerprint instance by converting the given bytes to its content.
func FingerprintFromBytes(data []byte) Fingerprint {
	typeSize := int(reflect.TypeOf(float32(0)).Size())
	output := make([]float32, len(data)/typeSize)

	for i := range output {
		bits := binary.LittleEndian.Uint32(data[i*typeSize : (i+1)*typeSize])
		output[i] = math.Float32frombits(bits)
	}
	return output
}
