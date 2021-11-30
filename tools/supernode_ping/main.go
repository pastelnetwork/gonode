package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
	"google.golang.org/grpc"
)

func main() {
	var (
		serverAddr = flag.String("server_addr", "", "The server address in the format of host:port, like localhost:4444")
	)

	flag.Parse()
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}
	// validate input
	if *serverAddr == "" {
		Usage()
		return
	}

	// Prepare the client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, *serverAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect")
		return
	}
	defer conn.Close()
	client := pb.NewHealthCheckClient(conn)

	// Send ping request
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.Ping(ctx, &pb.PingRequest{Msg: "hello"})
	if err != nil {
		log.WithError(err).Fatal("Failed to ping")
		return
	} else {
		log.Infof("Received ping reply: %s", res.Reply)
	}
}
