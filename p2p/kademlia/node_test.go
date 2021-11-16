package kademlia

import (
	"sort"
)

func (s *testSuite) initIDWithValues(v byte) []byte {
	id := [20]byte{}
	for i := 0; i < 20; i++ {
		id[i] = v
	}
	return id[:]
}

func (s *testSuite) initZeroIDWithNthByte(n int, v byte) []byte {
	id := s.initIDWithValues(0)
	id[n] = v
	return id
}

func (s *testSuite) TestNodeList() {
	nl := &NodeList{}

	comparator := s.initIDWithValues(0)
	n1 := &Node{ID: s.initZeroIDWithNthByte(19, 1)}
	n2 := &Node{ID: s.initZeroIDWithNthByte(18, 1)}
	n3 := &Node{ID: s.initZeroIDWithNthByte(17, 1)}
	n4 := &Node{ID: s.initZeroIDWithNthByte(16, 1)}

	nl.Nodes = []*Node{n1, n2, n3, n4}
	nl.Comparator = comparator

	sort.Sort(nl)

	s.Equal(n1, nl.Nodes[0])
	s.Equal(n2, nl.Nodes[1])
	s.Equal(n3, nl.Nodes[2])
	s.Equal(n4, nl.Nodes[3])
}

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
