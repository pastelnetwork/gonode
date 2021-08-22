package artworkregister

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/artwork"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/p2p"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqmock "github.com/pastelnetwork/gonode/raptorq/node/test"
	"github.com/pastelnetwork/gonode/supernode/node"
	nodeMock "github.com/pastelnetwork/gonode/supernode/node/mocks"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	type args struct {
		config        *Config
		pastelClient  pastel.Client
		nodeClient    node.Client
		p2pClient     p2p.Client
		raptorQClient rqnode.Client
	}

	config := NewConfig()
	pastelClient := pastelMock.NewMockClient(t)
	p2pClient := p2pMock.NewMockClient(t)
	nodeClient := &nodeMock.Client{}
	raptorQClient := rqmock.NewMockClient(t)

	testCases := []struct {
		args args
		want *Service
	}{
		{
			args: args{
				config:        config,
				pastelClient:  pastelClient.Client,
				p2pClient:     p2pClient.Client,
				raptorQClient: raptorQClient.Client,
				nodeClient:    nodeClient,
			},
			want: &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				p2pClient:    p2pClient.Client,
				rqClient:     raptorQClient.Client,
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			service := NewService(testCase.args.config, nil, testCase.args.pastelClient, testCase.args.nodeClient,
				testCase.args.p2pClient, testCase.args.raptorQClient, nil)
			assert.Equal(t, testCase.want.config, service.config)
			assert.Equal(t, testCase.want.pastelClient, service.pastelClient)
			assert.Equal(t, testCase.want.p2pClient, service.p2pClient)
			assert.Equal(t, testCase.want.rqClient, service.rqClient)
		})
	}
}

func TestServiceRun(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			config := &Config{
				Config: common.Config{
					PastelID: "pastelID",
				},
			}
			pastelClient := pastelMock.NewMockClient(t)
			p2pClient := p2pMock.NewMockClient(t)
			raptorQClient := rqmock.NewMockClient(t)
			service := &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				p2pClient:    p2pClient.Client,
				rqClient:     raptorQClient.Client,
				Worker:       task.NewWorker(),
				Storage:      artwork.NewStorage(nil),
			}
			ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
			defer cancel()
			err := service.Run(ctx)
			assert.Equal(t, testCase.want, err)
		})
	}
}

func TestServiceNewTask(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	testCases := []struct {
		args args
	}{
		{
			args: args{
				ctx: context.Background(),
			},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			t.Parallel()

			config := &Config{
				Config: common.Config{
					PastelID: "pastelID",
				},
			}
			pastelClient := pastelMock.NewMockClient(t)
			p2pClient := p2pMock.NewMockClient(t)
			raptorQClient := rqmock.NewMockClient(t)

			service := &Service{
				config:       config,
				pastelClient: pastelClient.Client,
				rqClient:     raptorQClient.Client,
				p2pClient:    p2pClient.Client,
				Worker:       task.NewWorker(),
			}
			ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
			defer cancel()
			go service.Run(ctx)
			task := service.NewTask()
			assert.Equal(t, service, task.Service)
		})
	}
}
