package kademlia

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/pastelnetwork/gonode/common/utils"
)

func (s *testSuite) TestHasBit() {
	for i := 0; i < 8; i++ {
		s.Equal(true, hasBit(byte(255), uint(i)))
	}
}

func (s *testSuite) newNodeList(n int) (*NodeList, error) {
	nl := &NodeList{}

	for i := 0; i < n; i++ {
		id, err := newRandomID()
		if err != nil {
			return nil, err
		}
		nl.AddNodes([]*Node{
			{
				ID:   id,
				IP:   s.IP,
				Port: s.bootstrapPort + i + 1,
			},
		})
	}
	return nl, nil
}

func (s *testSuite) TestAddNodes() {
	number := 5

	nl, err := s.newNodeList(number)
	if err != nil {
		s.T().Fatalf("new node list: %v", err)
	}
	s.Equal(number, nl.Len())
}

func (s *testSuite) TestDelNodes() {
	number := 5

	nl, err := s.newNodeList(number)
	if err != nil {
		s.T().Fatalf("new node list: %v", err)
	}
	s.Equal(number, nl.Len())

	ids := [][]byte{}
	for _, node := range nl.Nodes {
		ids = append(ids, node.ID)
	}
	for _, id := range ids {
		nl.DelNode(&Node{ID: id})
	}
	s.Equal(0, nl.Len())
}

func (s *testSuite) TestExists() {
	number := 5

	nl, err := s.newNodeList(number)
	if err != nil {
		s.T().Fatalf("new node list: %v", err)
	}
	s.Equal(number, nl.Len())

	for _, node := range nl.Nodes {
		s.Equal(true, nl.Exists(node))
	}
}

func computeExpectedDistances(nodes []*Node, comparator []byte) [][]byte {
	// Computes the distance between two hashed byte slices.
	distance := func(id1, id2 []byte) *big.Int {
		o1 := new(big.Int).SetBytes(id1)
		o2 := new(big.Int).SetBytes(id2)
		return new(big.Int).Xor(o1, o2)
	}

	// Sort based on distances and return the sorted IDs.
	type nodeDist struct {
		node *Node
		dist *big.Int
	}

	var nodeDists []nodeDist
	for _, node := range nodes {
		dist := distance(node.ID, comparator)
		nodeDists = append(nodeDists, nodeDist{node, dist})
	}

	// Sorting based on distances.
	expected := make([][]byte, len(nodes))
	for i, nd := range nodeDists {
		expected[i] = nd.node.ID
	}

	return expected
}

func TestNodeListSort(t *testing.T) {
	hashedComparator, err := utils.Sha3256hash([]byte{0, 1, 2})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		nodes       []*Node
		comparator  []byte
		expectedIds [][]byte
	}{
		{
			name:       "basic sort",
			comparator: hashedComparator,
		},
		{
			name:       "reversed order",
			comparator: hashedComparator,
		},
	}

	// Pre-compute the hashed IDs for the nodes.
	for _, tt := range tests {
		nodes := []*Node{
			{ID: []byte{}}, // placeholder for hashed ID
			{ID: []byte{}},
			{ID: []byte{}},
		}

		ids := [][]byte{
			[]byte{1, 2, 3},
			[]byte{2, 3, 4},
			[]byte{3, 4, 5},
		}

		for i, rawID := range ids {
			hashedID, err := utils.Sha3256hash(rawID)
			if err != nil {
				t.Fatal(err)
			}
			nodes[i].ID = hashedID
		}

		tt.nodes = nodes
		tt.expectedIds = computeExpectedDistances(tt.nodes, tt.comparator)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nl := &NodeList{
				Nodes:      tt.nodes,
				Comparator: tt.comparator,
			}
			nl.Sort()

			for i, expectedID := range tt.expectedIds {
				if !bytes.Equal(nl.Nodes[i].ID, expectedID) {
					t.Errorf("expected ID at index %d to be %v but got %v", i, expectedID, nl.Nodes[i].ID)
				}
			}
		})
	}
}
