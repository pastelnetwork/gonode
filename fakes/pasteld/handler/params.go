package handler

import "errors"

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
		str := param.(string)
		p.Params = append(p.Params, str)
	}

	return nil
}
