package kademlia

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

// Node is the over-the-wire representation of a node
type Node struct {
	// id is a 32 byte unique identifier
	ID []byte `json:"id,omitempty"`

	// ip address of the node
	IP string `json:"ip,omitempty"`

	// port of the node
	Port int `json:"port,omitempty"`
}

func (s *Node) String() string {
	return fmt.Sprintf("%v-%v:%d", base58.Encode(s.ID), s.IP, s.Port)
}

// NodeList is used in order to sort a list of nodes
type NodeList struct {
	Nodes []*Node

	// Comparator is the id to compare to
	Comparator []byte
}

// String returns the dump information for node list
func (s *NodeList) String() string {
	nodes := []string{}
	for _, node := range s.Nodes {
		nodes = append(nodes, node.String())
	}
	return strings.Join(nodes, ",")
}

// DelNode deletes a node from list
func (s *NodeList) DelNode(node *Node) {
	for i := 0; i < s.Len(); i++ {
		if bytes.Equal(s.Nodes[i].ID, node.ID) {
			newList := s.Nodes[:i]
			if i+1 < s.Len() {
				newList = append(newList, s.Nodes[i+1:]...)
			}

			s.Nodes = newList

			return
		}
	}
}

// Exists return true if the node is already there
func (s *NodeList) Exists(node *Node) bool {
	for _, item := range s.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			return true
		}
	}
	return false
}

// AddNodes appends the nodes to node list if it's not existed
func (s *NodeList) AddNodes(nodes []*Node) {
	for _, node := range nodes {
		if !s.Exists(node) {
			s.Nodes = append(s.Nodes, node)
		}
	}
}

// Len check length
func (s *NodeList) Len() int {
	return len(s.Nodes)
}

// Swap swap two nodes
func (s *NodeList) Swap(i, j int) {
	s.Nodes[i], s.Nodes[j] = s.Nodes[j], s.Nodes[i]
}

// Less compare two nodes
func (s *NodeList) Less(i, j int) bool {
	id := s.distance(s.Nodes[i].ID, s.Comparator)
	jd := s.distance(s.Nodes[j].ID, s.Comparator)

	return id.Cmp(jd) == -1
}

func (s *NodeList) distance(id1, id2 []byte) *big.Int {
	o1 := new(big.Int).SetBytes(id1)
	o2 := new(big.Int).SetBytes(id2)

	return new(big.Int).Xor(o1, o2)
}
