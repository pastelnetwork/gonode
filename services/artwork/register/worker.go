package register

import (
	"context"
	"fmt"
)

type Worker struct {
	jobCh chan *Job
}

func (worker *Worker) AddJob(ctx context.Context, job *Job) {
	select {
	case <-ctx.Done():
		fmt.Println("stop add job")
	case worker.jobCh <- job:
	}
}

func (worker *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case job := <-worker.jobCh:
			fmt.Println(job)
			return nil
		}
	}
	return nil
}

func NewWorker() *Worker {
	return &Worker{
		jobCh: make(chan *Job),
	}
}
