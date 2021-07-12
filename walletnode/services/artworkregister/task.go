package artworkregister

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rq "github.com/pastelnetwork/gonode/raptorq"
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
	fingerprints                 []byte
	fingerprintsHash             []byte
	imageEncodedWithFingerprints *artwork.File
	previewHash                  []byte
	mediumThumbnailHash          []byte
	smallThumbnailHash           []byte
	datahash                     []byte
	rqids                        []string

	// TODO: call cNodeAPI to get the reall signature instead of the fake one
	fingerprintSignature []byte

	// TODO: need to update rqservice code to return the following info
	rqoti []byte

	// TODO: call cNodeAPI to get the following info
	blockTxID string
	blockNum  int
	burnTxId  pastel.TxIDType

	rarenessScore int
	nSFWScore     int
	seenScore     int
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
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return err
	}
	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

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
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return nodes.WaitConnClose(ctx)
	})
	task.UpdateStatus(StatusConnected)
	task.nodes = nodes

	// Create a thumbnail copy of the image.
	log.WithContext(ctx).WithField("filename", task.Request.Image.Name()).Debugf("Copy image")
	thumbnail, err := task.Request.Image.Copy()
	if err != nil {
		return err
	}
	defer thumbnail.Remove()

	log.WithContext(ctx).WithField("filename", thumbnail.Name()).Debugf("Resize image to %dx%d pixeles", task.config.thumbnailSize, task.config.thumbnailSize)
	if err := thumbnail.ResizeImage(thumbnailSize, thumbnailSize); err != nil {
		return err
	}

	// Send thumbnail to supernodes for probing.
	if err := nodes.ProbeImage(ctx, thumbnail); err != nil {
		return err
	}

	// Match fingerprints received from supernodes.
	if err := nodes.MatchFingerprintAndScores(); err != nil {
		task.UpdateStatus(StatusErrorFingerprintsNotMatch)
		return err
	}
	task.UpdateStatus(StatusImageProbed)
	task.fingerprints = nodes.Fingerprint()
	if task.fingerprintsHash, err = sha3256hash(task.fingerprints); err != nil {
		return errors.Errorf("failed to hash fingerprints %w", err)
	}

	task.rarenessScore = nodes.RarenessScore()
	task.nSFWScore = nodes.NSFWScore()
	task.seenScore = nodes.SeenScore()

	finalImage, err := task.Request.Image.Copy()
	if err != nil {
		return errors.Errorf("copy image to encode failed %w", err)
	}
	log.WithContext(ctx).WithField("FileName", finalImage.Name()).Debugf("final image")
	task.Request.Image = finalImage

	imgBytes, err := finalImage.Bytes()
	if err != nil {
		return errors.Errorf("failed to convert image to byte stream %w", err)
	}
	if task.datahash, err = sha3256hash(imgBytes); err != nil {
		return errors.Errorf("failed to hash encoded image %w", err)
	}
	// defer finalImage.Remove()

	if err := task.encodeFingerprint(ctx, task.fingerprints, finalImage); err != nil {
		return errors.Errorf("encode image with fingerprint %w", err)
	}
	task.imageEncodedWithFingerprints = finalImage

	// Upload image with pqgsinganature and its thumb to supernodes
	if err := nodes.UploadImageWithThumbnail(ctx, finalImage, task.Request.Thumbnail); err != nil {
		return errors.Errorf("upload encoded image and thumbnail coordinate failed %w", err)
	}
	// Match thumbnail hashes receiveed from supernodes
	if err := nodes.MatchThumbnailHashes(); err != nil {
		task.UpdateStatus(StatusErrorThumbnailHashsesNotMatch)
		return errors.Errorf("thumbnail hash returns by supenodes not mached %w", err)
	}
	task.UpdateStatus(StatusImageAndThumbnailUploaded)

	task.previewHash = nodes.PreviewHash()
	task.mediumThumbnailHash = nodes.MediumThumbnailHash()
	task.smallThumbnailHash = nodes.SmallThumbnailHash()

	// Connect to rq serivce to get rq symbols identifier
	rqSymbolIdFiles, encoderParams, err := task.genRQIdentifiersFiles(ctx)
	if err != nil {
		task.UpdateStatus(StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers failed %w", err)
	}
	if len(rqSymbolIdFiles) < 1 {
		return errors.Errorf("nuber of raptorq symbol identifiers files must be greter than 1")
	}
	if task.rqids, err = rqSymbolIdFiles.FileIdentifers(); err != nil {
		return errors.Errorf("failed to get rq symbols' identifier file's identifier %w", err)
	}

	ticket, err := task.createTicket(ctx)
	if err != nil {
		return errors.Errorf("failed to create ticket %w", err)
	}

	buf, err := pastel.EncodeArtTicket(ticket)
	if err != nil {
		return errors.Errorf("failed to marshal ticket %w", err)
	} else {
		log.Debug(string(buf))
	}

	// sign ticket with artist signature
	artistSignature, err := task.signTicket(ctx, ticket)
	if err != nil {
		return errors.Errorf("failed to sign ticket %w", err)
	}

	// send signed ticket to supernodes to calculate registration fee
	symbolsIdFilesMap, err := rqSymbolIdFiles.ToMap()
	if err != nil {
		return errors.Errorf("failed to create rq symbol identifiers files map %w", err)
	}
	if err := nodes.UploadSignedTicket(ctx, buf, artistSignature, symbolsIdFilesMap, *encoderParams); err != nil {
		return errors.Errorf("failed to upload signed ticket %w", err)
	}

	if err := nodes.MatchRegistrationFee(); err != nil {
		return errors.Errorf("registration fees don't matched %w", err)
	}

	// burn 10 % of registration fee by sending to unspendable address with has the format of PtPasteLBurnAddressXXXXXXXXXTWPm3E
	// TODO: make this as configuration
	burnedAmount := nodes.RegistrationFee() / 10
	if task.burnTxId, err = task.pastelClient.SendToAddress(ctx, "PtPasteLBurnAddressXXXXXXXXXTWPm3E", burnedAmount); err != nil {
		return errors.Errorf("failed to burn 10% of transaction fee %w", err)
	}

	// send the txid of the preburn transaction to super nodes
	if err := nodes.SendPreBurntFeeTxId(ctx, task.burnTxId); err != nil {
		return errors.Errorf("failed to send txId of preburnt fee transaction %w", err)
	}
	fmt.Println("OK")

	// Wait for all connections to disconnect.
	return groupConnClose.Wait()
}

func (task *Task) encodeFingerprint(ctx context.Context, fingerprint []byte, img *artwork.File) error {
	// Sign fingerprint
	ed448PubKey := []byte(task.Request.ArtistPastelID)
	ed448Signature, err := task.pastelClient.Sign(ctx, fingerprint, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase)
	if err != nil {
		return errors.Errorf("sign fingerprint with pastelId and pastelPassphrase failed %w", err)
	}

	// TODO: Should be replaced with real data from the Pastel API.
	pqPubKey := []byte("pqPubKey")
	pqSignature := []byte("pqSignature")

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
	task.fingerprintSignature = pqSignature

	// Decode data from the image, to make sure their integrity.
	decSig := qrsignature.New()
	copyImage, _ := img.Copy()
	if err := copyImage.Decode(decSig); err != nil {
		return err
	}

	if !bytes.Equal(fingerprint, decSig.Fingerprint()) {
		return errors.Errorf("fingerprints do not match, original len:%d, decoded len:%d\n", len(fingerprint), len(decSig.Fingerprint()))
	}
	if !bytes.Equal(pqSignature, decSig.PostQuantumSignature()) {
		return errors.Errorf("post quantum signatures do not match, original len:%d, decoded len:%d\n", len(pqSignature), len(decSig.PostQuantumSignature()))
	}
	if !bytes.Equal(pqPubKey, decSig.PostQuantumPubKey()) {
		return errors.Errorf("post quantum public keys do not match, original len:%d, decoded len:%d\n", len(pqPubKey), len(decSig.PostQuantumPubKey()))
	}
	if !bytes.Equal(ed448Signature, decSig.Ed448Signature()) {
		return errors.Errorf("ed448 signatures do not match, original len:%d, decoded len:%d\n", len(ed448Signature), len(decSig.Ed448Signature()))
	}
	if !bytes.Equal(ed448PubKey, decSig.Ed448PubKey()) {
		return errors.Errorf("ed448 public keys do not match, original len:%d, decoded len:%d\n", len(ed448PubKey), len(decSig.Ed448PubKey()))
	}
	return nil
}

// meshNodes establishes communication between supernodes.
func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
	var meshNodes node.List

	primary := nodes[primaryIndex]
	if err := primary.Connect(ctx, task.config.connectTimeout); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
		return nil, err
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

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

					if err := node.Connect(ctx, task.config.connectTimeout); err != nil {
						return
					}
					if err := node.Session(ctx, false); err != nil {
						return
					}
					secondaries.Add(node)

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

	meshNodes.Add(primary)
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

func (task *Task) genRQIdentifiersFiles(ctx context.Context) (rq.SymbolIdFiles, *rqnode.EncoderParameters, error) {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return nil, nil, errors.Errorf("failed to connect to raptorQ service %w", err)
	}
	defer conn.Close()

	content, err := task.Request.Image.Bytes()
	if err != nil {
		return nil, nil, errors.Errorf("failed to read image contents")
	}

	rqService := conn.RaptorQ()
	encodeInfo, err := rqService.EncodeInfo(ctx, content)
	if err != nil {
		return nil, nil, errors.Errorf("failed to generate RaptorQ symbols' identifiers %w", err)
	}

	files := rq.SymbolIdFiles{}
	for i := 0; i < task.config.NumberRQIDSFiles; i++ {
		f, err := task.createRQIDSFile(ctx, encodeInfo.SymbolIds, "BLOCK_HASH")
		if err != nil {
			return nil, nil, errors.Errorf("failed to create rqids file %w", err)
		}

		files = append(files, f)
	}

	return files, &encodeInfo.EncoderParam, nil
}

func (task *Task) createRQIDSFile(ctx context.Context, symbols []string, blockHash string) (*rq.SymbolIdFile, error) {
	symbolIdFile := rq.SymbolIdFile{
		Id:                uuid.New().String(),
		BlockHash:         blockHash,
		SymbolIdentifiers: symbols,
		Signature:         nil,
	}

	js, err := json.Marshal(&symbolIdFile)
	if err != nil {
		return nil, errors.Errorf("failed to marshal identifiers file %w", err)
	}

	symbolIdFile.Signature, err = task.pastelClient.Sign(ctx, js, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase)
	if err != nil {
		return nil, errors.Errorf("failed to sign identifier file %w", err)
	}

	js, err = json.Marshal(&symbolIdFile)
	if err != nil {
		return nil, errors.Errorf("failed to marshal identifiers file with signature %w", err)
	}

	hasher := sha3.New256()
	src := bytes.NewReader(js)
	if _, err := io.Copy(hasher, src); err != nil {
		return nil, errors.Errorf("failed to hash identifiers file %w", err)
	}
	symbolIdFile.FileIdentifer = string(hasher.Sum(nil))

	return &symbolIdFile, nil
}

func (task *Task) createTicket(ctx context.Context) (*pastel.ArtTicket, error) {
	if task.fingerprints == nil {
		return nil, errors.Errorf("empty fingerprints")
	}
	if task.fingerprintsHash == nil {
		return nil, errors.Errorf("empty fingerprints hash")
	}
	if task.fingerprintSignature == nil {
		return nil, errors.Errorf("empty fingerprint signature")
	}
	if task.datahash == nil {
		return nil, errors.Errorf("empty data hash")
	}
	if task.previewHash == nil {
		return nil, errors.Errorf("empty preview hash")
	}
	if task.mediumThumbnailHash == nil {
		return nil, errors.Errorf("empty medium thumbnail hash")
	}
	if task.smallThumbnailHash == nil {
		return nil, errors.Errorf("empty small thumbnail hash")
	}
	if task.rqids == nil {
		return nil, errors.Errorf("empty RaptorQ symbols identifiers")
	}

	pastelID := base58.Decode(task.Request.ArtistPastelID)
	if pastelID == nil {
		return nil, errors.Errorf("base58 decode artist PastelID failed")
	}

	// Get block num
	blockNum, err := task.pastelClient.GetBlockCount(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get block num: %w", err)
	}

	// Get block hash string
	blockInfo, err := task.pastelClient.GetBlockVerbose1(ctx, blockNum)
	if err != nil {
		return nil, errors.Errorf("failed to get block info with given block num %d: %w", blockNum, err)
	}

	// Decode hash string to byte
	blockHash, err := base64.StdEncoding.DecodeString(blockInfo.Hash)
	if err != nil {
		return nil, errors.Errorf("failed to convert hash string %s to bytes: %w", blockInfo.Hash, err)
	}

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.ArtTicket{
		Version:   1,
		Author:    pastelID,
		BlockNum:  int(blockNum),
		BlockHash: blockHash,
		Copies:    task.Request.IssuedCopies,
		Royalty:   0,  // Not supported yet by cNode
		Green:     "", // Not supported yet by cNode
		AppTicketData: pastel.AppTicket{
			AuthorPastelID:                 task.Request.ArtistPastelID,
			BlockTxID:                      "TBD",
			BlockNum:                       0,
			ArtistName:                     task.Request.ArtistName,
			ArtistWebsite:                  safeString(task.Request.ArtistWebsiteURL),
			ArtistWrittenStatement:         safeString(task.Request.Description),
			ArtworkCreationVideoYoutubeURL: safeString(task.Request.YoutubeURL),
			ArtworkKeywordSet:              safeString(task.Request.Keywords),
			TotalCopies:                    task.Request.IssuedCopies,
			PreviewHash:                    task.previewHash,
			Thumbnail1Hash:                 task.mediumThumbnailHash,
			Thumbnail2Hash:                 task.smallThumbnailHash,
			DataHash:                       task.datahash,
			Fingerprints:                   task.fingerprints,
			FingerprintsHash:               task.fingerprintsHash,
			FingerprintsSignature:          task.fingerprintSignature,
			RarenessScore:                  task.rarenessScore,
			NSFWScore:                      task.nSFWScore,
			SeenScore:                      task.seenScore,
			RQIDs:                          task.rqids,
			RQOti:                          task.rqoti,
		},
	}

	return ticket, nil
}

func (task *Task) signTicket(ctx context.Context, ticket *pastel.ArtTicket) ([]byte, error) {
	js, err := json.Marshal(ticket)
	if err != nil {
		return nil, errors.Errorf("failed to encode ticket %w", err)
	}

	signature, err := task.pastelClient.Sign(ctx, js, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase)
	if err != nil {
		return nil, errors.Errorf("failed to sign ticket %w", err)
	}

	return signature, nil
}

//
// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Request) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Request: Ticket,
	}
}
