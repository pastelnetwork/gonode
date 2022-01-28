package senseregister

import (
	"context"
	"fmt"
	"testing"

	test "github.com/pastelnetwork/gonode/supernode/node/test/sense_register"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		node                    *SenseRegistrationNode
		address                 string
		args                    args
		err                     error
		numberConnectCall       int
		numberRegisterSenseCall int
		assertion               assert.ErrorAssertionFunc
	}{
		{
			node:                    &SenseRegistrationNode{Address: "127.0.0.1:4444"},
			address:                 "127.0.0.1:4444",
			args:                    args{context.Background()},
			err:                     nil,
			numberConnectCall:       1,
			numberRegisterSenseCall: 1,
			assertion:               assert.NoError,
		}, {
			node:                    &SenseRegistrationNode{Address: "127.0.0.1:4445"},
			address:                 "127.0.0.1:4445",
			args:                    args{context.Background()},
			err:                     fmt.Errorf("connection timeout"),
			numberConnectCall:       1,
			numberRegisterSenseCall: 0,
			assertion:               assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create client mocks
			clientMock := test.NewMockClient(t)

			//listen needed method
			clientMock.ListenOnConnect("", testCase.err).ListenOnRegisterSense()

			//set up node client only
			testCase.node.client = clientMock.Client

			//assertion error
			testCase.assertion(t, testCase.node.connect(testCase.args.ctx))
			//mock assertion
			clientMock.Client.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address)
			clientMock.AssertRegisterSenseCall(testCase.numberRegisterSenseCall)
		})
	}
}

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *SenseRegistrationNode
	}
	testCases := []struct {
		nodes SenseRegistrationNodes
		args  args
		want  SenseRegistrationNodes
	}{
		{
			nodes: SenseRegistrationNodes{},
			args:  args{node: &SenseRegistrationNode{Address: "127.0.0.1"}},
			want: SenseRegistrationNodes{
				&SenseRegistrationNode{Address: "127.0.0.1"},
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
		nodes SenseRegistrationNodes
		args  args
		want  *SenseRegistrationNode
	}{
		{
			nodes: SenseRegistrationNodes{
				&SenseRegistrationNode{ID: "1"},
				&SenseRegistrationNode{ID: "2"},
			},
			args: args{"2"},
			want: &SenseRegistrationNode{ID: "2"},
		}, {
			nodes: SenseRegistrationNodes{
				&SenseRegistrationNode{ID: "1"},
				&SenseRegistrationNode{ID: "2"},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Add(&SenseRegistrationNode{ID: "4"})
			testCase.nodes.Remove("4")
			assert.Equal(t, testCase.want, testCase.nodes.ByID(testCase.args.id))
		})
	}
}
