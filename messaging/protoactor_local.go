package messaging

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/pastelnetwork/gonode/common/context"
)

type localActor struct {
	mapPID  map[string]*actor.PID
	context *actor.RootContext
}

// Send func
func (r *localActor) Send(ctx context.Context, pid *actor.PID, message proto.Message) error {
	r.context.Send(pid, message)
	return nil
}

// RegisterActor func
func (r *localActor) RegisterActor(a actor.Actor, name string) (*actor.PID, error) {
	pid, err := r.context.SpawnNamed(actor.PropsFromProducer(func() actor.Actor { return a }), name)
	if err != nil {
		return nil, err
	}
	r.mapPID[name] = pid
	return pid, nil
}

// DeregisterActor func
func (r *localActor) DeregisterActor(name string) {
	if r.mapPID[name] != nil {
		r.context.Stop(r.mapPID[name])
	}
	delete(r.mapPID, name)
}

// Stop func
func (r *localActor) Stop() {
	for name, pid := range r.mapPID {
		r.context.StopFuture(pid).Wait()
		delete(r.mapPID, name)
	}
}

var lActor *localActor

// NewActor func
func NewActor(system *actor.ActorSystem) Actor {
	if lActor == nil {
		lActor = &localActor{
			context: system.Root,
			mapPID:  make(map[string]*actor.PID),
		}
	}

	return lActor
}
