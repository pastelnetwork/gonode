package artworkregister

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/node"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Request *Request

	// information of nodes
	nodes node.List

	// task data to create RegArt ticket
	creatorBlockHeight           int
	creatorBlockHash             string
	fingerprintAndScores         *pastel.DDAndFingerprints
	fingerprintsHash             []byte
	fingerprint                  []byte
	imageEncodedWithFingerprints *artwork.File
	previewHash                  []byte
	mediumThumbnailHash          []byte
	smallThumbnailHash           []byte
	datahash                     []byte
	rqids                        []string
	rqEncodeParams               rqnode.EncoderParameters

	// signatures from SN1, SN2 & SN3 over dd_and_fingerprints data from dd-server
	signatures [][]byte

	ddAndFingerprintsIc uint32
	rqIDsIc             uint32

	// ddAndFingerprintsIDs are Base58(SHA3_256(compressed(Base64URL(dd_and_fingerprints).
	// Base64URL(signatureSN1).Base64URL(signatureSN2).Base64URL(signatureSN3).dd_and_fingerprints_ic)))
	ddAndFingerprintsIDs []string

	// ddAndFpFile is Base64(compressed(Base64(dd_and_fingerprints).Base64(signatureSN1).Base64(signatureSN3)))
	ddAndFpFile []byte
	// rqIDsFile is Base64(compress(Base64(rq_ids).Base64(signature)))
	rqIDsFile []byte

	// TODO: call cNodeAPI to get the following info
	burnTxid        string
	regNFTTxid      string
	registrationFee int64

	// ticket
	creatorSignature []byte
	ticket           *pastel.NFTTicket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	defer task.removeArtifacts()
	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Error("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

func (task *Task) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.UpdateStatus(ErrorInsufficientFee)
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
	if err := task.getBlock(ctx); err != nil {
		return errors.Errorf("get current block heigth: %w", err)
	}

	// probe image for rareness, nsfw and seen score
	if err := task.probeImage(ctx); err != nil {
		return errors.Errorf("probe image: %w", err)
	}

	// generateredundantIDs generate   number of RQ IDs
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
	if err := task.checkCurrentBalance(ctx, float64(task.registrationFee)); err != nil {
		return errors.Errorf("check current balance: %w", err)
	}

	// send preburn-txid to master node(s)
	// master node will create reg-art ticket and returns transaction id
	task.UpdateStatus(StatusPreburntRegistrationFee)
	if err := task.preburntRegistrationFee(ctx); err != nil {
		return errors.Errorf("pre-burnt ten percent of registration fee: %w", err)
	}

	task.UpdateStatus(StatusTicketAccepted)

	log.WithContext(ctx).Debug("close connections to supernodes")
	close(nodesDone)
	for i := range task.nodes {
		if err := task.nodes[i].Connection.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": task.nodes[i].PastelID(),
				"addr":     task.nodes[i].String(),
			}).WithError(err).Errorf("close supernode connection failed")
		}
	}

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.waitTxidValid(newCtx, task.regNFTTxid, int64(task.config.RegArtTxMinConfirmations), 15*time.Second); err != nil {
		return errors.Errorf("wait reg-nft ticket valid: %w", err)
	}

	task.UpdateStatus(StatusTicketRegistered)

	// activate reg-art ticket at previous step
	actTxid, err := task.registerActTicket(newCtx)
	if err != nil {
		return errors.Errorf("register act ticket: %w", err)
	}
	log.Debugf("reg-act-txid: %s", actTxid)

	// Wait until actTxid is valid
	err = task.waitTxidValid(newCtx, actTxid, int64(task.config.RegActTxMinConfirmations), 15*time.Second)
	if err != nil {
		return errors.Errorf("wait reg-act ticket valid: %w", err)
	}
	task.UpdateStatus(StatusTicketActivated)
	log.Debugf("reg-act-tixd is confirmed")

	return nil
}

// generateDDAndFingerprintsIDs generates redundant IDs and assigns to task.redundantIDs
func (task *Task) generateDDAndFingerprintsIDs() error {
	ddDataJson, err := json.Marshal(task.fingerprintAndScores)
	if err != nil {
		return errors.Errorf("failed to marshal dd-data: %w", err)
	}

	ddEncoded := base64.StdEncoding.EncodeToString(ddDataJson)

	var buffer bytes.Buffer
	buffer.WriteString(ddEncoded)
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[0])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[1])
	buffer.WriteByte(pastel.SeparatorByte)
	buffer.Write(task.signatures[2])
	ddFpFile := buffer.Bytes()

	task.ddAndFingerprintsIc = rand.Uint32()
	task.ddAndFingerprintsIDs, err = pastel.GetIDFiles(ddFpFile, task.ddAndFingerprintsIc, task.config.DDAndFingerprintsMax)
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

func (task *Task) checkCurrentBalance(ctx context.Context, registrationFee float64) error {
	if registrationFee == 0 {
		return errors.Errorf("invalid registration fee: %f", registrationFee)
	}

	if registrationFee > task.Request.MaximumFee {
		return errors.Errorf("registration fee is to expensive - maximum-fee (%f) < registration-fee(%f)", task.Request.MaximumFee, registrationFee)
	}

	balance, err := task.pastelClient.GetBalance(ctx, task.Request.SpendableAddress)
	if err != nil {
		return errors.Errorf("get balance of address(%s): %w", task.Request.SpendableAddress, err)
	}

	if balance < registrationFee {
		return errors.Errorf("not enough PSL - balance(%f) < registration-fee(%f)", balance, registrationFee)
	}

	return nil
}

func (task *Task) waitTxidValid(ctx context.Context, txID string, expectedConfirms int64, interval time.Duration) error {
	log.WithContext(ctx).Debugf("Need %d confirmation for txid %s", expectedConfirms, txID)
	blockTracker := blocktracker.New(task.pastelClient)
	baseBlkCnt, err := blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done: %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := task.pastelClient.GetRawTransactionVerbose1(subCtx, txID)
				if err != nil {
					return errors.Errorf("get transaction: %w", err)
				}

				if result.Confirmations >= expectedConfirms {
					return nil
				}

				return errors.Errorf("not enough confirmations: expected %d, got %d", expectedConfirms, result.Confirmations)
			}

			err := checkConfirms()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("check confirmations failed")
			} else {
				return nil
			}

			currentBlkCnt, err := blockTracker.GetBlockCount()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("failed to get block count")
				continue
			}

			if currentBlkCnt-baseBlkCnt >= int32(expectedConfirms)+2 {
				return errors.Errorf("timeout when wating for confirmation of transaction %s", txID)
			}
		}
	}
}

func (task *Task) encodeFingerprint(ctx context.Context, fingerprint []byte, img *artwork.File) error {
	// Sign fingerprint
	ed448PubKey, err := getPubKey(task.Request.ArtistPastelID)
	if err != nil {
		return fmt.Errorf("encodeFingerprint: %v", err)
	}

	ed448Signature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign fingerprint: %w", err)
	}

	ticket, err := task.pastelClient.FindTicketByID(ctx, task.Request.ArtistPastelID)
	if err != nil {
		return errors.Errorf("find register ticket of artist pastel id(%s):%w", task.Request.ArtistPastelID, err)
	}

	pqSignature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, pastel.SignAlgorithmLegRoast)
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
	if err := img.Encode(encSig); err != nil {
		return err
	}

	return nil
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
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

					if err := node.ConnectTo(ctx, primary.PastelID(), primary.SessID()); err != nil {
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

func (task *Task) isSuitableStorageFee(ctx context.Context) (bool, error) {
	fee, err := task.pastelClient.StorageNetworkFee(ctx)
	if err != nil {
		return false, err
	}
	return fee <= task.Request.MaximumFee, nil
}

func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Request.MaximumFee {
			continue
		}

		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = task.pastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// determine current block height & hash of it
func (task *Task) getBlock(ctx context.Context) error {
	// Get block num
	blockNum, err := task.pastelClient.GetBlockCount(ctx)
	task.creatorBlockHeight = int(blockNum)
	if err != nil {
		return errors.Errorf("get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return errors.Errorf("get block info blocknum=%d: %w", blockNum, err)
	}

	// Decode hash string to byte
	task.creatorBlockHash = blockInfo.Hash
	if err != nil {
		return errors.Errorf("convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	return nil
}

func (task *Task) genRQIdentifiersFiles(ctx context.Context) error {
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
		RqFilesDir: task.Service.config.RqFilesDir,
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
			task.UpdateStatus(StatusErrorGenRaptorQSymbolsFailed)
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

func (task *Task) generateRQIDs(ctx context.Context, rawFile rqnode.RawSymbolIDFile) error {
	var content []byte
	appendStr := func(b []byte, s string) []byte {
		b = append(b, []byte(s)...)
		b = append(b, '\n')
		return b
	}

	content = appendStr(content, rawFile.BlockHash)
	content = appendStr(content, rawFile.PastelID)
	for _, id := range rawFile.SymbolIdentifiers {
		content = appendStr(content, id)
	}

	signature, err := task.pastelClient.Sign(ctx, content, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign identifiers file: %w", err)
	}

	encfile := utils.B64Encode(content)

	var buffer bytes.Buffer
	buffer.Write(encfile)
	buffer.WriteString(".")
	buffer.Write(signature)
	rqIDFile := buffer.Bytes()

	task.rqIDsIc = rand.Uint32()
	task.rqids, err = pastel.GetIDFiles(rqIDFile, task.rqIDsIc, task.config.RQIDsMax)
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

func (task *Task) createArtTicket(_ context.Context) error {
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

func (task *Task) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeNFTTicket(task.ticket)
	if err != nil {
		return errors.Errorf("encode ticket %w", err)
	}

	task.creatorSignature, err = task.pastelClient.Sign(ctx, data, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket %w", err)
	}
	return nil
}

func (task *Task) registerActTicket(ctx context.Context) (string, error) {
	return task.pastelClient.RegisterActTicket(ctx, task.regNFTTxid, task.creatorBlockHeight, task.registrationFee, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase)
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Request) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Request: Ticket,
	}
}

func (task *Task) connectToTopRankNodes(ctx context.Context) error {
	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return errors.Errorf("call masternode top: %w", err)
	}

	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	// For testing
	// // TODO: Remove this when releaset because in localnet there is no need to start 10 gonode
	// // and chase the log
	// pinned := []string{"127.0.0.1:4444", "127.0.0.1:4445", "127.0.0.1:4446"}
	// pinnedMn := make([]*node.Node, 0)
	// for i := range topNodes {
	// 	for j := range pinned {
	// 		if topNodes[i].Address() == pinned[j] {
	// 			pinnedMn = append(pinnedMn, topNodes[i])
	// 		}
	// 	}
	// }
	// log.WithContext(ctx).Debugf("%v", pinnedMn)

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
	task.UpdateStatus(StatusConnected)
	task.nodes = nodes

	return nil
}

func (task *Task) probeImage(ctx context.Context) error {
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debug("probe image")

	// Send image to supernodes for probing.
	if err := task.nodes.ProbeImage(ctx, task.Request.Image); err != nil {
		return errors.Errorf("send image: %w", err)
	}
	signatures := [][]byte{}
	// Match signatures received from supernodes.
	for i := 0; i < len(task.nodes); i++ {
		ddData, err := json.Marshal(task.nodes[i].FingerprintAndScores)
		if err != nil {
			return errors.Errorf("probeImage: marshal fingerprintAndScores %w", err)
		}

		verified, err := task.pastelClient.Verify(ctx, ddData, string(task.nodes[i].Signature), task.nodes[i].PastelID(), "ed448")
		if err != nil {
			return errors.Errorf("probeImage: pastelClient.Verify %w", err)
		}

		if !verified {
			task.UpdateStatus(StatusErrorSignaturesNotMatch)
			return errors.Errorf("node[%s] signature doesn't match", task.nodes[i].PastelID())
		}

		signatures = append(signatures, task.nodes[i].Signature)
	}
	task.signatures = signatures

	// Match fingerprints received from supernodes.
	if err := task.nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}
	task.fingerprintAndScores = task.nodes.FingerAndScores()
	task.UpdateStatus(StatusImageProbed)

	// As we are going to store the the compressed figerprint to kamedila
	// so we calculated the hash based on the compressed fingerprints as well
	fingerprintsHash, err := utils.Sha3256hash(task.fingerprintAndScores.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("hash zstd commpressed fingerprints: %w", err)
	}
	task.fingerprintsHash = []byte(base58.Encode(fingerprintsHash))

	// Decompress the fingerprint as bytes for later use
	task.fingerprint, err = zstd.Decompress(nil, task.fingerprintAndScores.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("decompress finger failed")
	}

	return nil
}

func (task *Task) uploadImage(ctx context.Context) error {
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
		task.UpdateStatus(StatusErrorThumbnailHashsesNotMatch)
		return errors.Errorf("thumbnail hash returns by supenodes not mached: %w", err)
	}
	task.UpdateStatus(StatusImageAndThumbnailUploaded)

	task.previewHash = task.nodes.PreviewHash()
	task.mediumThumbnailHash = task.nodes.MediumThumbnailHash()
	task.smallThumbnailHash = task.nodes.SmallThumbnailHash()

	return nil
}

func (task *Task) sendSignedTicket(ctx context.Context) error {
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

func (task *Task) preburntRegistrationFee(ctx context.Context) error {
	if task.registrationFee <= 0 {
		return errors.Errorf("invalid registration fee")
	}

	burnedAmount := float64(task.registrationFee) / 10
	burnTxid, err := task.pastelClient.SendFromAddress(ctx, task.Request.SpendableAddress, task.config.BurnAddress, burnedAmount)
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

func (task *Task) removeArtifacts() {
	removeFn := func(file *artwork.File) {
		if file != nil {
			log.Debugf("remove file: %s", file.Name())
			if err := file.Remove(); err != nil {
				log.Debugf("remove file failed: %s", err.Error())
			}
		}
	}
	if task.Request != nil {
		removeFn(task.Request.Image)
		removeFn(task.imageEncodedWithFingerprints)
	}
}
