package credentials

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/net/credentials/apis/greeter"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	address = "localhost:8080"
)

func TestSecretConnWithGRPC(t *testing.T) {
	altsTCServer := NewServerCreds(DefaultServerOptions())
	s := &greeter.Greeter{}
	server := grpc.NewServer(grpc.Creds(altsTCServer))
	greeter.RegisterGreeterServiceServer(server, s)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net listen: %v", err)
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wg.Done()
		// serve the incoming connections
		server.Serve(ln)
	}()

	altsTCClient := NewClientCreds(DefaultClientOptions())
	conn, err := grpc.Dial(
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(altsTCClient),
	)
	if err != nil {
		t.Fatalf("grpc dial: %v", err)
	}
	defer conn.Close()

	data := make([]byte, 1024)
	rand.Read(data)
	name := hex.EncodeToString(data)
	client := greeter.NewGreeterServiceClient(conn)
	reply, err := client.SayHi(context.Background(), &greeter.SayHiRequest{Name: name})
	if err != nil {
		t.Fatalf("grpc call: %v", err)
	}
	assert.Equal(t, name, reply.Message)

	server.GracefulStop()
	wg.Wait()
}
