package nftregister

import (
	"context"
	"fmt"
	"testing"

	test "github.com/pastelnetwork/gonode/supernode/node/test/nft_register"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNftNodeConnect(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		node                  *NftRegistrationNode
		address               string
		args                  args
		err                   error
		numberConnectCall     int
		numberRegisterNftCall int
		assertion             assert.ErrorAssertionFunc
	}{
		{
			node:                  &NftRegistrationNode{Address: "127.0.0.1:4444"},
			address:               "127.0.0.1:4444",
			args:                  args{context.Background()},
			err:                   nil,
			numberConnectCall:     1,
			numberRegisterNftCall: 1,
			assertion:             assert.NoError,
		}, {
			node:                  &NftRegistrationNode{Address: "127.0.0.1:4445"},
			address:               "127.0.0.1:4445",
			args:                  args{context.Background()},
			err:                   fmt.Errorf("connection timeout"),
			numberConnectCall:     1,
			numberRegisterNftCall: 0,
			assertion:             assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			//create client mocks
			clientMock := test.NewMockClient(t)

			//listen needed method
			clientMock.ListenOnConnect("", testCase.err).ListenOnRegisterNft()

			//set up node client only
			testCase.node.client = clientMock.Client

			//assertion error
			testCase.assertion(t, testCase.node.connect(testCase.args.ctx))
			//mock assertion
			clientMock.Client.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address)
			clientMock.AssertRegisterNftCall(testCase.numberRegisterNftCall)
		})
	}
}

func TestNftNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *NftRegistrationNode
	}
	testCases := []struct {
		nodes NftRegistrationNodes
		args  args
		want  NftRegistrationNodes
	}{
		{
			nodes: NftRegistrationNodes{},
			args:  args{node: &NftRegistrationNode{Address: "127.0.0.1"}},
			want: NftRegistrationNodes{
				&NftRegistrationNode{Address: "127.0.0.1"},
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

func TestNftByID(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}
	testCases := []struct {
		nodes NftRegistrationNodes
		args  args
		want  *NftRegistrationNode
	}{
		{
			nodes: NftRegistrationNodes{
				&NftRegistrationNode{ID: "1"},
				&NftRegistrationNode{ID: "2"},
			},
			args: args{"2"},
			want: &NftRegistrationNode{ID: "2"},
		}, {
			nodes: NftRegistrationNodes{
				&NftRegistrationNode{ID: "1"},
				&NftRegistrationNode{ID: "2"},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Add(&NftRegistrationNode{ID: "4"})
			testCase.nodes.Remove("4")
			assert.Equal(t, testCase.want, testCase.nodes.ByID(testCase.args.id))
		})
	}
}
