package generator

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/node"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

// NFTRegGenerator consume an nft file and return an activation ticket
type NFTRegGenerator interface {
	Gen(ctx context.Context, nft string) (*pastel.ActTicket, error)
}

type nftRegGenerator struct {
	wallet *node.Wallet
	miner  *node.Miner
}

func (g *nftRegGenerator) Gen(ctx context.Context, nft string) (*pastel.ActTicket, error) {
	imageID, err := g.wallet.UploadNFT(ctx, nft)
	if err != nil {
		return nil, errors.Errorf("failed to upload nft %s: %w", nft, err)
	}

	taskID, err := g.wallet.RegisterNFT(ctx, imageID)
	if err != nil {
		return nil, errors.Errorf("failed to register nft %s: %w", nft, err)
	}

	stateReceiver, err := g.wallet.SubscribeTaskState(ctx, taskID)
	if err != nil {
		return nil, errors.Errorf("failed to subscribe state of task %s: %w", taskID, err)
	}

	preburntConfirmation := 10
	registerNFTConfirmation := 10
	activateNFTConfirmation := 5
	for {
		state, err := stateReceiver.Recv()
		if err != nil {
			if err == io.EOF {
				log.WithContext(ctx).Debugf("EOF received")
				break
			} else {
				return nil, errors.Errorf("failed to receive state of task %s: %w", taskID, err)
			}
		}

		genFn := func(ctx context.Context, state *artworks.TaskState, expectedState string, amount int, waitTime time.Duration) error {
			if strings.Contains(state.Status, expectedState) {
				time.Sleep(waitTime)
				if err := g.miner.GenBlock(ctx, amount); err != nil {
					return errors.Errorf("failed to gen %d block for state %s: %w", amount, expectedState, err)
				}
			}
			return nil
		}

		log.WithContext(ctx).Debugf("%s", state.Status)
		waitTime := time.Duration(time.Second * 7)
		if err := genFn(ctx, state, artworkregister.StatusPreburntRegistrationFee.String(), preburntConfirmation, waitTime); err != nil {
			return nil, errors.Errorf("preburn transaction failed: %w", err)
		}
		if err := genFn(ctx, state, artworkregister.StatusTicketAccepted.String(), registerNFTConfirmation, waitTime); err != nil {
			return nil, errors.Errorf("register nft failed: %w", err)
		}
		if err := genFn(ctx, state, artworkregister.StatusTicketRegistered.String(), activateNFTConfirmation, waitTime); err != nil {
			return nil, errors.Errorf("activate nft failed: %w", err)
		}
		if strings.Contains(state.Status, artworkregister.StatusTicketActivated.String()) {
			log.WithContext(ctx).Debugf("ticket activated")
			break
		}
	}

	blockCount, err := g.wallet.GetBlockCount(ctx)
	if err != nil {
		return nil, errors.Errorf("get block count failed: %w", err)
	}

	blockDiff := preburntConfirmation + registerNFTConfirmation + activateNFTConfirmation
	actTicket, err := g.wallet.FindActTicketByCreatorHeight(ctx, blockCount-int32(blockDiff))
	if err != nil {
		return nil, errors.Errorf("get nft-act failed: %w", err)
	}

	return actTicket, nil
}

func NewNFTRegGenerator(wallet *node.Wallet, miner *node.Miner) *nftRegGenerator {
	return &nftRegGenerator{
		wallet: wallet,
		miner:  miner,
	}
}
