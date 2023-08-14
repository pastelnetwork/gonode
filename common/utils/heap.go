package utils

import (
	"math/big"
)

// DistanceEntry struct to hold XOR distance
type DistanceEntry struct {
	distance *big.Int
	str      string
}

//MaxHeap to efficiently sort XOR distances
type MaxHeap []DistanceEntry

//Len returns the length of heap
func (h MaxHeap) Len() int {
	return len(h)
}

// Less returns true if i > j
func (h MaxHeap) Less(i, j int) bool {
	return h[i].distance.Cmp(h[j].distance) > 0
}

// Swap is used to swap the values
func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push is used to push the results
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(DistanceEntry))
}

// Pop is used to pop the results
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
