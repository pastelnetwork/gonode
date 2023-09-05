package kademlia

import (
	"fmt"
	"testing"
)

func TestGenerateKeyFromNode(t *testing.T) {
	tests := map[string]struct {
		input    *Node
		expected string
	}{
		"simple node": {
			input:    &Node{ID: []byte("123"), IP: "192.168.1.1", Port: 8080},
			expected: "123^192.168.1.1^8080",
		},
		"empty ID": {
			input:    &Node{ID: []byte(""), IP: "192.168.1.1", Port: 8080},
			expected: "^192.168.1.1^8080",
		},
		"IPv6 address": {
			input:    &Node{ID: []byte("456"), IP: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", Port: 8081},
			expected: "456^2001:0db8:85a3:0000:0000:8a2e:0370:7334^8081",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := generateKeyFromNode(tc.input)
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestGetNodeFromKey(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected *Node
		err      error
	}{
		"valid key": {
			input:    "123^192.168.1.1^8080",
			expected: &Node{ID: []byte("123"), IP: "192.168.1.1", Port: 8080},
			err:      nil,
		},
		"missing port": {
			input: "123^192.168.1.1",
			err:   fmt.Errorf("invalid key"),
		},
		"IPv6 address": {
			input:    "456^2001:0db8:85a3:0000:0000:8a2e:0370:7334^8081",
			expected: &Node{ID: []byte("456"), IP: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", Port: 8081},
			err:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := getNodeFromKey(tc.input)

			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}

			if string(got.ID) != string(tc.expected.ID) || got.IP != tc.expected.IP || got.Port != tc.expected.Port {
				t.Errorf("expected node %+v, got %+v", tc.expected, got)
			}
		})
	}
}
