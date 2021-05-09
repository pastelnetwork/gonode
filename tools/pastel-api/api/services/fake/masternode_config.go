package fake

import "github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"

func newMasterNodeListConfByNode(node *models.MasterNode) *models.MasterNodeListConf {
	if node == nil {
		// try to replicate exactly the real Pastel API
		return &models.MasterNodeListConf{}
	}

	return &models.MasterNodeListConf{
		"masternode": models.MasterNodeConfig{
			Address:    node.IPPort,
			ExtAddress: node.ExtAddress,
			ExtKey:     node.ExtKey,
			ExtCfg:     node.ExtCfg,
		},
	}
}
