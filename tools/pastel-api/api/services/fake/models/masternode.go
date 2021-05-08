package models

import "net"

// MasterNode represents a single masternode from top list.
type MasterNode struct {
	Rank          string `json:"rank"`
	IPPort        string `json:"IP:Port"`
	Protocol      int    `json:"protocol"`
	Outpoint      string `json:"outpoint"`
	Payee         string `json:"payee"`
	LastSeen      int    `json:"lastseen"`
	ActiveSeconds int    `json:"activeseconds"`
	ExtAddress    string `json:"extAddress"`
	ExtKey        string `json:"extKey"`
	ExtCfg        string `json:"extCfg"`
}

// MasterNodes represents the API response that can be retrieved using the command `masternode top`.
type MasterNodes []MasterNode

func (nodes MasterNodes) ByPort(port string) *MasterNode {
	for _, node := range nodes {
		_, nodePort, _ := net.SplitHostPort(node.ExtAddress)
		if nodePort == port {
			return &node
		}
	}
	return nil
}
