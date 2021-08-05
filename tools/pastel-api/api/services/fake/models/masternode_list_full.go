package models

// FullListMasterNodes represents pastel top masternodes, that can be retrieved using the command `masternode top`.
type FullListMasterNodes map[string]MasterNodes

// LastBlock returns a list of the masternodes from the last block.
func (blocks FullListMasterNodes) LastBlock() MasterNodes {
	for _, nodes := range blocks {
		return nodes
	}
	return nil
}
