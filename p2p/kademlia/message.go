package kademlia

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"
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
	// Replicate is replicate data request towards a network node
	Replicate
	// BatchFindValues finds values in kademlia network
	BatchFindValues
	// StoreBatchData iterative stores data in batches
	BatchStoreData
	// BatchFindNode iterative find node in kademlia network
	BatchFindNode
	// BatchGetValues finds values in kademlia network
	BatchGetValues
)

func init() {
	gob.Register(&ResponseStatus{})
	gob.Register(&FindNodeRequest{})
	gob.Register(&FindNodeResponse{})
	gob.Register(&FindValueRequest{})
	gob.Register(&FindValueResponse{})
	gob.Register(&StoreDataRequest{})
	gob.Register(&StoreDataResponse{})
	gob.Register(&ReplicateDataRequest{})
	gob.Register(&ReplicateDataResponse{})
	gob.Register(&BatchFindValuesRequest{})
	gob.Register(&BatchFindValuesResponse{})
	gob.Register(&BatchStoreDataRequest{})
	gob.Register(&BatchFindNodeRequest{})
	gob.Register(&BatchFindNodeResponse{})
	gob.Register(&BatchGetValuesRequest{})
	gob.Register(&BatchGetValuesResponse{})
}

type MessageWithError struct {
	Message *Message
	Error   error
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

// ResultType specify success of message request
type ResultType int

const (
	// ResultOk means request is ok
	ResultOk ResultType = 0
	// ResultFailed meas request got failed
	ResultFailed ResultType = 1
)

// ReplicateDataRequest ...
type ReplicateDataRequest struct {
	Keys []string
}

// ReplicateDataResponse ...
type ReplicateDataResponse struct {
	Status ResponseStatus
}

// ResponseStatus defines the result of request
type ResponseStatus struct {
	Result ResultType
	ErrMsg string
}

// FindNodeRequest defines the request data for find node
type FindNodeRequest struct {
	Target []byte
}

// FindNodeResponse defines the response data for find node
type FindNodeResponse struct {
	Status  ResponseStatus
	Closest []*Node
}

// FindValueRequest defines the request data for find value
type FindValueRequest struct {
	Target []byte
}

// FindValueResponse defines the response data for find value
type FindValueResponse struct {
	Status  ResponseStatus
	Closest []*Node
	Value   []byte
}

// StoreDataRequest defines the request data for store data
type StoreDataRequest struct {
	Data []byte
	Type int
}

// StoreDataResponse defines the response data for store data
type StoreDataResponse struct {
	Status ResponseStatus
}

// BatchFindValuesRequest defines the request data for find value
type BatchFindValuesRequest struct {
	Keys []string
}

// BatchFindValuesResponse defines the response data for find value
type BatchFindValuesResponse struct {
	Status   ResponseStatus
	Response []byte
	Done     bool
}

type KeyValWithClosest struct {
	Value   []byte
	Closest []*Node
}

// BatchGetValuesRequest defines the request data for find value
type BatchGetValuesRequest struct {
	Data map[string]KeyValWithClosest //(compressed) // keys are hex encoded
	//Data []byte
}

type BatchGetValuesResponse struct {
	Data map[string]KeyValWithClosest // keys are hex encoded
	//Data   []byte
	Status ResponseStatus
}

// encode the message
func encode(message *Message) ([]byte, error) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	// encode the message with gob library
	if err := encoder.Encode(message); err != nil {
		return nil, err
	}

	if utils.BytesIntToMB(buf.Len()) > defaultMaxPayloadSize {
		return nil, errors.New("payload too big")
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
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	// parse the length of message
	length, err := binary.ReadUvarint(bytes.NewBuffer(header))
	if err != nil {
		return nil, errors.Errorf("parse header length: %w", err)
	}

	if utils.BytesToMB(length) > defaultMaxPayloadSize {
		return nil, errors.New("payload too big")
	}

	// read the message body
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
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

// BatchStoreDataRequest defines the request data for store data
type BatchStoreDataRequest struct {
	Data [][]byte
	Type int
}

// BatchFindNodeRequest defines the request data for find node
type BatchFindNodeRequest struct {
	HashedTarget [][]byte
}

// BatchFindNodeResponse defines the response data for find node
type BatchFindNodeResponse struct {
	Status       ResponseStatus
	ClosestNodes map[string][]*Node
}
