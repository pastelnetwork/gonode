package artworkregister

// SuperNodes is muptiple SuperNode
type SuperNodes []*SuperNode

// SuperNode represents supernode.
type SuperNode struct {
	Address   string  `json:"extAddress"`
	PastelKey string  `json:"extKey"`
	Fee       float64 `json:"fee"`
}
