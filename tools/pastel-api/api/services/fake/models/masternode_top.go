package models

// TopMasterNodes represents the API response that can be retrieved using the command `masternode top`.
type TopMasterNodes map[string]MasterNodes

func (blocks TopMasterNodes) LastBlock() MasterNodes {
	for _, nodes := range blocks {
		return nodes
	}
	return nil
}
