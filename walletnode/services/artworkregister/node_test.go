package artworkregister

import (
	"context"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/walletnode/node/mocks"
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

//TODO: Implement sendImage test with mock
func TestNodesSendImage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		file *artwork.File
	}
	tests := []struct {
		nodes     *Nodes
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			tt.assertion(t, tt.nodes.sendImage(tt.args.ctx, tt.args.file))
		})
	}
}

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type fields struct {
		client          *mocks.Client
		conn            *mocks.Connection
		registerArtwork *mocks.RegisterArtwork
		Address         string
		PastelID        string
	}

	type args struct {
		ctx context.Context
		err error
	}

	type methods struct {
		registerArtwork string
		connect         string
	}

	type methodCall struct {
		connect         int
		registerArtWork int
	}

	type mockArgs struct {
		ctx     interface{}
		address interface{}
	}

	testCases := []struct {
		fields     fields
		args       args
		methods    methods
		methodCall methodCall
		mockArgs   mockArgs
		assertion  assert.ErrorAssertionFunc
	}{
		{
			fields: fields{
				conn:            &mocks.Connection{},
				client:          &mocks.Client{},
				registerArtwork: &mocks.RegisterArtwork{},
				Address:         "127.0.0.1:4444",
				PastelID:        "1",
			},
			methods: methods{
				registerArtwork: "RegisterArtwork",
				connect:         "Connect",
			},
			methodCall: methodCall{
				connect:         1,
				registerArtWork: 1,
			},
			mockArgs: mockArgs{mock.Anything, mock.AnythingOfType("string")},
			args: args{
				ctx: context.Background(),
				err: nil,
			},
			assertion: assert.NoError,
		},
		{
			fields: fields{
				conn:            &mocks.Connection{},
				client:          &mocks.Client{},
				registerArtwork: &mocks.RegisterArtwork{},
				Address:         "127.0.0.1:4445",
				PastelID:        "2",
			},
			methods: methods{
				registerArtwork: "RegisterArtwork",
				connect:         "Connect",
			},
			methodCall: methodCall{
				connect:         1,
				registerArtWork: 0,
			},
			mockArgs: mockArgs{mock.Anything, mock.AnythingOfType("string")},
			args: args{
				ctx: context.Background(),
				err: fmt.Errorf("connection timeout"),
			},
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.fields.conn.On(testCase.methods.registerArtwork).Return(testCase.fields.registerArtwork)
			testCase.fields.client.On(testCase.methods.connect, testCase.mockArgs.ctx, testCase.mockArgs.address).Return(testCase.fields.conn, testCase.args.err)

			node := &Node{
				client:   testCase.fields.client,
				Address:  testCase.fields.Address,
				PastelID: testCase.fields.PastelID,
			}

			testCase.assertion(t, node.connect(testCase.args.ctx))
			testCase.fields.client.AssertExpectations(t)
			testCase.fields.client.AssertCalled(t, testCase.methods.connect, testCase.mockArgs.ctx, testCase.fields.Address)
			testCase.fields.client.AssertNumberOfCalls(t, testCase.methods.connect, testCase.methodCall.connect)
			testCase.fields.conn.AssertNumberOfCalls(t, testCase.methods.registerArtwork, testCase.methodCall.registerArtWork)
		})
	}
}
