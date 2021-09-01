package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	viper.SetConfigFile("example.yaml")
	err := viper.ReadInConfig()
	assert.Nil(t, err, "read configure failed: %s: %s", viper.ConfigFileUsed(), err)

	config := struct {
		MinerNode  MinerNode  `json:"miner_node" mapstructure:"miner_node"`
		WalletNode WalletNode `json:"wallet_node" mapstructure:"wallet_node"`
	}{}

	err = viper.Unmarshal(&config)
	assert.Nil(t, err, "marshal miner node failed: %s", err)
	assert.Equal(t, config.MinerNode.PastelAPI.Hostname, "localhost")
	assert.Equal(t, config.MinerNode.PastelAPI.Port, 12169)
	assert.Equal(t, config.MinerNode.PastelAPI.Username, "rt")
	assert.Equal(t, config.MinerNode.PastelAPI.Passphrase, "rt")

	wallet := config.WalletNode

	assert.Equal(t, wallet.PastelAPI.Hostname, "localhost")
	assert.Equal(t, wallet.PastelAPI.Port, 12170)
	assert.Equal(t, wallet.PastelAPI.Username, "rt")
	assert.Equal(t, wallet.PastelAPI.Passphrase, "rt")

	assert.Equal(t, wallet.RestAPI.Hostname, "localhost")
	assert.Equal(t, wallet.RestAPI.Port, 8080)

	assert.Equal(t, wallet.Artist.PastelID, "jXa8KpmY4cf2PMXbHFL6wYkS7MaD1FjqHV1MUJPnVv2Ayua8TQ3PUu2Dtyr9mTUEW1BGZbCrdQ54mPp6gsaF1s")
	assert.Equal(t, wallet.Artist.Passphrase, "passphrase")
	assert.Equal(t, wallet.Artist.SpendableAddr, "tPSFddCX7intdWJpig5931VeEAYPnkt4CFL")
}
