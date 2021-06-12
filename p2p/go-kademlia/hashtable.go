package kademlia

import (
	"bytes"
	"math/big"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const (
	IterateStore = iota
	IterateFindNode
	IterateFindValue
)

const (
	// a small number representing the degree of parallelism in network calls
	Alpha = 3

	// the size in bits of the keys used to identify nodes and store and
	// retrieve data; in basic Kademlia this is 160, the length of a SHA1
	B = 160

	// the maximum number of contacts stored in a bucket
	K = 20
)

// HashTable represents the hashtable state
type HashTable struct {
	// The ID of the local node
	Self *Node

	// Route table a list of all known nodes in the network
	// Nodes within buckets are sorted by least recently seen e.g.
	// [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
	//  ^                                                           ^
	//  └ Least recently seen                    Most recently seen ┘
	RouteTable [B][]*Node // 160x20

	// mutex for route table
	mutex sync.Mutex

	refreshers [B]time.Time
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewHashTable returns a new hashtable
func NewHashTable(options *Options) (*HashTable, error) {
	ht := &HashTable{
		Self: &Node{
			IP:   options.IP,
			Port: options.Port,
		},
	}
	// init the id for hashtable
	if options.ID != nil {
		ht.Self.ID = options.ID
	} else {
		id, err := newRandomID()
		if err != nil {
			return nil, err
		}
		ht.Self.ID = id
	}

	// reset the refresh time for every bucket
	for i := 0; i < B; i++ {
		ht.ResetRefreshTime(i)
	}

	return ht, nil
}

// ResetRefreshTime - reset the refresh time
func (ht *HashTable) ResetRefreshTime(bucket int) {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()
	ht.refreshers[bucket] = time.Now()
}

// RefreshNode makes the node to the end
func (ht *HashTable) RefreshNode(id []byte) {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	// bucket index of the node
	index := BucketIndex(ht.Self.ID, id)
	// point to the bucket
	bucket := ht.RouteTable[index]

	var offset int
	// find the position of the node
	for i, v := range bucket {
		if bytes.Equal(v.ID, id) {
			offset = i
			break
		}
	}

	// makes the node to the end
	current := bucket[offset]
	bucket = append(bucket[:offset], bucket[offset+1:]...)
	bucket = append(bucket, current)
	ht.RouteTable[index] = bucket
}

// NodeExists check if the node id is existed
func (ht *HashTable) NodeExists(bucket int, id []byte) bool {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	for _, node := range ht.RouteTable[bucket] {
		if bytes.Equal(node.ID, id) {
			return true
		}
	}
	return false
}

// ClosestContacts returns the closest contacts of target
func (ht *HashTable) ClosestContacts(num int, target []byte, ignoredNodes []*Node) *NodeList {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	// find the bucket index in local route tables
	index := BucketIndex(ht.Self.ID, target)
	indexList := []int{index}
	i := index - 1
	j := index + 1
	for len(indexList) < B {
		if j < B {
			indexList = append(indexList, j)
		}
		if i >= 0 {
			indexList = append(indexList, i)
		}
		i--
		j++
	}

	nl := &NodeList{}

	left := num
	// select alpha contacts and add them to the node list
	for left > 0 && len(indexList) > 0 {
		index, indexList = indexList[0], indexList[1:]
		for i := 0; i < len(ht.RouteTable[index]); i++ {
			node := ht.RouteTable[index][i]

			ignored := false
			for j := 0; j < len(ignoredNodes); j++ {
				if bytes.Equal(node.ID, ignoredNodes[j].ID) {
					ignored = true
				}
			}
			if !ignored {
				// add the node to list
				nl.AddNodes([]*Node{node})

				left--
				if left == 0 {
					break
				}
			}
		}
	}

	// sort the node list
	sort.Sort(nl)

	return nl
}

// CloserNodes returns the closer nodes in bucket for id comparing the local node
func (ht *HashTable) CloserNodes(bucket int, id []byte) [][]byte {
	var nodes [][]byte
	for _, node := range ht.RouteTable[bucket] {
		// distance between id and self
		d1 := ht.Distance(id, ht.Self.ID)
		// distance between id and node in bucket
		d2 := ht.Distance(id, node.ID)

		// if 2-self > 2-id
		if d1.Sub(d1, d2).Sign() >= 0 {
			nodes = append(nodes, node.ID)
		}
	}

	return nodes
}

// GetBucketNodes returns the nodes count of the bucket
func (ht *HashTable) GetBucketNodes(bucket int) int {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	return len(ht.RouteTable[bucket])
}

// Distance returns the distance between two ids
func (ht *HashTable) Distance(id1 []byte, id2 []byte) *big.Int {
	var dst [K]byte
	for i := 0; i < K; i++ {
		dst[i] = id1[i] ^ id2[i]
	}
	return big.NewInt(0).SetBytes(dst[:])
}

// BucketIndex return the bucket index from two node ids
func BucketIndex(id1 []byte, id2 []byte) int {
	// look at each byte from left to right
	for j := 0; j < len(id1); j++ {
		// xor the byte
		xor := id1[j] ^ id2[j]

		// check each bit on the xored result from left to right in order
		for i := 0; i < 8; i++ {
			if hasBit(xor, uint(i)) {
				byteIndex := j * 8
				bitIndex := i
				return B - (byteIndex + bitIndex) - 1
			}
		}
	}

	// only happen during bootstrapping
	return 0
}

// Count returns the number of nodes in route table
func (ht *HashTable) TotalCount() int {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	var num int
	for _, v := range ht.RouteTable {
		num += len(v)
	}
	return num
}

// newRandomID returns a new random id
func newRandomID() ([]byte, error) {
	id := make([]byte, 20)
	_, err := rand.Read(id)
	return id, err
}

// Simple helper function to determine the value of a particular
// bit in a byte by index

// Example:
// number:  1
// bits:    00000001
// pos:     01234567
func hasBit(n byte, pos uint) bool {
	pos = 7 - pos
	val := n & (1 << pos)
	return (val > 0)
}
