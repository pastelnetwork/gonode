package pastel

// MasterNodeStatus represensts pastel masternode status.
type MasterNodeStatus struct {
	Outpoint string `json:"outpoint"`
	Service  string `json:"service"`
	Payee    string `json:"payee"`
	Status   string `json:"status"`
}
