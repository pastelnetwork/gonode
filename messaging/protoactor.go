package messaging

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/pastelnetwork/gonode/common/context"
)

// Messaging interface
type Messaging interface {
	Send(ctx context.Context, pid *actor.PID, message proto.Message) error
	RegisterActor(a actor.Actor, name string) (*actor.PID, error)
	DeregisterActor(name string)
}

// Actor interface
type Actor interface {
	Messaging
	Stop()
}

// Remoter interface
type Remoter interface {
	Messaging
	Start()
	GracefulStop()
}
