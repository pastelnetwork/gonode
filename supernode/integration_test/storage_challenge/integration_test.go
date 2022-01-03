package storagechallenge

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pastelnetwork/gonode/common/net/credentials"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/utils"
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
	suite.conns = append(suite.conns, conn1)
	suite.gRPCClients = append(suite.gRPCClients, pb.NewStorageChallengeClient(conn1))
	suite.httpClient = http.DefaultClient
	suite.dockerClient = dockerClient
}

func (suite *testSuite) TearDownTest() {
	(&http.Client{}).Get("http://192.168.100.53:8088/sts/reset")
	for _, conn := range suite.conns {
		conn.Close()
	}
	suite.gRPCClients = make([]pb.StorageChallengeClient, 0)
}

func (suite *testSuite) TestFullChallengeSucceeded() {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		for _, client := range suite.gRPCClients {
			client.GenerateStorageChallenges(ctx, &pb.GenerateStorageChallengesRequest{ChallengesPerMasternodePerBlock: 1})
		}
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
	suite.T().Logf("challengeStatus %+v", outputChalengeStatusResult)

	sumaryChallengesSent := 0
	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		sumaryChallengesSent += challengeStatus.NumberOfChallengeSent
		suite.Require().Greaterf(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumberOfChallengeSuccess+challengeStatus.NumberOfChallengeFailed+challengeStatus.NumberOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' equals to challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeSuccess, "comparing challenges 'sent' equals to challenges 'succeeded' of node %s", challengingMasternodeID)
	}
	suite.Require().Equalf(10, sumaryChallengesSent, "sumary of challenge sent must be 10")
}

func (suite *testSuite) TestDishonestNodeByFailedOrTimeout() {
	err := suite.dockerClient.StartContainer(dishonestNode.ID, nil)
	suite.Require().NoError(err, "failed start dishonest node")
	defer suite.dockerClient.StopContainer(dishonestNode.ID, 30)

	_, err = removeKeys(suite.httpClient, "192.168.100.16")
	suite.Require().NoError(err, "failed remove 20% of files to prepare dishonest node")
	time.Sleep(3 * time.Second)

	verifyNodeID := base32.StdEncoding.EncodeToString([]byte("mn1key"))
	challengingNodeID := base32.StdEncoding.EncodeToString([]byte("mn6key"))

	conn, err := grpc.Dial("192.168.100.16:14444",
		grpc.WithTransportCredentials(credentials.NewClientCreds(&mockSecClient{}, &alts.SecInfo{PastelID: challengingNodeID, PassPhrase: "mock passphrase"})),
		grpc.WithBlock(),
	)
	suite.Require().NoError(err, "cannot setup grpc client connection to node1")
	suite.conns = append(suite.conns, conn)
	grpcClient := pb.NewStorageChallengeClient(conn)
	suite.gRPCClients = append(suite.gRPCClients, grpcClient)

	keylist, err := getKeyList(suite.httpClient)
	suite.Require().NoError(err, "failed retrive list of symbol file hash node")

	blkHeight, err := getCurrentBlock(suite.httpClient)
	suite.Require().NoError(err, "failed get current block")

	ctx := context.Background()
	for idx, fileHash := range keylist {
		challengeID := utils.GetHashFromString(fmt.Sprintf("%f%d%s%s", rand.Float64(), idx, fileHash, "mock-challenge-id"))
		err = prepareChallengeStatus(suite.httpClient, challengeID, challengingNodeID, fmt.Sprint(blkHeight))
		suite.Require().NoError(err, "could not prepare challenge status")
		grpcClient.ProcessStorageChallenge(ctx, &pb.ProcessStorageChallengeRequest{
			Data: &pb.StorageChallengeData{
				MessageId:                   utils.GetHashFromString(fmt.Sprintf("%f%d%s%s", rand.Float64(), idx, fileHash, "mock-message-id")),
				MessageType:                 pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE,
				ChallengeStatus:             pb.StorageChallengeData_Status_PENDING,
				MerklerootWhenChallengeSent: utils.GetHashFromString(fmt.Sprintf("%f%d%s%s", rand.Float64(), idx, fileHash, "mock-merkle-root")),
				BlockNumChallengeSent:       blkHeight,
				ChallengingMasternodeId:     challengingNodeID,
				RespondingMasternodeId:      verifyNodeID,
				ChallengeFile: &pb.StorageChallengeDataChallengeFile{
					FileHashToChallenge:      fileHash,
					ChallengeSliceStartIndex: 0,
					ChallengeSliceEndIndex:   20,
				},
				ChallengeId: challengeID,
			},
		})
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

	challengeStatus, ok := outputChalengeStatusResult[challengingNodeID]
	suite.T().Logf("challengeStatus %+v", challengeStatus)
	suite.Require().True(ok, "challenge status of dishonest node must be availabe")

	suite.Require().Greater(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
	suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeSuccess+challengeStatus.NumberOfChallengeTimeout, "comparing challenges 'sent' and sumary of challenges 'success, timeout' of node %s", challengingNodeID)
	suite.Require().Greater(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeSuccess, "comparing challenges 'sent' greater than challenges 'succeeded' of node %s", challengingNodeID)
	suite.Require().Greater(challengeStatus.NumberOfChallengeTimeout, 0, "challenge timeout from must be greater than 0")
}

func (suite *testSuite) TestIncrementalRQSymbolFileKey() {
	// insert incrementals files
	_, err := suite.httpClient.Get("http://192.168.100.53:8088/store/incrementals")
	suite.Require().NoError(err, "failed to call http request to add incrementals symnbol file hashes")

	keylist, err := getKeyList(suite.httpClient)
	suite.Require().NoError(err, "failed retrive list of symbol file hash node")
	suite.Require().Equal(14, len(keylist), "key list increasing from 10 to 14")
	ctx := context.Background()
	for i := 0; i < 14; i++ {
		for _, client := range suite.gRPCClients {
			client.GenerateStorageChallenges(ctx, &pb.GenerateStorageChallengesRequest{ChallengesPerMasternodePerBlock: 1})
		}
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
	suite.T().Logf("challengeStatus %+v", outputChalengeStatusResult)

	sumaryChallengesSent := 0
	for challengingMasternodeID, challengeStatus := range outputChalengeStatusResult {
		sumaryChallengesSent += challengeStatus.NumberOfChallengeSent
		suite.Require().Greater(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
		suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumberOfChallengeSuccess+challengeStatus.NumberOfChallengeFailed+challengeStatus.NumberOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' equals to challenges 'responded' of node %s", challengingMasternodeID)
		suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeSuccess, "comparing challenges 'sent' equals to challenges 'succeeded' of node %s", challengingMasternodeID)
	}
	suite.Require().Equalf(14, sumaryChallengesSent, "sumary of challenge sent must be 14")
}

func (suite *testSuite) TestIncrementalMasternode() {
	err := suite.dockerClient.StartContainer(incrementalNode.ID, nil)
	suite.Require().NoError(err, "failed start incremental node")
	defer suite.dockerClient.StopContainer(incrementalNode.ID, 30)

	verifyNodeID := base32.StdEncoding.EncodeToString([]byte("mn1key"))
	challengingNodeID := base32.StdEncoding.EncodeToString([]byte("mn5key"))

	conn, err := grpc.Dial("192.168.100.15:14444",
		grpc.WithTransportCredentials(credentials.NewClientCreds(&mockSecClient{}, &alts.SecInfo{PastelID: challengingNodeID, PassPhrase: "mock passphrase"})),
		grpc.WithBlock(),
	)
	suite.Require().NoError(err, "cannot setup grpc client connection to node1")
	suite.conns = append(suite.conns, conn)
	grpcClient := pb.NewStorageChallengeClient(conn)
	suite.gRPCClients = append(suite.gRPCClients, grpcClient)

	keylist, err := getKeyList(suite.httpClient)
	suite.Require().NoError(err, "failed retrive list of symbol file hash node")

	blkHeight, err := getCurrentBlock(suite.httpClient)
	suite.Require().NoError(err, "failed get current block")

	ctx := context.Background()
	for idx, fileHash := range keylist {
		challengeID := utils.GetHashFromString(fmt.Sprintf("%f%d%s%s", rand.Float64(), idx, fileHash, "mock-challenge-id"))
		err = prepareChallengeStatus(suite.httpClient, challengeID, challengingNodeID, fmt.Sprint(blkHeight))
		suite.Require().NoError(err, "could not prepare challenge status")
		grpcClient.ProcessStorageChallenge(ctx, &pb.ProcessStorageChallengeRequest{
			Data: &pb.StorageChallengeData{
				MessageId:                   utils.GetHashFromString(fmt.Sprintf("%f%d%s%s", rand.Float64(), idx, fileHash, "mock-message-id")),
				MessageType:                 pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE,
				ChallengeStatus:             pb.StorageChallengeData_Status_PENDING,
				MerklerootWhenChallengeSent: utils.GetHashFromString(fmt.Sprintf("%f%d%s%s", rand.Float64(), idx, fileHash, "mock-merkle-root")),
				BlockNumChallengeSent:       blkHeight,
				ChallengingMasternodeId:     challengingNodeID,
				RespondingMasternodeId:      verifyNodeID,
				ChallengeFile: &pb.StorageChallengeDataChallengeFile{
					FileHashToChallenge:      fileHash,
					ChallengeSliceStartIndex: 0,
					ChallengeSliceEndIndex:   20,
				},
				ChallengeId: challengeID,
			},
		})
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

	challengeStatus, ok := outputChalengeStatusResult[challengingNodeID]
	suite.T().Logf("challengeStatus %+v", challengeStatus)
	suite.Require().True(ok, "challenge status of dishonest node must be availabe")

	sumaryChallengesSent := 0
	sumaryChallengesSent += challengeStatus.NumberOfChallengeSent
	suite.Require().Greater(challengeStatus.NumberOfChallengeSent, 0, "challenge sent must be greater than 0")
	suite.Require().Equalf(challengeStatus.NumberOfChallengeRespondedTo, challengeStatus.NumberOfChallengeSuccess+challengeStatus.NumberOfChallengeFailed+challengeStatus.NumberOfChallengeTimeout, "comparing challenges 'responded' and sumary of challenges 'success, failed and timeout' of node %s", challengingNodeID)
	suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeRespondedTo, "comparing challenges 'sent' and challenges 'responded' of node %s", challengingNodeID)
	suite.Require().Equalf(challengeStatus.NumberOfChallengeSent, challengeStatus.NumberOfChallengeSuccess, "comparing challenges 'sent' equal to than challenges 'succeeded' of node %s", challengingNodeID)
}

func TestStorageChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
