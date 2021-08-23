package node

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	walletConfig := config.WalletNode{
		PastelAPI: config.PastelAPI{
			Hostname:   "localhost",
			Port:       12170,
			Username:   "rt",
			Passphrase: "rt",
		},
		RestAPI: config.RestAPI{
			Hostname: "localhost",
			Port:     8080,
		},
		Artist: config.Artist{
			PastelID:      "jXa8KpmY4cf2PMXbHFL6wYkS7MaD1FjqHV1MUJPnVv2Ayua8TQ3PUu2Dtyr9mTUEW1BGZbCrdQ54mPp6gsaF1s",
			Passphrase:    "passphrase",
			SpendableAddr: "tPSFddCX7intdWJpig5931VeEAYPnkt4CFL",
		},
	}

	ctx := context.Background()

	wallet := NewWallet(walletConfig)
	img := "../data/test-image1.jpg"
	imgID, err := wallet.UploadNFT(ctx, img)

	assert.Nil(t, err, "upload nft error is not nil")
	assert.NotEmpty(t, imgID, "image id is empty")

	taskID, err := wallet.RegisterNFT(ctx, imgID)
	assert.Nil(t, err, "register nft error is not nil")
	assert.NotEmpty(t, taskID, "register nft error is not nil")

	time.Sleep(time.Millisecond * 100)
	stateReceiver, err := wallet.SubscribeTaskState(ctx, taskID)
	assert.Nil(t, err, fmt.Sprintf("%s", err))
	assert.NotNil(t, stateReceiver, "task state receiver is nil")

	minerConfig := config.MinerNode{
		PastelAPI: config.PastelAPI{
			Hostname:   "localhost",
			Port:       12169,
			Username:   "rt",
			Passphrase: "rt",
		},
	}

	miner := NewMiner(minerConfig)
	for {
		state, err := stateReceiver.Recv()
		if err != nil && err != io.EOF {
			panic(err)
		}
		fmt.Println(state)
		if strings.Contains(state.Status, artworkregister.StatusPreburntRegistrationFee.String()) {
			time.Sleep(7 * time.Second)
			err := miner.GenBlock(ctx, 10)
			assert.Nil(t, err, fmt.Sprintf("%s", err))
			continue
		}
		if strings.Contains(state.Status, artworkregister.StatusTicketAccepted.String()) {
			time.Sleep(7 * time.Second)
			err := miner.GenBlock(ctx, 10)
			assert.Nil(t, err, "generate block error is not nil")
			continue
		}
		if strings.Contains(state.Status, artworkregister.StatusTicketRegistered.String()) {
			time.Sleep(7 * time.Second)
			err := miner.GenBlock(ctx, 5)
			assert.Nil(t, err, "generate block error is not nil")
			continue
		}
		if strings.Contains(state.Status, artworkregister.StatusTicketActivated.String()) {
			break
		}
	}

	blockCount, err := wallet.GetBlockCount(ctx)
	assert.Nil(t, err, "get block count: %s", err)
	assert.NotEqual(t, 0, blockCount, "block count is zero")

	actTicket, err := wallet.FindActTicketByCreatorHeight(ctx, blockCount-25)
	assert.Nil(t, err, "get arc-ticket by height: %s", err)
	fmt.Println(actTicket)
}
