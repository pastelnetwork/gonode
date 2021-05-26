package kademlia

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeNetMsg(t *testing.T) {
	netMsgInit()
	var conn bytes.Buffer

	node := newNode(&NetworkNode{})
	id, _ := newID()
	node.ID = id
	node.Port = 3000
	node.IP = net.ParseIP("0.0.0.0")

	msg := &message{}
	msg.Type = messageTypeFindNode
	msg.Receiver = node.NetworkNode
	msg.Data = &queryDataFindNode{
		Target: id,
	}

	serialized, err := serializeMessage(msg)
	assert.NoError(t, err)

	conn.Write(serialized)

	deserialized, err := deserializeMessage(&conn)
	assert.NoError(t, err)
	assert.Equal(t, msg, deserialized)
}
