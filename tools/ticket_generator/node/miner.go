package node

import (
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"
)

type Miner struct {
	config       config.MinderNode
	pastelClient pastel.Client
}
