package node

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"
)

type Miner struct {
	config       config.MinerNode
	pastelClient pastel.Client
}

func NewMiner(config config.MinerNode) *Miner {
	pastelCfg := &pastel.Config{
		Hostname: ptrFromStr(config.PastelAPI.Hostname),
		Port:     ptrFromInt(config.PastelAPI.Port),
		Username: ptrFromStr(config.PastelAPI.Username),
		Password: ptrFromStr(config.PastelAPI.Passphrase),
	}
	pastelClient := pastel.NewClient(pastelCfg)

	return &Miner{
		config:       config,
		pastelClient: pastelClient,
	}
}

func (miner *Miner) GenBlock(ctx context.Context, amount int) error {
	if _, err := miner.pastelClient.GenBlock(ctx, amount); err != nil {
		return errors.Errorf("failed to generate %d blocks: %w", amount, err)
	}
	return nil
}
