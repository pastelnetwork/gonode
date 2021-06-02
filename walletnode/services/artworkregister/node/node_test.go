package node

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/walletnode/node/test"
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
			node:                      &Node{address: "127.0.0.1:4444"},
			address:                   "127.0.0.1:4444",
			args:                      args{context.Background()},
			err:                       nil,
			numberConnectCall:         1,
			numberRegisterArtWorkCall: 1,
			assertion:                 assert.NoError,
		}, {
			node:                      &Node{address: "127.0.0.1:4445"},
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
			testCase.node.Client = client.ClientMock

			//assertion error
			testCase.assertion(t, testCase.node.Connect(testCase.args.ctx, time.Second))
			//mock assertion
			client.ClientMock.AssertExpectations(t)
			client.ClientMock.AssertCalled(t, "Connect", mock.Anything, testCase.address)
			client.ClientMock.AssertNumberOfCalls(t, "Connect", testCase.numberConnectCall)
			client.ConnectionMock.AssertNumberOfCalls(t, "RegisterArtwork", testCase.numberRegisterArtWorkCall)

		})
	}
}
