package kademlia

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
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
	Alpha = 6

	// B - the size in bits of the keys used to identify nodes and store and
	// retrieve data; in basic Kademlia this is 256, the length of a SHA3-256
	B = 256

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
	routeTable [][]*Node // 256x20

	// mutex for route table
	mutex sync.RWMutex

	// refresh time for every bucket
	refreshers []time.Time
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
		return nil, errors.New("id is nil")
	}
	ht.self.SetHashedID()

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
	// bucket index of the node
	index := ht.bucketIndex(ht.self.HashedID, id)
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
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

	return ht.refreshers[bucket]
}

// randomIDFromBucket returns a random id based on the bucket index and self id
func (ht *HashTable) randomIDFromBucket(bucket int) []byte {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

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
			nBig, _ := rand.Int(rand.Reader, big.NewInt(2))
			n := nBig.Int64()

			bit = n == 1
		}
		if bit {
			first += byte(math.Pow(2, float64(7-i)))
		}
	}
	id = append(id, first)

	// randomize each remaining byte
	for i := index + 1; i < 20; i++ {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(256))
		n := nBig.Int64()

		id = append(id, byte(n))
	}

	return id
}

// hasBucketNode check if the node id is existed in the bucket
func (ht *HashTable) hasBucketNode(bucket int, id []byte) bool {
	for _, node := range ht.routeTable[bucket] {
		if bytes.Equal(node.ID, id) {
			return true
		}
	}

	return false
}

// hasNode check if the node id is exists in the hash table
func (ht *HashTable) hasNode(id []byte) bool {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

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
func (ht *HashTable) closestContacts(num int, target []byte, ignoredNodes []*Node) (*NodeList, int) {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

	// Convert ignoredNodes slice to a map for faster lookup
	ignoredMap := make(map[string]bool)
	for _, node := range ignoredNodes {
		ignoredMap[string(node.ID)] = true
	}

	nl := &NodeList{
		Comparator: target,
	}

	counter := 0
	// Flatten the routeTable and add nodes to nl if they're not in the ignoredMap
	for _, bucket := range ht.routeTable {
		for _, node := range bucket {
			counter++
			if !ignoredMap[string(node.ID)] {
				nl.AddNodes([]*Node{node})
			}
		}
	}

	// Sort the node list and get the top 'num' nodes
	nl.Sort()
	nl.TopN(num)

	return nl, counter
}

// bucketIndex return the bucket index from two node ids
func (*HashTable) bucketIndex(id1 []byte, id2 []byte) int {
	// look at each byte from left to right
	for j := 0; j < len(id1); j++ {
		// xor the byte
		xor := id1[j] ^ id2[j]

		// check each bit on the xored result from left to right in order
		for i := 0; i < 8; i++ {
			if hasBit(xor, uint(i)) {
				byteIndex := j * 8 // Convert byte position to bit position
				bitIndex := i
				// Return bucket index based on position of differing bit (B - total bit position - 1)
				return B - (byteIndex + bitIndex) - 1
			}
		}
	}

	// only happen during bootstrap ping
	return 0
}

// Count returns the number of nodes in route table
func (ht *HashTable) totalCount() int {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

	var num int
	for _, v := range ht.routeTable {
		num += len(v)
	}
	return num
}

// nodes returns nodes in table
func (ht *HashTable) nodes() []*Node {
	nodeList := []*Node{}
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

	for _, v := range ht.routeTable {
		nodeList = append(nodeList, v...)
	}
	return nodeList
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
	val := n & (1 << (7 - pos)) // check bits from left to right (7 - pos)
	return (val > 0)
}

func (ht *HashTable) closestContactsWithInlcudingNode(num int, target []byte, ignoredNodes []*Node, includeNode *Node) *NodeList {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()

	// Convert ignoredNodes slice to a map for faster lookup
	ignoredMap := make(map[string]bool)
	for _, node := range ignoredNodes {
		ignoredMap[hex.EncodeToString(node.ID)] = true
	}

	nl := &NodeList{
		Comparator: target,
	}

	// Flatten the routeTable and add nodes to nl if they're not in the ignoredMap
	counter := 0
	for _, bucket := range ht.routeTable {
		for _, node := range bucket {
			if !ignoredMap[string(node.ID)] {
				counter++
				nl.AddNodes([]*Node{node})
			}
		}
	}

	// Add the included node (if any) to the list
	if includeNode != nil {
		nl.AddNodes([]*Node{includeNode})
	}

	// Sort the node list and get the top 'num' nodes
	nl.Sort()
	nl.TopN(num)

	return nl
}

// RemoveNode removes node from bucket
func (ht *HashTable) RemoveNode(index int, id []byte) bool {
	ht.mutex.Lock()
	defer ht.mutex.Unlock()

	bucket := ht.routeTable[index]
	for i, node := range bucket {
		if bytes.Equal(node.ID, id) {
			if i+1 < len(bucket) {
				bucket = append(bucket[:i], bucket[i+1:]...)
			} else {
				bucket = bucket[:i]
			}
			ht.routeTable[index] = bucket
			return true
		}
	}

	return false
}
