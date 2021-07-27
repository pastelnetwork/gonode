package credentials

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/net/credentials/apis/greeter"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	address = "localhost:8080"
)

type FakePastelClient struct {
	signatures map[string][]byte
	data       map[string][]byte
}

func (c *FakePastelClient) Sign(ctx context.Context, data []byte, pastelID, passPhrase string) ([]byte, error) {
	signature := make([]byte, 20)
	rand.Read(signature)
	c.signatures[pastelID] = signature
	c.data[pastelID] = data
	return signature, nil
}

func (c *FakePastelClient) Verify(ctx context.Context, data []byte, signature, pastelID string) (ok bool, err error) {
	ret := true
	if sig, ok := c.signatures[pastelID]; ok {
		ret = ret && bytes.Equal([]byte(signature), sig)
	}

	if d, ok := c.data[pastelID]; ok {
		ret = ret && bytes.Equal([]byte(data), d)
	}
	return ret, nil
}
func TestSecretConnWithGRPC(t *testing.T) {
	signInfoClient := &alts.SignInfo{
		PastelID:   "client_pastel_id",
		PassPhrase: "client_pass_phrase",
	}

	signInfoServer := &alts.SignInfo{
		PastelID:   "server_pastel_id",
		PassPhrase: "server_pass_phrase",
	}

	auth := &FakePastelClient{
		signatures: make(map[string][]byte),
		data:       make(map[string][]byte),
	}

	altsTCServer := NewServerCreds(auth, signInfoServer)
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

	altsTCClient := NewClientCreds(auth, signInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
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
