package nftregister

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/walletnode/services/common"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
)

// Registers an NFT on the blockchain
// Follow instructions from : https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration//
// NftRegistrationTask is Run from NftRegisterService.Run(), which eventually calls run, below
type NftRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler         *common.MeshHandler
	FingerprintsHandler *mixins.FingerprintsHandler
	ImageHandler        *mixins.NftImageHandler
	RqHandler           *mixins.RQHandler

	service *NftRegistrationService
	Request *NftRegisterRequest

	// task data to create RegArt ticket
	creatorBlockHeight int
	creatorBlockHash   string
	dataHash           []byte
	registrationFee    int64

	// ticket
	creatorSignature      []byte
	nftRegistrationTicket *pastel.NFTTicket
	serializedTicket      []byte

	regNFTTxid string
	burnTxid   string
}

// Run starts the task
func (task *NftRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

// Run sets up a connection to a mesh network of supernodes, then controls the communications to the mesh of nodes.
//	Task here will abstract away the individual node communications layer, and instead operate at the mesh control layer.
//  For individual communcations control, see node/grpc/nft_register.go
func (task *NftRegistrationTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.UpdateStatus(common.StatusErrorInsufficientFee)
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Request.MaximumFee)
	}

	// Setup mesh with supernodes with the highest ranks.
	creatorBlockHeight, creatorBlockHash, err := task.MeshHandler.SetupMeshOfNSupernodesNodes(ctx)
	if err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := task.MeshHandler.ConnectionsSupervisor(ctx, cancel)

	// send registration metadata
	if err := task.sendRegMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}

	// probe ORIGINAL image for average rareness, nsfw and seen score
	if err := task.probeImage(ctx, task.Request.Image, task.Request.Image.Name()); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// generateDDAndFingerprintsIDs generates dd & fp IDs
	if err := task.FingerprintsHandler.GenerateDDAndFingerprintsIDs(ctx, task.service.config.DDAndFingerprintsMax); err != nil {
		return errors.Errorf("DD and/or Fingerprint ID error: %w", err)
	}

	// Create copy of original image and embed fingerprints into it
	// result is in the - task.NftImageHandler
	if err := task.ImageHandler.CreateCopyWithEncodedFingerprint(ctx,
		task.Request.ArtistPastelID, task.Request.ArtistPastelID,
		task.FingerprintsHandler.FinalFingerprints, task.Request.Image); err != nil {
		return errors.Errorf("encode image with fingerprint: %w", err)
	}

	if task.dataHash, err = task.ImageHandler.GetHash(); err != nil {
		return errors.Errorf("get image hash: %w", err)
	}

	// upload image (from task.NftImageHandler) and thumbnail coordinates to supernode(s)
	// SN will return back: hashes, previews, thumbnails,
	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("upload image: %w", err)
	}

	// connect to rq serivce to get rq symbols identifier
	if err := task.RqHandler.GenRQIdentifiersFiles(ctx, task.ImageHandler.ImageEncodedWithFingerprints,
		task.creatorBlockHash, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase); err != nil {
		task.UpdateStatus(common.StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}

	if err := task.createNftTicket(ctx); err != nil {
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with artist signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("sign NFT ticket: %w", err)
	}

	// send signed ticket to supernodes to calculate registration fee
	if err := task.sendSignedTicket(ctx); err != nil {
		return errors.Errorf("send signed NFT ticket: %w", err)
	}

	// validate if address has enough psl
	if err := task.service.pastelHandler.CheckBalanceToPayRegistrationFee(ctx,
		task.Request.SpendableAddress,
		float64(task.registrationFee),
		task.Request.MaximumFee); err != nil {
		return errors.Errorf("check current balance: %w", err)
	}

	// send preburn-txid to master node(s)
	// master node will create reg-nft ticket and returns transaction id
	task.UpdateStatus(common.StatusPreburntRegistrationFee)
	if err := task.preburnRegistrationFeeGetTicketTxid(ctx); err != nil {
		return errors.Errorf("pre-burnt ten percent of registration fee: %w", err)
	}

	task.UpdateStatus(common.StatusTicketAccepted)

	// don't need SNs anymore
	_ = task.MeshHandler.CloseSNsConnections(ctx, nodesDone)

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regNFTTxid, int64(task.service.config.NFTRegTxMinConfirmations), 15*time.Second); err != nil {
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)

	// activate reg-nft ticket at previous step
	actTxid, err := task.registerActTicket(newCtx)
	if err != nil {
		return errors.Errorf("register act ticket: %w", err)
	}
	log.Debugf("reg-act-txid: %s", actTxid)

	// Wait until actTxid is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, actTxid, int64(task.service.config.NFTActTxMinConfirmations), 15*time.Second)
	if err != nil {
		return errors.Errorf("wait reg-act ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	log.Debugf("reg-act-tixd is confirmed")

	return nil
}

func (task *NftRegistrationTask) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.service.pastelHandler.PastelClient.StorageNetworkFee(ctx)
	if err != nil {
		return false, err
	}
	return fee <= task.Request.MaximumFee, nil
}

func (task *NftRegistrationTask) sendRegMetadata(ctx context.Context) error {
	if task.creatorBlockHash == "" {
		return errors.New("empty current block hash")
	}
	if task.Request.ArtistPastelID == "" {
		return errors.New("empty creator pastelID")
	}

	regMetadata := &types.NftRegMetadata{
		BlockHash:       task.creatorBlockHash,
		CreatorPastelID: task.Request.ArtistPastelID,
	}

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			err = nftRegNode.SendRegMetadata(ctx, regMetadata)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send registration metadata failed")
				return errors.Errorf("node %s: %w", someNode.String(), err)
			}

			return nil
		})
	}
	return group.Wait()
}

func (task *NftRegistrationTask) probeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	task.FingerprintsHandler.Clear()

	// Send image to supernodes for probing.
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() (err error) {
			compress, stateOk, err := nftRegNode.ProbeImage(ctx, file)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed")
				return errors.Errorf("node %s: probe failed :%w", someNode.String(), err)
			}

			someNode.SetRemoteState(stateOk)
			if !stateOk {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("probe image failed")
				return errors.Errorf("remote node %s: indicated processing error", someNode.String())
			}

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

	if err := task.FingerprintsHandler.Match(ctx); err != nil {
		task.UpdateStatus(common.StatusErrorSignaturesNotMatch)
		log.WithContext(ctx).WithError(err).WithField("filename", fileName).Error("probe image failed")
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}
	task.UpdateStatus(common.StatusImageProbed)

	return nil
}

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
	group, _ := errgroup.WithContext(ctx)

	task.ImageHandler.ClearHashes()

	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() error {
			hash1, hash2, hash3, err := nftRegNode.UploadImageWithThumbnail(ctx, file, thumbnail)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", someNode).Error("upload image with thumbnail failed")
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
	if task.RqHandler.IsEmpty() {
		return common.ErrEmptyRaptorQSymbols
	}
	if task.FingerprintsHandler.IsEmpty() {
		return common.ErrEmptyFingerprints
	}
	if task.ImageHandler.IsEmpty() {
		return common.ErrEmptyPreviewHash
	}

	nftType := pastel.NFTTypeImage

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.NFTTicket{
		Version:   1,
		Author:    task.Request.ArtistPastelID,
		BlockNum:  task.creatorBlockHeight,
		BlockHash: task.creatorBlockHash,
		Copies:    task.Request.IssuedCopies,
		Royalty:   0,     // Not supported yet by cNode
		Green:     false, // Not supported yet by cNode
		AppTicketData: pastel.AppTicket{
			CreatorName:                task.Request.ArtistName,
			CreatorWebsite:             utils.SafeString(task.Request.ArtistWebsiteURL),
			CreatorWrittenStatement:    utils.SafeString(task.Request.Description),
			NFTTitle:                   utils.SafeString(&task.Request.Name),
			NFTSeriesName:              utils.SafeString(task.Request.SeriesName),
			NFTCreationVideoYoutubeURL: utils.SafeString(task.Request.YoutubeURL),
			NFTKeywordSet:              utils.SafeString(task.Request.Keywords),
			NFTType:                    nftType,
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
		},
	}

	task.nftRegistrationTicket = ticket
	return nil
}

func (task *NftRegistrationTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeNFTTicket(task.nftRegistrationTicket)
	if err != nil {
		return errors.Errorf("encode ticket %w", err)
	}

	task.creatorSignature, err = task.service.pastelHandler.PastelClient.Sign(ctx, data, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket %w", err)
	}
	task.serializedTicket = data
	return nil
}

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
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		// TODO: Fix this when method to generate key1 and key2 are finalized
		key1 := "key1-" + uuid.New().String()
		key2 := "key2-" + uuid.New().String()
		group.Go(func() error {
			fee, err := nftRegNode.SendSignedTicket(ctx, task.serializedTicket, task.creatorSignature, key1, key2, rqidsFile, ddFpFile, encoderParams)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send signed ticket failed")
				return err
			}
			fees = append(fees, fee)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return errors.Errorf("uploading ticket has failed: %w", err)
	}

	if fees[0] != fees[1] || fees[0] != fees[2] || fees[1] != fees[2] {
		return errors.Errorf("registration fees don't match")
	}

	// check if fee is over-expection
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
		task.Request.ArtistPastelID,
		task.Request.ArtistPastelIDPassphrase)
}

func (task *NftRegistrationTask) preburnRegistrationFeeGetTicketTxid(ctx context.Context) error {

	task.service.pastelHandler.BurnSomeCoins(ctx, task.Request.SpendableAddress, task.registrationFee, 10)

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegisterNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegisterNode", someNode.String())
		}
		group.Go(func() error {
			ticketTxid, err := nftRegNode.SendPreBurntFeeTxid(ctx, task.burnTxid)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send pre-burnt fee txid failed")
				return err
			}
			if !someNode.IsPrimary() && ticketTxid != "" {
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

func (task *NftRegistrationTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
		task.RemoveFile(task.ImageHandler.ImageEncodedWithFingerprints)
	}
}

// NewNFTRegistrationTask returns a new Task instance.
func NewNFTRegistrationTask(service *NftRegistrationService, request *NftRegisterRequest) *NftRegistrationTask {
	task := &NftRegistrationTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		Request:        request,
	}

	task.ImageHandler = mixins.NewImageHandler(service.pastelHandler)

	task.MeshHandler = common.NewMeshHandler(task.WalletNodeTask,
		service.nodeClient, &RegisterNftNodeMaker{},
		service.pastelHandler,
		request.ArtistPastelID, request.ArtistPastelIDPassphrase,
		service.config.NumberSuperNodes, service.config.ConnectToNodeTimeout,
		service.config.AcceptNodesTimeout, service.config.ConnectToNextNodeDelay,
	)
	task.FingerprintsHandler = mixins.NewFingerprintsHandler(service.pastelHandler)
	task.RqHandler = mixins.NewRQHandler(service.rqClient,
		service.config.RaptorQServiceAddress, service.config.RqFilesDir, service.config.RQIDsMax,
		service.config.NumberRQIDSFiles)

	return task
}
