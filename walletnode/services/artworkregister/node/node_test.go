package node

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/artwork_register"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		node                      *NftRegisterNodeClient
		address                   string
		args                      args
		err                       error
		numberConnectCall         int
		numberRegisterArtWorkCall int
		assertion                 assert.ErrorAssertionFunc
	}{
		{
			node:                      &NftRegisterNodeClient{address: "127.0.0.1:4444", mtx: &sync.RWMutex{}},
			address:                   "127.0.0.1:4444",
			args:                      args{context.Background()},
			err:                       nil,
			numberConnectCall:         1,
			numberRegisterArtWorkCall: 1,
			assertion:                 assert.NoError,
		}, {
			node:                      &NftRegisterNodeClient{address: "127.0.0.1:4445", mtx: &sync.RWMutex{}},
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
			clientMock.ListenOnConnect("", testCase.err).ListenOnRegisterArtwork()

			//set up node client only
			testCase.node.ClientInterface = clientMock.Client

			//assertion error
			testCase.assertion(t, testCase.node.Connect(testCase.args.ctx, time.Second, &alts.SecInfo{}))
			//mock assertion
			clientMock.Client.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address, mock.Anything)
			clientMock.AssertRegisterArtworkCall(testCase.numberRegisterArtWorkCall)
		})
	}
}
