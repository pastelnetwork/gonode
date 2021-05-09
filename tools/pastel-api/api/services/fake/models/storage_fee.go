package models

// StorageFeeGetNetworkFee represents pastel masternode storage fee, that can be retrieved using the command `storagefee getnetworkfee`.
type StorageFeeGetNetworkFee struct {
	NetworkFee int `json:"networkfee"`
}
