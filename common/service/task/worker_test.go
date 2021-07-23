package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerTasks(t *testing.T) {
	t.Parallel()

	type fields struct {
		tasks []Task
	}
	tests := []struct {
		name   string
		fields fields
		want   []Task
	}{
		{
			name: "retrieve tasks",
			fields: fields{
				tasks: []Task{&task{id: "1"}, &task{id: "2"}},
			},
			want: []Task{&task{id: "1"}, &task{id: "2"}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{
				tasks: tt.fields.tasks,
			}
			assert.Equal(t, tt.want, worker.Tasks())
		})
	}
}

func TestWorkerTask(t *testing.T) {
	t.Parallel()

	type fields struct {
		tasks []Task
	}
	type args struct {
		taskID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Task
	}{
		{
			name: "get task with id 1",
			fields: fields{
				tasks: []Task{&task{id: "1"}, &task{id: "2"}},
			},
			args: args{"2"},
			want: &task{id: "2"},
		},
		{
			name: "get not exist task",
			fields: fields{
				tasks: []Task{&task{id: "1"}, &task{id: "2"}},
			},
			args: args{"3"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{
				tasks: tt.fields.tasks,
			}
			assert.Equal(t, tt.want, worker.Task(tt.args.taskID))
		})
	}
}

func TestWorkerAddTask(t *testing.T) {
	t.Parallel()

	type args struct {
		task Task
	}
	tests := []struct {
		name string
		args args
		want []Task
	}{
		{
			name: "add task",
			args: args{&task{id: "1"}},
			want: []Task{&task{id: "1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{
				taskCh: make(chan Task),
			}

			go func() {
				worker.AddTask(tt.args.task)
			}()

			<-worker.taskCh
			tasks := worker.tasks
			assert.Equal(t, tt.want, tasks)

		})
	}
}

func TestWorkerRemoveTask(t *testing.T) {
	t.Parallel()

	type args struct {
		subTask Task
	}
	tests := []struct {
		name string
		args args
		want []Task
	}{
		{
			name: "removed task",
			args: args{&task{id: "1"}},
			want: []Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{
				tasks: []Task{tt.args.subTask},
			}

			worker.RemoveTask(tt.args.subTask)
			assert.Equal(t, tt.want, worker.tasks)
		})
	}
}
