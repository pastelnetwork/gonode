package selfhealing

import (
	"context"
	json "github.com/json-iterator/go"
	"golang.org/x/crypto/sha3"
	"testing"

	fuzz "github.com/google/gofuzz"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqmock "github.com/pastelnetwork/gonode/raptorq/node/test"
	shtest "github.com/pastelnetwork/gonode/supernode/node/test/self_healing"
	"github.com/stretchr/testify/require"
)

func TestListSymbolFileKeysForNFTAndActionTickets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config := NewConfig()
	pastelClient := pastelMock.NewMockClient(t)
	p2pClient := p2pMock.NewMockClient(t)
	raptorQClient := rqmock.NewMockClient(t)
	var nodeClient *shtest.Client
	f := fuzz.New()

	file := []byte("test file")
	//rqDecode := &rqnode.Decode{
	//	File: file,
	//}
	fileHash := sha3.Sum256(file)

	ticket := pastel.RegTicket{TXID: "test-tx-id-nft"}
	f.Fuzz(&ticket)
	ticket.TXID = "test-tx-id-nft"
	ticket.RegTicketData.NFTTicketData.AppTicketData.DataHash = fileHash[:]
	ticket.RegTicketData.NFTTicketData.AppTicketData.RQIDs = []string{"file-hash-to-challenge"}
	nftTicket, err := pastel.EncodeNFTTicket(&ticket.RegTicketData.NFTTicketData)
	require.Nil(t, err)
	ticket.RegTicketData.NFTTicket = nftTicket
	tickets := pastel.RegTickets{ticket}

	cascadeActionRegTicket := pastel.ActionRegTicket{
		TXID: "test-tx-id-cascade",
		ActionTicketData: pastel.ActionTicketData{
			ActionType: pastel.ActionTypeCascade,
		},
	}

	cascadeActionTicket := &pastel.ActionTicket{
		APITicketData: pastel.APICascadeTicket{
			RQIDs:    []string{"cascade-file-hash-to-challenge"},
			DataHash: fileHash[:],
		},
		ActionType: pastel.ActionTypeCascade,
	}

	cascadeTicket, err := pastel.EncodeActionTicket(cascadeActionTicket)
	require.Nil(t, err)
	cascadeActionRegTicket.ActionTicketData.ActionTicket = cascadeTicket

	senseActionRegTicket := pastel.ActionRegTicket{
		TXID: "test-tx-id-sense",
		ActionTicketData: pastel.ActionTicketData{
			ActionType: pastel.ActionTypeSense,
		},
	}

	senseActionTicket := &pastel.ActionTicket{
		APITicketData: pastel.APISenseTicket{
			DDAndFingerprintsIDs: []string{"sense-file-hash-to-challenge"},
			DataHash:             fileHash[:],
		},
		ActionType: pastel.ActionTypeSense,
	}

	senseTicket, err := pastel.EncodeActionTicket(senseActionTicket)
	require.Nil(t, err)
	senseActionRegTicket.ActionTicketData.ActionTicket = senseTicket

	actionTickets := pastel.ActionTicketDatas{cascadeActionRegTicket, senseActionRegTicket}

	//rqIDsData, _ := fakeRQIDsData()

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

	//symbol, err := json.Marshal("24wvWw6zhpaCDwpBAjWXprsnnnB4HApKPkAyArDSi94z")
	//require.NoError(t, err)

	var encodeResp rqnode.Encode
	encodeResp.Symbols = make(map[string][]byte)
	encodeResp.Symbols["GfzPCcSwyyvz5faMjrjPk9rnL5QcSbW1MV1FSPuWXvS3"] = bytes

	tests := []struct {
		testcase string
		setup    func()
		expect   func(*testing.T, map[string]SymbolFileKeyDetails, error)
	}{
		{
			testcase: "when reg & action tickets found, list all the keys and create map with ticket mapping",
			setup: func() {
				pastelClient.ListenOnRegTickets(tickets, nil).ListenOnActionTickets(actionTickets, nil).ListenOnMasterNodesExtra(nodes, nil)
			},
			expect: func(t *testing.T, symbolFileKeyTicketMap map[string]SymbolFileKeyDetails, err error) {
				require.Nil(t, err)

				require.Equal(t, len(symbolFileKeyTicketMap), 3)

				nftKeyMap := symbolFileKeyTicketMap["file-hash-to-challenge"]
				require.Equal(t, nftKeyMap.TicketType, nftTicketType)
				require.Equal(t, nftKeyMap.TicketTxID, "test-tx-id-nft")

				senseKeyMap := symbolFileKeyTicketMap["sense-file-hash-to-challenge"]
				require.Equal(t, senseKeyMap.TicketType, senseTicketType)
				require.Equal(t, senseKeyMap.TicketTxID, "test-tx-id-sense")

				cascadeKeyMap := symbolFileKeyTicketMap["cascade-file-hash-to-challenge"]
				require.Equal(t, cascadeKeyMap.TicketType, cascadeTicketType)
				require.Equal(t, cascadeKeyMap.TicketTxID, "test-tx-id-cascade")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.testcase, func(t *testing.T) {
			// Run the setup for the testcase
			tt.setup()

			service := NewService(config, nil, pastelClient, nodeClient,
				p2pClient, nil)
			task := NewSHTask(service)
			task.StorageHandler.RqClient = raptorQClient
			// call the function to get return values
			_, mapSymbolFileKeys, err := task.ListSymbolFileKeysFromNFTAndActionTickets(ctx)
			// handle the test case's assertions with the provided func
			tt.expect(t, mapSymbolFileKeys, err)
		})
	}

}

func TestCreateClosestNodeMapAgainstKeys(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config := NewConfig()
	pastelClient := pastelMock.NewMockClient(t)
	p2pClient := p2pMock.NewMockClient(t)
	raptorQClient := rqmock.NewMockClient(t)
	var nodeClient *shtest.Client

	nodes := pastel.MasterNodes{}
	nodes = append(nodes, pastel.MasterNode{ExtKey: "PrimaryID"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "A"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "B"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "C"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "D"})
	nodes = append(nodes, pastel.MasterNode{ExtKey: "E"})

	tests := []struct {
		testcase string
		keys     []string
		setup    func()
		expect   func(*testing.T, map[string][]string)
	}{
		{
			testcase: "when keys are listed, should create a closestNodes map against the keys",
			keys:     []string{"file-hash-to-challenge", "cascade-file-hash-to-challenge", "sense-file-hash-to-challenge"},
			setup: func() {
				p2pClient.ListenOnNClosestNodes([]string{"A", "B", "C", "D", "E", "F"}, nil)
			},
			expect: func(t *testing.T, closestNodesMap map[string][]string) {
				require.Equal(t, len(closestNodesMap), 3)

				nftClosestNodes := closestNodesMap["file-hash-to-challenge"]
				senseClosestNodes := closestNodesMap["sense-file-hash-to-challenge"]
				cascadeClosestNodes := closestNodesMap["cascade-file-hash-to-challenge"]

				require.Equal(t, len(nftClosestNodes), 6)
				require.Equal(t, len(senseClosestNodes), 6)
				require.Equal(t, len(cascadeClosestNodes), 6)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.testcase, func(t *testing.T) {
			// Run the setup for the testcase
			tt.setup()

			service := NewService(config, nil, pastelClient, nodeClient,
				p2pClient, nil)
			task := NewSHTask(service)
			task.StorageHandler.RqClient = raptorQClient
			// call the function to get return values
			closestNodesMap := task.createClosestNodesMapAgainstKeys(ctx, tt.keys)
			// handle the test case's assertions with the provided func
			tt.expect(t, closestNodesMap)
		})
	}

}
