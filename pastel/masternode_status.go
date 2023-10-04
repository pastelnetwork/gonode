package pastel

// MasterNodeStatus represents pastel masternode status.
type MasterNodeStatus struct {
	Outpoint string `json:"outpoint"`
	Service  string `json:"service"`
	Payee    string `json:"payee"`
	Status   string `json:"status"`
	ExtKey   string `json:"extKey"`
}
