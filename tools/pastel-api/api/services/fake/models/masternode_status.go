package models

// MasterNodeStatus represents pastel masternode status, that can be retrieved using the command `masternode status`.
type MasterNodeStatus struct {
	Outpoint string `json:"outpoint"`
	Service  string `json:"service"`
	Payee    string `json:"payee"`
	Status   string `json:"status"`
}
