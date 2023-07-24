package kademlia

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"

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

	Mux sync.RWMutex
}

// String returns the dump information for node list
func (s *NodeList) String() string {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	var nodes []string
	for _, node := range s.Nodes {
		nodes = append(nodes, node.String())
	}
	return strings.Join(nodes, ",")
}

// DelNode deletes a node from list
func (s *NodeList) DelNode(node *Node) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

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
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	return s.exists(node)
}

func (s *NodeList) exists(node *Node) bool {
	for _, item := range s.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			return true
		}
	}
	return false
}

// AddNodes appends the nodes to node list if it's not existed
func (s *NodeList) AddNodes(nodes []*Node) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	for _, node := range nodes {
		if !s.exists(node) {
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
	if i >= 0 && i < s.Len() && j >= 0 && j < s.Len() {
		s.Nodes[i], s.Nodes[j] = s.Nodes[j], s.Nodes[i]
	}
}

// Sort sorts nodes
func (s *NodeList) Sort() {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	sort.Sort(s)
}

// Less compare two nodes
func (s *NodeList) Less(i, j int) bool {
	if i >= 0 && i < s.Len() && j >= 0 && j < s.Len() {
		id := s.distance(s.Nodes[i].ID, s.Comparator)
		jd := s.distance(s.Nodes[j].ID, s.Comparator)

		return id.Cmp(jd) == -1
	}
	return false
}

func (s *NodeList) distance(id1, id2 []byte) *big.Int {
	o1 := new(big.Int).SetBytes(id1)
	o2 := new(big.Int).SetBytes(id2)

	return new(big.Int).Xor(o1, o2)
}

// AddFirst adds a node to the first position of the list.
func (s *NodeList) AddFirst(node *Node) {
	s.Mux.Lock()         // lock for writing
	defer s.Mux.Unlock() // ensure the lock is released no matter what

	// Prepend the node to the slice
	s.Nodes = append([]*Node{node}, s.Nodes...)
}

// TopN keeps top N nodes
func (s *NodeList) TopN(n int) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	if s.Len() > n {
		s.Nodes = s.Nodes[:n]
	}
}
