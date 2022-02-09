package test

import (
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddNode(t *testing.T) {
	t.Parallel()

	type args struct {
		node *common.SuperNodeClient
	}
	testCases := []struct {
		nodes common.SuperNodeList
		args  args
		want  common.SuperNodeList
	}{
		{
			nodes: common.SuperNodeList{},
			args:  args{node: common.NewSuperNode(nil, "127.0.0.1:4444", "", nil)},
			want: common.SuperNodeList{
				common.NewSuperNode(nil, "127.0.0.1:4444", "", nil),
				common.NewSuperNode(nil, "127.0.0.1:4445", "", nil),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Add(testCase.args.node)
			testCase.nodes.AddNewNode(nil, "127.0.0.1:4445", "", nil)
			assert.Equal(t, testCase.want, testCase.nodes)
		})
	}
}

func TestNodesActivate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes common.SuperNodeList
	}{
		{
			nodes: common.SuperNodeList{
				common.NewSuperNode(nil, "127.0.0.1:4444", "", nil),
				common.NewSuperNode(nil, "127.0.0.1:4445", "", nil),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Activate()
			for _, n := range testCase.nodes {
				assert.True(t, n.IsActive())
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
		nodes common.SuperNodeList
		args  args
		want  *common.SuperNodeClient
	}{
		{
			nodes: common.SuperNodeList{
				common.NewSuperNode(nil, "", "1", nil),
				common.NewSuperNode(nil, "", "2", nil),
			},
			args: args{"2"},
			want: common.NewSuperNode(nil, "", "2", nil),
		},
		{
			nodes: common.SuperNodeList{
				common.NewSuperNode(nil, "", "1", nil),
				common.NewSuperNode(nil, "", "2", nil),
			},
			args: args{"3"},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, testCase.nodes.FindByPastelID(testCase.args.id))
		})
	}
}
