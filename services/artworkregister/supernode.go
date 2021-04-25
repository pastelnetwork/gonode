package artworkregister

type SuperNodes []*SuperNode

type SuperNode struct {
	Address   string  `json:"extAddress"`
	PastelKey string  `json:"extKey"`
	Fee       float64 `json:"fee"`
}
