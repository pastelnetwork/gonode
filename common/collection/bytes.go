package collection

import (
	"bytes"
	"math"
)

// SplitBytesBySize splits the given `data` into chunks with equal `size`, the last chunk can be less.
func SplitBytesBySize(data []byte, size int) [][]byte {
	if len(data) <= size {
		return [][]byte{data}
	}

	total := math.Ceil(float64(len(data)) / float64(size))
	chunks := make([][]byte, int64(total))

	buf := bytes.NewBuffer(data)
	for i := 0; ; i++ {
		chunk := buf.Next(size)
		if len(chunk) == 0 {
			break
		}
		chunks[i] = chunk
	}
	return chunks
}
