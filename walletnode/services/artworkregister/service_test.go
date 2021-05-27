package artworkregister

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/stretchr/testify/assert"
)

func TestService_Run(t *testing.T) {
	type fields struct {
		Worker       *task.Worker
		config       *Config
		db           storage.KeyValue
		pastelClient pastel.Client
		nodeClient   node.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker:       tt.fields.Worker,
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			tt.assertion(t, service.Run(tt.args.ctx))
		})
	}
}

func TestService_Tasks(t *testing.T) {
	type fields struct {
		Worker       *task.Worker
		config       *Config
		db           storage.KeyValue
		pastelClient pastel.Client
		nodeClient   node.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker:       tt.fields.Worker,
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			assert.Equal(t, tt.want, service.Tasks())
		})
	}
}

func TestService_Task(t *testing.T) {
	type fields struct {
		Worker       *task.Worker
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
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			assert.Equal(t, tt.want, service.Task(tt.args.id))
		})
	}
}

func TestService_AddTask(t *testing.T) {
	type fields struct {
		Worker       *task.Worker
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
	type args struct {
		config       *Config
		db           storage.KeyValue
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
			assert.Equal(t, tt.want, NewService(tt.args.config, tt.args.db, tt.args.pastelClient, tt.args.nodeClient))
		})
	}
}
