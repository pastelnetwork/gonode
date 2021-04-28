package legroast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeygen(t *testing.T) {
	pk, sk := Keygen()

	assert.Equal(t, len(pk), PKBytes, "Public key length doesn't match.")
	assert.Equal(t, len(sk), SKBytes, "Private key length doesn't match.")
}

func TestSignVerify(t *testing.T) {
	pk, sk := Keygen()

	msg := "Message to sign."
	m := ([]byte)(msg[:])

	signed := Sign(pk, sk, m)
	valid := Verify(pk, m, signed)
	assert.True(t, valid > 0, "Signature is not valid.")
}

func TestVerifyFail(t *testing.T) {
	pk, sk := Keygen()

	msg := "Message to sign."
	m := ([]byte)(msg[:])

	invalidMsg := "Wrong message."
	iM := ([]byte)(invalidMsg[:])

	signed := Sign(pk, sk, m)
	valid := Verify(pk, iM, signed)
	assert.True(t, valid <= 0, "Invalid Verify result is expected.")

	valid = Verify(pk, m, iM)
	assert.True(t, valid <= 0, "Invalid Verify result is expected.")
}
