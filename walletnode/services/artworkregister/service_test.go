package artworkregister

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/test"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/stretchr/testify/assert"
)

func TestServiceRun(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker task.Worker
	}

	tests := []struct {
		fields fields

		assertion assert.ErrorAssertionFunc
	}{
		{
			fields:    fields{task.NewWorker()},
			assertion: assert.NoError,
		},
	}

	for i, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {
			//only need file storage mock without further action
			fileStorage := test.NewMockFileStorage()
			storage := artwork.NewStorage(fileStorage.FileStorageMock)

			service := &Service{
				Storage: storage,
				Worker:  tt.fields.Worker,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			defer cancel()

			tt.assertion(t, service.Run(ctx))
		})
	}
}

func TestServiceTasks(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker task.Worker
	}
	tests := []struct {
		fields fields
		want   []*Task
	}{
		// {
		// 	fields: fields{},
		// 	want:   nil,
		// },
	}
	for i, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("testCase-%d", i), func(t *testing.T) {

			tt.fields.Worker.AddTask(&Task{})
			service := &Service{
				Worker: tt.fields.Worker,
			}
			assert.Equal(t, tt.want, service.Tasks())
		})
	}
}

func TestServiceTask(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker       task.Worker
		Storage      *artwork.Storage
		config       *Config
		db           storage.KeyValue
		pastelClient pastel.Client
		nodeClient   node.Client
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker:       tt.fields.Worker,
				Storage:      tt.fields.Storage,
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			assert.Equal(t, tt.want, service.Task(tt.args.id))
		})
	}
}

func TestServiceAddTask(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker       task.Worker
		Storage      *artwork.Storage
		config       *Config
		db           storage.KeyValue
		pastelClient pastel.Client
		nodeClient   node.Client
	}
	type args struct {
		ticket *Ticket
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker:       tt.fields.Worker,
				Storage:      tt.fields.Storage,
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			assert.Equal(t, tt.want, service.AddTask(tt.args.ticket))
		})
	}
}

func TestNewService(t *testing.T) {
	t.Parallel()

	type args struct {
		config       *Config
		db           storage.KeyValue
		fileStorage  storage.FileStorage
		pastelClient pastel.Client
		nodeClient   node.Client
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewService(tt.args.config, tt.args.db, tt.args.fileStorage, tt.args.pastelClient, tt.args.nodeClient))
		})
	}
}
