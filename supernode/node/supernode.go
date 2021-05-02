package node

// SuperNodes represents muptiple SuperNodes
type SuperNodes []*SuperNode

// Add adds a new node to the list
func (nodes *SuperNodes) Add(node *SuperNode) {
	*nodes = append(*nodes, node)
}

// FindByKey returns a node from the list by the given key.
func (nodes SuperNodes) FindByKey(key string) *SuperNode {
	for _, node := range nodes {
		if node.Key == key {
			return node
		}
	}
	return nil
}

// Remove removes a node from the list by the given key.
func (nodes *SuperNodes) Remove(key string) {
	for i, node := range *nodes {
		if node.Key == key {
			*nodes = append((*nodes)[:i], (*nodes)[:i+1]...)
			break
		}
	}
}

// SuperNode represents a single supernode
type SuperNode struct {
	Address string  `json:"extAddress"`
	Key     string  `json:"extKey"`
	Fee     float64 `json:"fee"`
}
