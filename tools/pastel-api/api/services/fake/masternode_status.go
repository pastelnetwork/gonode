package fake

import "github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"

func newMasterNodeStatusByNode(node *models.MasterNode) *models.MasterNodeStatus {
	return &models.MasterNodeStatus{
		Outpoint: node.Outpoint,
		Service:  node.IPPort,
		Payee:    node.Payee,
		Status:   "Masternode successfully started",
	}
}
