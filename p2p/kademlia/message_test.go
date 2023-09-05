package kademlia

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagePayloadTooBig(t *testing.T) {
	// generate message
	data := make([]byte, 200*1024*1024+1)
	rand.Read(data)
	msg := &Message{
		Sender:      nil,
		Receiver:    nil, // the receiver node
		MessageType: 0,   // the message type
		Data: &StoreDataRequest{
			Data: data,
		},
	}

	// try to encode
	_, err := encode(msg)
	assert.NotNil(t, err)
	assert.Equal(t, "payload too big", err.Error())
}
