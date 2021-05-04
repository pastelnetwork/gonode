package methods

type MasterNodes []*MasterNode

func (nodes *MasterNodes) Top() MasterNodes {
	return nil
}

type MasterNode struct {
	Rank          string `json:"rank"`
	IPPort        string `json:"IP:Port"`
	Protocol      int    `json:"protocol"`
	Outpoint      string `json:"outpoint"`
	Payee         string `json:"payee"`
	Lastseen      int    `json:"lastseen"`
	Activeseconds int    `json:"activeseconds"`
	ExtAddress    string `json:"extAddress"`
	ExtKey        string `json:"extKey"`
	ExtCfg        string `json:"extCfg"`
}
