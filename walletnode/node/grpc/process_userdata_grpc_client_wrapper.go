package grpc

import (
	"context"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"google.golang.org/grpc"
)

type processUserdataClientWrapper struct {
	pb.ProcessUserdataClient
	receiver    *actor.PID
	actorSystem *actor.ActorSystem
}

func (p *processUserdataClientWrapper) AcceptedNodes(ctx context.Context, in *pb.AcceptedNodesRequest, opts ...grpc.CallOption) (*pb.AcceptedNodesReply, error) {
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

func (p *processUserdataClientWrapper) ConnectTo(ctx context.Context, in *pb.ConnectToRequest, opts ...grpc.CallOption) (*pb.ConnectToReply, error) {
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

func (p *processUserdataClientWrapper) SendUserdata(ctx context.Context, in *pb.UserdataRequest, opts ...grpc.CallOption) (*pb.UserdataReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.UserdataReply)
	if ok {
		return reply, nil
	}
	return &pb.UserdataReply{}, nil
}

func (p *processUserdataClientWrapper) ReceiveUserdata(ctx context.Context, in *pb.RetrieveRequest, opts ...grpc.CallOption) (*pb.UserdataRequest, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.UserdataRequest)
	if ok {
		return reply, nil
	}
	return &pb.UserdataRequest{}, nil
}

func newProcessUserdataClientWrapper(client pb.ProcessUserdataClient, actorSystem *actor.ActorSystem, actorReceiver *actor.PID) pb.ProcessUserdataClient {
	return &processUserdataClientWrapper{
		ProcessUserdataClient: client,
		actorSystem:           actorSystem,
		receiver:              actorReceiver,
	}
}
