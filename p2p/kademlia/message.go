package kademlia

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

const (
	// Ping the target to check if it's online
	Ping = iota
	// StoreData iterative store data in kademlia network
	StoreData
	// FindNode iterative find node in kademlia network
	FindNode
	// FindValue iterative find value in kademlia network
	FindValue
)

func init() {
	gob.Register(&FindNodeRequest{})
	gob.Register(&FindValueRequest{})
	gob.Register(&StoreDataRequest{})
	gob.Register(&FindNodeResponse{})
	gob.Register(&FindValueResponse{})
}

// Message structure for kademlia network
type Message struct {
	Sender      *Node       // the sender node
	Receiver    *Node       // the receiver node
	MessageType int         // the message type
	Data        interface{} // the real data for the request
}

func (m *Message) String() string {
	return fmt.Sprintf("type: %v, sender: %v, receiver: %v, data type: %T", m.MessageType, m.Sender.String(), m.Receiver.String(), m.Data)
}

// FindNodeRequest defines the request data for find node
type FindNodeRequest struct {
	Target []byte
}

// FindValueRequest defines the request data for find value
type FindValueRequest struct {
	Target []byte
}

// StoreDataRequest defines the request data for store data
type StoreDataRequest struct {
	Data []byte
}

// FindNodeResponse defines the response data for find node
type FindNodeResponse struct {
	Closest []*Node
}

// FindValueResponse defines the response data for find value
type FindValueResponse struct {
	Closest []*Node
	Value   []byte
}

// encode the message
func encode(message *Message) ([]byte, error) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	// encode the message with gob library
	if err := encoder.Encode(message); err != nil {
		return nil, err
	}

	var header [8]byte
	// prepare the header
	binary.PutUvarint(header[:], uint64(buf.Len()))

	var data []byte
	data = append(data, header[:]...)
	data = append(data, buf.Bytes()...)

	return data, nil
}

// decode the message
func decode(conn io.Reader) (*Message, error) {
	// read the header
	header := make([]byte, 8)
	if _, err := conn.Read(header); err != nil {
		return nil, err
	}

	// parse the length of message
	length, err := binary.ReadUvarint(bytes.NewBuffer(header))
	if err != nil {
		return nil, err
	}

	// read the message body
	data := make([]byte, length)
	if _, err := conn.Read(data); err != nil {
		return nil, err
	}

	// new a decoder
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	// decode the message structure
	message := &Message{}
	if err = decoder.Decode(message); err != nil {
		return nil, err
	}

	return message, nil
}
