//go:generate mockery --name=DDServerClient

package ddclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/pastelnetwork/gonode/dupedetection"
	"google.golang.org/grpc"
)

type testService struct {
	*pb.UnimplementedDupeDetectionServerServer
	retReply *pb.ImageRarenessScoreReply
	retError error
}

func (srv *testService) ImageRarenessScore(_ context.Context, _ *pb.RarenessScoreRequest) (*pb.ImageRarenessScoreReply, error) {
	return srv.retReply, srv.retError
}

func startGrpcServer(_ *testing.T, wg *sync.WaitGroup, port int, serviceHandler *testService) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	srv := grpc.NewServer()
	pb.RegisterDupeDetectionServerServer(srv, serviceHandler)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(lis); err != nil {
			panic(fmt.Sprintf("failed to serve: %v", err))
		}
	}()

	return srv
}

// TestDDClientError
func TestDDClientError(t *testing.T) {
	testErr := errors.New("error")
	testReply := (*pb.ImageRarenessScoreReply)(nil)
	testPort := 50053
	serviceHandler := &testService{
		retError: testErr,
		retReply: testReply,
	}

	wg := &sync.WaitGroup{}
	server := startGrpcServer(t, wg, testPort, serviceHandler)
	defer func() {
		server.Stop()
		wg.Wait()
	}()

	config := &Config{
		Host:       "localhost",
		Port:       testPort,
		DDFilesDir: "/tmp",
	}
	client := NewDDServerClient(config)
	data := []byte{1, 2, 3, 4, 5}

	//ctx context.Context, img []byte, format string, blockHash string, blockHeight int, timestamp string, pastelID string, sn_1 string, sn_2 string, sn_3 string, openapi_request bool, open_api_subset_id string
	_, err := client.ImageRarenessScore(context.Background(), data, "png", "testBlockHash", "240669", "2022-03-31 16:55:28", "testPastelID", "jXYiHNqO9B7psxFQZb1thEgDNykZjL8GkHMZNPZx3iCYre1j3g0zHynlTQ9TdvY6dcRlYIsNfwIQ6nVXBSVJis", "jXpDb5K6S81ghCusMOXLP6k0RvqgFhkBJSFf6OhjEmpvCWGZiptRyRgfQ9cTD709sA58m5czpipFnvpoHuPX0F", "jXS9NIXHj8pd9mLNsP2uKgIh1b3EH2aq5dwupUF7hoaltTE8Zlf6R7Pke0cGr071kxYxqXHQmfVO5dA4jH0ejQ", false, "")
	assert.NotNil(t, err)
}

// TestDDClientSuccess
func TestDDClientSuccess(t *testing.T) {
	testErr := error(nil)
	testReply := &pb.ImageRarenessScoreReply{

		RarenessScores:        &pb.RarenessScores{},
		InternetRareness:      &pb.InternetRareness{},
		AlternativeNsfwScores: &pb.AltNsfwScores{},
		FingerprintsStat:      &pb.FingerprintsStat{},

		PerceptualImageHashes: &pb.PerceptualImageHashes{},

		Maxes:      &pb.Maxes{},
		Percentile: &pb.Percentile{},
	}
	testPort := 50054
	serviceHandler := &testService{
		retError: testErr,
		retReply: testReply,
	}

	wg := &sync.WaitGroup{}
	server := startGrpcServer(t, wg, testPort, serviceHandler)
	defer func() {
		server.Stop()
		wg.Wait()
	}()

	config := &Config{
		Host:       "localhost",
		Port:       testPort,
		DDFilesDir: "/tmp",
	}
	client := NewDDServerClient(config)
	data := []byte{1, 2, 3, 4, 5}

	_, err := client.ImageRarenessScore(context.Background(), data, "png", "testBlockHash", "240669", "2022-03-31 16:55:28", "testPastelID", "jXYiHNqO9B7psxFQZb1thEgDNykZjL8GkHMZNPZx3iCYre1j3g0zHynlTQ9TdvY6dcRlYIsNfwIQ6nVXBSVJis", "jXpDb5K6S81ghCusMOXLP6k0RvqgFhkBJSFf6OhjEmpvCWGZiptRyRgfQ9cTD709sA58m5czpipFnvpoHuPX0F", "jXS9NIXHj8pd9mLNsP2uKgIh1b3EH2aq5dwupUF7hoaltTE8Zlf6R7Pke0cGr071kxYxqXHQmfVO5dA4jH0ejQ", false, "")
	assert.Nil(t, err)
}
