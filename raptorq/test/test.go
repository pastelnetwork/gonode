package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/raptorq/node/grpc"
)

func main() {
	// create connection to service
	client := grpc.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	connection, err := client.Connect(ctx, "127.0.0.1:50051")
	if err != nil {
		fmt.Printf("Connect error:%v\n", err)
		return
	}

	data := []byte("ababcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdcd")

	service := connection.RaptorQ(&node.Config{
		RqFilesDir: "/tmp/",
	})

	_, err = service.Encode(context.Background(), data)
	if err != nil {
		fmt.Printf("Encode error:%v\n", err)
		return
	}

	_, err = service.EncodeInfo(context.Background(), data, 2, "block_hash", "pastelID")
	if err != nil {
		fmt.Printf("EncodeInfo error:%v\n", err)
		return
	}
}
