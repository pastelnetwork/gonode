package common

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

// SuperNodeList keeps list of SuperNodeClient for networking
type SuperNodeList []*SuperNodeClient

// AddNewNode created and adds a new node to the list.
func (nodes *SuperNodeList) AddNewNode(client node.ClientInterface, address string, pastelID string, nodeMaker node.RealNodeMaker) {
	someNode := NewSuperNode(client, address, pastelID, nodeMaker)
	*nodes = append(*nodes, someNode)
}

// Add adds a new node to the list.
func (nodes *SuperNodeList) Add(node *SuperNodeClient) {
	*nodes = append(*nodes, node)
}

// Activate marks all nodes as activated.
// Since any node can be present in the same time in several List and Node is a pointer, this is reflected in all lists.
func (nodes *SuperNodeList) Activate() {
	for _, someNode := range *nodes {
		someNode.SetActive(true)
	}
}

// WaitConnClose waits for the connection closing by any supernodes.
func (nodes *SuperNodeList) WaitConnClose(ctx context.Context, done <-chan struct{}) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, node := range *nodes {
		node := node
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case <-node.ConnectionInterface.Done():
				return errors.Errorf("%q unexpectedly closed the connection", node)
			case <-done:
				return nil
			}
		})
	}

	return group.Wait()
}

// FindByPastelID returns node by its patstelID.
func (nodes *SuperNodeList) FindByPastelID(id string) *SuperNodeClient {
	for _, someNode := range *nodes {
		fmt.Println("someNode.PastelID():", someNode.PastelID())
		fmt.Println("id", id)
		if someNode.PastelID() == id {
			return someNode
		}
	}
	return nil
}

// String returns the comma separated addresses of nodes
func (nodes *SuperNodeList) String() string {
	var commaSeparatedList string
	for _, someNode := range *nodes {
		commaSeparatedList += someNode.String() + ","
	}

	return commaSeparatedList
}

func (nodes *SuperNodeList) distance(id1, id2 []byte) *big.Int {
	o1 := new(big.Int).SetBytes(id1)
	o2 := new(big.Int).SetBytes(id2)

	return new(big.Int).Xor(o1, o2)
}

// Sort sorts nodes
func (nodes *SuperNodeList) Sort(key []byte) {
	for i := 0; i < nodes.Len(); i++ {
		(*nodes)[i].Comparator = key
		hash, _ := utils.Sha3256hash([]byte((*nodes)[i].PastelID()))
		(*nodes)[i].HashedID = hash
	}

	// Sort using the precomputed distances
	sort.Sort(nodes)
}

// Swap swap two nodes
func (nodes *SuperNodeList) Swap(i, j int) {
	if i >= 0 && i < nodes.Len() && j >= 0 && j < nodes.Len() {
		(*nodes)[i], (*nodes)[j] = (*nodes)[j], (*nodes)[i]
	}
}

// Less compare two nodes
func (nodes *SuperNodeList) Less(i, j int) bool {
	if i >= 0 && i < nodes.Len() && j >= 0 && j < nodes.Len() {
		id := nodes.distance((*nodes)[i].HashedID, (*nodes)[i].Comparator)
		jd := nodes.distance((*nodes)[j].HashedID, (*nodes)[j].Comparator)

		return id.Cmp(jd) == -1
	}

	return false
}

func (nodes *SuperNodeList) Len() int {
	return len(*nodes)
}
