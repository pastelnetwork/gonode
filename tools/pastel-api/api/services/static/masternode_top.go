package static

// TopMasterNodes represents the API response that can be retrieved using the command `masternode top`.
type TopMasterNodes map[int][]TopMasterNode

// TopMasterNode represents a single masternode from top list.
type TopMasterNode struct {
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
