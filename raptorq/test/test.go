package main

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/raptorq/node/grpc"
)

func main() {
	// create connection to service
	client := grpc.NewClient()
	connection, err := client.Connect(context.Background(), "127.0.0.1:50051")
	if err != nil {
		fmt.Printf("Connect error:%v\n", err)
		return
	}

	data := []byte("ababcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdcd")

	service := connection.RaptorQ()
	res, err := service.Encode(context.Background(), data)

	if err != nil {
		fmt.Printf("Encode error:%v\n", err)

		return
	}

	// print result
	for i, _ := range res {
		fmt.Println(string(res[i]))
	}
}
