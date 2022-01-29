package senseregister

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"testing"
	"time"

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
		nodeClient    node.ClientInterface
		p2pClient     p2p.Client
		raptorQClient rqnode.ClientInterface
	}

	config := NewConfig()
	pastelClient := pastelMock.NewMockClient(t)
	p2pClient := p2pMock.NewMockClient(t)
	nodeClient := &nodeMock.ClientInterface{}
	raptorQClient := rqmock.NewMockClient(t)

	testCases := []struct {
		args args
		want *SenseRegistrationService
	}{
		{
			args: args{
				config:        config,
				pastelClient:  pastelClient.Client,
				p2pClient:     p2pClient.Client,
				raptorQClient: raptorQClient.Client,
				nodeClient:    nodeClient,
			},
			want: &SenseRegistrationService{
				SuperNodeService: &common.SuperNodeService{
					PastelClient: pastelClient.Client,
					P2PClient:    p2pClient.Client,
					RQClient:     raptorQClient.Client,
					Worker:       task.NewWorker(),
				},
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
			assert.Equal(t, testCase.want.PastelClient, service.PastelClient)
			assert.Equal(t, testCase.want.P2PClient, service.P2PClient)
			assert.Equal(t, testCase.want.RQClient, service.RQClient)
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
			service := &SenseRegistrationService{
				config: config,
				SuperNodeService: &common.SuperNodeService{
					PastelClient: pastelClient.Client,
					P2PClient:    p2pClient.Client,
					RQClient:     raptorQClient.Client,
					Worker:       task.NewWorker(),
					Storage:      files.NewStorage(nil),
				},
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

			service := &SenseRegistrationService{
				config: config,
				SuperNodeService: &common.SuperNodeService{
					PastelClient: pastelClient.Client,
					P2PClient:    p2pClient.Client,
					RQClient:     raptorQClient.Client,
					Worker:       task.NewWorker(),
				},
			}
			ctx, cancel := context.WithTimeout(testCase.args.ctx, time.Second)
			defer cancel()
			go service.Run(ctx)
			task := service.NewTask()
			assert.Equal(t, service, task.SenseRegistrationService)
		})
	}
}
