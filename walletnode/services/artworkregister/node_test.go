package artworkregister

import (
	"context"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *Node
	}
	testCases := []struct {
		nodes Nodes
		args  args
		want  Nodes
	}{
		{
			nodes: Nodes{},
			args:  args{node: &Node{Address: "127.0.0.1"}},
			want: Nodes{
				&Node{Address: "127.0.0.1"},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.add(testCase.args.node)
			assert.Equal(t, testCase.want, testCase.nodes)
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

	type args struct {
		id string
	}
	testCases := []struct {
		nodes Nodes
		args  args
		want  *Node
	}{
		{
			nodes: Nodes{
				&Node{PastelID: "1"},
				&Node{PastelID: "2"},
			},
			args: args{"2"},
			want: &Node{PastelID: "2"},
		}, {
			nodes: Nodes{
				&Node{PastelID: "1"},
				&Node{PastelID: "2"},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, testCase.nodes.findByPastelID(testCase.args.id))
		})
	}
}

func TestNodesSendImage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		file *artwork.File
	}

	type nodeAttribute struct {
		address   string
		returnErr error
	}

	testCases := []struct {
		nodes                 []nodeAttribute
		args                  args
		err                   error
		numberUploadImageCall int
	}{
		{
			nodes:                 []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", nil}},
			args:                  args{context.Background(), &artwork.File{}},
			err:                   nil,
			numberUploadImageCall: 1,
		},
		{
			nodes:                 []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", fmt.Errorf("failed to open stream")}},
			args:                  args{context.Background(), &artwork.File{}},
			err:                   fmt.Errorf("failed to open stream"),
			numberUploadImageCall: 1,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			nodes := Nodes{}
			clients := []*test.Client{}

			for _, a := range testCase.nodes {
				//client mock
				client := test.NewMockClient()
				//listen on uploadImage call
				client.ListenOnUploadImage(testCase.err)
				clients = append(clients, client)

				nodes.add(&Node{
					Address:         a.address,
					RegisterArtwork: client.RegArtWorkMock,
				})
			}

			err := nodes.sendImage(testCase.args.ctx, testCase.args.file)
			assert.Equal(t, testCase.err, err)

			//mock assertion each client
			for _, client := range clients {
				client.RegArtWorkMock.AssertExpectations(t)
				client.RegArtWorkMock.AssertCalled(t, "UploadImage", testCase.args.ctx, testCase.args.file)
				client.RegArtWorkMock.AssertNumberOfCalls(t, "UploadImage", testCase.numberUploadImageCall)
			}
		})
	}
}

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		node                      *Node
		address                   string
		args                      args
		err                       error
		numberConnectCall         int
		numberRegisterArtWorkCall int
		assertion                 assert.ErrorAssertionFunc
	}{
		{
			node:                      &Node{Address: "127.0.0.1:4444"},
			address:                   "127.0.0.1:4444",
			args:                      args{context.Background()},
			err:                       nil,
			numberConnectCall:         1,
			numberRegisterArtWorkCall: 1,
			assertion:                 assert.NoError,
		}, {
			node:                      &Node{Address: "127.0.0.1:4445"},
			address:                   "127.0.0.1:4445",
			args:                      args{context.Background()},
			err:                       fmt.Errorf("connection timeout"),
			numberConnectCall:         1,
			numberRegisterArtWorkCall: 0,
			assertion:                 assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create client mocks
			client := test.NewMockClient()

			//listen needed method
			client.ListenOnConnect(testCase.err).ListenOnRegisterArtwork()

			//set up node client only
			testCase.node.client = client.ClientMock

			//assertion error
			testCase.assertion(t, testCase.node.connect(testCase.args.ctx))
			//mock assertion
			client.ClientMock.AssertExpectations(t)
			client.ClientMock.AssertCalled(t, "Connect", mock.Anything, testCase.address)
			client.ClientMock.AssertNumberOfCalls(t, "Connect", testCase.numberConnectCall)
			client.ConnectionMock.AssertNumberOfCalls(t, "RegisterArtwork", testCase.numberRegisterArtWorkCall)

		})
	}
}
