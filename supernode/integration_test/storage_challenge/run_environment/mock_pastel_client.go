package main

import (
	"context"
	"encoding/base32"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

var rqSymbolFileKeys []string
var mc = &mocks.Client{}

func init() {
	mc.On("Sign", mock.Anything, mock.MatchedBy(func(data []byte) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign data param:", string(data))
		return true
	}), mock.MatchedBy(func(pastelID string) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign pastel id param:", pastelID)
		return true
	}), mock.MatchedBy(func(passphrase string) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign passphrase param:", passphrase)
		return true
	}), mock.MatchedBy(func(algorithm string) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign algorithm param:", algorithm)
		return true
	})).Return([]byte("mock signature"), nil)

	mc.On("Verify", mock.Anything, mock.MatchedBy(func(data []byte) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify data param:", string(data))
		return true
	}), mock.MatchedBy(func(signature string) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify signature param:", signature)
		return true
	}), mock.MatchedBy(func(pastelID string) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify pastel id param:", pastelID)
		return true
	}), mock.MatchedBy(func(algorithm string) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify algorithm param:", algorithm)
		return true
	})).Return(func(_ context.Context, _ []byte, signature, _, _ string) bool {
		return signature == "mock signature"
	}, nil)

	mc.On("GetBlockCount", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- GetBlockHash")
	}).Return(func(_ context.Context) int32 {
		var ret int
		res, _ := http.DefaultClient.Get("http://helper:8088/getblockcount")
		b, _ := io.ReadAll(res.Body)
		ret, _ = strconv.Atoi(string(b))
		return int32(ret)
	}, func(_ context.Context) error { return nil })

	mc.On("GetBlockHash", mock.Anything, mock.MatchedBy(func(blockHeight int32) bool {
		log.Println("MOCK PASTEL CLIENT -- GetBlockHash blockHeight param:", blockHeight)
		return true
	})).Return(func(ctx context.Context, blockHeight int32) string { return getBlockHash(blockHeight) }, func(ctx context.Context, blockHeight int32) error { return nil })

	mc.On("GetBlockVerbose1", mock.Anything, mock.MatchedBy(func(blockHeight int32) bool {
		log.Println("MOCK PASTEL CLIENT -- GetBlockVerbose1 blockHeight param:", blockHeight)
		return true
	})).Return(func(_ context.Context, blockHeight int32) *pastel.GetBlockVerbose1Result {
		return getBlockVerbose1(blockHeight)
	}, func(_ context.Context) error { return nil })

	mc.On("MasterNodesExtra", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- MasterNodesExtra")
	}).Return(func(_ context.Context) pastel.MasterNodes {
		return getMasternodeList()
	}, func(_ context.Context) error { return nil })

	mc.On("MasterNodeTop", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- MasterNodeTop")
	}).Return(func(_ context.Context) pastel.MasterNodes {
		return getMasternodeList()
	}, func(_ context.Context) error { return nil })

	mc.On("RegTickets", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- RegTickets")
	}).Return(func(_ context.Context) pastel.RegTickets {
		return getRegTickets()
	}, func(_ context.Context) error { return nil })
}

func getBlockVerbose1(blkCount int32) *pastel.GetBlockVerbose1Result {
	return &pastel.GetBlockVerbose1Result{
		Hash:          getBlockHash(blkCount),
		Confirmations: 3,
		Height:        int64(blkCount),
		MerkleRoot:    utils.GetHashFromString(fmt.Sprintf("merkle_root_%d", blkCount)),
	}
}

func getBlockHash(blkCount int32) string {
	return utils.GetHashFromString(fmt.Sprintf("hash_%d", blkCount))
}

func getMasternodeList() pastel.MasterNodes {
	return pastel.MasterNodes{
		{
			Rank:       "1",
			IPAddress:  "192.168.100.11:18232",
			ExtAddress: "192.168.100.11:14444",
			ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn1key")),
			ExtP2P:     "192.168.100.11:14445",
		},
		{
			Rank:       "2",
			IPAddress:  "192.168.100.12:18232",
			ExtAddress: "192.168.100.12:14444",
			ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn2key")),
			ExtP2P:     "192.168.100.12:14445",
		},
		{
			Rank:       "3",
			IPAddress:  "192.168.100.13:18232",
			ExtAddress: "192.168.100.13:14444",
			ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn3key")),
			ExtP2P:     "192.168.100.13:14445",
		},
		{
			Rank:       "4",
			IPAddress:  "192.168.100.14:18232",
			ExtAddress: "192.168.100.14:14444",
			ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn4key")),
			ExtP2P:     "192.168.100.14:14445",
		},
		{
			Rank:       "5",
			IPAddress:  "192.168.100.15:18232",
			ExtAddress: "192.168.100.15:14444",
			ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn5key")),
			ExtP2P:     "192.168.100.15:14445",
		},
	}
}

func AddMockRQSymbolFileKeys(keys []string) {
	rqSymbolFileKeys = append(rqSymbolFileKeys, keys...)
}

func getRegTickets() pastel.RegTickets {
	return pastel.RegTickets{
		{
			RegTicketData: pastel.RegTicketData{
				NFTTicketData: pastel.NFTTicket{
					AppTicketData: pastel.AppTicket{
						RQIDs: rqSymbolFileKeys,
					},
				},
			},
		},
	}
}

func newMockPastelClient() pastel.Client {
	return mc
}
