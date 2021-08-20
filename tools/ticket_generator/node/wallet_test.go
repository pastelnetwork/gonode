package node

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"
	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	config := config.WalletNode{
		PastelAPI: config.PastelAPI{
			Hostname:   "localhost",
			Port:       11170,
			Username:   "rt",
			Passphrase: "rt",
		},
		RestAPI: config.RestAPI{
			Hostname: "localhost",
			Port:     8080,
		},
		Artist: config.Artist{
			PastelID:      "jXa8KpmY4cf2PMXbHFL6wYkS7MaD1FjqHV1MUJPnVv2Ayua8TQ3PUu2Dtyr9mTUEW1BGZbCrdQ54mPp6gsaF1s",
			Passphrase:    "passphrase",
			SpendableAddr: "tPSFddCX7intdWJpig5931VeEAYPnkt4CFL",
		},
	}

	wallet := NewWallet(config)
	img := "../data/test-image1.jpg"
	imgID, err := wallet.UploadImage(context.Background(), img)

	assert.Nil(t, err)
	assert.NotEmpty(t, imgID)
}
