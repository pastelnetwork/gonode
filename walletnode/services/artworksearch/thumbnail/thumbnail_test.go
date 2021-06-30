package thumbnail

import (
	"context"
	"errors"
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
		nodes = append(nodes, pastel.MasterNode{})
	}

	pastelClientMock := pastelMock.NewMockClient(t)
	pastelClientMock.ListenOnMasterNodesTop(nodes, nil)

	nodeClientMock := nodeMock.NewMockClient(t)
	nodeClientMock.ListenOnConnect(nil).ListenOnRegisterArtwork().ListenOnClose(nil)

	helper := New(pastelClientMock, nodeClientMock, 2*time.Second)

	tests := map[string]struct {
		helper      Helper
		connections uint
		err         error
	}{
		"one":           {helper: helper, connections: 1, err: nil},
		"max":           {helper: helper, connections: 10, err: nil},
		"more-than-max": {helper: helper, connections: maxConnections + 1, err: errors.New(maxConnectionsErr)},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.helper.Connect(context.Background(), tc.connections)
			assert.Equal(t, tc.err, err)
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
	nodeClientMock.ListenOnConnect(nil).ListenOnRegisterArtwork().ListenOnClose(nil)

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

			_, err = tc.helper.Fetch(ctx, "key")
			assert.Nil(t, err)
		})
	}
}
