package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/pastel/mocks"
	"github.com/stretchr/testify/mock"
)

var mc = &mocks.Client{}

func init() {
	mc.On("Sign", mock.Anything, mock.MatchedBy(func(data []byte) bool {
		return true
	}), mock.MatchedBy(func(pastelID string) bool {
		return true
	}), mock.MatchedBy(func(passphrase string) bool {
		return true
	}), mock.MatchedBy(func(algorithm string) bool {
		return true
	})).Return([]byte("mock signature"), nil)

	mc.On("Verify", mock.Anything, mock.MatchedBy(func(data []byte) bool {
		return true
	}), mock.MatchedBy(func(signature string) bool {
		return true
	}), mock.MatchedBy(func(pastelID string) bool {
		return true
	}), mock.MatchedBy(func(algorithm string) bool {
		return true
	})).Return(true, nil)

	mc.On("GetBlockCount", mock.Anything).Run(func(args mock.Arguments) {
	}).Return(func(_ context.Context) int32 {
		var ret int
		res, _ := http.DefaultClient.Get("http://helper:8088/getblockcount")
		b, _ := io.ReadAll(res.Body)
		ret, _ = strconv.Atoi(string(b))
		return int32(ret)
	}, func(ctx context.Context) error { return nil })

	mc.On("GetBlockHash", mock.Anything, mock.MatchedBy(func(blockHeight int32) bool {
		return true
	})).Return(func(ctx context.Context, blockHeight int32) string { return getBlockHash(blockHeight) },
		func(ctx context.Context, blockHeight int32) error { return nil })

	mc.On("GetBlockVerbose1", mock.Anything, mock.MatchedBy(func(blockHeight int32) bool {
		return true
	})).Return(func(ctx context.Context, blockHeight int32) *pastel.GetBlockVerbose1Result {
		return getBlockVerbose1(blockHeight)
	}, func(_ context.Context, blockHeight int32) error { return nil })

	mc.On("MasterNodesExtra", mock.Anything).Run(func(args mock.Arguments) {
	}).Return(func(_ context.Context) pastel.MasterNodes {
		return getMasternodeList()
	}, func(_ context.Context) error { return nil })

	mc.On("MasternodeList", mock.Anything).Run(func(args mock.Arguments) {
	}).Return(func(_ context.Context) pastel.MasterNodes {
		return getMasternodeList()
	}, func(ctx context.Context) error { return nil })

	mc.On("MasterNodeTop", mock.Anything).Run(func(args mock.Arguments) {
	}).Return(func(_ context.Context) pastel.MasterNodes {
		return getMasternodeList()
	}, func(ctx context.Context) error { return nil })

	mc.On("RegTickets", mock.Anything).Run(func(args mock.Arguments) {
	}).Return(func(_ context.Context) pastel.RegTickets {
		return getRegTickets()
	}, func(ctx context.Context) error { return nil })
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
	var ret pastel.MasterNodes
	res, _ := http.DefaultClient.Get("http://helper:8088/mnlist")
	json.NewDecoder(res.Body).Decode(&ret)
	return ret
}

func getRegTickets() pastel.RegTickets {
	var ret []string
	res, _ := http.DefaultClient.Get("http://helper:8088/store/keys")
	json.NewDecoder(res.Body).Decode(&ret)

	return pastel.RegTickets{
		{
			RegTicketData: pastel.RegTicketData{
				NFTTicketData: pastel.NFTTicket{
					AppTicketData: pastel.AppTicket{
						RQIDs: ret,
					},
				},
			},
		},
	}
}

func newMockPastelClient() pastel.Client {
	return mc
}
