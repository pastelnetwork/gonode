package externaldupedetection

import (
	"context"
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/supernode/node/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
			clientMock := test.NewMockClient(t)

			//listen needed method
			clientMock.ListenOnConnect("", testCase.err).ListenOnExternalDupeDetection()

			//set up node client only
			testCase.node.client = clientMock.Client

			//assertion error
			testCase.assertion(t, testCase.node.connect(testCase.args.ctx))
			//mock assertion
			clientMock.Client.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address)
			clientMock.AssertExternalDupeDetectionCall(testCase.numberRegisterArtWorkCall)
		})
	}
}

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

			testCase.nodes.Add(testCase.args.node)
			assert.Equal(t, testCase.want, testCase.nodes)
		})
	}
}

func TestByID(t *testing.T) {
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
				&Node{ID: "1"},
				&Node{ID: "2"},
			},
			args: args{"2"},
			want: &Node{ID: "2"},
		}, {
			nodes: Nodes{
				&Node{ID: "1"},
				&Node{ID: "2"},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Add(&Node{ID: "4"})
			testCase.nodes.Remove("4")
			assert.Equal(t, testCase.want, testCase.nodes.ByID(testCase.args.id))
		})
	}
}
