package rqstore

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func TestStoreSymbols(t *testing.T) {
	store := SetupTestDB(t)
	defer store.db.Close()

	testCases := []struct {
		description string
		id          string
		symbols     map[string][]byte
		expectError bool
	}{
		{
			description: "Single symbol",
			id:          "tx1",
			symbols:     map[string][]byte{"hash1": []byte("data1")},
			expectError: false,
		},
		{
			description: "Multiple symbols",
			id:          "tx2",
			symbols: map[string][]byte{
				"hash2": []byte("data2"),
				"hash3": []byte("data3"),
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := store.StoreSymbols(tc.id, tc.symbols)

			if (err != nil) != tc.expectError {
				t.Errorf("Unexpected error status. Got %v, expected error: %v", err, tc.expectError)
			}
		})
	}
}

func TestLoadSymbols(t *testing.T) {
	store := SetupTestDB(t)
	defer store.db.Close()

	// Pre-insert some symbols to load
	preInsertSymbols := map[string][]byte{
		"hash1": []byte("data1"),
		"hash2": []byte("data2"),
	}
	store.StoreSymbols("tx1", preInsertSymbols)

	testCases := []struct {
		description string
		id          string
		symbols     map[string][]byte // Symbols to load
		expected    map[string][]byte // Expected result after loading
		expectError bool
	}{
		{
			description: "Load existing symbol",
			id:          "tx1",
			symbols:     map[string][]byte{"hash1": nil},
			expected:    map[string][]byte{"hash1": []byte("data1")},
			expectError: false,
		},
		{
			description: "Load multiple symbols",
			id:          "tx1",
			symbols:     map[string][]byte{"hash1": nil, "hash2": nil},
			expected:    preInsertSymbols,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			loadedSymbols, err := store.LoadSymbols(tc.id, tc.symbols)

			if (err != nil) != tc.expectError {
				t.Fatalf("Unexpected error status. Got %v, expected error: %v", err, tc.expectError)
			}

			if !reflect.DeepEqual(loadedSymbols, tc.expected) {
				t.Errorf("Loaded symbols do not match expected. Got %v, want %v", loadedSymbols, tc.expected)
			}
		})
	}
}
