package configs

import "encoding/json"

// Start contains config of the Start command
type Start struct {
	InteractiveMode bool   `json:"interactive_mode" mapstructure:"interactive_mode"`
	Restart         bool   `json:"restart" mapstructure:"restart"`
	Name            string `json:"name" mapstructure:"name"`
	IsTestNet       bool   `json:"is_test_net" mapstructure:"is_test_net"`
	IsCreate        bool   `json:"is_create" mapstructure:"is_create"`
	IsUpdate        bool   `json:"is_update" mapstructure:"is_update"`
	TxID            string `json:"tx_id" mapstructure:"tx_id"`
	IND             string `json:"ind" mapstructure:"ind"`
	IP              string `json:"ip" mapstructure:"ip"`
	Port            int    `json:"port" mapstructure:"port"`
	PrivateKey      string `json:"private_key" mapstructure:"private_key"`
	PastelID        string `json:"pastel_id" mapstructure:"pastel_id"`
	PassPhrase      string `json:"pass_phrase" mapstructure:"pass_phrase"`
	RPCIP           string `json:"rpcip" mapstructure:"rpcip"`
	RPCPort         int    `json:"rpc_port" mapstructure:"rpc_port"`
	NodeIP          string `json:"node_ip" mapstructure:"node_ip"`
	NodePort        int    `json:"node_port" mapstructure:"node_port"`
}

func (s *Start) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// NewStart returns a new Start instance.
func NewStart() *Start {
	return &Start{}
}
