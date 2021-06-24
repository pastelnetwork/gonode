package secret

import (
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/otrv4/ed448"
	"github.com/stretchr/testify/assert"
)

func TestED448(t *testing.T) {
	curve := ed448.NewCurve()

	privA, pubA, ok := curve.GenerateKeys()
	if !ok {
		t.Fatal("curve generate keys")
	}

	privB, pubB, ok := curve.GenerateKeys()
	if !ok {
		t.Fatal("curve generate keys")
	}

	secretA := curve.ComputeSecret(privA, pubB)
	secretB := curve.ComputeSecret(privB, pubA)

	assert.Equal(t, secretA, secretB)
}

func TestPastelID(t *testing.T) {
	key := "jXYAvhwhYhanP1S8KjhVHYWRQbwhWXH56DLKDdWWiUfT5tWBjmH4aTtdFmJcLFnRbr1EGUUYW3gif9cs9KTeGa"
	decoded := base58.Decode(key)
	encoded := base58.Encode(decoded)
	assert.Equal(t, key, encoded)
}
