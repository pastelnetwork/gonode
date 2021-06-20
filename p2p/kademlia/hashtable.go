package kademlia

import (
	"bytes"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const (
	// IterateStore - iterative store the data
	IterateStore = iota
	// IterateFindNode - iterative find node
	IterateFindNode
	// IterateFindValue - iterative find value
	IterateFindValue
)

const (
	// Alpha - a small number representing the degree of parallelism in network calls
	Alpha = 3

	// B - the size in bits of the keys used to identify nodes and store and
	// retrieve data; in basic Kademlia this is 160, the length of a SHA1
	B = 160

	// K - the maximum number of contacts stored in a bucket
	K = 20
)

// HashTable represents the hashtable state
type HashTable struct {
	// The ID of the local node
	self *Node

	// Route table a list of all known nodes in the network
	// Nodes within buckets are sorted by least recently seen e.g.
	// [ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ]
	//  ^                                                           ^
	//  └ Least recently seen                    Most recently seen ┘
	routeTable [][]*Node // 160x20

	// mutex for route table
	mutex sync.Mutex

	// refresh time for every bucket
	refreshers []time.Time
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewHashTable returns a new hashtable
func NewHashTable(options *Options) (*HashTable, error) {
	ht := &HashTable{
		self: &Node{
			IP:   options.IP,
			Port: options.Port,
		},
		refreshers: make([]time.Time, B),
		routeTable: make([][]*Node, B),
	}
	// init the id for hashtable
	if options.ID != nil {
		ht.self.ID = options.ID
	} else {
		id, err := newRandomID()
		if err != nil {
			return nil, err
		}
		ht.self.ID = id
	}

	// reset the refresh time for every bucket
	for i := 0; i < B; i++ {
		ht.resetRefreshTime(i)
	}

	return ht, nil
}

// resetRefreshTime - reset the refresh time
func (ht *HashTable) resetRefreshTime(bucket int) {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()
	ht.refreshers[bucket] = time.Now()
}

// refreshNode makes the node to the end
func (ht *HashTable) refreshNode(id []byte) {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	// bucket index of the node
	index := ht.bucketIndex(ht.self.ID, id)
	// point to the bucket
	bucket := ht.routeTable[index]

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
	ht.routeTable[index] = bucket
}

// refreshTime returns the refresh time for bucket
func (ht *HashTable) refreshTime(bucket int) time.Time {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	return ht.refreshers[bucket]
}

// randomIDFromBucket returns a random id based on the bucket index and self id
func (ht *HashTable) randomIDFromBucket(bucket int) []byte {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	// set the new ID to to be equal in every byte up to
	// the byte of the first different bit in the bucket
	index := bucket / 8
	var id []byte
	for i := 0; i < index; i++ {
		id = append(id, ht.self.ID[i])
	}
	start := bucket % 8

	var first byte
	// check each bit from left to right in order
	for i := 0; i < 8; i++ {
		var bit bool
		if i < start {
			bit = hasBit(ht.self.ID[index], uint(i))
		} else {
			bit = rand.Intn(2) == 1
		}
		if bit {
			first += byte(math.Pow(2, float64(7-i)))
		}
	}
	id = append(id, first)

	// randomize each remaining byte
	for i := index + 1; i < 20; i++ {
		id = append(id, byte(rand.Intn(256)))
	}

	return id
}

// hasBucketNode check if the node id is existed in the bucket
func (ht *HashTable) hasBucketNode(bucket int, id []byte) bool {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	for _, node := range ht.routeTable[bucket] {
		if bytes.Equal(node.ID, id) {
			return true
		}
	}
	return false
}

// hasNode check if the node id is exists in the hash table
func (ht *HashTable) hasNode(id []byte) bool {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	for _, bucket := range ht.routeTable {
		for _, node := range bucket {
			if bytes.Equal(node.ID, id) {
				return true
			}
		}
	}

	return false
}

// closestContacts returns the closest contacts of target
func (ht *HashTable) closestContacts(num int, target []byte, ignoredNodes []*Node) *NodeList {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	// find the bucket index in local route tables
	index := ht.bucketIndex(ht.self.ID, target)
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
		for i := 0; i < len(ht.routeTable[index]); i++ {
			node := ht.routeTable[index][i]

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

// bucketIndex return the bucket index from two node ids
func (ht *HashTable) bucketIndex(id1 []byte, id2 []byte) int {
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
func (ht *HashTable) totalCount() int {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	var num int
	for _, v := range ht.routeTable {
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
//
// number:  1
// bits:    00000001
// pos:     01234567
func hasBit(n byte, pos uint) bool {
	pos = 7 - pos
	val := n & (1 << pos)
	return (val > 0)
}
