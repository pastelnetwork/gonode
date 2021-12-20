package client

import (
	"context"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"google.golang.org/grpc"
)

type registerArtworkClientWrapper struct {
	pb.RegisterArtworkClient
	receiver    *actor.PID
	actorSystem *actor.ActorSystem
}

func (p *registerArtworkClientWrapper) SendArtTicketSignature(ctx context.Context, in *pb.SendArtTicketSignatureRequest, opts ...grpc.CallOption) (*pb.SendArtTicketSignatureReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.SendArtTicketSignatureReply)
	if ok {
		return reply, nil
	}
	return &pb.SendArtTicketSignatureReply{}, nil
}

func newRegisterArtworkClientWrapper(client pb.RegisterArtworkClient, actorSystem *actor.ActorSystem, actorReceiver *actor.PID) pb.RegisterArtworkClient {
	return &registerArtworkClientWrapper{
		RegisterArtworkClient: client,
		actorSystem:           actorSystem,
		receiver:              actorReceiver,
	}
}
