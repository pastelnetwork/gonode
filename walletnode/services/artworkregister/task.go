package artworkregister

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/node"
	"golang.org/x/crypto/sha3"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Request *Request

	// information of nodes
	nodes node.List

	// task data to create RegArt ticket
	artistBlockHeight            int
	artistBlockHash              string
	fingerprintAndScores         *pastel.FingerAndScores
	fingerprintsHash             []byte
	fingerprint                  []byte
	imageEncodedWithFingerprints *artwork.File
	previewHash                  []byte
	mediumThumbnailHash          []byte
	smallThumbnailHash           []byte
	datahash                     []byte
	rqids                        []string

	// TODO: call cNodeAPI to get the reall signature instead of the fake one
	fingerprintSignature []byte

	// TODO: need to update rqservice code to return the following info
	rqSymbolIDFiles map[string][]byte
	rqEncodeParams  rqnode.EncoderParameters

	// TODO: call cNodeAPI to get the following info
	blockTxID       string
	burnTxid        string
	regArtTxid      string
	registrationFee int64

	// ticket
	artistSignature []byte
	ticket          *pastel.ArtTicket
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debugf("States updated")
	})

	log.WithContext(ctx).Debugf("Start task")
	defer log.WithContext(ctx).Debugf("End task")

	if err := task.run(ctx); err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Warnf("Task is rejected")
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
		return errors.Errorf("failed to connect to top rank nodes: %w", err)
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
		return errors.Errorf("failed to get current block heigth %w", err)
	}

	// probe image for rareness, nsfw and seen score
	if err := task.probeImage(ctx); err != nil {
		return errors.Errorf("failed to probe image: %w", err)
	}

	// upload image and thumbnail coordinates to supernode(s)
	if err := task.uploadImage(ctx); err != nil {
		return errors.Errorf("failed to upload image: %w", err)
	}

	// connect to rq serivce to get rq symbols identifier
	if err := task.genRQIdentifiersFiles(ctx); err != nil {
		return errors.Errorf("gen RaptorQ symbols' identifiers failed %w", err)
	}

	if err := task.createArtTicket(ctx); err != nil {
		return errors.Errorf("failed to create ticket %w", err)
	}

	// sign ticket with artist signature
	if err := task.signTicket(ctx); err != nil {
		return errors.Errorf("failed to sign art ticket %w", err)
	}

	// send signed ticket to supernodes to calculate registration fee
	if err := task.sendSignedTicket(ctx); err != nil {
		return errors.Errorf("failed to send signed art ticket: %w", err)
	}

	// send preburn-txid to master node(s)
	// master node will create reg-art ticket and returns transaction id
	if err := task.preburntRegistrationFee(ctx); err != nil {
		return errors.Errorf("faild to pre-burnt ten percent of registration fee: %w", err)
	}

	log.WithContext(ctx).Debugf("close connection")
	close(nodesDone)
	for i := range task.nodes {
		if err := task.nodes[i].Connection.Close(); err != nil {
			log.WithContext(ctx).WithError(err).Debugf("failed to close connection to node %s", task.nodes[i].PastelID())
		}
	}

	// new context because the old context already cancelled
	newCtx := context.Background()
	if err := task.waitTxidValid(newCtx, task.regArtTxid, int64(task.config.RegArtTxMinConfirmations), task.config.RegArtTxTimeout, 15*time.Second); err != nil {
		return errors.Errorf("failed to wait reg-art ticket valid %w", err)
	}

	// activate reg-art ticket at previous step
	actTxid, err := task.registerActTicket(newCtx)
	if err != nil {
		return errors.Errorf("failed to register act ticket %w", err)
	}
	log.Debugf("reg-act-txid: %s", actTxid)

	// Wait until actTxid is valid
	err = task.waitTxidValid(newCtx, actTxid, int64(task.config.RegActTxMinConfirmations), task.config.RegActTxTimeout, 15*time.Second)
	if err != nil {
		return errors.Errorf("failed to reg-act ticket valid %w", err)
	}
	log.Debugf("reg-act-tixd is confirmed")

	return nil
}

func (task *Task) waitTxidValid(ctx context.Context, txID string, expectedConfirms int64, timeout time.Duration, interval time.Duration) error {
	log.WithContext(ctx).Debugf("Need %d confirmation for txid %s, timeout %v", expectedConfirms, txID, timeout)

	checkTimeout := func(checked chan<- struct{}) {
		time.Sleep(timeout)
		close(checked)
	}

	timeoutCh := make(chan struct{})
	go checkTimeout(timeoutCh)

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("context done %w", ctx.Err())
		case <-time.After(interval):
			checkConfirms := func() error {
				subCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				defer cancel()

				result, err := task.pastelClient.GetRawTransactionVerbose1(subCtx, txID)
				log.WithContext(ctx).Debugf("getrawtransaction result: %s %v", txID, result)
				if err != nil {
					return errors.Errorf("failed to get transaction %w", err)
				}

				if result.Confirmations >= expectedConfirms {
					return nil
				}

				return errors.Errorf("confirmations not meet: expected %d, got %d", expectedConfirms, result.Confirmations)
			}

			err := checkConfirms()
			if err != nil {
				log.WithContext(ctx).Errorf("check confirmations failed : %v", err)
			} else {
				return nil
			}
		case <-timeoutCh:
			return errors.Errorf("timeout")
		}
	}
}

func (task *Task) encodeFingerprint(ctx context.Context, fingerprint []byte, img *artwork.File) error {
	// Sign fingerprint
	ed448PubKey := []byte(task.Request.ArtistPastelID)
	ed448Signature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, "ed448")
	if err != nil {
		return errors.Errorf("sign fingerprint with pastelId and pastelPassphrase failed %w", err)
	}

	// TODO: Should be replaced with real data from the Pastel API.
	ticket, err := task.pastelClient.FindTicketByID(ctx, task.Request.ArtistPastelID)
	if err != nil {
		return errors.Errorf("failed to find register ticket of artist pastel id:%s :%w", task.Request.ArtistPastelID, err)
	}
	pqSignature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, "legroast")
	if err != nil {
		return errors.Errorf("failed to sign fingerprint with legroats: %w", err)
	}

	// Encode data to the image.
	encSig := qrsignature.New(
		qrsignature.Fingerprint(fingerprint),
		qrsignature.PostQuantumSignature(pqSignature),
		qrsignature.PostQuantumPubKey([]byte(ticket.Signature)),
		qrsignature.Ed448Signature(ed448Signature),
		qrsignature.Ed448PubKey(ed448PubKey),
	)
	if err := img.Encode(encSig); err != nil {
		return err
	}
	task.fingerprintSignature = ed448Signature

	// TODO: check with the change in legroast
	// Decode data from the image, to make sure their integrity.
	// decSig := qrsignature.New()
	// copyImage, _ := img.Copy()
	// if err := copyImage.Decode(decSig); err != nil {
	// 	return err
	// }

	// if !bytes.Equal(fingerprint, decSig.Fingerprint()) {
	// 	return errors.Errorf("fingerprints do not match, original len:%d, decoded len:%d\n", len(fingerprint), len(decSig.Fingerprint()))
	// }
	// if !bytes.Equal(pqSignature, decSig.PostQuantumSignature()) {
	// 	return errors.Errorf("post quantum signatures do not match, original len:%d, decoded len:%d\n", len(pqSignature), len(decSig.PostQuantumSignature()))
	// }
	// if !bytes.Equal(pqPubKey, decSig.PostQuantumPubKey()) {
	// 	return errors.Errorf("post quantum public keys do not match, original len:%d, decoded len:%d\n", len(pqPubKey), len(decSig.PostQuantumPubKey()))
	// }
	// if !bytes.Equal(ed448Signature, decSig.Ed448Signature()) {
	// 	return errors.Errorf("ed448 signatures do not match, original len:%d, decoded len:%d\n", len(ed448Signature), len(decSig.Ed448Signature()))
	// }
	// if !bytes.Equal(ed448PubKey, decSig.Ed448PubKey()) {
	// 	return errors.Errorf("ed448 public keys do not match, original len:%d, decoded len:%d\n", len(ed448PubKey), len(decSig.Ed448PubKey()))
	// }
	return nil
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
	var meshNodes node.List

	primary := nodes[primaryIndex]
	if err := primary.Connect(ctx, task.config.ConnectTimeout); err != nil {
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

					if err := node.Connect(ctx, task.config.ConnectTimeout); err != nil {
						return
					}
					if err := node.Session(ctx, false); err != nil {
						return
					}
					go func() {
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
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// determine current block height & hash of it
func (task *Task) getBlock(ctx context.Context) error {
	// Get block num
	blockNum, err := task.pastelClient.GetBlockCount(ctx)
	task.artistBlockHeight = int(blockNum)
	if err != nil {
		return errors.Errorf("failed to get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return errors.Errorf("failed to get block info with given block num %d: %w", blockNum, err)
	}

	// Decode hash string to byte
	task.artistBlockHash = blockInfo.Hash
	if err != nil {
		return errors.Errorf("failed to convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	return nil
}

func (task *Task) genRQIdentifiersFiles(ctx context.Context) error {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("failed to connect to raptorQ service %w", err)
	}
	defer conn.Close()

	content, err := task.imageEncodedWithFingerprints.Bytes()
	if err != nil {
		return errors.Errorf("failed to read image contents: %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	})

	log.WithContext(ctx).Debugf("Image hash %x", sha3.Sum256(content))
	// FIXME :
	// - check format of artis block hash should be base58 or not
	encodeInfo, err := rqService.EncodeInfo(ctx, content, task.config.NumberRQIDSFiles, task.artistBlockHash, task.Request.ArtistPastelID)
	if err != nil {
		return errors.Errorf("failed to generate RaptorQ symbols' identifiers %w", err)
	}

	files := make(map[string][]byte)
	for _, rawSymbolIDFile := range encodeInfo.SymbolIDFiles {
		name, content, err := task.convertToSymbolIDFile(ctx, rawSymbolIDFile)
		if err != nil {
			task.UpdateStatus(StatusErrorGenRaptorQSymbolsFailed)
			return errors.Errorf("failed to create rqids file %w", err)
		}
		files[name] = content
	}

	if len(files) < 1 {
		return errors.Errorf("nuber of raptorq symbol identifiers files must be greater than 1")
	}
	for k := range files {
		log.WithContext(ctx).Debugf("symbol file identifier %s", k)
		task.rqids = append(task.rqids, k)
	}

	task.rqSymbolIDFiles = files
	task.rqEncodeParams = encodeInfo.EncoderParam
	return nil
}

func (task *Task) convertToSymbolIDFile(ctx context.Context, rawFile rqnode.RawSymbolIDFile) (string, []byte, error) {
	var content []byte
	appendStr := func(b []byte, s string) []byte {
		b = append(b, []byte(s)...)
		b = append(b, '\n')
		return b
	}

	content = appendStr(content, rawFile.ID)
	content = appendStr(content, rawFile.BlockHash)
	content = appendStr(content, rawFile.PastelID)
	for _, id := range rawFile.SymbolIdentifiers {
		content = appendStr(content, id)
	}

	// log.WithContext(ctx).Debugf("%s", content)
	signature, err := task.pastelClient.Sign(ctx, content, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, "ed448")
	if err != nil {
		return "", nil, errors.Errorf("failed to sign identifiers file: %w", err)
	}
	content = append(content, signature...)

	// read all content of the data b
	// compress the contente of the file
	compressedData, err := zstd.CompressLevel(nil, content, 22)
	if err != nil {
		return "", nil, errors.Errorf("failed to compress: ")
	}
	hasher := sha3.New256()
	src := bytes.NewReader(compressedData)
	if _, err := io.Copy(hasher, src); err != nil {
		return "", nil, errors.Errorf("failed to hash identifiers file %w", err)
	}
	return base58.Encode(hasher.Sum(nil)), compressedData, nil
}

func (task *Task) createArtTicket(_ context.Context) error {
	if task.fingerprint == nil {
		return errEmptyFingerprints
	}
	if task.fingerprintsHash == nil {
		return errEmptyFingerprintsHash
	}
	if task.fingerprintSignature == nil {
		return errEmptyFingerprintSignature
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

	pastelID := base58.Decode(task.Request.ArtistPastelID)
	if pastelID == nil {
		return errDecodePastelID
	}

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.ArtTicket{
		Version:   1,
		Author:    task.Request.ArtistPastelID,
		BlockNum:  task.artistBlockHeight,
		BlockHash: task.artistBlockHash,
		Copies:    task.Request.IssuedCopies,
		Royalty:   0,  // Not supported yet by cNode
		Green:     "", // Not supported yet by cNode
		AppTicketData: pastel.AppTicket{
			AuthorPastelID:                 task.Request.ArtistPastelID,
			BlockTxID:                      task.blockTxID,
			BlockNum:                       task.artistBlockHeight,
			ArtistName:                     task.Request.ArtistName,
			ArtistWebsite:                  safeString(task.Request.ArtistWebsiteURL),
			ArtistWrittenStatement:         safeString(task.Request.Description),
			ArtworkTitle:                   safeString(&task.Request.Name),
			ArtworkSeriesName:              safeString(task.Request.SeriesName),
			ArtworkCreationVideoYoutubeURL: safeString(task.Request.YoutubeURL),
			ArtworkKeywordSet:              safeString(task.Request.Keywords),
			TotalCopies:                    task.Request.IssuedCopies,
			PreviewHash:                    task.previewHash,
			Thumbnail1Hash:                 task.mediumThumbnailHash,
			Thumbnail2Hash:                 task.smallThumbnailHash,
			DataHash:                       task.datahash,
			FingerprintsHash:               task.fingerprintsHash,
			FingerprintsSignature:          task.fingerprintSignature,
			DupeDetectionSystemVer:         task.fingerprintAndScores.DupeDectectionSystemVersion,
			MatchesFoundOnFirstPage:        task.fingerprintAndScores.MatchesFoundOnFirstPage,
			NumberOfResultPages:            task.fingerprintAndScores.NumberOfPagesOfResults,
			FirstMatchURL:                  task.fingerprintAndScores.URLOfFirstMatchInPage,
			PastelRarenessScore:            task.fingerprintAndScores.OverallAverageRarenessScore,
			InternetRarenessScore:          float64(task.fingerprintAndScores.IsRareOnInternet),
			OpenNSFWScore:                  task.fingerprintAndScores.OpenNSFWScore,
			AlternateNSFWScores: pastel.AlternateNSFWScores{
				Drawing: task.fingerprintAndScores.AlternativeNSFWScore.Drawing,
				Hentai:  task.fingerprintAndScores.AlternativeNSFWScore.Hentai,
				Neutral: task.fingerprintAndScores.AlternativeNSFWScore.Neutral,
				Porn:    task.fingerprintAndScores.AlternativeNSFWScore.Porn,
				Sexy:    task.fingerprintAndScores.AlternativeNSFWScore.Sexy,
			},
			ImageHashes: pastel.ImageHashes{
				PerceptualHash: "",
				AverageHash:    "",
				DifferenceHash: "",
			},
			RQIDs: task.rqids,
			RQOti: task.rqEncodeParams.Oti,
		},
	}

	task.ticket = ticket
	return nil
}

func (task *Task) signTicket(ctx context.Context) error {
	data, err := pastel.EncodeArtTicket(task.ticket)
	if err != nil {
		return errors.Errorf("failed to encode ticket %w", err)
	}

	task.artistSignature, err = task.pastelClient.Sign(ctx, data, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase, "ed448")
	if err != nil {
		return errors.Errorf("failed to sign ticket %w", err)
	}
	return nil
}

func (task *Task) registerActTicket(ctx context.Context) (string, error) {
	return task.pastelClient.RegisterActTicket(ctx, task.regArtTxid, task.artistBlockHeight, task.registrationFee, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase)
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
		return errors.Errorf("failed to call masternode top: %w", err)
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
			if errors.IsContextCanceled(err) {
				return err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warnf("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(nodes) < task.config.NumberSuperNodes {
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
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debugf("Copy image")
	thumbnail, err := task.Request.Image.Copy()
	if err != nil {
		return errors.Errorf("failed to copy image: %w", err)
	}
	defer thumbnail.Remove()

	log.WithContext(ctx).WithField("filename", thumbnail.Name()).Debugf("Resize image to %dx%d pixeles", task.config.thumbnailSize, task.config.thumbnailSize)
	if err := thumbnail.ResizeImage(thumbnailSize, thumbnailSize); err != nil {
		return errors.Errorf("failed to resize image: %w", err)
	}

	// Send thumbnail to supernodes for probing.
	if err := task.nodes.ProbeImage(ctx, thumbnail); err != nil {
		return errors.Errorf("failed to probe image: %w", err)
	}

	// Match fingerprints received from supernodes.
	if err := task.nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(StatusErrorFingerprintsNotMatch)
		return errors.Errorf("fingerprints aren't matched :%w", err)
	}
	task.fingerprintAndScores = task.nodes.FingerAndScores()
	task.UpdateStatus(StatusImageProbed)

	// As we are going to store the the compressed figerprint to kamedila
	// so we calculated the hash base on the compressed fingerprint print also
	fingerprintsHash, err := sha3256hash(task.fingerprintAndScores.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("failed to hash zstd commpressed fingerprints %w", err)
	}
	task.fingerprintsHash = []byte(base58.Encode(fingerprintsHash))

	// Decompress the fingerprint as bytes for later use
	task.fingerprint, err = zstd.Decompress(nil, task.fingerprintAndScores.ZstdCompressedFingerprint)
	if err != nil {
		return errors.Errorf("failed to decompress finger")
	}
	return nil
}

func (task *Task) uploadImage(ctx context.Context) error {
	img1, err := task.Request.Image.Copy()
	if err != nil {
		return errors.Errorf("copy image to encode failed %w", err)
	}

	log.WithContext(ctx).WithField("FileName", img1.Name()).Debugf("final image")
	task.Request.Image = img1

	if err := task.encodeFingerprint(ctx, task.fingerprint, img1); err != nil {
		return errors.Errorf("encode image with fingerprint %w", err)
	}
	task.imageEncodedWithFingerprints = img1

	imgBytes, err := img1.Bytes()
	if err != nil {
		return errors.Errorf("failed to convert image to byte stream %w", err)
	}

	if task.datahash, err = sha3256hash(imgBytes); err != nil {
		return errors.Errorf("failed to hash encoded image: %w", err)
	}

	// Upload image with pqgsinganature and its thumb to supernodes
	if err := task.nodes.UploadImageWithThumbnail(ctx, img1, task.Request.Thumbnail); err != nil {
		return errors.Errorf("upload encoded image and thumbnail coordinate failed %w", err)
	}
	// Match thumbnail hashes receiveed from supernodes
	if err := task.nodes.MatchThumbnailHashes(); err != nil {
		task.UpdateStatus(StatusErrorThumbnailHashsesNotMatch)
		return errors.Errorf("thumbnail hash returns by supenodes not mached %w", err)
	}
	task.UpdateStatus(StatusImageAndThumbnailUploaded)

	task.previewHash = task.nodes.PreviewHash()
	task.mediumThumbnailHash = task.nodes.MediumThumbnailHash()
	task.smallThumbnailHash = task.nodes.SmallThumbnailHash()

	return nil
}

func (task *Task) sendSignedTicket(ctx context.Context) error {
	buf, err := pastel.EncodeArtTicket(task.ticket)
	if err != nil {
		return errors.Errorf("failed to marshal ticket %w", err)
	}
	log.Debug(string(buf))

	if err := task.nodes.UploadSignedTicket(ctx, buf, task.artistSignature, task.rqSymbolIDFiles, task.rqEncodeParams); err != nil {
		return errors.Errorf("failed to upload signed ticket %w", err)
	}

	if err := task.nodes.MatchRegistrationFee(); err != nil {
		return errors.Errorf("registration fees don't matched %w", err)
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
		return errors.Errorf("failed to burn 10 percent of transaction fee %w", err)
	}
	task.burnTxid = burnTxid
	log.WithContext(ctx).Debugf("preburn txid: %s", task.burnTxid)

	if err := task.nodes.SendPreBurntFeeTxid(ctx, task.burnTxid); err != nil {
		return errors.Errorf("failed to send pre-burn-txid: %s to supernode(s): %w", task.burnTxid, err)
	}
	task.regArtTxid = task.nodes.RegArtTicketID()
	if task.regArtTxid == "" {
		return errors.Errorf("regArtTxid is empty")
	}

	return nil
}
