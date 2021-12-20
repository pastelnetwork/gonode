package client

import (
	"context"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func headerFromContextMetadata(ctx context.Context) map[string]string {
	var header = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			header[k] = strings.Join(v, ",")
		}
	}
	return header
}

type processUserdataClientWrapper struct {
	pb.ProcessUserdataClient
	receiver    *actor.PID
	actorSystem *actor.ActorSystem
}

func (p *processUserdataClientWrapper) SendUserdataToPrimary(ctx context.Context, in *pb.SuperNodeRequest, _ ...grpc.CallOption) (*pb.SuperNodeReply, error) {
	header := headerFromContextMetadata(ctx)
	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	snReply, ok := res.(*pb.SuperNodeReply)
	if ok {
		return snReply, nil
	}
	return &pb.SuperNodeReply{}, nil
}

func (p *processUserdataClientWrapper) SendUserdataToLeader(ctx context.Context, in *pb.UserdataRequest, _ ...grpc.CallOption) (*pb.SuperNodeReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.SuperNodeReply)
	if ok {
		return reply, nil
	}
	return &pb.SuperNodeReply{}, nil
}

func newProcessUserdataClientWrapper(client pb.ProcessUserdataClient, actorSystem *actor.ActorSystem, actorReceiver *actor.PID) pb.ProcessUserdataClient {
	return &processUserdataClientWrapper{
		ProcessUserdataClient: client,
		actorSystem:           actorSystem,
		receiver:              actorReceiver,
	}
}
