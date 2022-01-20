package node

import (
	"fmt"
	"testing"

	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
)

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *NftSearchNodeClient
	}
	testCases := []struct {
		nodes List
		args  args
		want  List
	}{
		{
			nodes: List{},
			args:  args{node: &NftSearchNodeClient{address: "127.0.0.1"}},
			want: List{
				&NftSearchNodeClient{address: "127.0.0.1"},
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

func TestNodesActivate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes List
	}{
		{
			nodes: List{
				&NftSearchNodeClient{address: "127.0.0.1"},
				&NftSearchNodeClient{address: "127.0.0.2"},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			testCase.nodes.Activate()
			for _, n := range testCase.nodes {
				assert.True(t, n.activated)
			}
		})
	}
}

func TestNodesDisconnectInactive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes List
		conn  []struct {
			client       *test.Client
			activated    bool
			numberOfCall int
		}
	}{
		{
			nodes: List{},
			conn: []struct {
				client       *test.Client
				activated    bool
				numberOfCall int
			}{
				{
					client:       test.NewMockClient(t),
					numberOfCall: 1,
					activated:    false,
				},
				{
					client:       test.NewMockClient(t),
					numberOfCall: 0,
					activated:    true,
				},
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			for _, c := range testCase.conn {
				c.client.ListenOnClose(nil)

				node := &NftSearchNodeClient{
					ConnectionInterface: c.client.Connection,
					activated:           c.activated,
				}

				testCase.nodes = append(testCase.nodes, node)
			}

			testCase.nodes.DisconnectInactive()

			for j, c := range testCase.conn {
				c := c

				t.Run(fmt.Sprintf("close-called-%d", j), func(t *testing.T) {
					c.client.AssertCloseCall(c.numberOfCall)
				})

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
		nodes List
		args  args
		want  *NftSearchNodeClient
	}{
		{
			nodes: List{
				&NftSearchNodeClient{pastelID: "1"},
				&NftSearchNodeClient{pastelID: "2"},
			},
			args: args{"2"},
			want: &NftSearchNodeClient{pastelID: "2"},
		}, {
			nodes: List{
				&NftSearchNodeClient{pastelID: "1"},
				&NftSearchNodeClient{pastelID: "2"},
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
