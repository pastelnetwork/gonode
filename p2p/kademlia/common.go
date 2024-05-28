package kademlia

import (
	"fmt"
	"strconv"
	"strings"
)

// generateKeyFromNode to generate a key from a node object.
func generateKeyFromNode(node *Node) string {
	var builder strings.Builder
	builder.Grow(len(node.ID) + 1 + len(node.IP) + 1 + 5) // Assuming a maximum of 10 bytes for the port
	builder.WriteString(string(node.ID))
	builder.WriteByte('^')
	builder.WriteString(node.IP)
	builder.WriteByte('^')
	builder.WriteString(strconv.Itoa(node.Port))

	return builder.String()
}

// Function to retrieve a node object from a key.
func getNodeFromKey(key string) (*Node, error) {
	parts := strings.Split(key, "^")

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

/*
func compressKeysStr(keys []string) ([]byte, error) {
	buf, err := msgpack.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("msgpack encode error: %w", err)
	}

	compressed, err := utils.Compress(buf, 2)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	return compressed, nil
}

func decompressKeysStr(data []byte) ([]string, error) {
	decompressed, err := utils.Decompress(data)
	if err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	var keys []string
	if err := msgpack.Unmarshal(decompressed, &keys); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return keys, nil
}

/*
func compressSymbols(values [][]byte) ([]byte, error) {
	buf, err := msgpack.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("msgpack encode error: %w", err)
	}

	compressed, err := utils.Compress(buf, 2)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	return compressed, nil
}

func decompressSymbols(data []byte) (values [][]byte, err error) {
	decompressed, err := utils.Decompress(data)
	if err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	if err := msgpack.Unmarshal(decompressed, &values); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return values, nil
}
*/
