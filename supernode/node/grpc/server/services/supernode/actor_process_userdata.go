package supernode

import (
	"context"
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Receive represents actor mocel receive action
func (service *ProcessUserdata) Receive(ctx actor.Context) {
	var res proto.Message
	var err error
	var appCtx = metadata.NewIncomingContext(context.Background(), metadata.New(ctx.MessageHeader().ToMap()))
	switch msg := ctx.Message().(type) {
	case *pb.SuperNodeRequest:
		res, err = service.SendUserdataToPrimary(appCtx, msg)
	case *pb.UserdataRequest:
		res, err = service.SendUserdataToLeader(appCtx, msg)
	default:
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("unknown %T action", msg))
	}

	if err != nil {
		ctx.Respond(err)
		return
	}
	ctx.Respond(res)
}
