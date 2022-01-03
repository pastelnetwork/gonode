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
var dockerClient *docker.Client
var incrementalNode *docker.Container
var dishonestNode *docker.Container

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	dockerClient = pool.Client

	if err = dockerClient.BuildImage(docker.BuildImageOptions{
		Name:         "st-integration",
		Dockerfile:   "run_environment/Dockerfile",
		ContextDir:   ".",
		OutputStream: os.Stdout,
	}); err != nil {
		log.Fatalf("Could not build storage challenge image: %v", err)
	}

	if err = dockerClient.BuildImage(docker.BuildImageOptions{
		Name:         "helper",
		Dockerfile:   "test_helper/Dockerfile",
		ContextDir:   ".",
		OutputStream: os.Stdout,
	}); err != nil {
		log.Fatalf("Could not build helper image: %v", err)
	}

	networkResource, err := pool.CreateNetwork("st-integration", func(config *docker.CreateNetworkOptions) {
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
		log.Fatalf("Could not create network: %v", err)
	}
	defer networkResource.Close()
	network := networkResource.Network

	helper, _, err := createContainer(dockerClient, "helper", "helper:latest", network, []string{}, helperHost, false)
	if err != nil {
		log.Fatalf("Could not run helper container: %v", err)
	}

	// wait for helper container (which is containing mocking stuff as a simple ways for p2p and storage challenge challenge statictis status) is started first
	if err := pool.Retry(func() error {
		var err error
		_, err = http.DefaultClient.Get(fmt.Sprintf("http://%s:8088/getblockcount", helper.NetworkSettings.Networks[network.Name].IPAddress))
		if err != nil {
			log.Println("health check helper failed:", err)
		}
		return err
	}); err != nil {
		log.Fatalf("Could not wait for helper container successfully running: %v", err)
	}

	node1, remove1, err := createContainer(dockerClient, "mn1key", "st-integration:latest", network, []string{"SERVICE_HOST=192.168.100.11", "SERVICE_ID=mn1key"}, "192.168.100.11", false)
	defer remove1()
	if err != nil {
		log.Fatalf("Could not run node container 1 %v", err)
	}
	node2, remove2, err := createContainer(dockerClient, "mn2key", "st-integration:latest", network, []string{"SERVICE_HOST=192.168.100.12", "SERVICE_ID=mn2key"}, "192.168.100.12", false)
	defer remove2()
	if err != nil {
		log.Fatalf("Could not run node container 2: %v", err)
	}
	node3, remove3, err := createContainer(dockerClient, "mn3key", "st-integration:latest", network, []string{"SERVICE_HOST=192.168.100.13", "SERVICE_ID=mn3key"}, "192.168.100.13", false)
	defer remove3()
	if err != nil {
		log.Fatalf("Could not run node container 3: %v", err)
	}
	node4, remove4, err := createContainer(dockerClient, "mn4key", "st-integration:latest", network, []string{"SERVICE_HOST=192.168.100.14", "SERVICE_ID=mn4key"}, "192.168.100.14", false)
	defer remove4()
	if err != nil {
		log.Fatalf("Could not run node container 4: %v", err)
	}
	// use node5 as incremental node
	node5, remove5, err := createContainer(dockerClient, "mn5key", "st-integration:latest", network, []string{"SERVICE_HOST=192.168.100.15", "SERVICE_ID=mn5key"}, "192.168.100.15", true)
	defer remove5()
	if err != nil {
		log.Fatalf("Could not run node container 5: %v", err)
	}
	// use node6 as dishonest node
	node6, remove6, err := createContainer(dockerClient, "mn6key", "st-integration:latest", network, []string{"SERVICE_HOST=192.168.100.16", "SERVICE_ID=mn6key"}, "192.168.100.16", true)
	defer remove6()
	if err != nil {
		log.Fatalf("Could not run node container 5: %v", err)
	}
	time.Sleep(5 * time.Second)

	incrementalNode = node5
	dishonestNode = node6

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if errs := purgeDockerResources(dockerClient, helper, node1, node2, node3, node4, node5); len(errs) > 0 {
		log.Printf("Could not purge resources: %v", errs)
	}

	networkResource.Close()

	if code != 0 {
		log.Fatal("Storage challenge failed integration test")
	}
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

func createContainer(client *docker.Client, name string, image string, network *docker.Network, env []string, ipAddress string, noStart bool) (container *docker.Container, remove func(), err error) {
	remove = func() {}
	container, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: image,
			Env:   env,
		},
		HostConfig: &docker.HostConfig{
			AutoRemove:    true,
			RestartPolicy: docker.RestartPolicy{Name: "no"},
		},
		NetworkingConfig: &docker.NetworkingConfig{
			EndpointsConfig: map[string]*docker.EndpointConfig{
				network.Name: {
					Aliases:   []string{name},
					NetworkID: network.ID,
					IPAMConfig: &docker.EndpointIPAMConfig{
						IPv4Address: ipAddress,
					},
					IPAddress:   ipAddress,
					IPPrefixLen: 16,
					Gateway:     "192.168.100.1",
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

	if noStart {
		return
	}

	if err = client.StartContainer(container.ID, nil); err != nil {
		return
	}

	remove = func() {
		client.StopContainer(container.ID, 30)
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            container.ID,
			Force:         true,
			RemoveVolumes: true,
		})
	}

	container, err = client.InspectContainer(container.ID)

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
	conns        []*grpc.ClientConn
	gRPCClients  []pb.StorageChallengeClient
	httpClient   *http.Client
	dockerClient *docker.Client
}

func (suite *testSuite) SetupTest() {
	conn1, err := grpc.Dial("192.168.100.11:14444",
		grpc.WithTransportCredentials(credentials.NewClientCreds(&mockSecClient{}, &alts.SecInfo{PastelID: "test-pastel-id", PassPhrase: "test-passphrase"})),
		grpc.WithBlock(),
	)
	suite.Require().NoError(err, "cannot setup grpc client connection to node1")
	conn2, err := grpc.Dial("192.168.100.12:14444",
		grpc.WithTransportCredentials(credentials.NewClientCreds(&mockSecClient{}, &alts.SecInfo{PastelID: "test-pastel-id", PassPhrase: "test-passphrase"})),
		grpc.WithBlock(),
	)
	suite.Require().NoError(err, "cannot setup grpc client connection to node2")
	suite.conns = append(suite.conns, conn1, conn2)
	suite.gRPCClients = append(suite.gRPCClients, pb.NewStorageChallengeClient(conn1), pb.NewStorageChallengeClient(conn2))
	suite.httpClient = http.DefaultClient
	suite.dockerClient = dockerClient
}

func (suite *testSuite) TearDownTest() {
	(&http.Client{}).Get("http://192.168.100.53:8088/sts/reset")
	for _, conn := range suite.conns {
		conn.Close()
	}
}
func (suite *testSuite) TestFullChallengeSucceeded() {
	// collect status of each masternode process storage challenge result
	res, err := suite.httpClient.Get(fmt.Sprintf("http://%s:8088/sts/show", helperHost))
	suite.Require().NoError(err, "failed to call http request to retrive storage challenge status")
	defer res.Body.Close()
	suite.Require().Equal(http.StatusOK, res.StatusCode)
	var outputChalengeStatusResult map[string]*challengeStatusResult
	err = json.NewDecoder(res.Body).Decode(&outputChalengeStatusResult)
	suite.Require().NoError(err, "failed to decode storage challenge status response json")

	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		suite.Require().Greaterf(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumverOfChallengeSuccess+challengeStatus.NumverOfChallengeFailed+challengeStatus.NumverOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' and challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumverOfChallengeSuccess, "comparing challenges 'sent' and challenges 'succeeded' of node %s", challengingMasternodeID)
	}
}

func (suite *testSuite) TestDishonestNodeByFailedOrTimeout() {
	err := suite.dockerClient.StartContainer(dishonestNode.ID, nil)
	suite.Require().NoError(err, "failed start dishonest node")
	defer suite.dockerClient.StopContainer(dishonestNode.ID, 30)

	_, err = http.NewRequest("DELETE", fmt.Sprintf("http://%s:9090/p2p/0.2", "192.168.100.16"), nil)
	if err != nil {
		log.Fatalf("could not delete 20%% file hash to become dishonest node: %v", err)
	}
	time.Sleep(10 * time.Second)

	// collect status of each masternode process storage challenge result
	res, err := suite.httpClient.Get(fmt.Sprintf("http://%s:8088/sts/show", helperHost))
	suite.Require().NoError(err, "failed to call http request to retrive storage challenge status")
	defer res.Body.Close()
	suite.Require().Equal(http.StatusOK, res.StatusCode)
	var outputChalengeStatusResult map[string]*challengeStatusResult
	err = json.NewDecoder(res.Body).Decode(&outputChalengeStatusResult)
	suite.Require().NoError(err, "failed to decode storage challenge status response json")

	sumaryChallengesSent := 0
	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		sumaryChallengesSent += challengeStatus.NumberOfChallengeSent
		suite.Require().Greater(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumverOfChallengeSuccess+challengeStatus.NumverOfChallengeFailed+challengeStatus.NumverOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' and challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumverOfChallengeSuccess, "comparing challenges 'sent' greater than challenges 'succeeded' of node %s", challengingMasternodeID)
		suite.Require().Greaterf(challengeStatus.NumverOfChallengeTimeout, 0, "comparing challenges 'timeout' greater than 0 of node %s", challengingMasternodeID)
	}
}

func (suite *testSuite) TestIncrementalRQSymbolFileKey() {
	// // insert incrementals files
	// _, err := (&http.Client{}).Get("http://192.168.100.53:8088/store/incrementals")
	// suite.Require().NoError(err, "failed to call http request to add incrementals symnbol file hashes")
	time.Sleep(10 * time.Second)

	// collect status of each masternode process storage challenge result
	res, err := suite.httpClient.Get(fmt.Sprintf("http://%s:8088/sts/show", helperHost))
	suite.Require().NoError(err, "failed to call http request to retrive storage challenge status")
	defer res.Body.Close()
	suite.Require().Equal(http.StatusOK, res.StatusCode)
	var outputChalengeStatusResult map[string]*challengeStatusResult
	err = json.NewDecoder(res.Body).Decode(&outputChalengeStatusResult)
	suite.Require().NoError(err, "failed to decode storage challenge status response json")

	sumaryChallengesSent := 0
	log.Println(outputChalengeStatusResult)
	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		sumaryChallengesSent += challengeStatus.NumberOfChallengeSent
		suite.Require().Greater(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumverOfChallengeSuccess+challengeStatus.NumverOfChallengeFailed+challengeStatus.NumverOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' and challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumverOfChallengeSuccess, "comparing challenges 'sent' greater than challenges 'succeeded' of node %s", challengingMasternodeID)
	}
}

func (suite *testSuite) TestIncrementalMasternode() {
	err := suite.dockerClient.StartContainer(incrementalNode.ID, nil)
	suite.Require().NoError(err, "failed start dishonest node")
	defer suite.dockerClient.StopContainer(incrementalNode.ID, 30)

	// // wait for download and prepare images mock files to p2p network successfully
	// _, err = http.DefaultClient.Get(fmt.Sprintf("http://%s:8088/store/mocks", "192.168.100.53"))
	// if err != nil {
	// 	log.Fatalf("generate mock image failed from helper image: %v", err)
	// }

	time.Sleep(10 * time.Second)

	// collect status of each masternode process storage challenge result
	res, err := suite.httpClient.Get(fmt.Sprintf("http://%s:8088/sts/show", helperHost))
	suite.Require().NoError(err, "failed to call http request to retrive storage challenge status")
	defer res.Body.Close()
	suite.Require().Equal(http.StatusOK, res.StatusCode)
	var outputChalengeStatusResult map[string]*challengeStatusResult
	err = json.NewDecoder(res.Body).Decode(&outputChalengeStatusResult)
	suite.Require().NoError(err, "failed to decode storage challenge status response json")

	sumaryChallengesSent := 0
	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		sumaryChallengesSent += challengeStatus.NumberOfChallengeSent
		suite.Require().Greater(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumverOfChallengeSuccess+challengeStatus.NumverOfChallengeFailed+challengeStatus.NumverOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' and challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumverOfChallengeSuccess, "comparing challenges 'sent' greater than challenges 'succeeded' of node %s", challengingMasternodeID)
	}
}

func TestStorageChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
