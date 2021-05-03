package encryption

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/encryption/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net"
	"strings"
	"testing"
	"time"
)

type serverMock struct{
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
	grpcServer.Serve(listener)
	time.Sleep(1 * time.Second)
}

func TestClientAndServer(t *testing.T) {
	clientCipher := NewX448Cipher()
	serverCipher := NewX448Cipher()
	handshaker := NewHandshakerMock(clientCipher, serverCipher)
	transport := NewTransportCredentials(handshaker)
	go makeTestServer(t, transport)
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
	if response.Message != "egassem"{
		t.Fatalf("unexpected response from server %s", request.Message)
	}
}

type WriterMock struct {
	buff []byte
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	w.buff = append(w.buff, p...)
	return len(p), nil
}

func TestEncrypter(t *testing.T) {
	defKey := []byte("testsomekeytest1")
	//vi := []byte("smstsomekeytest1")
	w := new(WriterMock)
	encrypter, err := NewStreamEncrypter(defKey, defKey, w)
	if err != nil{
		t.Fatal(err)
	}
	testData := []byte("somebuffer")
	_, err = encrypter.Write(testData)
	if err != nil{
		t.Fatal(err)
	}
	//b, err := ioutil.ReadAll(strings.NewReader(string(w.buff)))
	//if err != nil{
	//	t.Fatal(err)
	//}
	sd, err := NewStreamDecrypter(defKey, defKey, strings.NewReader(string(w.buff)))
	if err != nil{
		t.Fatal(err)
	}
	d, err := ioutil.ReadAll(sd)
	if err != nil{
		t.Fatal(err)
	}
	if string(d) != string(testData){
		t.Fatalf("value is not expected %s %s", d, testData)
	}
}