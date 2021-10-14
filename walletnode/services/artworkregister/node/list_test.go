package node

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
)

func TestNodesAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		node *Node
	}
	testCases := []struct {
		nodes List
		args  args
		want  List
	}{
		{
			nodes: List{},
			args:  args{node: &Node{address: "127.0.0.1"}},
			want: List{
				&Node{address: "127.0.0.1"},
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
				&Node{address: "127.0.0.1"},
				&Node{address: "127.0.0.2"},
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

				node := &Node{
					Connection: c.client.Connection,
					activated:  c.activated,
					mtx:        &sync.RWMutex{},
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
		want  *Node
	}{
		{
			nodes: List{
				&Node{pastelID: "1"},
				&Node{pastelID: "2"},
			},
			args: args{"2"},
			want: &Node{pastelID: "2"},
		}, {
			nodes: List{
				&Node{pastelID: "1"},
				&Node{pastelID: "2"},
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

func TestNodesSendImage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		file *artwork.File
	}

	type nodeAttribute struct {
		address   string
		returnErr error
	}

	testCases := []struct {
		nodes                []nodeAttribute
		args                 args
		err                  error
		fingersAndScore      *pastel.FingerAndScores
		numberProbeImageCall int
	}{
		{
			nodes:                []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", nil}},
			args:                 args{context.Background(), &artwork.File{}},
			err:                  nil,
			fingersAndScore:      &pastel.FingerAndScores{},
			numberProbeImageCall: 1,
		},
		{
			nodes:                []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", fmt.Errorf("failed to open stream")}},
			args:                 args{context.Background(), &artwork.File{}},
			err:                  fmt.Errorf("failed to open stream"),
			fingersAndScore:      &pastel.FingerAndScores{},
			numberProbeImageCall: 1,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			nodes := List{}
			clients := []*test.Client{}

			for _, a := range testCase.nodes {
				//client mock
				client := test.NewMockClient(t)
				//listen on uploadImage call
				client.ListenOnRegisterArtwork_ProbeImage(testCase.fingersAndScore, testCase.err)
				clients = append(clients, client)

				nodes.Add(&Node{
					address:         a.address,
					RegisterArtwork: client.RegisterArtwork,
				})
			}

			err := nodes.ProbeImage(testCase.args.ctx, testCase.args.file)
			if err != nil {
				assert.Equal(t, errors.Errorf("failed to probe image: %w", testCase.err).Error(), err.Error())
			} else {
				assert.Equal(t, err, testCase.err)
			}

			//mock assertion each client
			for _, client := range clients {
				client.RegisterArtwork.AssertExpectations(t)
				client.AssertRegisterArtwork_ProbeImageCall(testCase.numberProbeImageCall, testCase.args.ctx, testCase.args.file)
			}
		})
	}
}
