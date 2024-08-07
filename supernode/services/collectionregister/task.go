package collectionregister

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	collectionRegFee = 1000
)

// CollectionRegistrationTask is the task of registering new Collection.
type CollectionRegistrationTask struct {
	*common.SuperNodeTask
	*common.RegTaskHelper

	*CollectionRegistrationService

	Ticket *pastel.CollectionTicket

	// signature of ticket data signed by this node's pastelID
	ownSignature []byte

	creatorSignature []byte
	registrationFee  int64
	burnTxID         string
}

// Run starts the task
func (task *CollectionRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.removeArtifacts)
}

// ValidatePreBurnTransaction validates the pre-burn txid
func (task *CollectionRegistrationTask) ValidatePreBurnTransaction(ctx context.Context, txid string) error {
	var err error
	log.WithContext(ctx).Debugf("pre-burn-txid: %s", txid)

	if err = task.CalculateFee(); err != nil {
		log.WithContext(ctx).WithError(err).Error("calculate fee failed")
		return errors.Errorf("calculate fee: %w", err)
	}
	log.WithContext(ctx).Info("collection reg fee has been computed")

	<-task.NewAction(func(ctx context.Context) error {
		task.burnTxID = txid

		confirmationChn := task.WaitConfirmation(ctx, txid, int64(task.config.PreburntTxMinConfirmations),
			5*time.Second, true, float64(task.registrationFee), 10)

		log.WithContext(ctx).Debug("waiting for confirmation")
		if err = <-confirmationChn; err != nil {
			log.WithContext(ctx).WithError(err).Errorf("validate pre-burn transaction validation")
			err = errors.Errorf("validate pre-burn transaction validation :%w", err)
			return err
		}
		log.WithContext(ctx).Debug("confirmation done")
		return nil
	})

	return err
}

// ValidateAndRegister will get signed ticket from fee txid, wait until it's confirmations meet expectation.
func (task *CollectionRegistrationTask) ValidateAndRegister(ctx context.Context, ticket []byte, creatorSignature []byte) (string, error) {
	var err error

	log.WithContext(ctx).Info("validating collection registration")
	task.creatorSignature = creatorSignature
	<-task.NewAction(func(ctx context.Context) error {
		if err = task.validateSignedTicketFromWN(ctx, ticket, creatorSignature); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("validate signed ticket failure")
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.NetworkHandler.ConnectedTo == nil)
		if err = task.signAndSendCollectionTicket(ctx, task.NetworkHandler.ConnectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("sign and send collection ticket")
			err = errors.Errorf("signed and send collection ticket")
			return nil
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// only primary node start this action
	var collectionRegTxid string
	if task.NetworkHandler.ConnectedTo == nil {
		<-task.NewAction(func(ctx context.Context) error {
			log.WithContext(ctx).Debug("waiting for signature from peers")
			for {
				select {
				case <-ctx.Done():
					err = ctx.Err()
					if err != nil {
						log.WithContext(ctx).Debug("waiting for signature from peers cancelled or timeout")
					}
					return nil
				case <-task.AllSignaturesReceivedChn:
					log.WithContext(ctx).Debug("all signature received so start validation")

					if err = task.VerifyPeersCollectionTicketSignature(ctx, task.Ticket); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' signature mismatched")
						err = errors.Errorf("peers' signature mismatched: %w", err)
						return nil
					}

					collectionRegTxid, err = task.registerCollection(ctx)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("register collection failed")
						err = errors.Errorf("register Collection: %w", err)
						return nil
					}

					return nil
				}
			}
		})
	}

	return collectionRegTxid, err
}

// CalculateFee calculates and assigns fee
func (task *CollectionRegistrationTask) CalculateFee() error {
	task.registrationFee = int64(collectionRegFee)
	task.ActionTicketRegMetadata = &types.ActionRegMetadata{EstimatedFee: task.registrationFee}
	task.RegTaskHelper.ActionTicketRegMetadata = &types.ActionRegMetadata{EstimatedFee: task.registrationFee}

	return nil
}

func (task *CollectionRegistrationTask) validateSignedTicketFromWN(ctx context.Context, ticket []byte, creatorSignature []byte) error {
	var err error
	task.creatorSignature = creatorSignature

	task.Ticket, err = pastel.DecodeCollectionTicket(ticket)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("ticket", string(ticket)).Errorf("decode collection ticket")
		return errors.Errorf("decode collection ticket: %w", err)
	}

	verified, err := task.PastelClient.VerifyCollectionTicket(ctx, ticket, string(creatorSignature), task.Ticket.Creator, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("verify ticket signature")
		return errors.Errorf("verify ticket signature %w", err)
	}

	if !verified {
		err = errors.New("ticket verification failed")
		log.WithContext(ctx).WithError(err).Errorf("verification failure")
		return err
	}

	return nil
}

// sign and send collection ticket if not primary
func (task *CollectionRegistrationTask) signAndSendCollectionTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeCollectionTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("serialize collection ticket: %w", err)
	}

	task.ownSignature, err = task.PastelClient.SignCollectionTicket(ctx, ticket, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket: %w", err)
	}

	if !isPrimary {
		log.WithContext(ctx).Debug("send signed collection ticket to primary node")

		collectionNode, ok := task.NetworkHandler.ConnectedTo.SuperNodePeerAPIInterface.(*CollectionRegistrationNode)
		if !ok {
			return errors.Errorf("node is not CollectionRegistrationNode")
		}

		if err := collectionNode.SendCollectionTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s: %w", task.NetworkHandler.ConnectedTo.ID, task.NetworkHandler.ConnectedTo.Address, err)
		}
	}
	return nil
}

func (task *CollectionRegistrationTask) registerCollection(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debug("all signature received so start validation")

	//ticketID := fmt.Sprintf("%s.%d.%s", task.Ticket.Caller, task.Ticket.BlockNum, hex.EncodeToString(task.dataHash))

	req := pastel.RegisterCollectionRequest{
		Ticket: &pastel.CollectionTicket{
			CollectionTicketVersion:                 task.Ticket.CollectionTicketVersion,
			CollectionName:                          task.Ticket.CollectionName,
			ItemType:                                task.Ticket.ItemType,
			Creator:                                 task.Ticket.Creator,
			ListOfPastelIDsOfAuthorizedContributors: task.Ticket.ListOfPastelIDsOfAuthorizedContributors,
			BlockNum:                                task.Ticket.BlockNum,
			BlockHash:                               task.Ticket.BlockHash,
			CollectionFinalAllowedBlockHeight:       task.Ticket.CollectionFinalAllowedBlockHeight,
			MaxCollectionEntries:                    task.Ticket.MaxCollectionEntries,
			CollectionItemCopyCount:                 task.Ticket.CollectionItemCopyCount,
			Royalty:                                 task.Ticket.Royalty,
			Green:                                   task.Ticket.Green,
			AppTicketData: pastel.AppTicket{
				MaxPermittedOpenNSFWScore:                      task.Ticket.AppTicketData.MaxPermittedOpenNSFWScore,
				MinimumSimilarityScoreToFirstEntryInCollection: task.Ticket.AppTicketData.MinimumSimilarityScoreToFirstEntryInCollection,
			},
		},
		Signatures: &pastel.CollectionTicketSignatures{
			Principal: map[string]string{
				task.Ticket.Creator: string(task.creatorSignature),
			},
			Mn2: map[string]string{
				task.NetworkHandler.Accepted[0].ID: string(task.PeersTicketSignature[task.NetworkHandler.Accepted[0].ID]),
			},
			Mn3: map[string]string{
				task.NetworkHandler.Accepted[1].ID: string(task.PeersTicketSignature[task.NetworkHandler.Accepted[1].ID]),
			},
		},
		Mn1PastelID: task.config.PastelID,
		Passphrase:  task.config.PassPhrase,
		Fee:         task.registrationFee,
		Label:       fmt.Sprintf("%s-%s", task.config.PastelID, task.Ticket.CollectionName),
	}

	collectionRegTxID, err := task.PastelClient.RegisterCollectionTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("register collection ticket: %w", err)
	}

	return collectionRegTxID, nil
}

func (task *CollectionRegistrationTask) removeArtifacts() {
	//nothing required here
}

// NewCollectionRegistrationTask returns a new Task instance.
func NewCollectionRegistrationTask(service *CollectionRegistrationService) *CollectionRegistrationTask {
	task := &CollectionRegistrationTask{
		SuperNodeTask:                 common.NewSuperNodeTask(logPrefix, service.historyDB),
		CollectionRegistrationService: service,
	}

	task.RegTaskHelper = common.NewRegTaskHelper(task.SuperNodeTask, task.config.PastelID, task.config.PassPhrase,
		common.NewNetworkHandler(task.SuperNodeTask, service.nodeClient,
			RegisterCollectionNodeMaker{}, service.PastelClient,
			task.config.PastelID,
			service.config.NumberConnectedNodes),
		service.PastelClient,
		service.config.PreburntTxMinConfirmations,
	)

	return task
}
