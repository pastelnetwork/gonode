package nftdownload

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/nft_download"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	type args struct {
		config       *Config
		pastelClient pastel.Client
		nodeClient   node.ClientInterface
	}

	config := &Config{}
	pastelClient := pastelMock.NewMockClient(t)
	nodeClient := test.NewMockClient(t)

	testCases := []struct {
		args args
		want *NftDownloadingService
	}{
		{
			args: args{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.ClientInterface,
			},
			want: &NftDownloadingService{
				config:        config,
				pastelHandler: &mixins.PastelHandler{PastelClient: pastelClient.Client},
				nodeClient:    nodeClient.ClientInterface,
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			service := NewNftDownloadService(testCase.args.config, testCase.args.pastelClient, testCase.args.nodeClient)
			assert.Equal(t, testCase.want.config, service.config)
			assert.Equal(t, testCase.want.pastelHandler.PastelClient, service.pastelHandler.PastelClient)
			assert.Equal(t, testCase.want.nodeClient, service.nodeClient)
		})
	}
}

func TestServiceRun(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		args args
		want error
	}{
		{
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			config := &Config{}
			pastelClient := pastelMock.NewMockClient(t)
			nodeClient := test.NewMockClient(t)
			service := &NftDownloadingService{
				config:        config,
				pastelHandler: &mixins.PastelHandler{PastelClient: pastelClient.Client},
				nodeClient:    nodeClient.ClientInterface,
				Worker:        task.NewWorker(),
			}
			ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
			defer cancel()
			err := service.Run(ctx)
			assert.Equal(t, testCase.want, err)
		})
	}
}

func TestServiceAddTask(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload *nft.NftDownloadPayload
	}
	payload := &nft.NftDownloadPayload{
		Txid: "txid",
		Pid:  "pastelid",
		Key:  "passphrase",
	}

	testCases := []struct {
		args args
		want *nft.NftDownloadPayload
	}{
		{
			args: args{
				ctx:     context.Background(),
				payload: payload,
			},
			want: payload,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			config := &Config{}
			pastelClient := pastelMock.NewMockClient(t)
			nodeClient := test.NewMockClient(t)
			service := &NftDownloadingService{
				config:        config,
				pastelHandler: &mixins.PastelHandler{PastelClient: pastelClient.Client},
				nodeClient:    nodeClient.ClientInterface,
				Worker:        task.NewWorker(),
			}
			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()
			go service.Run(ctx)
			taskID := service.AddTask(testCase.args.payload)
			task := service.GetTask(taskID)
			assert.Equal(t, testCase.want, task.Request)
		})
	}
}

func TestServiceGetTask(t *testing.T) {
	type args struct {
		ctx    context.Context
		ticket *NftDownloadingRequest
	}
	ticket := &NftDownloadingRequest{
		Txid:               "txid",
		PastelID:           "pastelid",
		PastelIDPassphrase: "passphrase",
	}

	testCases := []struct {
		args args
		want *NftDownloadingRequest
	}{
		{
			args: args{
				ctx:    context.Background(),
				ticket: ticket,
			},
			want: ticket,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			config := &Config{}
			pastelClient := pastelMock.NewMockClient(t)
			nodeClient := test.NewMockClient(t)
			service := &NftDownloadingService{
				config:        config,
				pastelHandler: &mixins.PastelHandler{PastelClient: pastelClient.Client},
				nodeClient:    nodeClient.ClientInterface,
				Worker:        task.NewWorker(),
			}
			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()
			go service.Run(ctx)
			task := NewNftDownloadTask(service, testCase.args.ticket)
			service.Worker.AddTask(task)
			taskID := task.ID()
			addedTask := service.GetTask(taskID)
			assert.Equal(t, testCase.want, addedTask.Request)
			time.Sleep(time.Second)
		})
	}
}

func TestServiceListTasks(t *testing.T) {
	type args struct {
		ctx      context.Context
		payloads []*nft.NftDownloadPayload
	}
	payload := &nft.NftDownloadPayload{
		Txid: "txid",
		Pid:  "pastelid",
		Key:  "passphrase",
	}
	var payloads []*nft.NftDownloadPayload
	payloads = append(payloads, payload)

	testCases := []struct {
		args args
		want []*nft.NftDownloadPayload
	}{
		{
			args: args{
				ctx:      context.Background(),
				payloads: payloads,
			},
			want: payloads,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			config := &Config{}
			pastelClient := pastelMock.NewMockClient(t)
			nodeClient := test.NewMockClient(t)
			service := &NftDownloadingService{
				config:        config,
				pastelHandler: &mixins.PastelHandler{PastelClient: pastelClient.Client},
				nodeClient:    nodeClient.ClientInterface,
				Worker:        task.NewWorker(),
			}
			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()
			go service.Run(ctx)
			var listTaskID []string
			for _, ticket := range testCase.args.payloads {
				listTaskID = append(listTaskID, service.AddTask(ticket))
			}

			for i := range listTaskID {
				task := service.GetTask(listTaskID[i])
				assert.Equal(t, testCase.want[i], task.Request)
			}
		})
	}
}
