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

func TestSecretConnWithGRPC(t *testing.T) {
	cSignature := make([]byte, 20)
	sSignature := make([]byte, 20)
	rand.Read(cSignature)
	rand.Read(sSignature)
	var sPubKey []byte
	var cPubKey []byte
	signInfoClient := &alts.SignInfo{
		PastelID:   "client_pastel_id",
		PassPhrase: "client_pass_phrase",
	}

	signInfoServer := &alts.SignInfo{
		PastelID:   "server_pastel_id",
		PassPhrase: "server_pass_phrase",
	}

	cSign := func(ctx context.Context, data []byte, pastelID, passPhrase string) ([]byte, error) {
		cPubKey = data
		return cSignature, nil
	}

	cVerify := func(ctx context.Context, data []byte, signature, pastelID string) (ok bool, err error) {
		if bytes.Equal(data, sPubKey) && bytes.Equal([]byte(signature), sSignature) &&
			pastelID == signInfoServer.PastelID {
			return true, nil
		}
		return false, fmt.Errorf("failed to verify client side")
	}

	sSign := func(ctx context.Context, data []byte, pastelID, passPhrase string) ([]byte, error) {
		sPubKey = data
		return sSignature, nil
	}

	sVerify := func(ctx context.Context, data []byte, signature, pastelID string) (ok bool, err error) {
		if bytes.Equal(data, cPubKey) && bytes.Equal([]byte(signature), cSignature) &&
			pastelID == signInfoClient.PastelID {
			return true, nil
		}
		return false, fmt.Errorf("failed to verify server side")
	}

	altsTCServer := NewServerCreds(sSign, sVerify, signInfoServer)
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

	altsTCClient := NewClientCreds(cSign, cVerify, signInfoClient)
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
