package node

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/walletnode/node"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/artwork_download"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeNewNode(t *testing.T) {
	t.Parallel()

	type args struct {
		client   node.ClientInterface
		address  string
		pastelID string
	}
	testCases := []struct {
		args args
		want *NftDownloadNodeClient
	}{
		{
			args: args{nil, "127.0.0.1:4444", "testID"},
			want: &NftDownloadNodeClient{
				ClientInterface: nil,
				address:         "127.0.0.1:4444",
				pastelID:        "testID",
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			node := NewNode(testCase.args.client, testCase.args.address, testCase.args.pastelID)
			assert.Equal(t, *testCase.want, *node)
		})
	}
}

func TestNodeString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		node *NftDownloadNodeClient
		want string
	}{
		{
			node: &NftDownloadNodeClient{address: "127.0.0.1:4444"},
			want: "127.0.0.1:4444",
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			name := testCase.node.String()
			assert.Equal(t, testCase.want, name)
		})
	}
}

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		node                      *NftDownloadNodeClient
		address                   string
		args                      args
		err                       error
		numberConnectCall         int
		numberDownloadArtWorkCall int
		assertion                 assert.ErrorAssertionFunc
	}{
		{
			node:                      &NftDownloadNodeClient{address: "127.0.0.1:4444"},
			address:                   "127.0.0.1:4444",
			args:                      args{context.Background()},
			err:                       nil,
			numberConnectCall:         1,
			numberDownloadArtWorkCall: 1,
			assertion:                 assert.NoError,
		}, {
			node:                      &NftDownloadNodeClient{address: "127.0.0.1:4445"},
			address:                   "127.0.0.1:4445",
			args:                      args{context.Background()},
			err:                       fmt.Errorf("connection timeout"),
			numberConnectCall:         1,
			numberDownloadArtWorkCall: 0,
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
			clientMock.ListenOnConnect(testCase.address, testCase.err).ListenOnDownloadArtwork()

			//set up node client only
			testCase.node.ClientInterface = clientMock.Client

			//assertion error
			testCase.assertion(t, testCase.node.Connect(testCase.args.ctx, time.Second, &alts.SecInfo{}))
			//mock assertion
			clientMock.Client.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address, mock.Anything)
			clientMock.AssertDownloadArtworkCall(testCase.numberDownloadArtWorkCall)
		})
	}
}
