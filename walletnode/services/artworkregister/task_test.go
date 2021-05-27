package artworkregister

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/stretchr/testify/assert"
)

func TestTaskRun(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
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
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			tt.assertion(t, task.Run(tt.args.ctx))
		})
	}
}

func TestTask_run(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
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
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			tt.assertion(t, task.run(tt.args.ctx))
		})
	}
}

func TestTaskMeshNodes(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
	}
	type args struct {
		ctx          context.Context
		nodes        Nodes
		primaryIndex int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      Nodes
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			got, err := task.meshNodes(tt.args.ctx, tt.args.nodes, tt.args.primaryIndex)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskisSuitableStorageFee(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			got, err := task.isSuitableStorageFee(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaskpastelTopNodes(t *testing.T) {
	type fields struct {
		Task    task.Task
		Service *Service
		Ticket  *Ticket
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      Nodes
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Task:    tt.fields.Task,
				Service: tt.fields.Service,
				Ticket:  tt.fields.Ticket,
			}
			got, err := task.pastelTopNodes(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewTask(t *testing.T) {
	type args struct {
		service *Service
		Ticket  *Ticket
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewTask(tt.args.service, tt.args.Ticket))
		})
	}
}
