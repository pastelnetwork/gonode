package grpc

import (
	"context"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"google.golang.org/grpc"
)

type registerArtworkClientWrapper struct {
	pb.RegisterArtworkClient
	receiver    *actor.PID
	actorSystem *actor.ActorSystem
}

func (p *registerArtworkClientWrapper) AcceptedNodes(ctx context.Context, in *pb.AcceptedNodesRequest, opts ...grpc.CallOption) (*pb.AcceptedNodesReply, error) {
	header := headerFromContextMetadata(ctx)
	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	snReply, ok := res.(*pb.AcceptedNodesReply)
	if ok {
		return snReply, nil
	}
	return &pb.AcceptedNodesReply{}, nil
}

func (p *registerArtworkClientWrapper) ConnectTo(ctx context.Context, in *pb.ConnectToRequest, opts ...grpc.CallOption) (*pb.ConnectToReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.ConnectToReply)
	if ok {
		return reply, nil
	}
	return &pb.ConnectToReply{}, nil
}

func (p *registerArtworkClientWrapper) SendSignedNFTTicket(ctx context.Context, in *pb.SendSignedNFTTicketRequest, opts ...grpc.CallOption) (*pb.SendSignedNFTTicketReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.SendSignedNFTTicketReply)
	if ok {
		return reply, nil
	}
	return &pb.SendSignedNFTTicketReply{}, nil
}

func (p *registerArtworkClientWrapper) SendPreBurntFeeTxid(ctx context.Context, in *pb.SendPreBurntFeeTxidRequest, opts ...grpc.CallOption) (*pb.SendPreBurntFeeTxidReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.SendPreBurntFeeTxidReply)
	if ok {
		return reply, nil
	}
	return &pb.SendPreBurntFeeTxidReply{}, nil
}

func (p *registerArtworkClientWrapper) SendTicket(ctx context.Context, in *pb.SendTicketRequest, opts ...grpc.CallOption) (*pb.SendTicketReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.SendTicketReply)
	if ok {
		return reply, nil
	}
	return &pb.SendTicketReply{}, nil
}

func newRegisterArtworkClientWrapper(client pb.RegisterArtworkClient, actorSystem *actor.ActorSystem, actorReceiver *actor.PID) pb.RegisterArtworkClient {
	return &registerArtworkClientWrapper{
		RegisterArtworkClient: client,
		actorSystem:           actorSystem,
		receiver:              actorReceiver,
	}
}
