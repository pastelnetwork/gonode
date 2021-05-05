package static

// StorageFeeGetNetworkFee represents the API response that can be retrieved using the command `storagefee getnetworkfee`.
type StorageFeeGetNetworkFee struct {
	NetworkFee int `json:"networkfee"`
}
