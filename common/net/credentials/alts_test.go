package credentials

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
	signatures  map[string][]byte
	data        map[string][]byte
	errorSign   bool
	errorVerify bool
}

func (c *FakePastelClient) SetErrorSign() {
	c.errorSign = true
}

func (c *FakePastelClient) SetErrorVerify() {
	c.errorVerify = true
}

func (c *FakePastelClient) Sign(_ context.Context, data []byte, pastelID, _ string, _ string) ([]byte, error) {
	signature := make([]byte, 20)
	rand.Read(signature)
	c.signatures[pastelID] = signature
	c.data[pastelID] = data
	if c.errorSign {
		return nil, fmt.Errorf("signing failed")
	}
	return signature, nil
}

func (c *FakePastelClient) Verify(_ context.Context, data []byte, signature, pastelID string, _ string) (ok bool, err error) {
	ret := true
	if sig, ok := c.signatures[pastelID]; ok {
		ret = ret && bytes.Equal([]byte(signature), sig)
	}

	if d, ok := c.data[pastelID]; ok {
		ret = ret && bytes.Equal([]byte(data), d)
	}

	if c.errorVerify {
		return false, fmt.Errorf("verifying failed")
	}
	return ret, nil
}
func TestSecretConnWithGRPC(t *testing.T) {
	secInfoClient := &alts.SecInfo{
		PastelID:   "client_pastel_id",
		PassPhrase: "client_pass_phrase",
	}

	secInfoServer := &alts.SecInfo{
		PastelID:   "server_pastel_id",
		PassPhrase: "server_pass_phrase",
	}

	secClient := &FakePastelClient{
		signatures: make(map[string][]byte),
		data:       make(map[string][]byte),
	}

	altsTCServer := NewServerCreds(secClient, secInfoServer)
	s := &greeter.Greeter{}
	server := grpc.NewServer(grpc.Creds(altsTCServer))
	greeter.RegisterGreeterServiceServer(server, s)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net listen: %v", err)
	}
	defer ln.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// serve the incoming connections
		server.Serve(ln)
	}()

	altsTCClient := NewClientCreds(secClient, secInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

func TestErrorSignInHandshake(t *testing.T) {
	secInfoClient := &alts.SecInfo{
		PastelID:   "client_pastel_id",
		PassPhrase: "client_pass_phrase",
	}

	secInfoServer := &alts.SecInfo{
		PastelID:   "server_pastel_id",
		PassPhrase: "server_pass_phrase",
	}

	secClient := &FakePastelClient{
		signatures: make(map[string][]byte),
		data:       make(map[string][]byte),
	}
	secClient.SetErrorSign()

	altsTCServer := NewServerCreds(secClient, secInfoServer)
	s := &greeter.Greeter{}
	server := grpc.NewServer(grpc.Creds(altsTCServer))
	greeter.RegisterGreeterServiceServer(server, s)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net listen: %v", err)
	}
	defer ln.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// serve the incoming connections
		server.Serve(ln)
	}()

	altsTCClient := NewClientCreds(secClient, secInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(altsTCClient),
	)

	assert.True(t, err != nil)
	assert.True(t, conn == nil)

	server.GracefulStop()
	wg.Wait()
}

func TestErrorVerifyInHandshake(t *testing.T) {
	secInfoClient := &alts.SecInfo{
		PastelID:   "client_pastel_id",
		PassPhrase: "client_pass_phrase",
	}

	secInfoServer := &alts.SecInfo{
		PastelID:   "server_pastel_id",
		PassPhrase: "server_pass_phrase",
	}

	secClient := &FakePastelClient{
		signatures: make(map[string][]byte),
		data:       make(map[string][]byte),
	}
	secClient.SetErrorVerify()

	altsTCServer := NewServerCreds(secClient, secInfoServer)
	s := &greeter.Greeter{}
	server := grpc.NewServer(grpc.Creds(altsTCServer))
	greeter.RegisterGreeterServiceServer(server, s)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net listen: %v", err)
	}
	defer ln.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// serve the incoming connections
		server.Serve(ln)
	}()

	altsTCClient := NewClientCreds(secClient, secInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(altsTCClient),
	)

	assert.True(t, err != nil)
	assert.True(t, conn == nil)

	server.GracefulStop()
	wg.Wait()
}

func TestMismatchTypeGrpcServer(t *testing.T) {
	secInfoClient := &alts.SecInfo{
		PastelID:   "client_pastel_id",
		PassPhrase: "client_pass_phrase",
	}

	client := &FakePastelClient{
		signatures: make(map[string][]byte),
		data:       make(map[string][]byte),
	}

	s := &greeter.Greeter{}
	server := grpc.NewServer()
	greeter.RegisterGreeterServiceServer(server, s)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net listen: %v", err)
	}
	defer ln.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// serve the incoming connections
		server.Serve(ln)
	}()

	altsTCClient := NewClientCreds(client, secInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(altsTCClient),
	)
	assert.True(t, err != nil)
	assert.True(t, conn == nil)
	server.GracefulStop()
	wg.Wait()
}

func TestMismatchTypeGrpcClient(t *testing.T) {
	secInfoServer := &alts.SecInfo{
		PastelID:   "server_pastel_id",
		PassPhrase: "server_pass_phrase",
	}

	secClient := &FakePastelClient{
		signatures: make(map[string][]byte),
		data:       make(map[string][]byte),
	}

	altsTCServer := NewServerCreds(secClient, secInfoServer)
	s := &greeter.Greeter{}
	server := grpc.NewServer(grpc.Creds(altsTCServer))
	greeter.RegisterGreeterServiceServer(server, s)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net listen: %v", err)
	}
	defer ln.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		// serve the incoming connections
		server.Serve(ln)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
	)
	assert.True(t, err != nil)
	assert.True(t, conn == nil)
	server.GracefulStop()
	wg.Wait()
}
