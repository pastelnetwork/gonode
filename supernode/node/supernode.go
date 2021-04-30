package node

type SuperNodes []*SuperNode

func (nodes *SuperNodes) Add(node *SuperNode) {
	*nodes = append(*nodes, node)
}

func (nodes SuperNodes) FindByKey(key string) *SuperNode {
	for _, node := range nodes {
		if node.Key == key {
			return node
		}
	}
	return nil
}

// SuperNode represents supernode
type SuperNode struct {
	Address string  `json:"extAddress"`
	Key     string  `json:"extKey"`
	Fee     float64 `json:"fee"`
}
