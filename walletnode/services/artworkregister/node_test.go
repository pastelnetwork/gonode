package artworkregister

import (
	"context"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		conn         *mocks.Connection
		client       *mocks.Client
		err          error
		address      string
		numberofCall int
	}{
		{
			conn:         &mocks.Connection{},
			client:       &mocks.Client{},
			err:          nil,
			address:      "127.0.0.1:4444",
			numberofCall: 1,
		},
		{
			conn:         &mocks.Connection{},
			client:       &mocks.Client{},
			err:          fmt.Errorf("connection timeout"),
			address:      "127.0.0.1:4445",
			numberofCall: 0,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.conn.On("RegisterArtwork").Return(&mocks.RegisterArtwork{})
			testCase.client.On("Connect", mock.Anything, mock.AnythingOfType("string")).Return(testCase.conn, testCase.err)

			node := &Node{
				client:  testCase.client,
				Address: testCase.address,
			}

			err := node.connect(context.Background())
			assert.Equal(t, testCase.err, err)

			//client assertion
			testCase.client.AssertExpectations(t)
			testCase.client.AssertCalled(t, "Connect", mock.Anything, testCase.address)
			testCase.client.AssertNumberOfCalls(t, "Connect", 1)

			testCase.conn.AssertNumberOfCalls(t, "RegisterArtwork", testCase.numberofCall)
		})
	}

}

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes        Nodes
		node         *Node
		expecteNodes Nodes
	}{
		{
			nodes: Nodes{},
			node: &Node{
				Address: "127.0.0.1",
			},
			expecteNodes: Nodes{
				&Node{
					Address: "127.0.0.1",
				},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.add(testCase.node)
			assert.Equal(t, testCase.expecteNodes, testCase.nodes)
		})
	}
}

func TestNodesActivate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes Nodes
	}{
		{
			nodes: Nodes{
				&Node{Address: "127.0.0.1"},
				&Node{Address: "127.0.0.2"},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.activate()

			for _, n := range testCase.nodes {
				assert.True(t, n.activated)
			}
		})
	}
}

func TestNodesDisconnectInactive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes Nodes
		conn  []struct {
			conn         *mocks.Connection
			activated    bool
			numberOfCall int
		}
	}{
		{
			nodes: Nodes{},
			conn: []struct {
				conn         *mocks.Connection
				activated    bool
				numberOfCall int
			}{
				{
					conn:         &mocks.Connection{},
					numberOfCall: 1,
					activated:    false,
				},
				{
					conn:         &mocks.Connection{},
					numberOfCall: 0,
					activated:    true,
				},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			for _, c := range testCase.conn {
				c.conn.On("Close").Return(nil)

				node := &Node{
					conn:      c.conn,
					activated: c.activated,
				}

				testCase.nodes = append(testCase.nodes, node)
			}

			testCase.nodes.disconnectInactive()

			for j, c := range testCase.conn {
				c := c

				t.Run(fmt.Sprintf("close-called-%d", j), func(t *testing.T) {
					c.conn.AssertNumberOfCalls(t, "Close", c.numberOfCall)
				})

			}

		})
	}

}

func TestNodesFindByPastelID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes        Nodes
		pastelId     string
		expectedNode *Node
	}{
		{
			nodes: Nodes{
				&Node{PastelID: "1"},
				&Node{PastelID: "2"},
			},
			pastelId:     "2",
			expectedNode: &Node{PastelID: "2"},
		}, {
			nodes: Nodes{
				&Node{PastelID: "1"},
				&Node{PastelID: "2"},
			},
			pastelId:     "3",
			expectedNode: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			node := testCase.nodes.findByPastelID(testCase.pastelId)
			assert.Equal(t, testCase.expectedNode, node)
		})
	}
}
