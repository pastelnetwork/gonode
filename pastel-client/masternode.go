package pastel

// MasterNodes represensts multiple MasterNode.
type MasterNodes []MasterNode

// MasterNode represensts pastel top masternode.
type MasterNode struct {
	Rank       string  `json:"rank"`
	Address    string  `json:"address"`
	Payee      string  `json:"payee"`
	Outpoint   string  `json:"outpoint"`
	Fee        float64 `json:"fee"`
	ExtAddress string  `json:"extAddress"`
	ExtKey     string  `json:"extKey"`
}
