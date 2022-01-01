package storagechallenge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

var helperHost = "192.168.100.53"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	client := pool.Client

	networkResource, err := pool.CreateNetwork("pastel-integration-test-net", func(config *docker.CreateNetworkOptions) {
		config.CheckDuplicate = true
		config.Driver = "bridge"
		config.IPAM = &docker.IPAMOptions{
			Driver: "default",
			Config: []docker.IPAMConfig{
				{
					Subnet:  "192.168.0.0/16",
					Gateway: "192.168.100.1",
				},
			},
		}
	})
	if err != nil {
		log.Fatalf("Could not run helper container: %s", err)
	}
	defer networkResource.Close()
	network := networkResource.Network

	helper, removeHelper, err := createContainer(client, "helper", "helper:latest", network, []string{}, helperHost)
	if err != nil {
		log.Fatalf("Could not run helper container: %s", err)
	}
	defer removeHelper()

	// wait for helper container (which is containing mocking stuff as a simple ways for p2p and storage challenge challenge statictis status) is started first
	if err := pool.Retry(func() error {
		var err error
		_, err = http.DefaultClient.Get(fmt.Sprintf("http://%s:8088", helper.NetworkSettings.IPAddress))
		return err
	}); err != nil {
		log.Fatalf("Could not run helper container: %s", err)
	}

	node1, remove1, err := createContainer(client, "mn1key", "storage-challenge:latest", network, []string{"SERVICE_HOST=192.168.100.11", "SERVICE_ID=mn1key"}, "192.168.100.11")
	if err != nil {
		log.Fatalf("Could not run node container 1 %s", err)
	}
	defer remove1()
	node2, remove2, err := createContainer(client, "mn2key", "storage-challenge:latest", network, []string{"SERVICE_HOST=192.168.100.12", "SERVICE_ID=mn2key"}, "192.168.100.12")
	if err != nil {
		log.Fatalf("Could not run node container 2: %s", err)
	}
	defer remove2()
	node3, remove3, err := createContainer(client, "mn3key", "storage-challenge:latest", network, []string{"SERVICE_HOST=192.168.100.13", "SERVICE_ID=mn3key"}, "192.168.100.13")
	if err != nil {
		log.Fatalf("Could not run node container 3: %s", err)
	}
	defer remove3()
	node4, remove4, err := createContainer(client, "mn4key", "storage-challenge:latest", network, []string{"SERVICE_HOST=192.168.100.14", "SERVICE_ID=mn4key"}, "192.168.100.14")
	if err != nil {
		log.Fatalf("Could not run node container 4: %s", err)
	}
	defer remove4()
	node5, remove5, err := createContainer(client, "mn5key", "storage-challenge:latest", network, []string{"SERVICE_HOST=192.168.100.15", "SERVICE_ID=mn5key"}, "192.168.100.15")
	if err != nil {
		log.Fatalf("Could not run node container 5: %s", err)
	}
	defer remove5()

	time.Sleep(5 * time.Second)

	// wait for download and prepare 500 images mock files to p2p network successfully
	if err := pool.Retry(func() error {
		var err error
		_, err = http.DefaultClient.Get(fmt.Sprintf("http://%s:8088/store/mocks", helper.NetworkSettings.IPAddress))
		return err
	}); err != nil {
		log.Fatalf("Could not run helper container: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if errs := purgeDockerResources(client, helper, node1, node2, node3, node4, node5); len(errs) > 0 {
		log.Fatalf("Could not purge resources: %v", errs)
	}

	os.Exit(code)
}

func purgeDockerResources(client *docker.Client, containers ...*docker.Container) []error {
	var errs = make([]error, 0)
	for _, container := range containers {
		if err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            container.ID,
			Force:         true,
			RemoveVolumes: true,
		}); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func createContainer(client *docker.Client, name string, image string, network *docker.Network, env []string, ipAddress string) (container *docker.Container, remove func(), err error) {
	container, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: image,
		},
		HostConfig: &docker.HostConfig{
			AutoRemove:    true,
			RestartPolicy: docker.RestartPolicy{Name: "no"},
		},
		NetworkingConfig: &docker.NetworkingConfig{
			EndpointsConfig: map[string]*docker.EndpointConfig{
				network.Name: {
					Aliases:   []string{name},
					IPAddress: ipAddress,
				},
			},
		},
	})
	if err != nil {
		return
	}

	remove = func() {
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            container.ID,
			Force:         true,
			RemoveVolumes: true,
		})
	}

	return
}

type mockSecClient struct{}

func (*mockSecClient) Sign(_ context.Context, _ []byte, _, _, _ string) ([]byte, error) {
	return []byte("mock signature"), nil
}

func (*mockSecClient) Verify(_ context.Context, _ []byte, signature, _, _ string) (bool, error) {
	return signature == "mock signature", nil
}

type challengeStatusResult struct {
	// number of challenge send to a challenging masternode from challenger masternode
	NumberOfChallengeSent int `json:"sent"`
	// number of challenge which is challenging masternode respond to verifying masternode
	NumberOfChallengeRespondedTo int `json:"respond"`
	// number of challenge which is challenging masternode respond to verifying masternode and verifying masternode successfully verified
	NumverOfChallengeSuccess int `json:"success"`
	// number of challenge which is challenging masternode respond to verifying masternode and verifying masternode says incorrect hash response
	NumverOfChallengeFailed int `json:"failed"`
	// number of challenge which is challenging masternode respond to verifying masternode and verifying masternode says the process take too long
	NumverOfChallengeTimeout int `json:"timeout"`
}

type testSuite struct {
	suite.Suite
	conn       *grpc.ClientConn
	gRPCClient pb.StorageChallengeClient
	httpClient *http.Client
}

func (suite *testSuite) SetupTest() {
	conn, err := grpc.Dial("192.168.100.11:14444",
		grpc.WithTransportCredentials(credentials.NewClientCreds(&mockSecClient{}, &alts.SecInfo{PastelID: "test-pastel-id", PassPhrase: "test-passphrase"})),
		grpc.WithBlock(),
	)
	suite.Require().NoError(err, "cannot setup grpc client connection")
	suite.conn = conn
	suite.gRPCClient = pb.NewStorageChallengeClient(suite.conn)
	suite.httpClient = http.DefaultClient
}

func (suite *testSuite) TearDownTest() {
	suite.conn.Close()
}
func (suite *testSuite) TestFullChallengeSucceeded() {
	for i := 0; i < 100; i++ {
		_, err := suite.gRPCClient.GenerateStorageChallenges(context.Background(), &pb.GenerateStorageChallengesRequest{
			ChallengesPerMasternodePerBlock: 1,
		})
		suite.Require().NoErrorf(err, "failed to call gRPC generate storage challenge '%d'", i)
	}
	time.Sleep(20 * time.Second)

	// collect status of each masternode process storage challenge result
	res, err := suite.httpClient.Get(fmt.Sprintf("http://%s:8088", helperHost))
	suite.Require().NoError(err, "failed to call http request to retrive storage challenge status")
	defer res.Body.Close()
	suite.Require().Equal(http.StatusOK, res.StatusCode)
	var outputChalengeStatusResult map[string]*challengeStatusResult
	err = json.NewDecoder(res.Body).Decode(&outputChalengeStatusResult)
	suite.Require().NoError(err, "failed to decode storage challenge status response json")

	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumverOfChallengeSuccess+challengeStatus.NumverOfChallengeFailed+challengeStatus.NumverOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' and challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumverOfChallengeSuccess, "comparing challenges 'sent' and challenges 'succeeded' of node %s", challengingMasternodeID)
	}
}

func (suite *testSuite) TestDishonestNodeByFailedOrTimeout() {
	// implement the test here
}

func (suite *testSuite) TestIncrementalRQSymbolFileKey() {
	// implement the test here
}

func (suite *testSuite) TestIncrementalMasternode() {
	// implement the test here
}

func TestStorageChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
