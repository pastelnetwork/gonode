package nftregister

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
)

const (
	//DateTimeFormat following the go convention for request timestamp
	DateTimeFormat = "2006:01:02 15:04:05"
)

// NftRegistrationTask an NFT on the blockchain
// Follow instructions from : https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration//
// NftRegistrationTask is Run from NftRegisterService.Run(), which eventually calls run, below
type NftRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler         *common.MeshHandler
	FingerprintsHandler *mixins.FingerprintsHandler
	ImageHandler        *mixins.NftImageHandler
	RqHandler           *mixins.RQHandler

	service         *NftRegistrationService
	Request         *NftRegistrationRequest
	downloadService *download.NftDownloadingService

	// task data to create RegArt ticket
	creatorBlockHeight      int
	creatorBlockHash        string
	creationTimestamp       string
	dataHash                []byte
	registrationFee         int64
	originalFileSizeInBytes int
	fileType                string

	// ticket
	creatorSignature      []byte
	nftRegistrationTicket *pastel.NFTTicket
	serializedTicket      []byte

	regNFTTxid string

	// only set to true for unit tests
	skipPrimaryNodeTxidVerify bool
}

// Run starts the task
func (task *NftRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *NftRegistrationTask) skipPrimaryNodeTxidCheck() bool {
	return task.skipPrimaryNodeTxidVerify || os.Getenv("INTEGRATION_TEST_ENV") == "true"
}

// Run sets up a connection to a mesh network of supernodes, then controls the communications to the mesh of nodes.
//
//		Task here will abstract away the individual node communications layer, and instead operate at the mesh control layer.
//	 For individual communcations control, see node/grpc/nft_register.go
func (task *NftRegistrationTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.WithContext(ctx).Info("checking storage fee")
	task.StatusLog[common.FieldTaskType] = "NFT Registration"
	task.StatusLog[common.FieldNFTRegisterTaskMaxFee] = task.Request.MaximumFee
	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking storage fee")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorCheckingStorageFee)
		return err
	} else if !ok {
		task.UpdateStatus(common.StatusErrorInsufficientFee)
		log.WithContext(ctx).Error("insufficient storage fee")
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Request.MaximumFee)
	}

	task.StatusLog[common.FieldSpendableAddress] = task.Request.SpendableAddress
	if err := task.service.pastelHandler.IsBalanceMoreThanMaxFee(ctx, task.Request.SpendableAddress, task.Request.MaximumFee); err != nil {
		log.WithContext(ctx).WithError(err).Error("current balance is not sufficient for the storage fee")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorInsufficientBalance)
		return errors.Errorf("current balance is less than max fee provided in the ticket: %v", task.Request.MaximumFee)
	}

	log.WithContext(ctx).Info("setting up mesh with top supernodes")

	// Setup mesh with supernodes with the highest ranks.
	creatorBlockHeight, creatorBlockHash, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error establishing a mesh of SNs")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorMeshSetupFailed)
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash
	//set timestamp for nft reg metadata
	task.creationTimestamp = time.Now().Format(DateTimeFormat)
	task.StatusLog[common.FieldBlockHeight] = creatorBlockHeight

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	log.WithContext(ctx).Info("uploading data to supernodes")

	// send registration metadata
	if err := task.sendRegMetadata(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending reg metadata")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSendingRegMetadata)
		return errors.Errorf("send registration metadata: %w", err)
	}

	nftBytes, err := task.Request.Image.Bytes()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error converting image to bytes")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorConvertingImageBytes)
		return errors.Errorf("convert image to byte stream %w", err)
	}
	task.originalFileSizeInBytes = len(nftBytes)

	//Detect the file type
	task.fileType = mimetype.Detect(nftBytes).String()

	// probe ORIGINAL image for average rareness, nsfw and seen score
	if err := task.probeImage(ctx, task.Request.Image, task.Request.Image.Name()); err != nil {
		log.WithContext(ctx).WithError(err).Error("error probing image")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorProbeImage)
		return errors.Errorf("probe image: %w", err)
	}

	log.WithContext(ctx).Info("image probed successfully")

	// Calculate IDs for redundant number of dd_and_fingerprints files
	// Step 5.4
	if err := task.FingerprintsHandler.GenerateDDAndFingerprintsIDs(ctx, task.service.config.DDAndFingerprintsMax); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Error generating FPs:%s", err.Error())
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorGenerateDDAndFPIds)
		return errors.Errorf("DD and/or Fingerprint ID error: %w", err)
	}
	log.WithContext(ctx).Info("fingerprint && dd IDs have been generated")

	// Create copy of original image and embed fingerprints into it
	// result is in the - task.NftImageHandler
	// Step 6
	if err := task.ImageHandler.CreateCopyWithEncodedFingerprint(ctx,
		task.Request.CreatorPastelID, task.Request.CreatorPastelIDPassphrase,
		task.FingerprintsHandler.FinalFingerprints, task.Request.Image); err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating copy with encoded FPs")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorEncodingImage)
		return errors.Errorf("encode image with fingerprint: %w", err)
	}
	log.WithContext(ctx).Info("copy created with encoded fingerprints")

	if task.dataHash, err = task.ImageHandler.GetHash(); err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving nft hash")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorGetImageHash)
		return errors.Errorf("get image hash: %w", err)
	}
	log.WithContext(ctx).Info("data hash has been generated")

	// upload image (from task.NftImageHandler) and thumbnail coordinates to supernode(s)
	// SN will return back: hashes, previews, thumbnails,
	// Step 6 - 8 for Walletnode
	if err := task.uploadImage(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error uploading signed image")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorUploadImageFailed)
		return errors.Errorf("upload image: %w", err)
	}
	log.WithContext(ctx).Info("signed image has been uploaded")

	// connect to rq serivce to get rq symbols identifier
	// All of step 9
	if err := task.RqHandler.GenRQIdentifiersFiles(ctx, task.ImageHandler.ImageEncodedWithFingerprints,
		task.creatorBlockHash, task.Request.CreatorPastelID, task.Request.CreatorPastelIDPassphrase); err != nil {
		log.WithContext(ctx).WithError(err).Error("error generating RQIDs")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}
	log.WithContext(ctx).Info("rq IDs have been generated")

	if err := task.createNftTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating NFT ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorCreatingTicket)
		return errors.Errorf("create ticket: %w", err)
	}
	log.WithContext(ctx).Info("nft-reg ticket has been created")

	// sign ticket with artist signature
	// Step 10
	if err := task.signTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSigningTicket)
		return errors.Errorf("sign NFT ticket: %w", err)
	}
	log.WithContext(ctx).Info("nft-reg ticket has been signed")

	// send signed ticket to supernodes to calculate registration fee
	// Step 11.A WalletNode - Upload Signed Ticket; RaptorQ IDs file and dd_and_fingerprints file to the SuperNode’s 1, 2 and 3
	if err := task.sendSignedTicket(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error sending signed ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorSendSignTicketFailed)
		return errors.Errorf("send signed NFT ticket: %w", err)
	}
	log.WithContext(ctx).Info("signed ticket has been sent")

	// validate if address has enough psl
	if err := task.service.pastelHandler.CheckBalanceToPayRegistrationFee(ctx,
		task.Request.SpendableAddress,
		float64(task.registrationFee),
		task.Request.MaximumFee); err != nil {
		log.WithContext(ctx).WithError(err).Error("error checking balance to pay reg fee")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorCheckBalance)
		return errors.Errorf("check current balance: %s", err.Error())
	}
	log.WithContext(ctx).Info("spendable address has been validated to check if have enough PSL")

	// send preburn-txid to master node(s)
	// master node will create reg-nft ticket and returns transaction id
	task.UpdateStatus(common.StatusPreburntRegistrationFee)
	if err := task.preburnRegistrationFeeGetTicketTxid(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error registering NFT ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorPreburnRegFeeGetTicketID)
		return errors.Errorf("pre-burnt ten percent of registration fee: %s", err.Error())
	}
	task.StatusLog[common.FieldRegTicketTxnID] = task.regNFTTxid
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validating NFT Reg TXID: ",
		StatusString:  task.regNFTTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	task.UpdateStatus(common.StatusTicketAccepted)
	log.WithContext(ctx).Debug("pre-burned 10% of fee, signed ticket uploaded to SNs")

	now := time.Now()
	if task.downloadService != nil {
		if err := common.DownloadWithRetry(ctx, task, now, now.Add(1*time.Minute)); err != nil {
			log.WithContext(ctx).WithField("reg_tx_id", task.regNFTTxid).WithError(err).Error("error validating nft ticket data")

			log.WithContext(ctx).WithField("reg_tx_id", task.regNFTTxid).Info("initiating the new registration request")
			request := &NftRegistrationRequest{
				Name:                      task.Request.Name,
				Description:               task.Request.Description,
				Keywords:                  task.Request.Keywords,
				SeriesName:                task.Request.SeriesName,
				IssuedCopies:              task.Request.IssuedCopies,
				YoutubeURL:                task.Request.YoutubeURL,
				CreatorPastelID:           task.Request.CreatorPastelID,
				CreatorPastelIDPassphrase: task.Request.CreatorPastelIDPassphrase,
				CreatorName:               task.Request.CreatorName,
				CreatorWebsiteURL:         task.Request.CreatorWebsiteURL,
				SpendableAddress:          task.Request.SpendableAddress,
				MaximumFee:                task.Request.MaximumFee,
				Green:                     task.Request.Green,
				Royalty:                   task.Request.Royalty,
				Thumbnail:                 task.Request.Thumbnail,
				MakePubliclyAccessible:    task.Request.MakePubliclyAccessible,
			}

			newTask := NewNFTRegistrationTask(task.service, request)
			task.service.Worker.AddTask(newTask)

			log.WithContext(ctx).WithField("task_id", newTask.ID()).Info("new NFT registration task has been initiated")

			task.UpdateStatus(common.StatusErrorDownloadFailed)
			task.UpdateStatus(&common.EphemeralStatus{
				StatusTitle:   "Error validating cascade ticket data",
				StatusString:  task.regNFTTxid,
				IsFailureBool: false,
				IsFinalBool:   false,
			})

			task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

			return nil
		}
	}

	log.WithContext(ctx).Infof("nft-reg ticket registered. NFT Registration Ticket txid: %s", task.regNFTTxid)
	log.WithContext(ctx).Info("Closing SNs connections")

	// don't need SNs anymore
	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	// new context because the old context already cancelled
	newCtx := log.ContextWithPrefix(context.Background(), "nft-reg")
	//Start Step 20

	log.WithContext(ctx).Info("waiting confirmations for nft-reg ticket")

	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regNFTTxid, int64(task.service.config.NFTRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		log.WithContext(ctx).WithError(err).Error("error validating nft-reg-txid")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorValidateRegTxnIDFailed)

		return errors.Errorf("wait reg-nft ticket valid: %s", err.Error())
	}
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Validated NFT Reg TXID: ",
		StatusString:  task.regNFTTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	task.UpdateStatus(common.StatusTicketRegistered)
	log.WithContext(ctx).Info("nft-reg ticket confirmed, activating nft-reg ticket")

	// activate reg-nft ticket at previous step
	// Send Activation ticket and Registration Fee (cNode API)
	activateTxID, err := task.registerActTicket(newCtx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error registering act ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorActivatingTicket)
		return errors.Errorf("register act ticket: %s", err.Error())
	}
	task.StatusLog[common.FieldActivateTicketTxnID] = activateTxID
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "NFT Activated - Ticket TXID: ",
		StatusString:  activateTxID,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("nft-reg ticket activated. Activation ticket txid: %s", activateTxID)
	log.WithContext(ctx).Infof("waiting confirmations for NFT Activation Ticket - Ticket txid: %s", activateTxID)

	// Wait until actTxid is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, activateTxID,
		int64(task.service.config.NFTActTxMinConfirmations), time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error validating nft-act ticket")
		task.StatusLog[common.FieldErrorDetail] = err.Error()
		task.UpdateStatus(common.StatusErrorValidateActivateTxnIDFailed)
		return errors.Errorf("wait reg-act ticket valid: %s", err.Error())
	}
	task.UpdateStatus(common.StatusTicketActivated)
	task.UpdateStatus(&common.EphemeralStatus{
		StatusTitle:   "Registered NFT Act Ticket TXID: ",
		StatusString:  task.regNFTTxid,
		IsFailureBool: false,
		IsFinalBool:   false,
	})
	log.WithContext(ctx).Infof("NFT Activation ticket is confirmed. Activation ticket txid: %s", activateTxID)

	return nil
}

// Download downloads the data from p2p for data validation before ticket activation
func (task *NftRegistrationTask) Download(ctx context.Context) error {
	// add the task to the worker queue, and worker will process the task in the background
	log.WithContext(ctx).WithField("nft_tx_id", task.regNFTTxid).Info("Downloading has been started")
	taskID := task.downloadService.AddTask(&nft.DownloadPayload{Txid: task.regNFTTxid, Pid: task.Request.CreatorPastelID, Key: task.Request.CreatorPastelIDPassphrase}, "")
	downloadTask := task.downloadService.GetTask(taskID)
	defer downloadTask.Cancel()

	sub := downloadTask.SubscribeStatus()
	log.WithContext(ctx).WithField("nft_tx_id", task.regNFTTxid).Info("Subscribed to status channel")

	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-sub():
			if status.IsFailure() {
				log.WithContext(ctx).WithField("nft_tx_id", task.regNFTTxid).WithError(task.Error())

				return errors.New("Download failed")
			}

			if status.IsFinal() {
				log.WithContext(ctx).WithField("nft_tx_id", task.regNFTTxid).Info("task has been downloaded successfully")
				return nil
			}
		case <-time.After(20 * time.Minute):
			log.WithContext(ctx).WithField("nft_tx_id", task.regNFTTxid).Info("Download request has been timed out")
			return errors.New("download request timeout, data validation failed")
		}
	}
}

func (task *NftRegistrationTask) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.service.pastelHandler.PastelClient.StorageNetworkFee(ctx)
	if err != nil {
		return false, err
	}
	task.StatusLog[common.FieldFee] = fee
	return fee <= task.Request.MaximumFee, nil
}

func (task *NftRegistrationTask) sendRegMetadata(ctx context.Context) error {
	if task.creatorBlockHash == "" {
		return errors.New("empty current block hash")
	}
	if task.Request.CreatorPastelID == "" {
		return errors.New("empty creator pastelID")
	}

	regMetadata := &types.NftRegMetadata{
		BlockHash:       task.creatorBlockHash,
		CreatorPastelID: task.Request.CreatorPastelID,
		BlockHeight:     strconv.Itoa(task.creatorBlockHeight),
		Timestamp:       task.creationTimestamp,
	}

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}
		group.Go(func() (err error) {
			err = nftRegNode.SendRegMetadata(gctx, regMetadata)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Controls walletnode interaction from steps 3-5
func (task *NftRegistrationTask) probeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	task.FingerprintsHandler.Clear()

	// Send image to supernodes for probing.
	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			// result is 4.B.5
			compress, stateOk, hashExists, err := nftRegNode.ProbeImage(gctx, file)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Errorf("probe image failed:%s", err.Error())
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			if hashExists {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Error("image already registered")
				return errors.Errorf("remote node %s: image already registered", someNode.String())
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed:")
				return errors.Errorf("remote node %s: indicated processing error:%w", someNode.String(), err)
			}
			// Step 5.1
			fingerprintAndScores, fingerprintAndScoresBytes, signature, err := pastel.ExtractCompressSignedDDAndFingerprints(compress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", someNode).Error("extract compressed signed DDAandFingerprints failed")
				return errors.Errorf("node %s: extract failed: %w", someNode.String(), err)
			}
			task.FingerprintsHandler.AddNew(fingerprintAndScores, fingerprintAndScoresBytes, signature, someNode.PastelID())

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}
	// Step 5.2 - 5.3
	if err := task.FingerprintsHandler.Match(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Error matching FPs:%s", err.Error())
		task.UpdateStatus(common.StatusErrorSignaturesNotMatch)
		log.WithContext(ctx).WithError(err).WithField("filename", fileName).Error("probe image failed")
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}
	task.UpdateStatus(common.StatusImageProbed)

	return nil
}

// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Step 6 - 8 for Walletnode
func (task *NftRegistrationTask) uploadImage(ctx context.Context) error {
	// Upload image with pqgsinganature and its thumb to supernodes
	if err := task.uploadImageWithThumbnail(ctx, task.ImageHandler.ImageEncodedWithFingerprints, task.Request.Thumbnail); err != nil {
		return errors.Errorf("upload encoded image and thumbnail coordinate: %w", err)
	}
	// Match thumbnail hashes receiveed from supernodes
	if err := task.ImageHandler.MatchThumbnailHashes(); err != nil {
		task.UpdateStatus(common.StatusErrorThumbnailHashesNotMatch)
		return errors.Errorf("thumbnail hash returns by supenodes not mached: %w", err)
	}
	task.UpdateStatus(common.StatusImageAndThumbnailUploaded)
	return nil
}

// uploadImageWithThumbnail uploads the image with pqsignatured appended and thumbnail's coordinate to super nodes
func (task *NftRegistrationTask) uploadImageWithThumbnail(ctx context.Context, file *files.File, thumbnail files.ThumbnailCoordinate) error {
	group, gctx := errgroup.WithContext(ctx)

	task.ImageHandler.ClearHashes()

	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			hash1, hash2, hash3, err := nftRegNode.UploadImageWithThumbnail(gctx, file, thumbnail)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", someNode).Error("upload image with thumbnail failed")
				return err
			}
			task.ImageHandler.AddNewHashes(hash1, hash2, hash3, someNode.PastelID())
			return nil
		})
	}
	return group.Wait()
}

func (task *NftRegistrationTask) createNftTicket(_ context.Context) error {
	if task.dataHash == nil {
		return common.ErrEmptyDatahash
	}
	if task.FingerprintsHandler.IsEmpty() {
		return common.ErrEmptyFingerprints
	}
	if task.RqHandler.IsEmpty() {
		return common.ErrEmptyRaptorQSymbols
	}
	if task.ImageHandler.IsEmpty() {
		return common.ErrEmptyPreviewHash
	}

	nftType := pastel.NFTTypeImage

	ticket := &pastel.NFTTicket{
		Version:   1,
		Author:    task.Request.CreatorPastelID,
		BlockNum:  task.creatorBlockHeight,
		BlockHash: task.creatorBlockHash,
		Copies:    task.Request.IssuedCopies,
		Royalty:   (task.Request.Royalty) / 100,
		Green:     task.Request.Green,
		AppTicketData: pastel.AppTicket{
			CreatorName:                task.Request.CreatorName,
			NFTTitle:                   utils.SafeString(&task.Request.Name),
			NFTSeriesName:              utils.SafeString(task.Request.SeriesName),
			NFTCreationVideoYoutubeURL: utils.SafeString(task.Request.YoutubeURL),
			NFTKeywordSet:              utils.SafeString(task.Request.Keywords),
			NFTType:                    nftType,
			CreatorWebsite:             utils.SafeString(task.Request.CreatorWebsiteURL),
			CreatorWrittenStatement:    utils.SafeString(task.Request.Description),
			TotalCopies:                task.Request.IssuedCopies,
			PreviewHash:                task.ImageHandler.PreviewHash,
			Thumbnail1Hash:             task.ImageHandler.MediumThumbnailHash,
			Thumbnail2Hash:             task.ImageHandler.SmallThumbnailHash,
			DataHash:                   task.dataHash,
			DDAndFingerprintsIc:        task.FingerprintsHandler.DDAndFingerprintsIc,
			DDAndFingerprintsMax:       task.service.config.DDAndFingerprintsMax,
			DDAndFingerprintsIDs:       task.FingerprintsHandler.DDAndFingerprintsIDs,
			RQIc:                       task.RqHandler.RQIDsIc,
			RQMax:                      task.service.config.RQIDsMax,
			RQIDs:                      task.RqHandler.RQIDs,
			RQOti:                      task.RqHandler.RQEncodeParams.Oti,
			OriginalFileSizeInBytes:    task.originalFileSizeInBytes,
			FileType:                   task.fileType,
			MakePubliclyAccessible:     task.Request.MakePubliclyAccessible,
		},
	}

	task.nftRegistrationTicket = ticket
	return nil
}

// Step 10
// WalletNode - Prepare and sign NFT Ticket
func (task *NftRegistrationTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeNFTTicket(task.nftRegistrationTicket)
	if err != nil {
		return errors.Errorf("encode ticket %w", err)
	}

	task.creatorSignature, err = task.service.pastelHandler.PastelClient.Sign(ctx, data, task.Request.CreatorPastelID, task.Request.CreatorPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket %w", err)
	}
	task.serializedTicket = data
	return nil
}

// Step 11.A WalletNode - Upload Signed Ticket; RaptorQ IDs file and dd_and_fingerprints file to the SuperNode’s 1, 2 and 3
// As part of step 12, WalletNode will check to make sure that the fee returned is not too high
func (task *NftRegistrationTask) sendSignedTicket(ctx context.Context) error {
	if task.serializedTicket == nil {
		return errors.Errorf("uploading ticket: serializedTicket is empty")
	}
	if task.creatorSignature == nil {
		return errors.Errorf("uploading ticket: creatorSignature is empty")
	}

	ddFpFile := task.FingerprintsHandler.DDAndFpFile
	rqidsFile := task.RqHandler.RQIDsFile
	encoderParams := task.RqHandler.RQEncodeParams

	var fees []int64
	var feesMtx sync.Mutex

	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		//ticketID := fmt.Sprintf("%s.%d.%s", task.Request.CreatorPastelID, task.creatorBlockHeight, hex.EncodeToString(task.dataHash))
		label := uuid.New().String()
		group.Go(func() error {
			fee, err := nftRegNode.SendSignedTicket(gctx, task.serializedTicket, task.creatorSignature, label, rqidsFile, ddFpFile, encoderParams)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Error("send signed ticket failed")
				return err
			}

			func() {
				feesMtx.Lock()
				defer feesMtx.Unlock()

				fees = append(fees, fee)
			}()

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("uploading ticket has failed: %w", err)
	}

	if len(fees) < 3 {
		return errors.Errorf("registration fees not received")
	}

	if fees[0] != fees[1] || fees[0] != fees[2] || fees[1] != fees[2] {
		return errors.Errorf("registration fees don't match")
	}

	// check if fee is over-expectation
	task.registrationFee = fees[0]

	if task.registrationFee > int64(task.Request.MaximumFee) {
		return errors.Errorf("fee too high: registration fee %d, maximum fee %d", task.registrationFee, int64(task.Request.MaximumFee))
	}

	return nil
}

func (task *NftRegistrationTask) registerActTicket(ctx context.Context) (string, error) {
	return task.service.pastelHandler.PastelClient.RegisterActTicket(ctx,
		task.regNFTTxid,
		task.creatorBlockHeight,
		task.registrationFee,
		task.Request.CreatorPastelID,
		task.Request.CreatorPastelIDPassphrase)
}

// Step 13 - 14
// WalletNode - Burn 10% of Registration fee
// WalletNode - Upload txid of burned transaction to the SuperNode’s 1, 2 and 3
func (task *NftRegistrationTask) preburnRegistrationFeeGetTicketTxid(ctx context.Context) error {
	burnTxid, err := task.service.pastelHandler.BurnSomeCoins(ctx, task.Request.SpendableAddress,
		task.registrationFee, 10)
	if err != nil {
		return fmt.Errorf("burn some coins: %w", err)
	}

	task.StatusLog[common.FieldBurnTxnID] = burnTxid
	group, gctx := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := nftRegNode.SendPreBurntFeeTxid(gctx, burnTxid)
			if err != nil {
				log.WithContext(gctx).WithError(err).WithField("node", nftRegNode).Error("send pre-burnt fee txid failed")
				return err
			}
			if !someNode.IsPrimary() && ticketTxid != "" && !task.skipPrimaryNodeTxidCheck() {

				return errors.Errorf("receive response %s from secondary node %s", ticketTxid, someNode.PastelID())
			}

			if someNode.IsPrimary() {
				if ticketTxid == "" {
					return errors.Errorf("primary node - %s, returned empty txid", someNode.PastelID())
				}
				task.regNFTTxid = ticketTxid
			}
			return nil
		})
	}
	return group.Wait()
}

// Error returns task err
func (task *NftRegistrationTask) Error() error {
	return task.WalletNodeTask.Error()
}

func (task *NftRegistrationTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
		task.RemoveFile(task.ImageHandler.ImageEncodedWithFingerprints)
	}
}

// NewNFTRegistrationTask returns a new Task instance.
func NewNFTRegistrationTask(service *NftRegistrationService, request *NftRegistrationRequest) *NftRegistrationTask {
	task := common.NewWalletNodeTask(logPrefix)
	meshHandlerOpts := common.MeshHandlerOpts{
		Task:          task,
		NodeMaker:     &RegisterNftNodeMaker{},
		PastelHandler: service.pastelHandler,
		NodeClient:    service.nodeClient,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay: service.config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   service.config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     service.config.AcceptNodesTimeout,
			MinSNs:                 service.config.NumberSuperNodes,
			PastelID:               request.CreatorPastelID,
			Passphrase:             request.CreatorPastelIDPassphrase,
			CheckDDDatabaseHashes:  true,
		},
	}

	return &NftRegistrationTask{
		WalletNodeTask:      task,
		service:             service,
		Request:             request,
		MeshHandler:         common.NewMeshHandler(meshHandlerOpts),
		ImageHandler:        mixins.NewImageHandler(service.pastelHandler),
		FingerprintsHandler: mixins.NewFingerprintsHandler(service.pastelHandler),
		downloadService:     service.downloadHandler,
		RqHandler: mixins.NewRQHandler(service.rqClient, service.pastelHandler,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir,
			service.config.NumberRQIDSFiles, service.config.RQIDsMax),
	}
}
