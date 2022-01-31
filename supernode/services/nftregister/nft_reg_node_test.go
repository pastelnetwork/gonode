package nftregister

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/supernode/services/common"
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
		node                  *common.SuperNodePeer
		address               string
		args                  args
		err                   error
		numberConnectCall     int
		numberRegisterNftCall int
		assertion             assert.ErrorAssertionFunc
	}{
		{
			node: &common.SuperNodePeer{
				Address:   "127.0.0.1:4444",
				NodeMaker: &RegisterNftNodeMaker{}},
			args:                  args{context.Background()},
			err:                   nil,
			numberConnectCall:     1,
			numberRegisterNftCall: 1,
			assertion:             assert.NoError,
		}, {
			node: &common.SuperNodePeer{
				Address:   "127.0.0.1:4445",
				NodeMaker: &RegisterNftNodeMaker{}},
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
			testCase.node.ClientInterface = clientMock.ClientInterface

			//assertion error
			testCase.assertion(t, testCase.node.Connect(testCase.args.ctx))
			//mock assertion
			clientMock.ClientInterface.AssertExpectations(t)
			clientMock.AssertConnectCall(testCase.numberConnectCall, mock.Anything, testCase.address)
			clientMock.AssertRegisterNftCall(testCase.numberRegisterNftCall)
		})
	}
}

func TestNftNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *common.SuperNodePeer
	}
	testCases := []struct {
		nodes common.SuperNodePeerList
		args  args
		want  common.SuperNodePeer
	}{
		{
			nodes: common.SuperNodePeerList{},
			args: args{node: &common.SuperNodePeer{
				Address:   "127.0.0.1",
				NodeMaker: &RegisterNftNodeMaker{}}},
			want: common.SuperNodePeer{
				Address:   "127.0.0.1",
				NodeMaker: &RegisterNftNodeMaker{}},
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
		nodes common.SuperNodePeerList
		args  args
		want  *common.SuperNodePeer
	}{
		{
			nodes: common.SuperNodePeerList{
				&common.SuperNodePeer{
					ID:        "1",
					NodeMaker: &RegisterNftNodeMaker{},
				},
				&common.SuperNodePeer{
					ID:        "1",
					NodeMaker: &RegisterNftNodeMaker{},
				},
			},
			args: args{"2"},
			want: &common.SuperNodePeer{
				ID:        "2",
				NodeMaker: &RegisterNftNodeMaker{}},
		}, {
			nodes: common.SuperNodePeerList{
				&common.SuperNodePeer{
					ID:        "1",
					NodeMaker: &RegisterNftNodeMaker{},
				},
				&common.SuperNodePeer{
					ID:        "1",
					NodeMaker: &RegisterNftNodeMaker{},
				},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Add(&common.SuperNodePeer{
				ID:        "4",
				NodeMaker: &RegisterNftNodeMaker{},
			})
			testCase.nodes.Remove("4")
			assert.Equal(t, testCase.want, testCase.nodes.ByID(testCase.args.id))
		})
	}
}
