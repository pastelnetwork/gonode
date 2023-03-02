package selfhealing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	fuzz "github.com/google/gofuzz"
	"golang.org/x/crypto/sha3"

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

func TestVerifySelfHealingChallenge(t *testing.T) {
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
	ticket.RegTicketData.NFTTicketData.AppTicket = base64.StdEncoding.EncodeToString(b)

	b, err = json.Marshal(ticket.RegTicketData.NFTTicketData)
	if err != nil {
		t.Fatalf("faied to marshal, err: %s", err)
	}
	ticket.RegTicketData.NFTTicket = b

	rqIDsData, _ := fakeRQIDsData()

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
		expect   func(*testing.T, *pb.SelfHealingData, error)
	}{
		{
			testcase: "when reconstruction is not required, should fail verification",
			message: &pb.SelfHealingData{
				MessageId:                   "test-message-1",
				MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE,
				MerklerootWhenChallengeSent: "previous-block-hash",
				ChallengingMasternodeId:     "challenging-node",
				RespondingMasternodeId:      "responding-node",
				ChallengeFile: &pb.SelfHealingDataChallengeFile{
					FileHashToChallenge: "file-hash-to-challenge",
				},
				ChallengeId: "test-challenge-1",
				RegTicketId: "reg-ticket-tx-id",
			},
			setup: func() {
				pastelClient.On("RegTicket", mock.Anything, mock.Anything).Return(ticket, nil)
				raptorQClient.ListenOnConnect(nil)
				raptorQClient.ListenOnRaptorQ().ListenOnClose(nil)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(rqIDsData, nil).Times(1)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(symbol, nil).Times(1)
				raptorQClient.ListenOnDone()
			},
			expect: func(t *testing.T, data *pb.SelfHealingData, err error) {
				require.Nil(t, err)

				require.Equal(t, data.MessageType, pb.SelfHealingData_MessageType_SELF_HEALING_RESPONSE_MESSAGE)
				require.Equal(t, data.RespondingMasternodeId, "challenging-node")
			},
		},
		{
			testcase: "when challenging node file hash matches with the verifying node generated file hash, should be successful",
			message: &pb.SelfHealingData{
				MessageId:                   "test-message-1",
				MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_VERIFICATION_MESSAGE,
				MerklerootWhenChallengeSent: "previous-block-hash",
				ChallengingMasternodeId:     "challenging-node",
				RespondingMasternodeId:      "responding-node",
				ChallengeFile: &pb.SelfHealingDataChallengeFile{
					FileHashToChallenge: "file-hash-to-challenge",
				},
				ChallengeId:           "test-challenge-1",
				RegTicketId:           "reg-ticket-tx-id",
				ReconstructedFileHash: fileHash[:],
			},
			setup: func() {
				pastelClient.ListenOnRegTicket(mock.Anything, ticket, nil)
				raptorQClient.ListenOnConnect(nil)
				raptorQClient.ListenOnRaptorQ().ListenOnClose(nil)
				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(rqIDsData, nil).Times(1)

				symbol, err := json.Marshal("GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3")
				require.NoError(t, err)

				p2pClient.On(p2pMock.RetrieveMethod, mock.Anything, mock.Anything, mock.Anything).Return(symbol, nil).Times(1)
				raptorQClient.RaptorQ.On(rqmock.DecodeMethod, mock.Anything, mock.Anything).Return(rqDecode, nil)

				raptorQClient.ListenOnDone()
			},
			expect: func(t *testing.T, data *pb.SelfHealingData, err error) {
				require.Nil(t, err)
				require.Equal(t, data.ChallengeStatus, pb.SelfHealingData_Status_SUCCEEDED)
				require.Equal(t, data.MessageType, pb.SelfHealingData_MessageType_SELF_HEALING_RESPONSE_MESSAGE)
				require.Equal(t, data.RespondingMasternodeId, "challenging-node")
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
			res, err := task.VerifySelfHealingChallenge(ctx, tt.message)

			// handle the test case's assertions with the provided func
			tt.expect(t, res, err)
		})
	}

}
