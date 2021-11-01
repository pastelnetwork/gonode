package node

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestNodesActive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes List
		want  int
	}{
		{
			nodes: List{
				&Node{address: "127.0.0.1", activated: true},
				&Node{address: "127.0.0.2", activated: false},
			},
			want: 1,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			nodes := testCase.nodes.Active()
			assert.Equal(t, testCase.want, len(nodes))
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

func TestNodesMatchFiles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes     List
		assertion assert.ErrorAssertionFunc
	}{
		{
			nodes: List{
				&Node{
					address: "127.0.0.1",
					file:    []byte("1234"),
				},
				&Node{
					address: "127.0.0.2",
					file:    []byte("1234"),
				},
			},
			assertion: assert.NoError,
		},
		{
			nodes: List{
				&Node{
					address: "127.0.0.1",
					file:    []byte("1234"),
				},
				&Node{
					address: "127.0.0.2",
					file:    []byte("1235"),
				},
			},
			assertion: assert.Error,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			err := testCase.nodes.MatchFiles()
			testCase.assertion(t, err)
		})
	}
}

func TestNodesFiles(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		nodes List
		want  []byte
	}{
		{
			nodes: List{
				&Node{
					address: "127.0.0.1",
					file:    []byte("1234"),
				},
			},
			want: []byte("1234"),
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			file := testCase.nodes.File()
			assert.Equal(t, testCase.want, file)
		})
	}
}

func TestNodesDownload(t *testing.T) {
	// FIXME: enable later
	t.Skip()
	t.Parallel()

	type args struct {
		ctx                               context.Context
		txid, timestamp, signature, ttxid string
		// downloadErrs error
	}

	type nodeAttribute struct {
		address   string
		returnErr error
	}

	testCases := []struct {
		nodes              []nodeAttribute
		args               args
		err                error
		file               []byte
		numberDownloadCall int
	}{
		{
			nodes:              []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", nil}},
			args:               args{context.Background(), "txid", "timestamp", "signature", "ttxid"},
			err:                nil,
			file:               []byte("test"),
			numberDownloadCall: 1,
		},
		{
			nodes:              []nodeAttribute{{"127.0.0.1:4444", nil}, {"127.0.0.1:4445", nil}},
			args:               args{context.Background(), "txid", "timestamp", "signature", "ttxid"},
			err:                nil,
			file:               nil,
			numberDownloadCall: 1,
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
				client.ListenOnConnect("", nil)
				client.ListenOnDownload(testCase.file, testCase.err)
				clients = append(clients, client)

				nodes.Add(&Node{
					address:         a.address,
					DownloadArtwork: client.DownloadArtwork,
				})
			}

			_, err := nodes.Download(testCase.args.ctx, testCase.args.txid, testCase.args.timestamp, testCase.args.signature, testCase.args.ttxid, 5*time.Second, nil)
			assert.Equal(t, testCase.err, err)

			//mock assertion each client
			for _, client := range clients {
				client.DownloadArtwork.AssertExpectations(t)
				client.AssertDownloadCall(testCase.numberDownloadCall, testCase.args.ctx,
					testCase.args.txid, testCase.args.timestamp, testCase.args.signature, testCase.args.ttxid)
			}
		})
	}
}
