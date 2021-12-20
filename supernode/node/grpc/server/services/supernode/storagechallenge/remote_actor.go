package storagechallenge

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

func (service *appSc) Receive(ctx actor.Context) {
	var res proto.Message
	var err error
	var appCtx = metadata.NewIncomingContext(context.Background(), metadata.New(ctx.MessageHeader().ToMap()))
	switch msg := ctx.Message().(type) {
	case *pb.GenerateStorageChallengesRequest:
		res, err = service.GenerateStorageChallenges(appCtx, msg)
	case *pb.ProcessStorageChallengeRequest:
		res, err = service.ProcessStorageChallenge(appCtx, msg)
	case *pb.VerifyStorageChallengeRequest:
		res, err = service.VerifyStorageChallenge(appCtx, msg)
	default:
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("unknown %T action", msg))
	}

	if err != nil {
		ctx.Respond(err)
		return
	}
	ctx.Respond(res)
}
