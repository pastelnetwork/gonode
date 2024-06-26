package selfhealing

/*
import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/supernode/services/download"

	fuzz "github.com/google/gofuzz"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/storage/rqstore"
	"github.com/pastelnetwork/gonode/common/types"
	p2pMock "github.com/pastelnetwork/gonode/p2p/test"
	"github.com/pastelnetwork/gonode/pastel"
	pastelMock "github.com/pastelnetwork/gonode/pastel/test"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqmock "github.com/pastelnetwork/gonode/raptorq/node/test"
	shtest "github.com/pastelnetwork/gonode/supernode/node/test/self_healing"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
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

	configD := download.Config{}
	downloadService := download.NewService(&configD, pastelClient, p2pClient, nil)

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
				p2pClient, nil, downloadService, rqstore.SetupTestDB(t))
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

	configD := download.Config{}
	downloadService := download.NewService(&configD, pastelClient, p2pClient, nil)

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
				p2pClient.ListenOnNClosestNodesWithIncludingNodelist([]string{"A", "B", "C", "D", "E", "F"}, nil)
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
				p2pClient, nil, downloadService, rqstore.SetupTestDB(t))
			task := NewSHTask(service)
			task.StorageHandler.RqClient = raptorQClient
			// call the function to get return values
			closestNodesMap := task.createClosestNodesMapAgainstKeys(ctx, tt.keys, nil)
			// handle the test case's assertions with the provided func
			tt.expect(t, closestNodesMap)
		})
	}

}

func TestCreateSelfHealingTicketsMap(t *testing.T) {
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

	closestNodesMap := make(map[string][]string)
	closestNodesMap["file-hash-to-challenge"] = []string{"A", "B", "C", "D", "E", "F"}
	closestNodesMap["cascade-file-hash-to-challenge"] = []string{"G", "H", "I", "J", "K", "L"}
	closestNodesMap["sense-file-hash-to-challenge"] = []string{"AG", "BH", "CI", "DJ", "EK", "FL"}

	watchlistPingInfo := []types.PingInfo{
		types.PingInfo{SupernodeID: "A"},
		types.PingInfo{SupernodeID: "B"},
		types.PingInfo{SupernodeID: "C"},
		types.PingInfo{SupernodeID: "D"},
		types.PingInfo{SupernodeID: "E"},
		types.PingInfo{SupernodeID: "F"},
	}

	symbolFileKeyMap := make(map[string]SymbolFileKeyDetails)
	symbolFileKeyMap["file-hash-to-challenge"] = SymbolFileKeyDetails{TicketTxID: "test-tx-id-nft", TicketType: nftTicketType}
	symbolFileKeyMap["file-hash-to-challenge-cascade"] = SymbolFileKeyDetails{TicketTxID: "test-tx-id-nft", TicketType: cascadeTicketType}
	symbolFileKeyMap["file-hash-to-challenge-sense"] = SymbolFileKeyDetails{TicketTxID: "test-tx-id-nft", TicketType: senseTicketType}

	configD := download.Config{}
	downloadService := download.NewService(&configD, pastelClient, p2pClient, nil)

	tests := []struct {
		testcase string
		keys     []string
		setup    func()
		expect   func(*testing.T, map[string]SymbolFileKeyDetails)
	}{
		{
			testcase: "when all the closest nodes are on watchlist, should include the ticket for self-healing",
			keys:     []string{"file-hash-to-challenge", "cascade-file-hash-to-challenge", "sense-file-hash-to-challenge"},
			setup: func() {
				p2pClient.ListenOnNClosestNodes([]string{"A", "B", "C", "D", "E", "F"}, nil)
			},
			expect: func(t *testing.T, selfHealingTicketsMap map[string]SymbolFileKeyDetails) {
				require.Equal(t, len(selfHealingTicketsMap), 1)

				ticketDetails := selfHealingTicketsMap["test-tx-id-nft"]
				require.Equal(t, ticketDetails.TicketType, nftTicketType)
				require.Equal(t, len(ticketDetails.Keys), 1)
				require.Equal(t, ticketDetails.Keys[0], "file-hash-to-challenge")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.testcase, func(t *testing.T) {
			// Run the setup for the testcase
			tt.setup()

			service := NewService(config, nil, pastelClient, nodeClient,
				p2pClient, nil, downloadService, rqstore.SetupTestDB(t))
			task := NewSHTask(service)
			task.StorageHandler.RqClient = raptorQClient
			// call the function to get return values
			selfHealingTicketsMap := task.identifySelfHealingTickets(ctx, watchlistPingInfo, closestNodesMap, symbolFileKeyMap)
			// handle the test case's assertions with the provided func
			tt.expect(t, selfHealingTicketsMap)
		})
	}

}

func TestShouldTriggerSelfHealing(t *testing.T) {
	t.Parallel()

	config := NewConfig()
	pastelClient := pastelMock.NewMockClient(t)
	p2pClient := p2pMock.NewMockClient(t)
	raptorQClient := rqmock.NewMockClient(t)
	var nodeClient *shtest.Client

	configD := download.Config{}
	downloadService := download.NewService(&configD, pastelClient, p2pClient, nil)

	tests := []struct {
		testcase string
		infos    types.PingInfos
		expect   func(*testing.T, bool, types.PingInfos)
	}{
		{
			testcase: "when more than 6 nodes that are on watchlist, have last_seen within 30 min time span",
			infos: []types.PingInfo{
				types.PingInfo{SupernodeID: "A", LastSeen: sql.NullTime{Time: time.Now().Add(-5 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "B", LastSeen: sql.NullTime{Time: time.Now().Add(-10 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "C", LastSeen: sql.NullTime{Time: time.Now().Add(-15 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "D", LastSeen: sql.NullTime{Time: time.Now().Add(-45 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "E", LastSeen: sql.NullTime{Time: time.Now().Add(-50 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "F", LastSeen: sql.NullTime{Time: time.Now().Add(-75 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "G", LastSeen: sql.NullTime{Time: time.Now().Add(-100 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "H", LastSeen: sql.NullTime{Time: time.Now().Add(-25 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "I", LastSeen: sql.NullTime{Time: time.Now().Add(-29 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "J", LastSeen: sql.NullTime{Time: time.Now().Add(-22 * time.Minute), Valid: true}},
			},
			expect: func(t *testing.T, shouldTrigger bool, infos types.PingInfos) {
				require.Equal(t, len(infos), 6)
				require.Equal(t, shouldTrigger, true)
			},
		},
		{
			testcase: "when less than 6 nodes that are on watchlist, have last_seen within 30 min time span",
			infos: []types.PingInfo{
				types.PingInfo{SupernodeID: "A", LastSeen: sql.NullTime{Time: time.Now().Add(-55 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "B", LastSeen: sql.NullTime{Time: time.Now().Add(-90 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "C", LastSeen: sql.NullTime{Time: time.Now().Add(-45 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "D", LastSeen: sql.NullTime{Time: time.Now().Add(-45 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "E", LastSeen: sql.NullTime{Time: time.Now().Add(-50 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "F", LastSeen: sql.NullTime{Time: time.Now().Add(-75 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "G", LastSeen: sql.NullTime{Time: time.Now().Add(-100 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "H", LastSeen: sql.NullTime{Time: time.Now().Add(-75 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "I", LastSeen: sql.NullTime{Time: time.Now().Add(-29 * time.Minute), Valid: true}},
				types.PingInfo{SupernodeID: "J", LastSeen: sql.NullTime{Time: time.Now().Add(-22 * time.Minute), Valid: true}},
			},
			expect: func(t *testing.T, shouldTrigger bool, infos types.PingInfos) {
				require.Equal(t, len(infos), 2)
				require.Equal(t, shouldTrigger, false)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.testcase, func(t *testing.T) {
			// Run the setup for the testcase
			service := NewService(config, nil, pastelClient, nodeClient,
				p2pClient, nil, downloadService, rqstore.SetupTestDB(t))
			task := NewSHTask(service)
			task.StorageHandler.RqClient = raptorQClient
			// call the function to get return values
			shouldTrigger, filteredPings := task.shouldTriggerSelfHealing(tt.infos)
			// handle the test case's assertions with the provided func
			tt.expect(t, shouldTrigger, filteredPings)
		})
	}

}
*/
