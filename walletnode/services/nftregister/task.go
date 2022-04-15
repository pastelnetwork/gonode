package nftregister

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"sync"
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

// NftRegistrationTask an NFT on the blockchain
// Follow instructions from : https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration//
// NftRegistrationTask is Run from NftRegisterService.Run(), which eventually calls run, below
type NftRegistrationTask struct {
	*common.WalletNodeTask

	MeshHandler         *common.MeshHandler
	FingerprintsHandler *mixins.FingerprintsHandler
	ImageHandler        *mixins.NftImageHandler
	RqHandler           *mixins.RQHandler

	service *NftRegistrationService
	Request *NftRegistrationRequest

	// task data to create RegArt ticket
	creatorBlockHeight int
	creatorBlockHash   string
	creationTimestamp  string
	dataHash           []byte
	registrationFee    int64

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
	//set timestamp for nft reg metadata
	task.creationTimestamp = time.Now().Format("YYYY-MM-DD hh:mm:ss")

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

	// Calculate IDs for redundant number of dd_and_fingerprints files
	// Step 5.4
	if err := task.FingerprintsHandler.GenerateDDAndFingerprintsIDs(ctx, task.service.config.DDAndFingerprintsMax); err != nil {
		return errors.Errorf("DD and/or Fingerprint ID error: %w", err)
	}

	// Create copy of original image and embed fingerprints into it
	// result is in the - task.NftImageHandler
	// Step 6
	if err := task.ImageHandler.CreateCopyWithEncodedFingerprint(ctx,
		task.Request.CreatorPastelID, task.Request.CreatorPastelIDPassphrase,
		task.FingerprintsHandler.FinalFingerprints, task.Request.Image); err != nil {
		return errors.Errorf("encode image with fingerprint: %w", err)
	}

	if task.dataHash, err = task.ImageHandler.GetHash(); err != nil {
		return errors.Errorf("get image hash: %w", err)
	}

	// upload image (from task.NftImageHandler) and thumbnail coordinates to supernode(s)
	// SN will return back: hashes, previews, thumbnails,
	// Step 6 - 8 for Walletnode
	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("upload image: %w", err)
	}

	// connect to rq serivce to get rq symbols identifier
	// All of step 9
	if err := task.RqHandler.GenRQIdentifiersFiles(ctx, task.ImageHandler.ImageEncodedWithFingerprints,
		task.creatorBlockHash, task.Request.CreatorPastelID, task.Request.CreatorPastelIDPassphrase); err != nil {
		task.UpdateStatus(common.StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}

	if err := task.createNftTicket(ctx); err != nil {
		return errors.Errorf("create ticket: %w", err)
	}

	// sign ticket with artist signature
	// Step 10
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("sign NFT ticket: %w", err)
	}

	// send signed ticket to supernodes to calculate registration fee
	// Step 11.A WalletNode - Upload Signed Ticket; RaptorQ IDs file and dd_and_fingerprints file to the SuperNode’s 1, 2 and 3
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
	//Start Step 20

	if err := task.service.pastelHandler.WaitTxidValid(newCtx, task.regNFTTxid, int64(task.service.config.NFTRegTxMinConfirmations),
		time.Duration(task.service.config.WaitTxnValidInterval)*time.Second); err != nil {
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketRegistered)

	// activate reg-nft ticket at previous step
	// Send Activation ticket and Registration Fee (cNode API)
	actTxid, err := task.registerActTicket(newCtx)
	if err != nil {
		return errors.Errorf("register act ticket: %w", err)
	}
	log.Debugf("reg-act-txid: %s", actTxid)

	// Wait until actTxid is valid
	err = task.service.pastelHandler.WaitTxidValid(newCtx, actTxid,
		int64(task.service.config.NFTActTxMinConfirmations), time.Duration(task.service.config.WaitTxnValidInterval)*time.Second)
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
	if task.Request.CreatorPastelID == "" {
		return errors.New("empty creator pastelID")
	}

	regMetadata := &types.NftRegMetadata{
		BlockHash:       task.creatorBlockHash,
		CreatorPastelID: task.Request.CreatorPastelID,
		BlockHeight:     strconv.Itoa(task.creatorBlockHeight),
		Timestamp:       task.creationTimestamp,
	}

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
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

// https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
// Controls walletnode interaction from steps 3-5
func (task *NftRegistrationTask) probeImage(ctx context.Context, file *files.File, fileName string) error {
	log.WithContext(ctx).WithField("filename", fileName).Debug("probe image")

	task.FingerprintsHandler.Clear()

	// Send image to supernodes for probing.
	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() (err error) {
			// result is 4.B.5
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
		task.UpdateStatus(common.StatusErrorSignaturesNotMatch)
		log.WithContext(ctx).WithError(err).WithField("filename", fileName).Error("probe image failed")
		return errors.Errorf("probing image %s failed: %w", fileName, err)
	}
	task.UpdateStatus(common.StatusImageProbed)

	return nil
}

//https://pastel.wiki/en/Architecture/Workflows/NewArtRegistration
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
	group, _ := errgroup.WithContext(ctx)

	task.ImageHandler.ClearHashes()

	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
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

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.NFTTicket{
		Version:   1,
		Author:    task.Request.CreatorPastelID,
		BlockNum:  task.creatorBlockHeight,
		BlockHash: task.creatorBlockHash,
		Copies:    task.Request.IssuedCopies,
		Royalty:   0,     // Not supported yet by cNode
		Green:     false, // Not supported yet by cNode
		AppTicketData: pastel.AppTicket{
			CreatorName:                task.Request.CreatorName,
			CreatorWebsite:             utils.SafeString(task.Request.CreatorWebsiteURL),
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

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			//TODO: use assert here
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		ticketID := fmt.Sprintf("%s.%d.%s", task.Request.CreatorPastelID, task.creatorBlockHeight, hex.EncodeToString(task.dataHash))
		key1 := ticketID
		key2 := uuid.New().String()
		group.Go(func() error {
			fee, err := nftRegNode.SendSignedTicket(ctx, task.serializedTicket, task.creatorSignature, key1, key2, rqidsFile, ddFpFile, encoderParams)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send signed ticket failed")
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

	group, _ := errgroup.WithContext(ctx)
	for _, someNode := range task.MeshHandler.Nodes {
		nftRegNode, ok := someNode.SuperNodeAPIInterface.(*NftRegistrationNode)
		if !ok {
			return errors.Errorf("node %s is not NftRegistrationNode", someNode.String())
		}

		someNode := someNode
		group.Go(func() error {
			ticketTxid, err := nftRegNode.SendPreBurntFeeTxid(ctx, burnTxid)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("node", nftRegNode).Error("send pre-burnt fee txid failed")
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
		},
	}

	return &NftRegistrationTask{
		WalletNodeTask:      task,
		service:             service,
		Request:             request,
		MeshHandler:         common.NewMeshHandler(meshHandlerOpts),
		ImageHandler:        mixins.NewImageHandler(service.pastelHandler),
		FingerprintsHandler: mixins.NewFingerprintsHandler(service.pastelHandler),
		RqHandler: mixins.NewRQHandler(service.rqClient, service.pastelHandler,
			service.config.RaptorQServiceAddress, service.config.RqFilesDir,
			service.config.NumberRQIDSFiles, service.config.RQIDsMax),
	}
}
