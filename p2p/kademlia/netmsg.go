package kademlia

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	messageTypePing = iota
	messageTypeStore
	messageTypeFindNode
	messageTypeFindValue
)

type message struct {
	Sender     *NetworkNode
	Receiver   *NetworkNode
	ID         int64
	Error      error
	Type       int
	IsResponse bool
	Data       interface{}
}

type queryDataFindNode struct {
	Target []byte
}

type queryDataFindValue struct {
	Target []byte
}

type queryDataStore struct {
	Data       []byte
	Publishing bool // Whether or not we are the original publisher
}

type responseDataFindNode struct {
	Closest []*NetworkNode
}

type responseDataFindValue struct {
	Closest []*NetworkNode
	Value   []byte
}

type responseDataStore struct {
	Success bool
}

func netMsgInit() {
	gob.Register(&queryDataFindNode{})
	gob.Register(&queryDataFindValue{})
	gob.Register(&queryDataStore{})
	gob.Register(&responseDataFindNode{})
	gob.Register(&responseDataFindValue{})
	gob.Register(&responseDataStore{})
}

func serializeMessage(q *message) ([]byte, error) {
	var msgBuffer bytes.Buffer
	enc := gob.NewEncoder(&msgBuffer)
	err := enc.Encode(q)
	if err != nil {
		return nil, errors.Errorf("failed to encode transmits: %w", err)
	}

	length := msgBuffer.Len()

	var lengthBytes [8]byte
	binary.PutUvarint(lengthBytes[:], uint64(length))

	var result []byte
	result = append(result, lengthBytes[:]...)
	result = append(result, msgBuffer.Bytes()...)

	return result, nil
}

func deserializeMessage(conn io.Reader) (*message, error) {
	lengthBytes := make([]byte, 8)
	if _, err := conn.Read(lengthBytes); err != nil {
		return nil, errors.Errorf("failed to read bytes: %w", err)
	}

	lengthReader := bytes.NewBuffer(lengthBytes)
	length, err := binary.ReadUvarint(lengthReader)
	if err != nil {
		return nil, errors.Errorf("failed to read an encoded buffer: %w", err)
	}

	msgBytes := make([]byte, length)
	if _, err := conn.Read(msgBytes); err != nil {
		return nil, errors.Errorf("failed to read bytes: %w", err)
	}

	reader := bytes.NewBuffer(msgBytes)
	msg := &message{}
	dec := gob.NewDecoder(reader)

	if err := dec.Decode(msg); err != nil {
		return nil, errors.Errorf("failed to decode message: %w", err)
	}
	return msg, nil
}
