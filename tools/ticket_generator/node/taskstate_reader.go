package node

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
)

type TaskStateReceiver struct {
	ctx    context.Context
	data   chan *artworks.TaskState
	err    error
	stream artworks.RegisterTaskStateClientStream
}

func (r *TaskStateReceiver) begin() {
	defer close(r.data)
	for {
		state, err := r.stream.Recv()
		if err == io.EOF {
			log.WithContext(r.ctx).Debug("EOF received")
			r.data <- state
			return
		}
		if err != nil {
			log.WithContext(r.ctx).Debugf("Receive error: %s", err.Error())
			r.err = errors.Errorf("received error: %w", err)
			return
		}
		r.data <- state
	}
}

func (r *TaskStateReceiver) Recv() (*artworks.TaskState, error) {
	select {
	case <-r.ctx.Done():
		log.WithContext(r.ctx).Debug("Contex done")
		return nil, r.ctx.Err()
	case state, ok := <-r.data:
		if !ok {
			return nil, r.err
		}
		return state, nil
	}
}

func NewTaskStateReceiver(ctx context.Context, stream artworks.RegisterTaskStateClientStream) *TaskStateReceiver {
	r := &TaskStateReceiver{
		ctx:    ctx,
		data:   make(chan *artworks.TaskState),
		err:    nil,
		stream: stream,
	}
	go r.begin()
	return r
}
