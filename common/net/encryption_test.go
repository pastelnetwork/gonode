package net

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/crypto"

	"github.com/pastelnetwork/gonode/common/net/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type serverMock struct {
	proto.UnsafeReverseServer
}

func (s *serverMock) Do(c context.Context, request *proto.Request) (response *proto.Response, err error) {
	n := 0
	rune := make([]rune, len(request.Message))
	for _, r := range request.Message {
		rune[n] = r
		n++
	}
	// Reverse using runes.
	rune = rune[0:n]
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	output := string(rune)
	response = &proto.Response{
		Message: output,
	}
	return response, nil
}

func makeTestServer(t *testing.T, transport credentials.TransportCredentials) {
	listener, err := net.Listen("tcp", ":5300")
	fmt.Printf("test server started\n")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(transport),
	}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterReverseServer(grpcServer, &serverMock{})
	fmt.Printf("test server listen\n")
	if err := grpcServer.Serve(listener); err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Printf("test server finished\n")

	time.Sleep(1 * time.Second)
}

func TestClientAndServer(t *testing.T) {
	clientCipher := crypto.NewCipher()
	serverCipher := crypto.NewCipher()
	handshaker := NewHandshakerMock(clientCipher, serverCipher, true)
	transport := NewTransportCredentials(handshaker)

	serverHandshaker := NewHandshakerMock(clientCipher, serverCipher, false)
	serverTransport := NewTransportCredentials(serverHandshaker)

	go makeTestServer(t, serverTransport)
	time.Sleep(1 * time.Second)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transport),
	}
	conn, err := grpc.Dial("127.0.0.1:5300", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewReverseClient(conn)
	request := &proto.Request{
		Message: "message",
	}
	response, err := client.Do(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	fmt.Println(response.Message)
	assert.Equal(t, "egassem", response.Message)
}
