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
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/node"
)

// NftRegistrationTask is the task of registering new artwork.
type NftRegistrationTask struct {
	*common.WalletNodeTask

	service *NftRegisterService
	Request *NftRegisterRequest

	meshHandler *mixins.MeshHandler
	fpHandler   *mixins.FingerprintsHandler

	// task data to create RegArt ticket
	creatorBlockHeight           int
	creatorBlockHash             string
	imageEncodedWithFingerprints *files.File
	previewHash                  []byte
	mediumThumbnailHash          []byte
	smallThumbnailHash           []byte
	datahash                     []byte

	// signatures from SN1, SN2 & SN3 over dd_and_fingerprints data from dd-server
	signatures [][]byte

	rqids          []string
	rqEncodeParams rqnode.EncoderParameters
	rqIDsFile      []byte
	rqIDsIc        uint32

	// ticket
	creatorSignature []byte
	ticket           *pastel.NFTTicket

	// TODO: call cNodeAPI to get the following info
	//burnTxid        string
	//regNFTTxid      string
	//registrationFee int64
}

// Run starts the task
func (task *NftRegistrationTask) Run(ctx context.Context) error {
	_ = task.RunHelper(ctx, task.run, task.removeArtifacts)
	return nil
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
	if err := task.connectToTopRankNodes(ctx); err != nil {
		return errors.Errorf("connect to top rank nodes: %w", err)
	}

	// supervise the connection to top rank nodes
	// cancel any ongoing context if the connections are broken
	nodesDone := make(chan struct{})
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return task.nodes.WaitConnClose(ctx, nodesDone)
	})

	// get block height + hash
	blockNum, blockHash, err := task.GetBlock(ctx)
	if err != nil {
		return errors.Errorf("get current block heigth: %v", err)
	}
	task.creatorBlockHeight = blockNum
	task.creatorBlockHash = blockHash

	// send registration metadata
	if err := task.sendRegMetadata(ctx); err != nil {
		return errors.Errorf("send registration metadata: %w", err)
	}

	// probe image for rareness, nsfw and seen score
	if err := task.probeImage(ctx); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// generateDDAndFingerprintsIDs generates dd & fp IDs
	if err := task.generateDDAndFingerprintsIDs(); err != nil {
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

// generateDDAndFingerprintsIDs generates redundant IDs and assigns to task.redundantIDs
func (task *NftRegistrationTask) generateDDAndFingerprintsIDs() error {
	ddDataJSON, err := json.Marshal(task.fingerprintAndScores)
	if err != nil {
		return errors.Errorf("failed to marshal dd-data: %w", err)
	}
	task.fingerprint = ddDataJSON

	ddEncoded := utils.B64Encode(ddDataJSON)

	var buffer bytes.Buffer
	buffer.Write(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[0])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[1])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[2])
	ddFpFile := buffer.Bytes()

	task.ddAndFingerprintsIc = rand.Uint32()
	task.ddAndFingerprintsIDs, _, err = pastel.GetIDFiles(ddFpFile, task.ddAndFingerprintsIc, task.config.DDAndFingerprintsMax)
	if err != nil {
		return fmt.Errorf("get ID Files: %w", err)
	}

	comp, err := zstd.CompressLevel(nil, ddFpFile, 22)
	if err != nil {
		return errors.Errorf("compress: %w", err)
	}
	task.ddAndFpFile = utils.B64Encode(comp)

	return nil
}

func (task *NftRegistrationTask) encodeFingerprint(ctx context.Context, fingerprint []byte, img *files.File) error {
	// Sign fingerprint
	ed448PubKey, err := getPubKey(task.Request.ArtistPastelID)
	if err != nil {
		return fmt.Errorf("encodeFingerprint: %v", err)
	}

	ed448Signature, err := task.PastelClient.Sign(ctx,
		fingerprint,
		task.Request.ArtistPastelID,
		task.Request.ArtistPastelIDPassphrase,
		pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign fingerprint: %w", err)
	}

	ticket, err := task.PastelClient.FindTicketByID(ctx, task.Request.ArtistPastelID)
	if err != nil {
		return errors.Errorf("find register ticket of artist pastel id(%s):%w", task.Request.ArtistPastelID, err)
	}

	pqSignature, err := task.PastelClient.Sign(ctx,
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

// meshNodes establishes communication between supernodes.
func (task *NftRegistrationTask) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
	var meshNodes node.List
	secInfo := &alts.SecInfo{
		PastelID:   task.Request.ArtistPastelID,
		PassPhrase: task.Request.ArtistPastelIDPassphrase,
		Algorithm:  "ed448",
	}

	primary := nodes[primaryIndex]
	log.WithContext(ctx).Debugf("Trying to connect to primary node %q", primary)
	if err := primary.Connect(ctx, task.config.ConnectToNodeTimeout, secInfo); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
		return nil, err
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	// FIXME: ugly hack here. Need to make the Node and List to be safer
	secondariesMtx := &sync.Mutex{}
	var secondaries node.List
	go func() {
		for i, node := range nodes {
			node := node

			if i == primaryIndex {
				continue
			}

			select {
			case <-nextConnCtx.Done():
				return
			case <-time.After(task.config.connectToNextNodeDelay):
				go func() {
					defer errors.Recover(log.Fatal)

					if err := node.Connect(ctx, task.config.ConnectToNodeTimeout, secInfo); err != nil {
						return
					}
					if err := node.Session(ctx, false); err != nil {
						return
					}
					// Should not run this code in go routine
					func() {
						secondariesMtx.Lock()
						defer secondariesMtx.Unlock()
						secondaries.Add(node)
					}()

					if err := node.ConnectTo(ctx, types.MeshedSuperNode{
						NodeID: primary.PastelID(),
						SessID: primary.SessID(),
					}); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Seconary %q connected to primary", node)
				}()
			}
		}
	}()

	acceptCtx, acceptCancel := context.WithTimeout(ctx, task.config.acceptNodesTimeout)
	defer acceptCancel()

	accepted, err := primary.AcceptedNodes(acceptCtx)
	if err != nil {
		return nil, err
	}

	primary.SetPrimary(true)
	meshNodes.Add(primary)

	secondariesMtx.Lock()
	defer secondariesMtx.Unlock()

	for _, pastelID := range accepted {
		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

		node := secondaries.FindByPastelID(pastelID)
		if node == nil {
			return nil, errors.New("not found accepted node")
		}
		meshNodes.Add(node)
	}

	return meshNodes, nil
}

func (task *NftRegistrationTask) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.PastelClient.StorageNetworkFee(ctx)
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
	if task.datahash == nil {
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
			DataHash:                   task.datahash,
			DDAndFingerprintsIc:        task.ddAndFingerprintsIc,
			DDAndFingerprintsMax:       task.config.DDAndFingerprintsMax,
			DDAndFingerprintsIDs:       task.ddAndFingerprintsIDs,
			RQIc:                       task.rqIDsIc,
			RQMax:                      task.config.RQIDsMax,
			RQIDs:                      task.rqids,
			RQOti:                      task.rqEncodeParams.Oti,
		},
	}

	task.ticket = ticket
	return nil
}

func (task *NftRegistrationTask) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeNFTTicket(task.ticket)
	if err != nil {
		return errors.Errorf("encode ticket %w", err)
	}

	task.creatorSignature, err = task.PastelClient.Sign(ctx, data, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket %w", err)
	}
	return nil
}

func (task *NftRegistrationTask) registerActTicket(ctx context.Context) (string, error) {
	return task.PastelClient.RegisterActTicket(ctx,
		task.regNFTTxid,
		task.creatorBlockHeight,
		task.registrationFee,
		task.Request.ArtistPastelID,
		task.Request.ArtistPastelIDPassphrase)
}

func (task *NftRegistrationTask) connectToTopRankNodes(ctx context.Context) error {
	// Retrieve supernodes with highest ranks.

	var topNodes node.List
	err := task.GetTopNodes(ctx, &topNodes, task.NftRegisterService.nodeClient, 0)
	if err != nil {
		return errors.Errorf("call masternode top: %w", err)
	}

	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(common.StatusErrorInsufficientFee)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	// Try to create mesh of supernodes, connecting to all supernodes in a different sequences.
	var nodes node.List
	var errs error

	for primaryRank := range topNodes {
		nodes, err = task.meshNodes(ctx, topNodes, primaryRank)
		if err != nil {
			// close connected connections
			topNodes.DisconnectAll()

			if errors.IsContextCanceled(err) {
				return err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warn("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(nodes) < task.config.NumberSuperNodes {
		// close connected connections
		topNodes.DisconnectAll()
		return errors.Errorf("Could not create a mesh of %d nodes: %w", task.config.NumberSuperNodes, errs)
	}

	// Activate supernodes that are in the mesh.
	nodes.Activate()
	// Disconnect supernodes that are not involved in the process.
	topNodes.DisconnectInactive()

	// Cancel context when any connection is broken.
	task.UpdateStatus(common.StatusConnected)
	task.nodes = nodes

	// Send all meshed supernode info to nodes - that will be used to node send info to other nodes
	meshedSNInfo := []types.MeshedSuperNode{}
	for _, node := range nodes {
		meshedSNInfo = append(meshedSNInfo, types.MeshedSuperNode{
			NodeID: node.PastelID(),
			SessID: node.SessID(),
		})
	}

	for _, node := range nodes {
		err = node.MeshNodes(ctx, meshedSNInfo)
		if err != nil {
			nodes.DisconnectAll()
			return errors.Errorf("could not send info of meshed nodes: %w", err)
		}
	}
	return nil
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

	return task.nodes.SendRegMetadata(ctx, regMetadata)
}

func (task *NftRegistrationTask) probeImage(ctx context.Context) error {
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debug("probe image")

	// Send image to supernodes for probing.
	if err := task.nodes.ProbeImage(ctx, task.Request.Image); err != nil {
		return errors.Errorf("send image: %w", err)
	}
	signatures := [][]byte{}
	// Match signatures received from supernodes.
	for i := 0; i < len(task.nodes); i++ {
		verified, err := task.PastelClient.Verify(ctx,
			task.nodes[i].FingerprintAndScoresBytes,
			string(task.nodes[i].Signature),
			task.nodes[i].PastelID(),
			pastel.SignAlgorithmED448)
		if err != nil {
			return errors.Errorf("probeImage: pastelClient.Verify %w", err)
		}

		if !verified {
			task.UpdateStatus(common.StatusErrorSignaturesNotMatch)
			return errors.Errorf("node[%s] signature doesn't match", task.nodes[i].PastelID())
		}

		signatures = append(signatures, task.nodes[i].Signature)
	}
	task.signatures = signatures

	// Match fingerprints received from supernodes.
	if err := task.nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(common.StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}

	task.fingerprintAndScores = task.nodes.FingerAndScores()
	task.UpdateStatus(common.StatusImageProbed)

	return nil
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

	if task.datahash, err = utils.Sha3256hash(imgBytes); err != nil {
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
	buf, err := pastel.EncodeNFTTicket(task.ticket)
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
