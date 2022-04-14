package nftregister

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

// NftRegistrationTask is the task of registering new Nft.
type NftRegistrationTask struct {
	*common.SuperNodeTask
	*common.DupeDetectionHandler
	*NftRegistrationService

	storage *common.StorageHandler

	nftRegMetadata *types.NftRegMetadata
	Ticket         *pastel.NFTTicket
	ResampledNft   *files.File
	Nft            *files.File
	imageSizeBytes int

	PreviewThumbnail *files.File
	MediumThumbnail  *files.File
	SmallThumbnail   *files.File

	Oti []byte

	// signature of ticket data signed by this node's pastelID
	ownSignature []byte

	creatorSignature []byte
	key1             string
	key2             string
	registrationFee  int64

	rawRqFile   []byte
	rqIDFiles   [][]byte
	rawDdFpFile []byte
	ddFpFiles   [][]byte
}

type tasker struct {
}

func (t *tasker) SendDDFBack(ctx context.Context, node node.SuperNodePeerAPIInterface, nodeInfo *types.MeshedSuperNode, pastelID string, data []byte) error {
	regNode, ok := node.(*NftRegistrationNode)
	if !ok {
		return errors.Errorf("node is not NFTRegistrationNode")
	}
	return regNode.SendSignedDDAndFingerprints(ctx, nodeInfo.SessID, pastelID, data)
}

// Run starts the task
func (task *NftRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.removeArtifacts)
}

// SendRegMetadata sends reg metadata NB this method should really be called SET regMetadata
func (task *NftRegistrationTask) SendRegMetadata(_ context.Context, regMetadata *types.NftRegMetadata) error {
	if err := task.RequiredStatus(common.StatusConnected); err != nil {
		return err
	}
	task.nftRegMetadata = regMetadata
	return nil
}

// ProbeImage sends the resampled image to dd-server and return a compression of pastel.DDAndFingerprints
//  https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration Step 4.A.3
func (task *NftRegistrationTask) ProbeImage(ctx context.Context, file *files.File) ([]byte, error) {
	if task.nftRegMetadata == nil || task.nftRegMetadata.BlockHash == "" || task.nftRegMetadata.CreatorPastelID == "" || task.nftRegMetadata.BlockHeight == "" || task.nftRegMetadata.Timestamp == "" {
		return nil, errors.Errorf("invalid nftRegMetadata")
	}
	task.ResampledNft = file
	return task.DupeDetectionHandler.ProbeImage(ctx, file,
		task.nftRegMetadata.BlockHash, task.nftRegMetadata.BlockHeight, task.nftRegMetadata.Timestamp, task.nftRegMetadata.CreatorPastelID, &tasker{})
}

// validates RQIDs and DdFp IDs file and its IDs
// Step 11.B.3 - 11.B.4
func (task *NftRegistrationTask) validateRqIDsAndDdFpIds(ctx context.Context, rq []byte, dd []byte) error {
	var err error

	pastelIDs := task.NetworkHandler.MeshedNodesPastelID()

	task.rawDdFpFile, task.ddFpFiles, err = task.ValidateIDFiles(ctx, dd,
		task.Ticket.AppTicketData.DDAndFingerprintsIc, task.Ticket.AppTicketData.DDAndFingerprintsMax,
		task.Ticket.AppTicketData.DDAndFingerprintsIDs, 3,
		pastelIDs,
		task.PastelClient)
	if err != nil {
		return errors.Errorf("validate dd_and_fingerprints: %w", err)
	}

	task.rawRqFile, task.rqIDFiles, err = task.ValidateIDFiles(ctx, rq,
		task.Ticket.AppTicketData.RQIc, task.Ticket.AppTicketData.RQMax,
		task.Ticket.AppTicketData.RQIDs, 1,
		[]string{task.Ticket.Author},
		task.PastelClient)
	if err != nil {
		return errors.Errorf("validate rq_ids: %w", err)
	}

	return nil
}

// GetNftRegistrationFee get the fee to register Nft to blockchain
// Step 11.B ALL SuperNode - Validate signature and IDs in the ticket
// Step 12. ALL SuperNode - Calculate final Registration Fee
//  Step 11.B.3 Validate dd_and_fingerprints file signature using PastelID’s of all 3 MNs -- Does this exist?
func (task *NftRegistrationTask) GetNftRegistrationFee(_ context.Context,
	ticket []byte, creatorSignature []byte, key1 string, key2 string,
	rqidFile []byte, ddFpFile []byte, oti []byte,
) (int64, error) {
	var err error
	if err = task.RequiredStatus(common.StatusImageAndThumbnailCoordinateUploaded); err != nil {
		return 0, errors.Errorf("require status %s not satisfied", common.StatusImageAndThumbnailCoordinateUploaded)
	}

	task.Oti = oti
	task.creatorSignature = creatorSignature
	task.key1 = key1
	task.key2 = key2

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(common.StatusRegistrationFeeCalculated)

		// TODO: fix this like how can we get the signature before calling cNode
		task.Ticket, err = pastel.DecodeNFTTicket(ticket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("decode NFT ticket")
			err = errors.Errorf("decode NFT ticket %w", err)
			return nil
		}

		verified, err := task.PastelClient.Verify(ctx, ticket, string(creatorSignature), task.Ticket.Author, pastel.SignAlgorithmED448)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("verify ticket signature")
			err = errors.Errorf("verify ticket signature %w", err)
			return nil
		}

		if !verified {
			err = errors.New("ticket verification failed")
			log.WithContext(ctx).WithError(err).Errorf("verification failure")
			return nil
		}

		if err := task.validateRqIDsAndDdFpIds(ctx, rqidFile, ddFpFile); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("validate rq & dd id files")
			err = errors.Errorf("validate rq & dd id files %w", err)
			return nil
		}

		// Assume passphrase is 16-bytes length

		getFeeRequest := pastel.GetRegisterNFTFeeRequest{
			Ticket: task.Ticket,
			Signatures: &pastel.RegTicketSignatures{
				Creator: map[string]string{
					task.Ticket.Author: string(creatorSignature),
				},
				Mn2: map[string]string{
					task.NftRegistrationService.config.PastelID: string(creatorSignature),
				},
				Mn3: map[string]string{
					task.NftRegistrationService.config.PastelID: string(creatorSignature),
				},
			},
			Mn1PastelID: task.NftRegistrationService.config.PastelID, // all ID has same length, so can use any id here
			Passphrase:  task.config.PassPhrase,
			Key1:        key1,
			Key2:        key2,
			Fee:         0, // fake data
			ImgSizeInMb: int64(task.imageSizeBytes) / (1024 * 1024),
		}

		task.registrationFee, err = task.PastelClient.GetRegisterNFTFee(ctx, getFeeRequest)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("get register NFT fee")
			err = errors.Errorf("get register NFT fee %w", err)
		}
		return nil
	})

	return task.registrationFee, err
}

// ValidatePreBurnTransaction will get pre-burnt transaction fee txid, wait until it's confirmations meet expectation.
// Step 15 - 17
// TODO verify Step 15 "Validate burned transaction (that 10% was really burned)"
func (task *NftRegistrationTask) ValidatePreBurnTransaction(ctx context.Context, txid string) error {
	var err error
	if err = task.RequiredStatus(common.StatusRegistrationFeeCalculated); err != nil {
		return errors.Errorf("require status %s not satisfied", common.StatusRegistrationFeeCalculated)
	}

	log.WithContext(ctx).Debugf("preburn-txid: %s", txid)
	<-task.NewAction(func(ctx context.Context) error {
		confirmationChn := task.WaitConfirmation(ctx, txid, int64(task.config.PreburntTxMinConfirmations), 15*time.Second)

		// compare rqsymbols
		if err = task.compareRQSymbolID(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("generate rqids")
			err = errors.Errorf("generate rqids: %w", err)
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.NetworkHandler.ConnectedTo == nil)
		if err = task.signAndSendNftTicket(ctx, task.NetworkHandler.ConnectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("signed and send NFT ticket")
			err = errors.Errorf("signed and send NFT ticket")
			return nil
		}

		log.WithContext(ctx).Debug("waiting for confimation")
		if err = <-confirmationChn; err != nil {
			log.WithContext(ctx).WithError(err).Errorf("validate preburn transaction validation")
			err = errors.Errorf("validate preburn transaction validation :%w", err)
			return nil
		}
		log.WithContext(ctx).Debug("confirmation done")

		return nil
	})
	return err
}

// ActivateAndStoreNft started by primary node only, it waits for signatures from secondaries and create
// activation ticket and stores everything into P2P
// Step 18 - 19
func (task *NftRegistrationTask) ActivateAndStoreNft(_ context.Context) (string, error) {
	var err error
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

					//Step 18
					if err = task.verifyPeersSignature(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("peers' singature mismatched: %w", err)
						return nil
					}

					nftRegTxid, err = task.registerNft(ctx)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("task register nft failure")
						err = errors.Errorf("register NFT: %w", err)
						return nil
					}

					// We aren't waiting because we are returned a transaction ID from registerNft.
					//	Once we see that on the blockchain, we can proceed.

					// confirmations := task.WaitConfirmation(ctx, nftRegTxid, 10, 30*time.Second, 55)
					// err = <-confirmations
					// if err != nil {
					// 	return errors.Errorf("wait for confirmation of reg-art ticket %w", err)
					// }

					// Step 19
					if err = task.storeRaptorQSymbols(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("store raptor symbols")
						err = errors.Errorf("store raptor symbols: %w", err)
						return nil
					}

					if err = task.storeThumbnails(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("store thumbnails")
						err = errors.Errorf("store thumbnails: %w", err)
						return nil
					}

					//This includes adding fingerprints to the dd-service fingerprint sqlite database
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

// sign and send NFT ticket if not primary
func (task *NftRegistrationTask) signAndSendNftTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeNFTTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("serialize NFT ticket: %w", err)
	}

	task.ownSignature, err = task.PastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket: %w", err)
	}

	if !isPrimary {
		log.WithContext(ctx).Debug("send signed articket to primary node")

		nftRegNode, ok := task.NetworkHandler.ConnectedTo.SuperNodePeerAPIInterface.(*NftRegistrationNode)
		if !ok {
			return errors.Errorf("node is not SenseRegistrationNode")
		}

		if err := nftRegNode.SendNftTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s: %w", task.NetworkHandler.ConnectedTo.ID, task.NetworkHandler.ConnectedTo.Address, err)
		}
	}
	return nil
}

func (task *NftRegistrationTask) verifyPeersSignature(ctx context.Context) error {
	log.WithContext(ctx).Debug("all signature received so start validation")

	data, err := pastel.EncodeNFTTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("encoded NFT ticket: %w", err)
	}
	return task.RegTaskHelper.VerifyPeersSignature(ctx, data)
}

func (task *NftRegistrationTask) registerNft(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debug("all signature received so start validation")

	req := pastel.RegisterNFTRequest{
		Ticket: &pastel.NFTTicket{
			Version:       task.Ticket.Version,
			Author:        task.Ticket.Author,
			BlockNum:      task.Ticket.BlockNum,
			BlockHash:     task.Ticket.BlockHash,
			Copies:        task.Ticket.Copies,
			Royalty:       task.Ticket.Royalty,
			Green:         task.Ticket.Green,
			AppTicketData: task.Ticket.AppTicketData,
		},
		Signatures: &pastel.RegTicketSignatures{
			Creator: map[string]string{
				task.Ticket.Author: string(task.creatorSignature),
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
		// TODO: fix this when how to get key1 and key2 are finalized
		Key1: task.key1,
		Key2: task.key2,
		Fee:  task.registrationFee,
	}

	nftRegTxid, err := task.PastelClient.RegisterNFTTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("register NFT work: %w", err)
	}
	return nftRegTxid, nil
}

// Step 19
func (task *NftRegistrationTask) storeRaptorQSymbols(ctx context.Context) error {
	data, err := task.Nft.Bytes()
	if err != nil {
		return errors.Errorf("read image data: %w", err)
	}

	return task.storage.StoreRaptorQSymbolsIntoP2P(ctx, data, task.Nft.Name())
}

// Step 19
func (task *NftRegistrationTask) storeThumbnails(ctx context.Context) error {
	if _, err := task.storage.StoreFileIntoP2P(ctx, task.PreviewThumbnail); err != nil {
		return errors.Errorf("store preview thumbnail into kademlia: %w", err)
	}
	if _, err := task.storage.StoreFileIntoP2P(ctx, task.MediumThumbnail); err != nil {
		return errors.Errorf("store medium thumbnail into kademlia: %w", err)
	}
	if _, err := task.storage.StoreFileIntoP2P(ctx, task.SmallThumbnail); err != nil {
		return errors.Errorf("store small thumbnail into kademlia: %w", err)
	}

	return nil
}

// Step 19
func (task *NftRegistrationTask) storeIDFiles(ctx context.Context) error {
	if err := task.storage.StoreListOfBytesIntoP2P(ctx, task.ddFpFiles); err != nil {
		return errors.Errorf("store ddAndFp files into kademlia: %w", err)
	}
	if err := task.storage.StoreListOfBytesIntoP2P(ctx, task.rqIDFiles); err != nil {
		return errors.Errorf("store rqIDFiles files into kademlia: %w", err)
	}

	return nil
}

// validates actual RQ Symbol IDs inside RQIDs file
func (task *NftRegistrationTask) compareRQSymbolID(ctx context.Context) error {

	content, err := task.Nft.Bytes()
	if err != nil {
		return errors.Errorf("read image contents: %w", err)
	}

	return task.storage.ValidateRaptorQSymbolIDs(ctx,
		content /*uint32(len(task.Ticket.AppTicketData.RQIDs))*/, 1,
		hex.EncodeToString([]byte(task.Ticket.BlockHash)), task.Ticket.Author,
		task.rawRqFile)
}

// UploadImageWithThumbnail uploads the image that contained image with pqsignature
// generate the image thumbnail from the coordinate provided for user and return
// the hash for the genreated thumbnail
func (task *NftRegistrationTask) UploadImageWithThumbnail(_ context.Context, file *files.File, coordinate files.ThumbnailCoordinate) ([]byte, []byte, []byte, error) {
	var err error
	if err = task.RequiredStatus(common.StatusImageProbed); err != nil {
		return nil, nil, nil, errors.Errorf("require status %s not satisfied", common.StatusImageProbed)
	}

	previewThumbnailHash := make([]byte, 0)
	mediumThumbnailHash := make([]byte, 0)
	smallThumbnailHash := make([]byte, 0)
	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(common.StatusImageAndThumbnailCoordinateUploaded)

		task.Nft = file

		// Determine file size
		// TODO: improve it by call stats on file
		var fileBytes []byte
		fileBytes, err = file.Bytes()
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("read image file")
			err = errors.Errorf("read image file: %w", err)
			return nil
		}
		task.imageSizeBytes = len(fileBytes)

		previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, err = task.createAndHashThumbnails(coordinate)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("create and hash thumbnail")
			return nil
		}

		return nil
	})

	return previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, err
}

func (task *NftRegistrationTask) removeArtifacts() {
	task.RemoveFile(task.ResampledNft)
	task.RemoveFile(task.Nft)
	task.RemoveFile(task.PreviewThumbnail)
	task.RemoveFile(task.MediumThumbnail)
	task.RemoveFile(task.SmallThumbnail)
}

// NewNftRegistrationTask returns a new Task instance.
func NewNftRegistrationTask(service *NftRegistrationService) *NftRegistrationTask {
	task := &NftRegistrationTask{
		SuperNodeTask:          common.NewSuperNodeTask(logPrefix),
		NftRegistrationService: service,
		storage: common.NewStorageHandler(service.P2PClient, service.RQClient,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir),
	}

	task.DupeDetectionHandler = common.NewDupeDetectionTaskHelper(task.SuperNodeTask, service.ddClient,
		task.config.PastelID, task.config.PassPhrase,
		common.NewNetworkHandler(task.SuperNodeTask, service.nodeClient,
			RegisterNftNodeMaker{}, service.PastelClient,
			task.config.PastelID,
			service.config.NumberConnectedNodes),
		service.PastelClient,
		task.config.PreburntTxMinConfirmations,
	)

	return task
}
