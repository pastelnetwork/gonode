package client

import (
	"context"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"google.golang.org/grpc"
)

type storageChallengeClientWrapper struct {
	pb.StorageChallengeClient
	receiver    *actor.PID
	actorSystem *actor.ActorSystem
}

func (p *storageChallengeClientWrapper) GenerateStorageChallenges(ctx context.Context, in *pb.GenerateStorageChallengesRequest, _ ...grpc.CallOption) (*pb.GenerateStorageChallengesReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.GenerateStorageChallengesReply)
	if ok {
		return reply, nil
	}
	return &pb.GenerateStorageChallengesReply{}, nil
}

func (p *storageChallengeClientWrapper) ProcessStorageChallenge(ctx context.Context, in *pb.ProcessStorageChallengeRequest, _ ...grpc.CallOption) (*pb.ProcessStorageChallengeReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.ProcessStorageChallengeReply)
	if ok {
		return reply, nil
	}
	return &pb.ProcessStorageChallengeReply{}, nil
}

func (p *storageChallengeClientWrapper) VerifyStorageChallenge(ctx context.Context, in *pb.VerifyStorageChallengeRequest, _ ...grpc.CallOption) (*pb.VerifyStorageChallengeReply, error) {
	header := headerFromContextMetadata(ctx)

	res, err := actor.NewRootContext(p.actorSystem, header).RequestFuture(p.receiver, in, time.Second*30).Result()
	if err != nil {
		return nil, err
	}
	err, ok := res.(error)
	if ok {
		return nil, err
	}
	reply, ok := res.(*pb.VerifyStorageChallengeReply)
	if ok {
		return reply, nil
	}
	return &pb.VerifyStorageChallengeReply{}, nil
}

func newStorageChallengeClientWrapper(client pb.StorageChallengeClient, actorSystem *actor.ActorSystem, actorReceiver *actor.PID) pb.StorageChallengeClient {
	return &storageChallengeClientWrapper{
		StorageChallengeClient: client,
		actorSystem:            actorSystem,
		receiver:               actorReceiver,
	}
}
