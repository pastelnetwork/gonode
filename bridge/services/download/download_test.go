package download

import (
	"context"
	"testing"
	"time"

	nodeMock "github.com/pastelnetwork/gonode/bridge/node/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()
	nodes := pastel.MasterNodes{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, pastel.MasterNode{ExtKey: "key", ExtAddress: "127.0.0.1:14445"})
	}

	tests := map[string]struct {
		wantthumb map[int][]byte
		want      []byte
		wantErr   error
	}{
		"a": {wantthumb: map[int][]byte{0: {2, 78, 67, 33}}, want: []byte{2, 78, 67, 83}, wantErr: nil},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	t.Run("a", func(t *testing.T) {

		pastelClientMock := pastelMock.NewMockClient(t)
		pastelClientMock.ListenOnMasterNodesTop(nodes, nil).ListenOnFindTicketByID(&pastel.IDTicket{TXID: "txid"}, nil).
			ListenOnGetBlockCount(10, nil)

		nodeClientMock := nodeMock.NewMockClient(t)
		nodeClientMock.ListenOnConnect("", nil).ListenOnDownloadThumbnail(tests["a"].wantthumb, nil).ListenOnClose(nil).
			ListenOnDownloadDDAndFP(tests["a"].want, nil)
		nodeClientMock.ConnectionInterface.On("DownloadData").Return(nodeClientMock.DownloadDataInterface)
		doneCh := make(<-chan struct{})
		nodeClientMock.ConnectionInterface.On("Done").Return(doneCh)

		svc := NewService(NewConfig(), pastelClientMock, nodeClientMock)
		go svc.Run(ctx)
		time.Sleep(5 * time.Second)
		data, err := svc.FetchThumbnail(ctx, "txid", 1)
		assert.Nil(t, err)

		assert.Equal(t, data[0], tests["a"].wantthumb[0])

		ddData, err := svc.FetchDupeDetectionData(ctx, "txid")
		assert.Nil(t, err)

		assert.Equal(t, ddData, tests["a"].want)
	})

	cancel()
}
