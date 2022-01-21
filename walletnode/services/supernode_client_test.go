package services

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"

	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	testSense "github.com/pastelnetwork/gonode/walletnode/node/test/register_sense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		node                      *common.SuperNodeClient
		address                   string
		args                      args
		err                       error
		numberConnectCall         int
		numberRegisterArtWorkCall int
		assertion                 assert.ErrorAssertionFunc
	}{
		{
			node:                      common.NewSuperNode(nil, "127.0.0.1:4444", "", nil),
			address:                   "127.0.0.1:4444",
			args:                      args{context.Background()},
			err:                       nil,
			numberConnectCall:         1,
			numberRegisterArtWorkCall: 1,
			assertion:                 assert.NoError,
		}, {
			node:                      common.NewSuperNode(nil, "127.0.0.1:4445", "", nil),
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
			senseClientMock := testSense.NewMockClient(t)

			//listen needed method
			senseClientMock.ListenOnConnect("", testCase.err).ListenOnRegisterSense()

			//set up node client only
			testCase.node.ClientInterface = senseClientMock.Client

			makers := []node.NodeMaker{senseregister.SenseRegisterNodeMaker{}}

			for _, m := range makers {

				testCase.node.NodeMaker = m

				//assertion error
				testCase.assertion(t, testCase.node.Connect(testCase.args.ctx, time.Second, &alts.SecInfo{}))
				//mock assertion
				senseClientMock.Client.AssertExpectations(t)
				senseClientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address, mock.Anything)
				senseClientMock.AssertRegisterSenseCall(testCase.numberRegisterArtWorkCall)
			}
		})
	}
}
