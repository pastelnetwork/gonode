package kademlia

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/zstd"
)

func compressKeys(keys [][]byte) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(keys); err != nil {
		return nil, fmt.Errorf("encode error: %w", err)
	}

	compressed, err := zstd.CompressLevel(nil, buf.Bytes(), 22)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	return compressed, nil
}

func decompressKeys(data []byte) ([][]byte, error) {
	decompressed, err := zstd.Decompress(nil, data)
	if err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(decompressed))

	var keys [][]byte
	if err := dec.Decode(&keys); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return keys, nil
}

// Function to generate a key from a node object.
func generateKeyFromNode(node *Node) string {
	return string(node.ID) + ":" + node.IP + ":" + fmt.Sprint(node.Port)
}

// Function to retrieve a node object from a key.
func getNodeFromKey(key string) (*Node, error) {
	parts := strings.Split(key, ":")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid key")
	}

	id := []byte(parts[0])
	ip := parts[1]

	// strconv.Atoi returns an int and an error, which we need to handle.
	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	return &Node{ID: id, IP: ip, Port: port}, nil
}

func compressKeysStr(keys []string) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(keys); err != nil {
		return nil, fmt.Errorf("encode error: %w", err)
	}

	compressed, err := zstd.CompressLevel(nil, buf.Bytes(), 22)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	return compressed, nil
}

func decompressKeysStr(data []byte) ([]string, error) {
	decompressed, err := zstd.Decompress(nil, data)
	if err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(decompressed))

	var keys []string
	if err := dec.Decode(&keys); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return keys, nil
}
