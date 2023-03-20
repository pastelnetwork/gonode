package senseregister

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/ddstore"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// SenseRegistrationTask is the task of registering new Sense.
type SenseRegistrationTask struct {
	*common.SuperNodeTask
	*common.DupeDetectionHandler
	*SenseRegistrationService

	storage *common.StorageHandler

	Ticket   *pastel.ActionTicket
	Asset    *files.File
	dataHash []byte

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
//
//	caller/creator PastelID; block when ticket registration has started; txid of the pre-burn fee
func (task *SenseRegistrationTask) SendRegMetadata(_ context.Context, regMetadata *types.ActionRegMetadata) error {
	if err := task.RequiredStatus(common.StatusConnected); err != nil {
		return err
	}
	task.ActionTicketRegMetadata = regMetadata

	return nil
}

// CalculateFee calculates and assigns fee
func (task *SenseRegistrationTask) CalculateFee(ctx context.Context, file *files.File) error {
	if task.ActionTicketRegMetadata == nil || task.ActionTicketRegMetadata.BlockHash == "" || task.ActionTicketRegMetadata.CreatorPastelID == "" {
		return errors.Errorf("invalid senseRegMetadata")
	}

	task.Asset = file

	var fileBytes []byte
	fileBytes, err := file.Bytes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("read image file")
		return errors.Errorf("read image file: %w", err)
	}

	fileDataInMb := utils.GetFileSizeInMB(fileBytes)
	fee, err := task.PastelHandler.GetEstimatedSenseFee(ctx, fileDataInMb)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("GetEstimatedSenseFeeFailed:%s", err.Error()))
		return errors.Errorf("getting estimated fee %w", err)
	}

	task.registrationFee = int64(fee)
	task.ActionTicketRegMetadata.EstimatedFee = task.registrationFee
	task.RegTaskHelper.ActionTicketRegMetadata.EstimatedFee = task.registrationFee

	return nil
}

// HashExists checks if hash exists in database
func (task *SenseRegistrationTask) HashExists(ctx context.Context, file *files.File) (bool, error) {
	var fileBytes []byte
	fileBytes, err := file.Bytes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("read image file")
		return false, errors.Errorf("read image file: %w", err)
	}

	dataHash, err := utils.Sha3256hash(fileBytes)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error converting bytes to hash")
		return false, errors.Errorf("hash encoded image: %w", err)
	}

	db, err := ddstore.NewSQLiteDDStore(task.config.DDDatabase)
	if err != nil {
		return false, err
	}

	exists, err := db.IfFingerprintExists(ctx, string(dataHash))
	if err != nil {
		return false, err
	}

	if err := db.Close(); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to close dd database")
	}

	if exists {
		task.UpdateStatus(common.StatusFileExists)
	}

	return exists, nil

}

// ProbeImage sends the original image to dd-server and return a compression of pastel.DDAndFingerprints
func (task *SenseRegistrationTask) ProbeImage(ctx context.Context, file *files.File) ([]byte, error) {

	return task.DupeDetectionHandler.ProbeImage(ctx, file,
		task.ActionTicketRegMetadata.BlockHash, task.ActionTicketRegMetadata.BlockHeight, task.ActionTicketRegMetadata.Timestamp, task.ActionTicketRegMetadata.CreatorPastelID, &tasker{})
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
	apiTicket, err := task.Ticket.APISenseTicket()
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

	task.dataHash = apiTicket.DataHash

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
			log.WithContext(ctx).WithError(err).Errorf("validate signed ticket failure")
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.NetworkHandler.ConnectedTo == nil)
		if err = task.signAndSendSenseTicket(ctx, task.NetworkHandler.ConnectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("sign and send Sense ticket")
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
						log.WithContext(ctx).WithError(err).Errorf("peers' signature mismatched")
						err = errors.Errorf("peers' signature mismatched: %w", err)
						return nil
					}

					nftRegTxid, err = task.registerAction(ctx)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("register action failed")
						err = errors.Errorf("register Action: %w", err)
						return nil
					}

					// Store dd_and_fingerprints into Kademlia
					log.WithContext(ctx).WithField("txid", nftRegTxid).Info("storing dd_and_fingerprints symbols")
					if err = task.storeIDFiles(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("store id files")
						err = errors.Errorf("store id files: %w", err)
						return nil
					}

					return nil
				}
			}
		})
	}

	return nftRegTxid, err
}

// ValidateActionActAndConfirm checks Action activation ticket and reStore if possible
func (task *SenseRegistrationTask) ValidateActionActAndConfirm( /*ctx*/ _ context.Context /*actionRegTxID*/, _ string) error {
	//var err error

	//// Wait for action ticket to be activated by walletnode
	//confirmations := task.waitActionActivation(ctx, actionRegTxID, 3, 30*time.Second)
	//err = <-confirmations
	//if err != nil {
	//	return errors.Errorf("wait for confirmation of sense ticket %w", err)
	//}

	return nil
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

	//ticketID := fmt.Sprintf("%s.%d.%s", task.Ticket.Caller, task.Ticket.BlockNum, hex.EncodeToString(task.dataHash))

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
		Label:       task.ActionTicketRegMetadata.BurnTxID,
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
