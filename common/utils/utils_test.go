package utils

import (
	"bytes"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		level  int
		expect func([]byte, []byte) bool
	}{
		{
			name:  "basic",
			input: []byte("Hello, world!"),
			level: 2,
			expect: func(orig, decompressed []byte) bool {
				return bytes.Equal(orig, decompressed)
			},
		},
		{
			name:  "basic",
			input: []byte("Hello, world Level 2!"),
			level: 1,
			expect: func(orig, decompressed []byte) bool {
				return bytes.Equal(orig, decompressed)
			},
		},
		{
			name:  "basic",
			input: []byte("Hello, world Now!"),
			level: 1,
			expect: func(orig, decompressed []byte) bool {
				return bytes.Equal(orig, decompressed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := Compress(tt.input, tt.level)
			if err != nil {
				t.Fatalf("Error compressing: %v", err)
			}

			decompressed, err := Decompress(compressed)
			if err != nil {
				t.Fatalf("Error decompressing: %v", err)
			}

			if !tt.expect(tt.input, decompressed) {
				t.Fatalf("Expected result not obtained. Original: %v, Decompressed: %v", string(tt.input), string(decompressed))
			}
		})
	}
}
