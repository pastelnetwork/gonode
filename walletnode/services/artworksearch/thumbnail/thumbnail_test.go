package thumbnail

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/pastel"

	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	nodeMock "github.com/pastelnetwork/gonode/walletnode/node/test"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtAddress: fmt.Sprint(i)})
	}

	pastelClientMock := pastelMock.NewMockClient(t)
	pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

	tests := map[string]struct {
		helper      Helper
		connections uint
		nodeErr     error
		err         error
	}{
		"one":              {connections: 1, err: nil},
		"max":              {connections: 10, err: nil},
		"more-than-max":    {connections: maxConnections + 1, err: errors.New(maxConnectionsErr)},
		"node-connect-err": {connections: 4, nodeErr: errors.New("test"), err: errors.New("test")},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nodeClientMock := nodeMock.NewMockClient(t)
			if tc.nodeErr == nil {
				nodeClientMock.ListenOnConnect("", nil)
			} else {
				nodeClientMock.ListenOnConnect("0", nil)
				nodeClientMock.ListenOnConnect("1", nil)
				nodeClientMock.ListenOnConnect("2", nil)
				nodeClientMock.ListenOnConnect("3", tc.nodeErr)
			}
			nodeClientMock.ListenOnClose(nil)

			helper := New(pastelClientMock, nodeClientMock, 2*time.Second)

			err := helper.Connect(context.Background(), tc.connections)
			assert.Equal(t, tc.err, err)
			helper.Close()
		})
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{})
	}

	pastelClientMock := pastelMock.NewMockClient(t)
	pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

	nodeClientMock := nodeMock.NewMockClient(t)
	nodeClientMock.ListenOnConnect("", nil).ListenOnRegisterArtwork().ListenOnClose(nil)

	helper := New(pastelClientMock, nodeClientMock, 2*time.Second)

	tests := map[string]struct {
		helper      Helper
		connections uint
		err         error
	}{
		"one-node":  {helper: helper, connections: 1, err: nil},
		"max-nodes": {helper: helper, connections: 10, err: nil},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			err := tc.helper.Connect(ctx, tc.connections)
			assert.Nil(t, err)

			_, err = tc.helper.Fetch(ctx, []byte("key"))
			assert.Nil(t, err)
		})
	}
}
