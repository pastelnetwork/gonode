package kademlia

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestBanList_Add(t *testing.T) {
	// Create a BanList
	ctx := context.Background()
	banList := NewBanList(ctx)

	// Create a Node
	node := &Node{
		ID:   []byte{1, 2, 3},
		IP:   "127.0.0.1",
		Port: 12345,
	}

	// Add the Node to the BanList
	banList.Add(node)

	// Check if the Node was added to the BanList
	found := false
	for _, item := range banList.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			found = true
			break
		}
	}

	if !found {
		t.Error("Failed to add Node to BanList")
	}
}

func TestBanList_IncrementCount(t *testing.T) {
	// Create a BanList
	ctx := context.Background()
	banList := NewBanList(ctx)

	// Create a Node
	node := &Node{
		ID:   []byte{1, 2, 3},
		IP:   "127.0.0.1",
		Port: 12345,
	}

	// Add the Node to the BanList
	banList.Add(node)

	// Increment the count for the Node in the BanList
	banList.IncrementCount(node)

	// Check if the count was incremented
	found := false
	for _, item := range banList.Nodes {
		if bytes.Equal(item.ID, node.ID) {
			if item.count != 2 { // Expected count after increment is 2
				t.Errorf("Failed to increment count for Node in BanList, expected: 2, got: %d", item.count)
			}
			found = true
			break
		}
	}

	if !found {
		t.Error("Failed to find Node in BanList")
	}
}

func TestBanList_Exists(t *testing.T) {
	// Create a BanList
	ctx := context.Background()
	banList := NewBanList(ctx)

	// Create a Node
	node := &Node{
		ID:   []byte{1, 2, 3},
		IP:   "127.0.0.1",
		Port: 12345,
	}

	// Add the Node to the BanList
	banList.AddWithCreatedAt(node, time.Now().UTC(), 6)

	// Check if the Node exists in the BanList
	exists := banList.Exists(node)
	if !exists {
		t.Error("Failed to detect existing Node in BanList")
	}

	// Create a Node that doesn't exist in the BanList
	nonExistingNode := &Node{
		ID:   []byte{4, 5, 6},
		IP:   "192.168.0.1",
		Port: 54321,
	}

	// Check if the non-existing Node doesn't exist in the BanList
	exists = banList.Exists(nonExistingNode)
	if exists {
		t.Error("Failed to detect non-existing Node in BanList")
	}
}

func TestBanList_Delete(t *testing.T) {
	ctx := context.Background()
	banList := NewBanList(ctx)

	// Add some nodes to the BanList
	node1 := &Node{ID: []byte{1}, IP: "127.0.0.1", Port: 8000}
	node2 := &Node{ID: []byte{2}, IP: "127.0.0.2", Port: 8001}
	node3 := &Node{ID: []byte{3}, IP: "127.0.0.3", Port: 8002}
	banList.AddWithCreatedAt(node1, time.Now().UTC(), 6)
	banList.AddWithCreatedAt(node2, time.Now().UTC(), 6)
	banList.AddWithCreatedAt(node3, time.Now().UTC(), 6)

	// Test deleting a node that exists in the BanList
	banList.Delete(node2)

	// Check if the deleted node is removed from the BanList
	if banList.Exists(node2) {
		t.Errorf("Expected node2 to be deleted from BanList, but it still exists")
	}

	// Check if the BanList still contains the other nodes
	if !banList.Exists(node1) {
		t.Errorf("Expected node1 to exist in BanList, but it doesn't")
	}
	if !banList.Exists(node3) {
		t.Errorf("Expected node3 to exist in BanList, but it doesn't")
	}

	// Test deleting a node that doesn't exist in the BanList
	node4 := &Node{ID: []byte{4}, IP: "127.0.0.4", Port: 8003}
	banList.Delete(node4)

	// Check if the BanList remains unchanged
	if !banList.Exists(node1) {
		t.Errorf("Expected node1 to exist in BanList, but it doesn't")
	}
	if !banList.Exists(node3) {
		t.Errorf("Expected node3 to exist in BanList, but it doesn't %d", len(banList.Nodes))
	}
}

func TestBanList_Purge(t *testing.T) {
	ctx := context.Background()
	banList := NewBanList(ctx)

	// Add some nodes to the BanList
	node1 := &Node{ID: []byte{1}, IP: "127.0.0.1", Port: 8000}
	node2 := &Node{ID: []byte{2}, IP: "127.0.0.2", Port: 8001}
	node3 := &Node{ID: []byte{3}, IP: "127.0.0.3", Port: 8002}
	banList.AddWithCreatedAt(node1, time.Now().UTC().Add(-2*time.Hour), 6)
	banList.AddWithCreatedAt(node2, time.Now().UTC().Add(-6*time.Hour), 6)
	banList.AddWithCreatedAt(node3, time.Now().UTC().Add(-1*time.Hour), 16)

	// Call Purge to remove expired nodes
	banList.Purge()

	// Check if the expired node is removed from the BanList
	if banList.Exists(node2) {
		t.Errorf("Expected node2 to be removed from BanList due to expiration, but it still exists")
	}

	// Check if the non-expired nodes still exist in the BanList
	if !banList.Exists(node1) {
		t.Errorf("Expected node1 to exist in BanList, but it doesn't")
	}
	if !banList.Exists(node3) {
		t.Errorf("Expected node3 to exist in BanList, but it doesn't")
	}
}
