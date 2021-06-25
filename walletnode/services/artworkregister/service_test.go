package artworkregister

import (
	"context"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/stretchr/testify/assert"
)

func parseTicketNames(tasks ...*Task) []string {
	var s []string
	for _, t := range tasks {
		s = append(s, t.Ticket.Name)
	}
	return s
}

func TestServiceRun(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker  *task.Worker
		Storage *artwork.Storage
		config  *Config
	}

	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "run service",
			fields: fields{
				Worker:  task.NewWorker(),
				Storage: artwork.NewStorage(fs.NewFileStorage("./")),
				config:  NewConfig(),
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker:  tt.fields.Worker,
				Storage: tt.fields.Storage,
				config:  tt.fields.config,
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			tt.assertion(t, service.Run(ctx))
		})
	}
}

func TestServiceTasks(t *testing.T) {
	t.Parallel()

	type fields struct {
		Worker *task.Worker
		config *Config
	}

	type args struct {
		ticketNames []string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "get tasks",
			fields: fields{
				Worker: task.NewWorker(),
				config: NewConfig(),
			},
			args: args{
				ticketNames: []string{"ticket 1"},
			},
			want: []string{"ticket 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				Worker: tt.fields.Worker,
				config: tt.fields.config,
			}

			group, ctx := errgroup.WithContext(context.Background())
			group.Go(func() error {
				for _, t := range tt.args.ticketNames {
					service.AddTask(&Ticket{Name: t})
				}

				tasks := service.Tasks()
				assert.Equal(t, tt.want, parseTicketNames(tasks...))

				return nil
			})

			service.Run(ctx)

			err := group.Wait()
			assert.NoError(t, err)
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
		want   []string
	}{
		{
			name: "get task with given id",
			fields: fields{
				Worker: task.NewWorker(),
				config: NewConfig(),
			},
			args: args{"Ticket 1"},
			want: []string{"Ticket 1"},
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
				taskID := service.AddTask(&Ticket{Name: tt.args.ticketName})
				task := service.Task(taskID)
				assert.Equal(t, tt.want, parseTicketNames(task))

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
				taskID := service.AddTask(tt.args.ticket)
				assert.NotEmpty(t, taskID)
				return nil
			})

			service.Run(ctx)

			err := group.Wait()
			assert.NoError(t, err)
		})
	}

}
