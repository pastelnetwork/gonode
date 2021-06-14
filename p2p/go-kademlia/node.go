package kademlia

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// Node is the over-the-wire representation of a node
type Node struct {
	// id is a 20 byte unique identifier
	ID []byte

	// ip address of the node
	IP string

	// port of the node
	Port int
}

func (s *Node) String() string {
	return fmt.Sprintf("%v-%v:%d", hex.EncodeToString(s.ID), s.IP, s.Port)
}

// NewNode returns a new node for bootstrapping
func NewNode(ip string, port int) *Node {
	return &Node{
		IP:   ip,
		Port: port,
	}
}

// NodeList is used in order to sort a list of nodes
type NodeList struct {
	Nodes []*Node

	// Comparator is the id to compare to
	Comparator []byte
}

func (s *NodeList) String() string {
	nodes := []string{}
	for _, node := range s.Nodes {
		nodes = append(nodes, node.String())
	}
	return strings.Join(nodes, ",")
}

// CompareNodes checks if two nodes are equal
func CompareNodes(n1 *Node, n2 *Node) bool {
	if n1 == nil || n2 == nil {
		return false
	}
	if n1.IP != n2.IP {
		return false
	}
	if n1.Port != n2.Port {
		return false
	}
	return true
}

// DelNode deletes a node from list
func (n *NodeList) DelNode(node *Node) {
	for i := 0; i < n.Len(); i++ {
		if bytes.Equal(n.Nodes[i].ID, node.ID) {
			n.Nodes = append(n.Nodes[:i], n.Nodes[i+1:]...)
			return
		}
	}
}

// Exists return true if the node is already there
func (n *NodeList) Exists(node *Node) bool {
	for _, item := range n.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			return true
		}
	}
	return false
}

// AddNode appends the nodes to node list if it's not existed
func (n *NodeList) AddNodes(nodes []*Node) {
	for _, node := range nodes {
		if !n.Exists(node) {
			n.Nodes = append(n.Nodes, node)
		}
	}
}

func (n *NodeList) Len() int {
	return len(n.Nodes)
}
func (n *NodeList) Swap(i, j int) {
	n.Nodes[i], n.Nodes[j] = n.Nodes[j], n.Nodes[i]
}
func (n *NodeList) Less(i, j int) bool {
	id := n.distance(n.Nodes[i].ID, n.Comparator)
	jd := n.distance(n.Nodes[j].ID, n.Comparator)

	return id.Cmp(jd) == -1
}

func (n *NodeList) distance(id1, id2 []byte) *big.Int {
	o1 := new(big.Int).SetBytes(id1)
	o2 := new(big.Int).SetBytes(id2)

	return new(big.Int).Xor(o1, o2)
}
