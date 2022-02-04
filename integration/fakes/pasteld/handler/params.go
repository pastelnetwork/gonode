package handler

import (
	"errors"
	"fmt"

	"github.com/pastelnetwork/gonode/integration/fakes/common/testconst"
)

// RpcParamsis used for unmarshaling the method params
type RpcParams struct {
	Params []string `json:"params"`
	N      int      `json:"n"`
}

// FromPositional will be passed an array of interfaces if positional parameters
// are passed in the rpc call
func (p *RpcParams) FromPositional(params []interface{}) error {
	if len(params) < 1 {
		return errors.New("atleast 1 arg is required")
	}

	for _, param := range params {
		str := fmt.Sprint(param)
		p.Params = append(p.Params, str)
	}

	return nil
}

type GetOperationStatusResult struct {
	ID            string                 `json:"id,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Error         map[string]interface{} `json:"error,omitempty"`
	Result        map[string]interface{} `json:"result"`
	ExecutionSecs float64                `json:"execution_secs"`
	Method        string                 `json:"method"`
	Params        map[string]interface{} `json:"params"`
	MinConf       int                    `json:"min_conf"`
	Fee           float64                `json:"fee"`
}

func isValidMNPastelID(id string) bool {
	for _, val := range testconst.SNsPastelIDs {
		if id == val {
			return true
		}
	}

	return false
}
