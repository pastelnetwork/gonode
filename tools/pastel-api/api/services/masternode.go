package services

import (
	"context"
	"encoding/json"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services/models"
	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services/models/data"
)

type MasterNode struct{}

func (service *MasterNode) Handle(ctx context.Context, params []string) (interface{}, error) {
	paramsLen := len(params)
	if paramsLen == 0 {
		return nil, errWrongParams
	}

	command := params[0]
	switch command {
	case "top":
		if paramsLen > 1 {
			return nil, errWrongParams
		}
		return service.topMasterNode()
	default:
		return nil, errWrongParams
	}

	return nil, nil
}

func (service *MasterNode) topMasterNode() (interface{}, error) {
	var topMasterNodes models.TopMasterNodes
	res, err := data.ReadFile("masternode_top.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(res, &topMasterNodes); err != nil {
		return nil, err
	}
	return &topMasterNodes, nil
}

func NewMasterNode() *MasterNode {
	return &MasterNode{}
}
