package artworkregister

import (
	"context"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/supernode/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *Node
	}

	testCases := []struct {
		name  string
		nodes *Nodes
		args  args
		want  *Nodes
	}{
		{
			name:  "add node",
			nodes: &Nodes{},
			args:  args{&Node{ID: "1", Address: "127.0.0.1"}},
			want:  &Nodes{&Node{ID: "1", Address: "127.0.0.1"}},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			testCase.nodes.Add(testCase.args.node)
			assert.Equal(t, testCase.want, testCase.nodes)
		})
	}
}

func TestNodesByID(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}

	testCases := []struct {
		name  string
		nodes Nodes
		args  args
		want  *Node
	}{
		{
			name:  "return node with id 1",
			nodes: Nodes{&Node{ID: "2"}, &Node{ID: "1"}},
			args:  args{"1"},
			want:  &Node{ID: "1"},
		}, {
			name:  "return nil node",
			nodes: Nodes{&Node{ID: "2"}, &Node{ID: "1"}},
			args:  args{"3"},
			want:  nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, testCase.nodes.ByID(testCase.args.id))
		})
	}
}

func TestNodesRemove(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}

	testCases := []struct {
		name  string
		nodes *Nodes
		args  args
		want  *Nodes
	}{
		{
			name:  "remove node with id 2",
			nodes: &Nodes{&Node{ID: "1"}, &Node{ID: "2"}},
			args:  args{"2"},
			want:  &Nodes{&Node{ID: "1"}},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			testCase.nodes.Remove(testCase.args.id)
			assert.Equal(t, testCase.want, testCase.nodes)
		})
	}
}

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type fields struct {
		ID      string
		Address string
	}
	type args struct {
		ctx       context.Context
		returnErr error
	}

	testCases := []struct {
		name                   string
		fields                 fields
		args                   args
		assertion              assert.ErrorAssertionFunc
		numConnectCall         int
		numArtworkRegisterCall int
	}{
		{
			name:                   "connected to supernode",
			fields:                 fields{"1", "127.0.0.1"},
			args:                   args{context.Background(), nil},
			assertion:              assert.NoError,
			numConnectCall:         1,
			numArtworkRegisterCall: 1,
		}, {
			name:                   "failed connect to supernode",
			fields:                 fields{"2", "127.0.0.2"},
			args:                   args{context.Background(), fmt.Errorf("failed connect to supernode")},
			assertion:              assert.Error,
			numConnectCall:         1,
			numArtworkRegisterCall: 0,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			client := test.NewMockClient(t)
			client.ListenOnConnectCall(testCase.args.returnErr).
				ListenOnRegisterArtworkCall()

			node := &Node{
				ID:      testCase.fields.ID,
				Address: testCase.fields.Address,
				client:  client.Client,
			}

			testCase.assertion(t, node.connect(testCase.args.ctx))

			client.AssertConnectCall(testCase.numConnectCall, mock.Anything, testCase.fields.Address)
			client.AssertRegisterArtworkCall(testCase.numArtworkRegisterCall)
		})
	}
}
