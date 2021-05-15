package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceForAppend(t *testing.T) {
	t.Parallel()
	const GcmTagSize int = 16

	plainText := "TestSliceForAppend"
	plainTextByte := []byte(plainText)

	dst, tail := SliceForAppend(plainTextByte, len(plainText)+GcmTagSize)

	assert.Equal(t, len(tail), len(plainText)+GcmTagSize)
	assert.Equal(t, len(dst), len(plainTextByte)+len(tail)) //Expected dst to have enough capacity to avoid a reallocation

	newB := make([]byte, len(plainTextByte)+GcmTagSize) // creating new byte that has size of plaintext and GcmTagSize
	assert.Equal(t, newB, tail)

	data := tail[:len(plainText)]
	copy(data, plainText)

	assert.Equal(t, data, plainTextByte) //Expected the data to be same with plaintext after the reallocation"
}
