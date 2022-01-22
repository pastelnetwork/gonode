package artworkregister

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	"math/rand"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// NftRegistrationTask is the task of registering new artwork.
type NftRegistrationTask struct {
	*common.WalletNodeTask

	nftHandler *NftRegisterHandler

	service *NftRegisterService
	Request *NftRegisterRequest

	// task data to create RegArt ticket
	creatorBlockHeight           int
	creatorBlockHash             string
	imageEncodedWithFingerprints *files.File
	previewHash                  []byte
	mediumThumbnailHash          []byte
	smallThumbnailHash           []byte
	dataHash                     []byte
	registrationFee              int64

	rqids          []string
	rqEncodeParams rqnode.EncoderParameters
	rqIDsFile      []byte
	rqIDsIc        uint32

	// ticket
	creatorSignature      []byte
	nftRegistrationTicket *pastel.NFTTicket

	regNFTTxid string
}

// Run starts the task
func (task *NftRegistrationTask) Run(ctx context.Context) error {
	return task.RunHelper(ctx, task.run, task.removeArtifacts)
}

func (task *NftRegistrationTask) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.UpdateStatus(common.StatusErrorInsufficientFee)
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Request.MaximumFee)
	}

	// Retrieve supernodes with highest ranks.
	creatorBlockHeight, creatorBlockHash, err := task.nftHandler.meshHandler.ConnectToTopRankNodes(ctx)
	if err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}
	task.creatorBlockHeight = creatorBlockHeight
	task.creatorBlockHash = creatorBlockHash

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return task.nftHandler.meshHandler.Nodes.WaitConnClose(ctx, nodesDone)
	})

	// send registration metadata
	if err := task.sendRegMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}

	// probe image for average rareness, nsfw and seen score
	if err := task.nftHandler.ProbeImage(ctx, task.Request.Image, task.Request.Image.Name()); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// generateDDAndFingerprintsIDs generates dd & fp IDs
	if err := task.nftHandler.fpHandler.GenerateDDAndFingerprintsIDs(ctx, task.service.config.DDAndFingerprintsMax); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// upload image and thumbnail coordinates to supernode(s)
	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("upload image: %w", err)
	}

	// connect to rq serivce to get rq symbols identifier
	if err := task.genRQIdentifiersFiles(ctx); err != nil {
		return errors.Errorf("gen RaptorQ symbols' identifiers: %w", err)
	}

	if err := task.createArtTicket(ctx); err != nil {
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
	if err := task.CheckRegistrationFee(ctx,
		task.Request.SpendableAddress,
		float64(task.registrationFee),
		task.Request.MaximumFee); err != nil {
		return errors.Errorf("check current balance: %w", err)
	}

	// send preburn-txid to master node(s)
	// master node will create reg-art ticket and returns transaction id
	task.UpdateStatus(common.StatusPreburntRegistrationFee)
	if err := task.preburntRegistrationFee(ctx); err != nil {
		return errors.Errorf("pre-burnt ten percent of registration fee: %w", err)
	}

	task.UpdateStatus(common.StatusTicketAccepted)

	log.WithContext(ctx).Debug("close connections to supernodes")
	close(nodesDone)
	for i := range task.nodes {
		if err := task.nodes[i].ConnectionInterface.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": task.nodes[i].PastelID(),
				"addr":     task.nodes[i].String(),
			}).WithError(err).Errorf("close supernode connection failed")
		}
	}

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.PastelHandler.WaitTxidValid(newCtx, task.regNFTTxid, int64(task.config.RegArtTxMinConfirmations), 15*time.Second); err != nil {
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}

	task.UpdateStatus(common.StatusTicketRegistered)

	// activate reg-art ticket at previous step
	actTxid, err := task.registerActTicket(newCtx)
	if err != nil {
		return errors.Errorf("register act ticket: %w", err)
	}
	log.Debugf("reg-act-txid: %s", actTxid)

	// Wait until actTxid is valid
	err = task.PastelHandler.WaitTxidValid(newCtx, actTxid, int64(task.config.RegActTxMinConfirmations), 15*time.Second)
	if err != nil {
		return errors.Errorf("wait reg-act ticket valid: %w", err)
	}
	task.UpdateStatus(common.StatusTicketActivated)
	log.Debugf("reg-act-tixd is confirmed")

	return nil
}

func (task *NftRegistrationTask) encodeFingerprint(ctx context.Context, fingerprint []byte, img *files.File) error {
	// Sign fingerprint
	ed448PubKey, err := getPubKey(task.Request.ArtistPastelID)
	if err != nil {
		return fmt.Errorf("encodeFingerprint: %v", err)
	}

	ed448Signature, err := task.service.pastelHandler.PastelClient.Sign(ctx,
		fingerprint,
		task.Request.ArtistPastelID,
		task.Request.ArtistPastelIDPassphrase,
		pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign fingerprint: %w", err)
	}

	ticket, err := task.service.pastelHandler.PastelClient.FindTicketByID(ctx, task.Request.ArtistPastelID)
	if err != nil {
		return errors.Errorf("find register ticket of artist pastel id(%s):%w", task.Request.ArtistPastelID, err)
	}

	pqSignature, err := task.service.pastelHandler.PastelClient.Sign(ctx,
		fingerprint,
		task.Request.ArtistPastelID,
		task.Request.ArtistPastelIDPassphrase,
		pastel.SignAlgorithmLegRoast)
	if err != nil {
		return errors.Errorf("sign fingerprint with legroats: %w", err)
	}

	pqPubKey, err := getPubKey(ticket.PqKey)
	if err != nil {
		return fmt.Errorf("encodeFingerprint: %v", err)
	}

	// Encode data to the image.
	encSig := qrsignature.New(
		qrsignature.Fingerprint(fingerprint),
		qrsignature.PostQuantumSignature(pqSignature),
		qrsignature.PostQuantumPubKey(pqPubKey),
		qrsignature.Ed448Signature(ed448Signature),
		qrsignature.Ed448PubKey(ed448PubKey),
	)

	return img.Encode(encSig)
}

func (task *NftRegistrationTask) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.service.pastelHandler.PastelClient.StorageNetworkFee(ctx)
	if err != nil {
		return false, err
	}
	return fee <= task.Request.MaximumFee, nil
}

func (task *NftRegistrationTask) genRQIdentifiersFiles(ctx context.Context) error {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("connect to raptorQ: %w", err)
	}
	defer conn.Close()

	content, err := task.imageEncodedWithFingerprints.Bytes()
	if err != nil {
		return errors.Errorf("read image content: %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.NftRegisterService.config.RqFilesDir,
	})

	// FIXME :
	// - check format of artis block hash should be base58 or not
	encodeInfo, err := rqService.EncodeInfo(ctx, content, task.config.NumberRQIDSFiles, task.creatorBlockHash, task.Request.ArtistPastelID)
	if err != nil {
		return errors.Errorf("generate RaptorQ symbols identifiers: %w", err)
	}

	files := 0
	for _, rawSymbolIDFile := range encodeInfo.SymbolIDFiles {
		err := task.generateRQIDs(ctx, rawSymbolIDFile)
		if err != nil {
			task.UpdateStatus(common.StatusErrorGenRaptorQSymbolsFailed)
			return errors.Errorf("create rqids file :%w", err)
		}
		files++
		break
	}

	if files != 1 {
		return errors.Errorf("number of raptorq symbol identifiers files must be greater than 1")
	}
	task.rqEncodeParams = encodeInfo.EncoderParam

	return nil
}

func (task *NftRegistrationTask) generateRQIDs(ctx context.Context, rawFile rqnode.RawSymbolIDFile) error {
	file, err := json.Marshal(rawFile)
	if err != nil {
		return fmt.Errorf("marshal rqID file")
	}

	signature, err := task.PastelClient.Sign(ctx,
		file,
		task.Request.ArtistPastelID,
		task.Request.ArtistPastelIDPassphrase,
		pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign identifiers file: %w", err)
	}

	encfile := utils.B64Encode(file)

	var buffer bytes.Buffer
	buffer.Write(encfile)
	buffer.WriteString(".")
	buffer.Write(signature)
	rqIDFile := buffer.Bytes()

	task.rqIDsIc = rand.Uint32()
	task.rqids, _, err = pastel.GetIDFiles(rqIDFile, task.rqIDsIc, task.config.RQIDsMax)
	if err != nil {
		return fmt.Errorf("get ID Files: %w", err)
	}

	comp, err := zstd.CompressLevel(nil, rqIDFile, 22)
	if err != nil {
		return errors.Errorf("compress: %w", err)
	}
	task.rqIDsFile = utils.B64Encode(comp)

	return nil
}

func (task *NftRegistrationTask) createArtTicket(_ context.Context) error {
	if task.fingerprint == nil {
		return errEmptyFingerprints
	}
	if task.dataHash == nil {
		return errEmptyDatahash
	}
	if task.previewHash == nil {
		return errEmptyPreviewHash
	}
	if task.mediumThumbnailHash == nil {
		return errEmptyMediumThumbnailHash
	}
	if task.smallThumbnailHash == nil {
		return errEmptySmallThumbnailHash
	}
	if task.rqids == nil {
		return errEmptyRaptorQSymbols
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
			PreviewHash:                task.previewHash,
			Thumbnail1Hash:             task.mediumThumbnailHash,
			Thumbnail2Hash:             task.smallThumbnailHash,
			DataHash:                   task.dataHash,
			DDAndFingerprintsIc:        task.ddAndFingerprintsIc,
			DDAndFingerprintsMax:       task.config.DDAndFingerprintsMax,
			DDAndFingerprintsIDs:       task.ddAndFingerprintsIDs,
			RQIc:                       task.rqIDsIc,
			RQMax:                      task.config.RQIDsMax,
			RQIDs:                      task.rqids,
			RQOti:                      task.rqEncodeParams.Oti,
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

	return task.nftHandler.SendRegMetadata(ctx, regMetadata)
}

func (task *NftRegistrationTask) uploadImage(ctx context.Context) error {
	img1, err := task.Request.Image.Copy()
	if err != nil {
		return errors.Errorf("copy image to encode: %w", err)
	}

	log.WithContext(ctx).WithField("FileName", img1.Name()).Debug("final image")
	if err := task.encodeFingerprint(ctx, task.fingerprint, img1); err != nil {
		return errors.Errorf("encode image with fingerprint: %w", err)
	}
	task.imageEncodedWithFingerprints = img1

	imgBytes, err := img1.Bytes()
	if err != nil {
		return errors.Errorf("convert image to byte stream %w", err)
	}

	if task.dataHash, err = utils.Sha3256hash(imgBytes); err != nil {
		return errors.Errorf("hash encoded image: %w", err)
	}

	// Upload image with pqgsinganature and its thumb to supernodes
	if err := task.nodes.UploadImageWithThumbnail(ctx, img1, task.Request.Thumbnail); err != nil {
		return errors.Errorf("upload encoded image and thumbnail coordinate: %w", err)
	}
	// Match thumbnail hashes receiveed from supernodes
	if err := task.nodes.MatchThumbnailHashes(); err != nil {
		task.UpdateStatus(common.StatusErrorThumbnailHashesNotMatch)
		return errors.Errorf("thumbnail hash returns by supenodes not mached: %w", err)
	}
	task.UpdateStatus(common.StatusImageAndThumbnailUploaded)

	task.previewHash = task.nodes.PreviewHash()
	task.mediumThumbnailHash = task.nodes.MediumThumbnailHash()
	task.smallThumbnailHash = task.nodes.SmallThumbnailHash()

	return nil
}

func (task *NftRegistrationTask) sendSignedTicket(ctx context.Context) error {
	buf, err := pastel.EncodeNFTTicket(task.nftRegistrationTicket)
	if err != nil {
		return errors.Errorf("marshal ticket: %w", err)
	}
	log.Debug(string(buf))

	if err := task.nodes.UploadSignedTicket(ctx, buf, task.creatorSignature, task.rqIDsFile, task.ddAndFpFile, task.rqEncodeParams); err != nil {
		return errors.Errorf("upload signed ticket: %w", err)
	}

	if err := task.nodes.MatchRegistrationFee(); err != nil {
		return errors.Errorf("registration fees don't matched: %w", err)
	}

	// check if fee is over-expection
	task.registrationFee = task.nodes.RegistrationFee()

	if task.registrationFee > int64(task.Request.MaximumFee) {
		return errors.Errorf("fee too high: registration fee %d, maximum fee %d", task.registrationFee, int64(task.Request.MaximumFee))
	}
	task.registrationFee = task.nodes.RegistrationFee()

	return nil
}

func (task *NftRegistrationTask) preburntRegistrationFee(ctx context.Context) error {
	if task.registrationFee <= 0 {
		return errors.Errorf("invalid registration fee")
	}

	burnedAmount := float64(task.registrationFee) / 10
	burnTxid, err := task.PastelClient.SendFromAddress(ctx, task.Request.SpendableAddress, task.config.BurnAddress, burnedAmount)
	if err != nil {
		return errors.Errorf("burn 10 percent of transaction fee: %w", err)
	}
	task.burnTxid = burnTxid
	log.WithContext(ctx).Debugf("preburn txid: %s", task.burnTxid)

	if err := task.nodes.SendPreBurntFeeTxid(ctx, task.burnTxid); err != nil {
		return errors.Errorf("send pre-burn-txid: %s to supernode(s): %w", task.burnTxid, err)
	}
	task.regNFTTxid = task.nodes.RegArtTicketID()
	if task.regNFTTxid == "" {
		return errors.Errorf("empty regNFTTxid")
	}

	return nil
}

func (task *NftRegistrationTask) removeArtifacts() {
	if task.Request != nil {
		task.RemoveFile(task.Request.Image)
		task.RemoveFile(task.imageEncodedWithFingerprints)
	}
}

// NewNFTRegistrationTask returns a new Task instance.
func NewNFTRegistrationTask(service *NftRegisterService, request *NftRegisterRequest) *NftRegistrationTask {
	return &NftRegistrationTask{
		WalletNodeTask: common.NewWalletNodeTask(logPrefix),
		service:        service,
		Request:        request,
		meshHandler:    mixins.NewMeshHandler(service.nodeClient),
	}
}
