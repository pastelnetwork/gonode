package main

import (
	"context"
	"encoding/base32"
	"fmt"
	"net"
	"os"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/client"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/middleware"
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
	"google.golang.org/grpc"
)

var (
	host  = os.Getenv("SERVICE_HOST")
	rawID = os.Getenv("SERVICE_ID")
)

func main() {
	pID := base32.StdEncoding.EncodeToString([]byte(rawID))
	pclient := newMockPastelClient()
	secInfo := &alts.SecInfo{
		PastelID:   pID,
		PassPhrase: rawID,
		Algorithm:  "ed448",
	}
	p2p := newMockP2P(pID)

	nodeClient := client.New(pclient, secInfo)
	domainSt, stopActorFunc := storagechallenge.NewService(&storagechallenge.Config{
		Config: common.Config{
			PastelID:   pID,
			PassPhrase: rawID,
		},
		StorageChallengeExpiredBlocks: 2,
		NumberOfChallengeReplicas:     3,
	}, nodeClient, p2p, pclient, &challengeStateStorage{})
	defer stopActorFunc()

	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	grpcServer := grpc.NewServer(
		middleware.UnaryInterceptor(),
		middleware.StreamInterceptor(),
		middleware.AltsCredential(pclient, secInfo),
	)
	go p2p.Run(ctx)
	go domainSt.Run(ctx)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:14444", host))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	grpcServer.Serve(listener)
	defer grpcServer.GracefulStop()
}
