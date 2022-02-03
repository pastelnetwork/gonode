package senseregister

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// Task is the task of registering new Sense.
type SenseRegistrationTask struct {
	*common.SuperNodeTask
	*common.DupeDetectionHandler
	*SenseRegistrationService

	storage *common.StorageHandler

	Ticket *pastel.ActionTicket
	Asset  *files.File

	// signature of ticket data signed by this node's pastelID
	ownSignature []byte

	creatorSignature []byte
	registrationFee  int64

	rawDdFpFile []byte
	ddFpFiles   [][]byte
}

type tasker struct {
}

func (t *tasker) SendDDFBack(ctx context.Context, node node.SuperNodePeerAPIInterface, nodeInfo *types.MeshedSuperNode, pastelID string, data []byte) error {
	senseNode, ok := node.(*SenseRegistrationNode)
	if !ok {
		return errors.Errorf("node is not SenseRegistrationNode")
	}
	return senseNode.SendSignedDDAndFingerprints(ctx, nodeInfo.SessID, pastelID, data)
}

// Run starts the task
func (task *SenseRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.removeArtifacts)
}

// SendRegMetadata receives registration metadata -
//		caller/creator PastelID; block when ticket registration has started; txid of the pre-burn fee
func (task *SenseRegistrationTask) SendRegMetadata(_ context.Context, regMetadata *types.ActionRegMetadata) error {
	if err := task.RequiredStatus(common.StatusConnected); err != nil {
		return err
	}
	task.SenseRegMetadata = regMetadata

	return nil
}

// ProbeImage sends the original image to dd-server and return a compression of pastel.DDAndFingerprints
func (task *SenseRegistrationTask) ProbeImage(ctx context.Context, file *files.File) ([]byte, error) {
	if task.SenseRegMetadata == nil || task.SenseRegMetadata.BlockHash == "" || task.SenseRegMetadata.CreatorPastelID == "" {
		return nil, errors.Errorf("invalid senseRegMetadata")
	}
	task.Asset = file
	return task.DupeDetectionHandler.ProbeImage(ctx, file,
		task.SenseRegMetadata.BlockHash, task.SenseRegMetadata.CreatorPastelID, &tasker{})
}

func (task *SenseRegistrationTask) validateDdFpIds(ctx context.Context, dd []byte) error {

	pastelIDs := task.NetworkHandler.MeshedNodesPastelID()

	apiSenseTicket, err := task.Ticket.APISenseTicket()
	if err != nil {
		return errors.Errorf("invalid sense ticket: %w", err)
	}

	task.rawDdFpFile, task.ddFpFiles, err = task.ValidateIDFiles(ctx, dd,
		apiSenseTicket.DDAndFingerprintsIc, apiSenseTicket.DDAndFingerprintsMax,
		apiSenseTicket.DDAndFingerprintsIDs, 3,
		pastelIDs,
		task.PastelClient)
	if err != nil {
		return errors.Errorf("validate dd_and_fingerprints: %w", err)
	}

	return nil
}

func (task *SenseRegistrationTask) validateSignedTicketFromWN(ctx context.Context, ticket []byte, creatorSignature []byte, ddFpFile []byte) error {
	var err error
	task.creatorSignature = creatorSignature

	// TODO: fix this like how can we get the signature before calling cNode
	task.Ticket, err = pastel.DecodeActionTicket(ticket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("decode action ticket")
		return errors.Errorf("decode action ticket: %w", err)
	}

	// Verify APISenseTicket
	_, err = task.Ticket.APISenseTicket()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("invalid api sense ticket")
		return errors.Errorf("invalid api sense ticket: %w", err)
	}

	verified, err := task.PastelClient.Verify(ctx, ticket, string(creatorSignature), task.Ticket.Caller, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("verify ticket signature")
		return errors.Errorf("verify ticket signature %w", err)
	}

	if !verified {
		err = errors.New("ticket verification failed")
		log.WithContext(ctx).WithError(err).Errorf("verification failure")
		return err
	}

	if err := task.validateDdFpIds(ctx, ddFpFile); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("validate rq & dd id files")

		return errors.Errorf("validate rq & dd id files %w", err)
	}
	return nil
}

// ValidateAndRegister will get signed ticket from fee txid, wait until it's confirmations meet expectation.
func (task *SenseRegistrationTask) ValidateAndRegister(_ context.Context, ticket []byte, creatorSignature []byte, ddFpFile []byte) (string, error) {
	var err error
	if err = task.RequiredStatus(common.StatusImageProbed); err != nil {
		return "", errors.Errorf("require status %s not satisfied", common.StatusImageProbed)
	}

	task.creatorSignature = creatorSignature

	<-task.NewAction(func(ctx context.Context) error {
		if err = task.validateSignedTicketFromWN(ctx, ticket, creatorSignature, ddFpFile); err != nil {
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.NetworkHandler.ConnectedTo == nil)
		if err = task.signAndSendSenseTicket(ctx, task.NetworkHandler.ConnectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("signed and send NFT ticket")
			err = errors.Errorf("signed and send NFT ticket")
			return nil
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// only primary node start this action
	var nftRegTxid string
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

					if err = task.VerifyPeersTicketSignature(ctx, task.Ticket); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("peers' singature mismatched: %w", err)
						return nil
					}

					nftRegTxid, err = task.registerAction(ctx)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("register NFT: %w", err)
						return nil
					}

					return nil
				}
			}
		})
	}

	return nftRegTxid, err
}

// ValidateActionActAndStore informs actionRegTxID to trigger store ID files in case of actionRegTxID was
func (task *SenseRegistrationTask) ValidateActionActAndStore(ctx context.Context, actionRegTxID string) error {
	var err error

	// Wait for action ticket to be activated by walletnode
	confirmations := task.waitActionActivation(ctx, actionRegTxID, 2, 30*time.Second)
	err = <-confirmations
	if err != nil {
		return errors.Errorf("wait for confirmation of reg-art ticket %w", err)
	}

	// Store dd_and_fingerprints into Kademlia
	if err = task.storeIDFiles(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("store id files")
		err = errors.Errorf("store id files: %w", err)
		return nil
	}

	return nil
}

func (task *SenseRegistrationTask) waitActionActivation(ctx context.Context, txid string, timeoutInBlock int64, interval time.Duration) <-chan error {
	ch := make(chan error)

	go func(ctx context.Context, txid string) {
		defer close(ch)
		blockTracker := blocktracker.New(task.PastelClient)
		baseBlkCnt, err := blockTracker.GetBlockCount()
		if err != nil {
			log.WithContext(ctx).WithError(err).Warn("failed to get block count")
			ch <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				// context cancelled or abort by caller so no need to return anything
				log.WithContext(ctx).Debugf("context done: %s", ctx.Err())
				ch <- ctx.Err()
				return
			case <-time.After(interval):
				txResult, err := task.PastelClient.FindActionActByActionRegTxid(ctx, txid)
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("FindActionActByActionRegTxid err")
				} else {
					if txResult != nil {
						log.WithContext(ctx).Debug("action reg is activated")
						ch <- nil
						return
					}
				}

				currentBlkCnt, err := blockTracker.GetBlockCount()
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("failed to get block count")
					continue
				}

				if currentBlkCnt-baseBlkCnt >= int32(timeoutInBlock)+2 {
					ch <- errors.Errorf("timeout when waiting for confirmation of transaction %s", txid)
					return
				}
			}

		}
	}(ctx, txid)
	return ch
}

// sign and send NFT ticket if not primary
func (task *SenseRegistrationTask) signAndSendSenseTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeActionTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("serialize NFT ticket: %w", err)
	}

	task.ownSignature, err = task.PastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket: %w", err)
	}

	if !isPrimary {
		log.WithContext(ctx).Debug("send signed sense ticket to primary node")

		senseNode, ok := task.NetworkHandler.ConnectedTo.SuperNodePeerAPIInterface.(*SenseRegistrationNode)
		if !ok {
			return errors.Errorf("node is not SenseRegistrationNode")
		}

		if err := senseNode.SendSenseTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s: %w", task.NetworkHandler.ConnectedTo.ID, task.NetworkHandler.ConnectedTo.Address, err)
		}
	}
	return nil
}

func (task *SenseRegistrationTask) registerAction(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debug("all signature received so start validation")

	req := pastel.RegisterActionRequest{
		Ticket: &pastel.ActionTicket{
			Version:       task.Ticket.Version,
			Caller:        task.Ticket.Caller,
			BlockNum:      task.Ticket.BlockNum,
			BlockHash:     task.Ticket.BlockHash,
			ActionType:    task.Ticket.ActionType,
			APITicketData: task.Ticket.APITicketData,
		},
		Signatures: &pastel.ActionTicketSignatures{
			Caller: map[string]string{
				task.Ticket.Caller: string(task.creatorSignature),
			},
			Mn1: map[string]string{
				task.config.PastelID: string(task.ownSignature),
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
	}

	nftRegTxid, err := task.PastelClient.RegisterActionTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("register action ticket: %w", err)
	}
	return nftRegTxid, nil
}

func (task *SenseRegistrationTask) storeIDFiles(ctx context.Context) error {
	if err := task.storage.StoreListOfBytesIntoP2P(ctx, task.ddFpFiles); err != nil {
		return errors.Errorf("store ID files into kademlia: %w", err)
	}
	return nil
}

func (task *SenseRegistrationTask) removeArtifacts() {
	task.RemoveFile(task.Asset)
}

// NewSenseRegistrationTask returns a new Task instance.
func NewSenseRegistrationTask(service *SenseRegistrationService) *SenseRegistrationTask {
	task := &SenseRegistrationTask{
		SuperNodeTask:            common.NewSuperNodeTask(logPrefix),
		SenseRegistrationService: service,
		storage:                  common.NewStorageHandler(service.P2PClient, nil, "", ""),
	}

	task.DupeDetectionHandler = common.NewDupeDetectionTaskHelper(task.SuperNodeTask, service.ddClient,
		task.config.PastelID, task.config.PassPhrase,
		common.NewNetworkHandler(task.SuperNodeTask, service.nodeClient,
			RegisterSenseNodeMaker{}, service.PastelClient,
			task.config.PastelID,
			service.config.NumberConnectedNodes),
		service.PastelClient,
		task.config.PreburntTxMinConfirmations,
	)

	return task
}
