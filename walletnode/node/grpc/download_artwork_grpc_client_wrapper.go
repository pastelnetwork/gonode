package grpc

import (
	"context"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
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

type downloadArtworkClientWrapper struct {
	pb.DownloadArtworkClient
	receiver    *actor.PID
	actorSystem *actor.ActorSystem
}

func (p *downloadArtworkClientWrapper) AcceptedNodes(ctx context.Context, in *pb.AcceptedNodesRequest, _ ...grpc.CallOption) (*pb.AcceptedNodesReply, error) {
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

func newDownloadArtworkClientWrapper(client pb.DownloadArtworkClient, actorSystem *actor.ActorSystem, actorReceiver *actor.PID) pb.DownloadArtworkClient {
	return &downloadArtworkClientWrapper{
		DownloadArtworkClient: client,
		actorSystem:           actorSystem,
		receiver:              actorReceiver,
	}
}
