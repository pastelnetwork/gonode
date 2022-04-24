package cascaderegister

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// CascadeRegistrationTask is the task of registering new Sense.
type CascadeRegistrationTask struct {
	*common.SuperNodeTask
	*common.RegTaskHelper
	*CascadeRegistrationService

	storage *common.StorageHandler

	Ticket         *pastel.ActionTicket
	Asset          *files.File
	assetSizeBytes int
	dataHash       []byte

	Oti []byte

	// signature of ticket data signed by this node's pastelID
	ownSignature []byte

	creatorSignature []byte
	registrationFee  int64

	rawRqFile []byte
	rqIDFiles [][]byte
}

// Run starts the task
func (task *CascadeRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.removeArtifacts)
}

// SendRegMetadata receives registration metadata -
//		caller/creator PastelID; block when ticket registration has started; txid of the pre-burn fee
func (task *CascadeRegistrationTask) SendRegMetadata(_ context.Context, regMetadata *types.ActionRegMetadata) error {
	if err := task.RequiredStatus(common.StatusConnected); err != nil {
		return err
	}
	task.ActionTicketRegMetadata = regMetadata

	return nil
}

// UploadAsset uploads the asset
func (task *CascadeRegistrationTask) UploadAsset(_ context.Context, file *files.File) error {
	var err error
	if err = task.RequiredStatus(common.StatusConnected); err != nil {
		return errors.Errorf("require status %s not satisfied", common.StatusConnected)
	}

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(common.StatusAssetUploaded)

		task.Asset = file

		// Determine file size
		// TODO: improve it by call stats on file
		var fileBytes []byte
		fileBytes, err = file.Bytes()
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("read image file")
			err = errors.Errorf("read image file: %w", err)
			return nil
		}
		task.assetSizeBytes = len(fileBytes)

		fileDataInMb := float64(len(fileBytes)) / (1024 * 1024)
		fee, err := task.PastelHandler.GetEstimatedCascadeFee(ctx, fileDataInMb)
		if err != nil {
			err = errors.Errorf("getting estimated fee %w", err)
			return nil
		}
		task.registrationFee = int64(fee)

		return nil
	})

	return err
}

// ValidateAndRegister will get signed ticket from fee txid, wait until it's confirmations meet expectation.
func (task *CascadeRegistrationTask) ValidateAndRegister(_ context.Context,
	ticket []byte, creatorSignature []byte,
	rqidFile []byte, _ /*oti*/ []byte,
) (string, error) {
	var err error
	if err = task.RequiredStatus(common.StatusAssetUploaded); err != nil {
		return "", errors.Errorf("require status %s not satisfied", common.StatusAssetUploaded)
	}

	task.creatorSignature = creatorSignature

	<-task.NewAction(func(ctx context.Context) error {
		if err = task.validateSignedTicketFromWN(ctx, ticket, creatorSignature, rqidFile); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("validateSignedTicketFromWN")
			err = errors.Errorf("validateSignedTicketFromWN: %w", err)
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.NetworkHandler.ConnectedTo == nil)
		if err = task.signAndSendCascadeTicket(ctx, task.NetworkHandler.ConnectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("signed and send Cascade ticket")
			err = errors.Errorf("signed and send NFT ticket: %w", err)
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

func (task *CascadeRegistrationTask) validateSignedTicketFromWN(ctx context.Context,
	ticket []byte, creatorSignature []byte, rqidFile []byte) error {
	var err error
	task.creatorSignature = creatorSignature

	// TODO: fix this like how can we get the signature before calling cNode
	task.Ticket, err = pastel.DecodeActionTicket(ticket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("decode action ticket")
		return errors.Errorf("decode action ticket: %w", err)
	}

	// Verify APICascadeTicket
	apiTicket, err := task.Ticket.APICascadeTicket()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("invalid api cascade ticket")
		return errors.Errorf("invalid api cascade ticket: %w", err)
	}

	verified, err := task.PastelClient.Verify(ctx, ticket, string(creatorSignature), task.Ticket.Caller, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("verify cascade ticket signature")
		return errors.Errorf("verify cascade ticket signature %w", err)
	}

	if !verified {
		err = errors.New("ticket verification failed")
		log.WithContext(ctx).WithError(err).Errorf("verification failure")
		return err
	}

	if err := task.validateRqIDs(ctx, rqidFile); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("validate rqids files")

		return errors.Errorf("validate rq & dd id files %w", err)
	}

	if err = task.validateRQSymbolID(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("valdate rq ids inside rqids file")
		err = errors.Errorf("generate rqids: %w", err)
		return nil
	}

	task.dataHash = apiTicket.DataHash

	return nil
}

// validates RQIDs file
func (task *CascadeRegistrationTask) validateRqIDs(ctx context.Context, dd []byte) error {
	pastelIDs := []string{task.Ticket.Caller}

	apiCascadeTicket, err := task.Ticket.APICascadeTicket()
	if err != nil {
		return errors.Errorf("invalid sense ticket: %w", err)
	}

	task.rawRqFile, task.rqIDFiles, err = task.ValidateIDFiles(ctx, dd,
		apiCascadeTicket.RQIc, apiCascadeTicket.RQMax,
		apiCascadeTicket.RQIDs, 1,
		pastelIDs,
		task.PastelClient)
	if err != nil {
		return errors.Errorf("validate rq_ids file: %w", err)
	}

	return nil
}

// validates actual RQ Symbol IDs inside RQIDs file
func (task *CascadeRegistrationTask) validateRQSymbolID(ctx context.Context) error {

	content, err := task.Asset.Bytes()
	if err != nil {
		return errors.Errorf("read image contents: %w", err)
	}

	return task.storage.ValidateRaptorQSymbolIDs(ctx,
		content /*uint32(len(task.Ticket.AppTicketData.RQIDs))*/, 1,
		hex.EncodeToString([]byte(task.Ticket.BlockHash)), task.Ticket.Caller,
		task.rawRqFile)
}

// sign and send NFT ticket if not primary
func (task *CascadeRegistrationTask) signAndSendCascadeTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeActionTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("serialize NFT ticket: %w", err)
	}

	task.ownSignature, err = task.PastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket: %w", err)
	}

	if !isPrimary {
		log.WithContext(ctx).Debug("send signed cascade ticket to primary node")

		senseNode, ok := task.NetworkHandler.ConnectedTo.SuperNodePeerAPIInterface.(*CascadeRegistrationNode)
		if !ok {
			return errors.Errorf("node is not SenseRegistrationNode")
		}

		if err := senseNode.SendCascadeTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s: %w", task.NetworkHandler.ConnectedTo.ID, task.NetworkHandler.ConnectedTo.Address, err)
		}
	}
	return nil
}

func (task *CascadeRegistrationTask) registerAction(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debug("all signature received so start validation")

	ticketID := fmt.Sprintf("%s.%d.%s", task.Ticket.Caller, task.Ticket.BlockNum, hex.EncodeToString(task.dataHash))

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
			Principal: map[string]string{
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
		Key1:        ticketID,
		Key2:        "key2-" + uuid.New().String(),
	}

	nftRegTxid, err := task.PastelClient.RegisterActionTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("register action ticket: %w", err)
	}
	return nftRegTxid, nil
}

// ValidateActionActAndStore informs actionRegTxID to trigger store ID files in case of actionRegTxID was
func (task *CascadeRegistrationTask) ValidateActionActAndStore(ctx context.Context, actionRegTxID string) error {
	var err error

	// Wait for action ticket to be activated by walletnode
	confirmations := task.waitActionActivation(ctx, actionRegTxID, 3, 30*time.Second)
	err = <-confirmations
	if err != nil {
		return errors.Errorf("wait for confirmation of reg-art ticket %w", err)
	}

	log.WithContext(ctx).WithField("txid", actionRegTxID).Info("storing rq symbols")
	if err = task.storeRaptorQSymbols(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("store raptor symbols")
		err = errors.Errorf("store raptor symbols: %w", err)
		return nil
	}

	// Store dd_and_fingerprints into Kademlia
	log.WithContext(ctx).WithField("txid", actionRegTxID).Info("storing id files")
	if err = task.storeIDFiles(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("store id files")
		err = errors.Errorf("store id files: %w", err)
		return nil
	}

	log.WithContext(ctx).WithField("txid", actionRegTxID).Info("ID files & symbols stored")

	return nil
}

func (task *CascadeRegistrationTask) waitActionActivation(ctx context.Context, txid string, timeoutInBlock int64, interval time.Duration) <-chan error {
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

func (task *CascadeRegistrationTask) storeRaptorQSymbols(ctx context.Context) error {
	data, err := task.Asset.Bytes()
	if err != nil {
		return errors.Errorf("read image data: %w", err)
	}
	return task.storage.StoreRaptorQSymbolsIntoP2P(ctx, data, task.Asset.Name())
}

func (task *CascadeRegistrationTask) storeIDFiles(ctx context.Context) error {
	if err := task.storage.StoreListOfBytesIntoP2P(ctx, task.rqIDFiles); err != nil {
		return errors.Errorf("store ID files into kademlia: %w", err)
	}
	return nil
}

func (task *CascadeRegistrationTask) removeArtifacts() {
	task.RemoveFile(task.Asset)
}

// NewCascadeRegistrationTask returns a new Task instance.
func NewCascadeRegistrationTask(service *CascadeRegistrationService) *CascadeRegistrationTask {

	task := &CascadeRegistrationTask{
		SuperNodeTask:              common.NewSuperNodeTask(logPrefix),
		CascadeRegistrationService: service,
		storage: common.NewStorageHandler(service.P2PClient, service.RQClient,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir),
	}

	task.RegTaskHelper = common.NewRegTaskHelper(task.SuperNodeTask,
		task.config.PastelID, task.config.PassPhrase,
		common.NewNetworkHandler(task.SuperNodeTask, service.nodeClient,
			RegisterCascadeNodeMaker{}, service.PastelClient, task.config.PastelID,
			service.config.NumberConnectedNodes),
		service.PastelClient, task.config.PreburntTxMinConfirmations)

	return task
}
