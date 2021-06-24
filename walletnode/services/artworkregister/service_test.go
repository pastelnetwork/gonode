package artworkregister

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/stretchr/testify/assert"
)

func TestServiceRun(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker       *task.Worker
		Storage      *artwork.Storage
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
				Storage:      tt.fields.Storage,
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			tt.assertion(t, service.Run(tt.args.ctx))
		})
	}
}

func TestServiceTasks(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker       *task.Worker
		Storage      *artwork.Storage
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
				Storage:      tt.fields.Storage,
				config:       tt.fields.config,
				db:           tt.fields.db,
				pastelClient: tt.fields.pastelClient,
				nodeClient:   tt.fields.nodeClient,
			}
			assert.Equal(t, tt.want, service.Tasks())
		})
	}
}

func TestServiceTask(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker *task.Worker
		config *Config
	}
	type args struct {
		ticketName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "get task with given id",
			fields: fields{
				Worker: task.NewWorker(),
				config: NewConfig(),
			},
			args: args{"Ticket 1"},
			want: "Ticket 1",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker: tt.fields.Worker,
				config: tt.fields.config,
			}

			group, ctx := errgroup.WithContext(context.Background())
			group.Go(func() error {
				tId := service.AddTask(&Ticket{Name: tt.args.ticketName})
				task := service.Task(tId)
				assert.Equal(t, tt.want, task.Ticket.Name)

				return nil
			})

			service.Run(ctx)

			err := group.Wait()
			assert.NoError(t, err)
		})
	}
}

func TestServiceAddTask(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker *task.Worker
		config *Config
	}
	type args struct {
		ticket *Ticket
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add task",
			fields: fields{
				Worker: task.NewWorker(),
				config: NewConfig(),
			},
			args: args{
				ticket: &Ticket{
					Name: "ticket 1",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker: tt.fields.Worker,
				config: tt.fields.config,
			}

			group, ctx := errgroup.WithContext(context.Background())
			group.Go(func() error {
				tId := service.AddTask(tt.args.ticket)
				assert.NotEmpty(t, tId)
				return nil
			})

			service.Run(ctx)

			err := group.Wait()
			assert.NoError(t, err)
		})
	}
}
