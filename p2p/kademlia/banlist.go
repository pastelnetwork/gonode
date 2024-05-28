package kademlia

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"
)

const (
	// banDuration - ban duration
	banDuration = 3 * time.Hour

	// threshold - threshold
	threshold = 3
)

// BanNode is the over-the-wire representation of a node
type BanNode struct {
	Node
	CreatedAt time.Time
	count     int
}

// BanList is used in order to sort a list of nodes
type BanList struct {
	Nodes []BanNode
	mtx   sync.RWMutex
}

// Add adds a node to the ban list
func (s *BanList) Add(node *Node) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	banNode := BanNode{
		Node: Node{
			ID:   node.ID,
			IP:   node.IP,
			Port: node.Port,
		},
		CreatedAt: time.Now().UTC(),
		count:     1,
	}

	s.Nodes = append(s.Nodes, banNode)
}

// IncrementCount increments the count of a node in the ban list
func (s *BanList) IncrementCount(node *Node) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	found := false
	for i := 0; i < len(s.Nodes); i++ {
		if bytes.Equal(s.Nodes[i].ID, node.ID) {
			s.Nodes[i].count++
			found = true

			break
		}
	}

	if !found {
		banNode := BanNode{
			Node: Node{
				ID:   node.ID,
				IP:   node.IP,
				Port: node.Port,
			},
			CreatedAt: time.Now().UTC(),
			count:     1,
		}

		s.Nodes = append(s.Nodes, banNode)
	}
}

// Banned return true if the node is banned
func (s *BanList) Banned(node *Node) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, item := range s.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			return item.count > threshold
		}
	}

	return false
}

// Exists return true if the node is already there
func (s *BanList) Exists(node *Node) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, item := range s.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			return true
		}
	}

	return false
}

// Delete deletes a node from list
func (s *BanList) Delete(node *Node) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	l := len(s.Nodes)
	for i := 0; i < l; i++ {
		if bytes.Equal(s.Nodes[i].ID, node.ID) {
			newNodes := s.Nodes[:i]
			if i+1 < l {
				newNodes = append(newNodes, s.Nodes[i+1:]...)
			}
			s.Nodes = newNodes

			return
		}
	}
}

// Purge removes all expired nodes from the ban list
func (s *BanList) Purge() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for i := 0; i < len(s.Nodes); i++ {
		if time.Now().UTC().After(s.Nodes[i].CreatedAt.Add(banDuration)) {

			newNodes := s.Nodes[:i]
			if i+1 < len(s.Nodes) {
				newNodes = append(newNodes, s.Nodes[i+1:]...)
				i--
			}

			s.Nodes = newNodes
		}
	}
}

// schedulePurge schedules a purge of the ban list
func (s *BanList) schedulePurge(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Purge()
		case <-ctx.Done():
			return
		}
	}
}

// ToNodeList returns the list of nodes
func (s *BanList) ToNodeList() []*Node {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	ret := make([]*Node, 0)

	for i := 0; i < len(s.Nodes); i++ {
		if s.Nodes[i].count > threshold {
			ret = append(ret, &s.Nodes[i].Node)
		}
	}

	return ret
}

// AddWithCreatedAt adds a node to the ban list with a specific creation date
func (s *BanList) AddWithCreatedAt(node *Node, createdAt time.Time, count int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	banNode := BanNode{
		Node: Node{
			ID:   node.ID,
			IP:   node.IP,
			Port: node.Port,
		},
		CreatedAt: createdAt,
		count:     count,
	}

	s.Nodes = append(s.Nodes, banNode)
}

// String returns the dump information for node list
func (s *BanList) String() string {
	nodes := []string{}
	for _, node := range s.Nodes {
		nodes = append(nodes, node.Node.String())
	}
	return strings.Join(nodes, ",")
}

// NewBanList creates a new ban list
func NewBanList(ctx context.Context) *BanList {
	list := &BanList{
		Nodes: make([]BanNode, 0),
	}

	go list.schedulePurge(ctx)

	return list
}
