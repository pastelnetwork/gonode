package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
	service "github.com/pastelnetwork/gonode/walletnode/services/senseregister"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/sense_register"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterSenseNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	var TestCases = []struct {
		node              *common.SuperNodeClient
		address           string
		args              args
		err               error
		numberConnectCall int
		numberAPICall     int
		assertion         assert.ErrorAssertionFunc
	}{
		{
			node:              common.NewSuperNode(nil, "127.0.0.1:4444", "", nil),
			address:           "127.0.0.1:4444",
			args:              args{context.Background()},
			err:               nil,
			numberConnectCall: 1,
			numberAPICall:     1,
			assertion:         assert.NoError,
		}, {
			node:              common.NewSuperNode(nil, "127.0.0.1:4445", "", nil),
			address:           "127.0.0.1:4445",
			args:              args{context.Background()},
			err:               fmt.Errorf("connection timeout"),
			numberConnectCall: 1,
			numberAPICall:     0,
			assertion:         assert.Error,
		},
	}
	for i, testCase := range TestCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create client mocks
			clientMock := test.NewMockClient(t)

			//listen needed method
			clientMock.ListenOnConnect("", testCase.err).ListenOnRegisterSense()

			//set up node client only
			testCase.node.ClientInterface = clientMock.ClientInterface

			testCase.node.RealNodeMaker = service.RegisterSenseNodeMaker{}

			//assertion error
			testCase.assertion(t, testCase.node.Connect(testCase.args.ctx, time.Second, &alts.SecInfo{}, ""))
			//mock assertion
			clientMock.ClientInterface.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address, mock.Anything)
			clientMock.AssertRegisterSenseCall(testCase.numberAPICall)
		})
	}
}
