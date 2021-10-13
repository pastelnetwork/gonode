package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
	"google.golang.org/grpc"
)

func main() {
	var (
		serverAddr = flag.String("server_addr", "", "The server address in the format of host:port, like localhost:4444")
		cmd        = flag.String("cmd", "", "one of value set/get/status")
		key        = flag.String("key", "", "set when action is `get`")
		value      = flag.String("value", "", "set when action is `set`")
	)

	flag.Parse()
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}
	// validate input
	if len(*serverAddr) == 0 {
		Usage()
		return
	}

	// Prepare the client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// entities
	pastelClient := pastel.NewClient(&pastel.Config{
		Hostname: "localhost",
		Port:     12170,
		Username: "rt",
		Password: "rt",
	})
	secInfo := &alts.SecInfo{
		PastelID:   "jXYX6qAEiQKvTmpLZvNKQuphrzKPLACkx73zo9mE3B1kRQ7sjvnJMit3fHVPGNCY7REwdnVB2H42FaZoG8keAi",
		PassPhrase: "passphrase",
		Algorithm:  "ed448",
	}

	altsTCClient := credentials.NewClientCreds(pastelClient, secInfo)
	conn, err := grpc.DialContext(ctx, *serverAddr,
		grpc.WithTransportCredentials(altsTCClient),
		grpc.WithBlock(),
	)

	if err != nil {
		log.WithError(err).Fatal("Failed to connect")
		return
	}
	defer conn.Close()
	client := pb.NewHealthCheckClient(conn)

	switch *cmd {
	case "set":
		if len(*value) == 0 {
			log.Error("value is empty")
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.P2PSet(ctx, &pb.P2PSetRequest{Value: []byte(*value)})
		if err != nil {
			log.WithError(err).Fatal("Failed to call P2PSet")
			return
		} else {
			log.Infof("Returned Key is: %s", res.Key)
		}

	case "get":
		if len(*key) == 0 {
			log.Error("key is empty")
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.P2PGet(ctx, &pb.P2PGetRequest{Key: *key})
		if err != nil {
			log.WithError(err).Fatal("Failed to call P2PGet")
			return
		} else {
			log.Infof("Returned Value is: %s", string(res.Value))
		}

	case "status":
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.Status(ctx, &pb.StatusRequest{})
		if err != nil {
			log.WithError(err).Fatal("Failed to get status")
			return
		} else {
			log.Infof("Received Status reply: %s", res.StatusInJson)
		}
	default:
		log.Error("cmd is not supported")
	}

}
