package selfhealing

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"strconv"
	"testing"

	"github.com/DataDog/zstd"
	fuzz "github.com/google/gofuzz"
	"github.com/pastelnetwork/gonode/common/b85"
	"github.com/pastelnetwork/gonode/common/utils"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqmock "github.com/pastelnetwork/gonode/raptorq/node/test"
	shtest "github.com/pastelnetwork/gonode/supernode/node/test/self_healing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessSelfHealingTest(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config := NewConfig()
	pastelClient := pastelMock.NewMockClient(t)
	p2pClient := p2pMock.NewMockClient(t)
	raptorQClient := rqmock.NewMockClient(t)
	var nodeClient *shtest.Client
	f := fuzz.New()

	file := []byte("test file")
	rqDecode := &rqnode.Decode{
		File: file,
	}
	fileHash := sha3.Sum256(file)

	ticket := pastel.RegTicket{}
	f.Fuzz(&ticket)
	ticket.RegTicketData.NFTTicketData.AppTicketData.DataHash = fileHash[:]
	b, err := json.Marshal(ticket.RegTicketData.NFTTicketData.AppTicketData)
	if err != nil {
		t.Fatalf("faied to marshal, err: %s", err)
	}
	ticket.RegTicketData.NFTTicketData.AppTicket = b85.Encode(b)

	b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
	if err != nil {
		t.Fatalf("faied to marshal, err: %s", err)
	}
	ticket.RegTicketData.NFTTicket = b
	tickets := pastel.RegTickets{ticket}

	rqIDsData, _ := fakeRQIDsData()

	nodes := pastel.MasterNodes{}
	nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "C"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "D"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "5072696d6172794944"})

	rqFile := rq.SymbolIDFile{ID: "GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3"}
	bytes, err := json.Marshal(rqFile)
	require.Nil(t, err)

	symbol, err := json.Marshal("24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z")
	require.NoError(t, err)

	var encodeResp rqnode.Encode
	encodeResp.Symbols = make(map[string][]byte)
	encodeResp.Symbols["GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3"] = bytes

	tests := []struct {
		testcase string
		message  *pb.SelfHealingData
		setup    func()
		expect   func(*testing.T, error)
	}{
		{
			testcase: "when symbols do not mismatch, does not proceed self-healing",
			message: &pb.SelfHealingData{
				MessageId:                   "test-message-1",
				MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE,
				MerklerootWhenChallengeSent: "previous-block-hash",
				ChallengingMasternodeId:     "challenging-node",
				RespondingMasternodeId:      "respondong-node",
				ChallengeFile: &pb.SelfHealingDataChallengeFile{
					FileHashToChallenge: "file-hash-to-challenge",
				},
				ChallengeId: "test-challenge-1",
				RegTicketId: "reg-ticket-tx-id",
			},
			setup: func() {
				pastelClient.ListenOnRegTickets(tickets, nil).ListenOnMasterNodesExtra(nodes, nil)
				raptorQClient.ListenOnConnect(nil)
				raptorQClient.ListenOnRaptorQ().ListenOnClose(nil)
				pastelClient.On("RegTicket", mock.Anything, mock.Anything).Return(ticket, nil).Times(1)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(rqIDsData, nil).Times(1)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(symbol, nil).Times(1)
				raptorQClient.ListenOnDone()

			},
			expect: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
		{
			testcase: "when symbols mismatch, should self-heal",
			message: &pb.SelfHealingData{
				MessageId:                   "test-message-1",
				MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE,
				MerklerootWhenChallengeSent: "previous-block-hash",
				ChallengingMasternodeId:     "challenging-node",
				RespondingMasternodeId:      "respondong-node",
				ChallengeFile: &pb.SelfHealingDataChallengeFile{
					FileHashToChallenge: "file-hash-to-challenge",
				},
				ChallengeId: "test-challenge-1",
				RegTicketId: "reg-ticket-tx-id",
			},
			setup: func() {
				pastelClient.ListenOnRegTickets(tickets, nil).ListenOnMasterNodesExtra(nodes, nil)
				raptorQClient.ListenOnConnect(nil)
				raptorQClient.ListenOnRaptorQ().ListenOnClose(nil)
				pastelClient.On("RegTicket", mock.Anything, mock.Anything).Return(ticket, nil).Times(1)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(rqIDsData, nil).Times(1)

				symbol, err := json.Marshal("GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3")
				require.NoError(t, err)

				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(symbol, nil).Times(1)
				raptorQClient.RaptorQ.On(rqmock.DecodeMethod, mock.Anything, mock.Anything).Return(rqDecode, nil)

				nodeClient = shtest.NewMockClient(t)
				nodeClient.ListenOnConnect("", nil).ListenOnSelfHealingChallengeInterface().ListenOnVerifySelfHealingChallengeFunc(&pb.SelfHealingData{ChallengeStatus: pb.SelfHealingData_Status_SUCCEEDED}, nil)

				raptorQClient.RaptorQ.On(rqmock.EncodeMethod, mock.Anything, mock.Anything).Return(&encodeResp, nil)

				p2pClient.ListenOnStore("", nil)
				raptorQClient.ListenOnDone()

			},
			expect: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
		{
			testcase: "when self-healing verification failed, should not store symbols into P2P",
			message: &pb.SelfHealingData{
				MessageId:                   "test-message-1",
				MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE,
				MerklerootWhenChallengeSent: "previous-block-hash",
				ChallengingMasternodeId:     "challenging-node",
				RespondingMasternodeId:      "respondong-node",
				ChallengeFile: &pb.SelfHealingDataChallengeFile{
					FileHashToChallenge: "file-hash-to-challenge",
				},
				ChallengeId: "test-challenge-1",
				RegTicketId: "reg-ticket-tx-id",
			},
			setup: func() {
				pastelClient.ListenOnRegTickets(tickets, nil).ListenOnMasterNodesExtra(nodes, nil)
				raptorQClient.ListenOnConnect(nil)
				raptorQClient.ListenOnRaptorQ().ListenOnClose(nil)
				pastelClient.On("RegTicket", mock.Anything, mock.Anything).Return(ticket, nil).Times(1)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(rqIDsData, nil).Times(1)

				symbol, err := json.Marshal("GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3")
				require.NoError(t, err)

				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(symbol, nil).Times(1)
				raptorQClient.RaptorQ.On(rqmock.DecodeMethod, mock.Anything, mock.Anything).Return(rqDecode, nil)
				nodeClient = shtest.NewMockClient(t)
				nodeClient.ListenOnConnect("", nil).ListenOnSelfHealingChallengeInterface().ListenOnVerifySelfHealingChallengeFunc(&pb.SelfHealingData{ChallengeStatus: pb.SelfHealingData_Status_FAILED_INCORRECT_RESPONSE}, nil)
				raptorQClient.ListenOnDone()
			},
			expect: func(t *testing.T, err error) {
				require.NotNil(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.testcase, func(t *testing.T) {
			// Run the setup for the testcase
			tt.setup()

			service := NewService(config, nil, pastelClient, nodeClient,
				p2pClient, raptorQClient)
			task := NewSHTask(service)
			// call the function to get return values
			err := task.ProcessSelfHealingChallenge(ctx, tt.message)

			// handle the test case's assertions with the provided func
			tt.expect(t, err)
		})
	}

}

func fakeRQIDsData() ([]byte, []string) {
	rqfile := &rqnode.RawSymbolIDFile{
		ID:                "09f6c459-ec2a-4db1-a8fe-0648fd97b5cb",
		PastelID:          "jXY1wJkRFt4hsPn6LnRqUtoRmBx5QTiGcbCXorKq7JuKVy4Zo89PmE8BoGjyujqj6NwfvfGsxhUH2ute6kW2gW",
		SymbolIdentifiers: []string{"GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3"},
	}

	dataJSON, _ := json.Marshal(rqfile)
	encoded := utils.B64Encode(dataJSON)

	var buffer bytes.Buffer
	buffer.Write(encoded)
	buffer.WriteByte(46)
	buffer.WriteString("test-signature")
	buffer.WriteByte(46)
	buffer.WriteString(strconv.Itoa(55))

	compressedData, _ := zstd.CompressLevel(nil, buffer.Bytes(), 22)

	return compressedData, []string{
		"GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3",
	}
}
