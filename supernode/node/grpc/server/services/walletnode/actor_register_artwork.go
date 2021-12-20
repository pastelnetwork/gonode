package walletnode

import (
	"context"
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Receive represents actor mocel receive action
func (service *RegisterArtwork) Receive(ctx actor.Context) {
	var res proto.Message
	var err error
	var appCtx = metadata.NewIncomingContext(context.Background(), metadata.New(ctx.MessageHeader().ToMap()))
	switch msg := ctx.Message().(type) {
	case *pb.AcceptedNodesRequest:
		res, err = service.AcceptedNodes(appCtx, msg)
	case *pb.ConnectToRequest:
		res, err = service.ConnectTo(appCtx, msg)
	case *pb.SendSignedNFTTicketRequest:
		res, err = service.SendSignedNFTTicket(appCtx, msg)
	case *pb.SendPreBurntFeeTxidRequest:
		res, err = service.SendPreBurntFeeTxid(appCtx, msg)
	case *pb.SendTicketRequest:
		res, err = service.SendTicket(appCtx, msg)
	default:
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("unknown %T action", msg))
	}

	if err != nil {
		ctx.Respond(err)
		return
	}
	ctx.Respond(res)
}
