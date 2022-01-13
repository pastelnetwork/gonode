package artworkdownload

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	"github.com/pastelnetwork/gonode/walletnode/node"
	test "github.com/pastelnetwork/gonode/walletnode/node/test/artwork_download"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	type args struct {
		config       *Config
		pastelClient pastel.Client
		nodeClient   node.Client
	}

	config := &Config{}
	pastelClient := pastelMock.NewMockClient(t)
	nodeClient := test.NewMockClient(t)

	testCases := []struct {
		args args
		want *Service
	}{
		{
			args: args{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
			},
			want: &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			service := NewService(testCase.args.config, testCase.args.pastelClient, testCase.args.nodeClient)
			assert.Equal(t, testCase.want.config, service.config)
			assert.Equal(t, testCase.want.pastelClient, service.pastelClient)
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
			service := &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
				Worker:       task.NewWorker(),
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
		ctx    context.Context
		ticket *Ticket
	}
	ticket := &Ticket{
		Txid:               "txid",
		PastelID:           "pastelid",
		PastelIDPassphrase: "passphrase",
	}

	testCases := []struct {
		args args
		want *Ticket
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
			service := &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
				Worker:       task.NewWorker(),
			}
			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()
			go service.Run(ctx)
			taskID := service.AddTask(testCase.args.ticket)
			task := service.GetTask(taskID)
			assert.Equal(t, testCase.want, task.Ticket)
		})
	}
}

func TestServiceGetTask(t *testing.T) {
	type args struct {
		ctx    context.Context
		ticket *Ticket
	}
	ticket := &Ticket{
		Txid:               "txid",
		PastelID:           "pastelid",
		PastelIDPassphrase: "passphrase",
	}

	testCases := []struct {
		args args
		want *Ticket
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
			service := &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
				Worker:       task.NewWorker(),
			}
			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()
			go service.Run(ctx)
			task := NewNftDownloadTask(service, testCase.args.ticket)
			service.Worker.AddTask(task)
			taskID := task.ID()
			addedTask := service.GetTask(taskID)
			assert.Equal(t, testCase.want, addedTask.Ticket)
			time.Sleep(time.Second)
		})
	}
}

func TestServiceListTasks(t *testing.T) {
	type args struct {
		ctx     context.Context
		tickets []*Ticket
	}
	ticket := &Ticket{
		Txid:               "txid",
		PastelID:           "pastelid",
		PastelIDPassphrase: "passphrase",
	}
	var tickets []*Ticket
	tickets = append(tickets, ticket)

	testCases := []struct {
		args args
		want []*Ticket
	}{
		{
			args: args{
				ctx:     context.Background(),
				tickets: tickets,
			},
			want: tickets,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			// t.Parallel()

			config := &Config{}
			pastelClient := pastelMock.NewMockClient(t)
			nodeClient := test.NewMockClient(t)
			service := &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				nodeClient:   nodeClient.Client,
				Worker:       task.NewWorker(),
			}
			ctx, cancel := context.WithCancel(testCase.args.ctx)
			defer cancel()
			go service.Run(ctx)
			var listTaskID []string
			for _, ticket := range testCase.args.tickets {
				listTaskID = append(listTaskID, service.AddTask(ticket))
			}

			for i := range listTaskID {
				task := service.GetTask(listTaskID[i])
				assert.Equal(t, testCase.want[i], task.Ticket)
			}
		})
	}
}
