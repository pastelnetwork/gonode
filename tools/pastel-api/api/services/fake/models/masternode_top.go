package models

// TopMasterNodes represents pastel top masternodes, that can be retrieved using the command `masternode top`.
type TopMasterNodes map[string]MasterNodes

// LastBlock returns a list of the masternodes from the last block.
func (blocks TopMasterNodes) LastBlock() MasterNodes {
	for _, nodes := range blocks {
		return nodes
	}
	return nil
}
