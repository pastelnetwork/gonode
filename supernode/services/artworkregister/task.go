package artworkregister

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"sort"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Ticket           *pastel.ArtTicket
	ResampledArtwork *artwork.File
	Artwork          *artwork.File
	imageSizeBytes   int

	fingerprintsData []byte
	fingerAndScores  *pastel.FingerAndScores
	fingerprints     pastel.Fingerprint

	PreviewThumbnail *artwork.File
	MediumThumbnail  *artwork.File
	SmallThumbnail   *artwork.File

	acceptedMu sync.Mutex
	accepted   Nodes

	RQIDS map[string][]byte
	Oti   []byte

	// signature of ticket data signed by node's pateslID
	ownSignature []byte

	artistSignature []byte
	key1            string
	key2            string
	registrationFee int64

	// valid only for a task run as primary
	peersArtTicketSignatureMtx *sync.Mutex
	peersArtTicketSignature    map[string][]byte
	allSignaturesReceived      chan struct{}

	// valid only for secondary node
	connectedTo *Node
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = task.context(ctx)
	defer log.WithContext(ctx).Debug("Task canceled")
	defer task.Cancel()

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debugf("States updated")
	})

	return task.RunAction(ctx)
}

// Session is handshake wallet to supernode
func (task *Task) Session(_ context.Context, isPrimary bool) error {
	if err := task.RequiredStatus(StatusTaskStarted); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if isPrimary {
			log.WithContext(ctx).Debugf("Acts as primary node")
			task.UpdateStatus(StatusPrimaryMode)
			return nil
		}

		log.WithContext(ctx).Debugf("Acts as secondary node")
		task.UpdateStatus(StatusSecondaryMode)

		return nil
	})
	return nil
}

// AcceptedNodes waits for connection supernodes, as soon as there is the required amount returns them.
func (task *Task) AcceptedNodes(serverCtx context.Context) (Nodes, error) {
	if err := task.RequiredStatus(StatusPrimaryMode); err != nil {
		return nil, err
	}

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("Waiting for supernodes to connect")

		sub := task.SubscribeStatus()
		for {
			select {
			case <-serverCtx.Done():
				return nil
			case <-ctx.Done():
				return nil
			case status := <-sub():
				if status.Is(StatusConnected) {
					return nil
				}
			}
		}
	})
	return task.accepted, nil
}

// SessionNode accepts secondary node
func (task *Task) SessionNode(_ context.Context, nodeID string) error {
	task.acceptedMu.Lock()
	defer task.acceptedMu.Unlock()

	if err := task.RequiredStatus(StatusPrimaryMode); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if node := task.accepted.ByID(nodeID); node != nil {
			return errors.Errorf("node %q is already registered", nodeID)
		}

		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			return errors.Errorf("failed to get node by extID %s %w", nodeID, err)
		}
		task.accepted.Add(node)

		log.WithContext(ctx).WithField("nodeID", nodeID).Debugf("Accept secondary node")

		if len(task.accepted) >= task.config.NumberConnectedNodes {
			task.UpdateStatus(StatusConnected)
		}
		return nil
	})
	return nil
}

// ConnectTo connects to primary node
func (task *Task) ConnectTo(_ context.Context, nodeID, sessID string) error {
	if err := task.RequiredStatus(StatusSecondaryMode); err != nil {
		return err
	}

	task.NewAction(func(ctx context.Context) error {
		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			return err
		}

		if err := node.connect(ctx); err != nil {
			return err
		}

		if err := node.Session(ctx, task.config.PastelID, sessID); err != nil {
			return err
		}

		task.connectedTo = node
		task.UpdateStatus(StatusConnected)
		return nil
	})
	return nil
}

// ProbeImage uploads the resampled image compute and return a fingerpirnt.
func (task *Task) ProbeImage(_ context.Context, file *artwork.File) (*pastel.FingerAndScores, error) {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		return nil, err
	}

	task.NewAction(func(ctx context.Context) error {
		task.ResampledArtwork = file
		defer task.ResampledArtwork.Remove()

		<-ctx.Done()
		return nil
	})

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusImageProbed)

		img, err := file.LoadImage()
		if err != nil {
			return errors.Errorf("failed to load image %s %w", file.Name(), err)
		}

		task.fingerAndScores, task.fingerprints, err = task.genFingerprintsData(ctx, img)
		if err != nil {
			return errors.Errorf("failed to generate fingerprints data %w", err)
		}
		// NOTE: for testing Kademlia and should be removed before releasing.
		// data, err := file.Bytes()
		// if err != nil {
		// 	return err
		// }

		// id, err := task.p2pClient.Store(ctx, data)
		// if err != nil {
		// 	return err
		// }
		// log.WithContext(ctx).WithField("id", id).Debugf("Image stored into Kademlia")
		return nil
	})

	return task.fingerAndScores, nil
}

// GetRegistrationFee get the fee to register artwork to bockchain
func (task *Task) GetRegistrationFee(_ context.Context, ticket []byte, artistSignature []byte, key1 string, key2 string, rqids map[string][]byte, oti []byte) (int64, error) {
	var err error
	if err = task.RequiredStatus(StatusImageAndThumbnailCoordinateUploaded); err != nil {
		return 0, errors.Errorf("require status %s not satisfied", StatusImageAndThumbnailCoordinateUploaded)
	}

	task.Oti = oti
	task.RQIDS = rqids
	task.artistSignature = artistSignature
	task.key1 = key1
	task.key2 = key2

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusRegistrationFeeCalculated)

		// TODO: fix this like how can we get the signature before calling cNode
		task.Ticket, err = pastel.DecodeArtTicket(ticket)
		if err != nil {
			return errors.Errorf("failed to decode art ticket %w", err)
		}

		// Assume passphase is 16-bytes length

		getFeeRequest := pastel.GetRegisterArtFeeRequest{
			Ticket: task.Ticket,
			Signatures: &pastel.TicketSignatures{
				Artist: map[string]string{
					task.Ticket.AppTicketData.AuthorPastelID: string(artistSignature),
				},
				Mn2: map[string]string{
					task.Service.config.PastelID: string(artistSignature),
				},
				Mn3: map[string]string{
					task.Service.config.PastelID: string(artistSignature),
				},
			},
			Mn1PastelID: task.Service.config.PastelID, // all ID has same lenght, so can use any id here
			Passphrase:  task.config.PassPhrase,
			Key1:        key1,
			Key2:        key2,
			Fee:         0, // fake data
			ImgSizeInMb: int64(task.imageSizeBytes) / (1024 * 1024),
		}

		task.registrationFee, err = task.pastelClient.GetRegisterArtFee(ctx, getFeeRequest)
		if err != nil {
			return errors.Errorf("failed to get register art fee %w", err)
		}
		return nil
	})

	return task.registrationFee, nil
}

// ValidatePreBurnTransaction will get pre-burnt transaction fee txid, wait until it's confirmations meet expectation.
func (task *Task) ValidatePreBurnTransaction(ctx context.Context, txid string) (string, error) {
	var err error
	if err = task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return "", errors.Errorf("require status %s not satisfied", StatusRegistrationFeeCalculated)
	}

	log.WithContext(ctx).Debugf("preburn-txid: %s", txid)
	<-task.NewAction(func(ctx context.Context) error {
		confirmationChn := task.waitConfirmation(ctx, txid, int64(task.config.PreburntTxMinConfirmations), task.config.PreburntTxConfirmationTimeout)

		if err := task.matchFingersPrintAndScores(ctx); err != nil {
			return errors.Errorf("fingerprints or scores don't matched")
		}

		// compare rqsymbols
		if err := task.compareRQSymbolID(ctx); err != nil {
			return errors.Errorf("failed to generate rqids %w", err)
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.connectedTo == nil)
		if err := task.signAndSendArtTicket(ctx, task.connectedTo == nil); err != nil {
			return errors.Errorf("failed to signed and send art ticket")
		}

		log.WithContext(ctx).Debugf("waiting for confimation")
		if err := <-confirmationChn; err != nil {
			return errors.Errorf("failed to validate preburn transaction validation %w", err)
		}
		log.WithContext(ctx).Debugf("confirmation done")

		return nil
	})

	// only primary node start this action
	var artRegTxid string
	if task.connectedTo == nil {
		<-task.NewAction(func(ctx context.Context) error {
			log.WithContext(ctx).Debug("waiting for signature from peers")
			for {
				select {
				case <-ctx.Done():
					err := ctx.Err()
					if err != nil {
						log.WithContext(ctx).Debug("waiting for signature from peers cancelled or timeout")
					}
					return err
				case <-task.allSignaturesReceived:
					log.WithContext(ctx).Debugf("all signature received so start validation")

					if err := task.verifyPeersSingature(ctx); err != nil {
						return errors.Errorf("peers' singature mismatched %w", err)
					}

					artRegTxid, err = task.registerArt(ctx)
					if err != nil {
						return errors.Errorf("failed to register art %w", err)
					}

					// confirmations := task.waitConfirmation(ctx, artRegTxid, 10, 30*time.Second, 55)
					// err = <-confirmations
					// if err != nil {
					// 	return errors.Errorf("failed to wait for confirmation of reg-art ticket %w", err)
					// }

					if err := task.storeRaptorQSymbols(ctx); err != nil {
						return errors.Errorf("failed to store raptor symbols %w", err)
					}

					if err := task.storeThumbnails(ctx); err != nil {
						return errors.Errorf("failed to store thumbnails %w", err)
					}

					if err := task.storeFingerprints(ctx); err != nil {
						return errors.Errorf("failed to store fingerprints %w", err)
					}

					return nil
				}
			}
		})
	}

	return artRegTxid, nil
}

func (task *Task) matchFingersPrintAndScores(ctx context.Context) error {
	log.WithContext(ctx).Debug("match fingerprints and scores")

	copyArtwork, err := task.Artwork.Copy()
	if err != nil {
		return errors.Errorf("failed to create copy of artwork %s %w", task.Artwork.Name(), err)
	}
	defer copyArtwork.Remove()

	if err := copyArtwork.ResizeImage(224, 224); err != nil {
		return errors.Errorf("failed to resize copy of artwork %w", err)
	}

	resizeImg, err := copyArtwork.LoadImage()
	if err != nil {
		return errors.Errorf("failed to load image from copied artwork %w", err)
	}

	fingerAndScores, _, err := task.genFingerprintsData(ctx, resizeImg)
	if err != nil {
		return errors.Errorf("failed to generate fingerprints data %w", err)
	}

	if err := pastel.CompareFingerPrintAndScore(task.fingerAndScores, fingerAndScores); err != nil {
		return errors.Errorf("fingerprint or scores are not equal: %w", err)
	}
	return nil
}

func (task *Task) waitConfirmation(ctx context.Context, txid string, minConfirmation int64, timeout time.Duration) <-chan error {
	ch := make(chan error)
	timeoutCh := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		close(timeoutCh)
	}()

	go func(ctx context.Context, txid string) {
		defer close(ch)
		retry := 0
		for {
			select {
			case <-ctx.Done():
				// context cancelled or abort by caller so no need to return anything
				log.WithContext(ctx).Debugf("context done: %s", ctx.Err())
				ch <- ctx.Err()
				return
			case <-time.After(15 * time.Second):
				log.WithContext(ctx).Debugf("retry: %d", retry)
				txResult, _ := task.pastelClient.GetRawTransactionVerbose1(ctx, txid)
				if txResult.Confirmations >= minConfirmation {
					ch <- nil
					return
				}
			case <-timeoutCh:
				ch <- errors.Errorf("timeout when wating for confirmation of transaction %s", txid)
				return
			}
		}
	}(ctx, txid)
	return ch
}

// sign and send art ticket if not primary
func (task *Task) signAndSendArtTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeArtTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("failed to serialize art ticket %w", err)
	}
	task.ownSignature, err = task.pastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase, "ed448")
	if err != nil {
		return errors.Errorf("failed to sign ticket %w", err)
	}
	if !isPrimary {
		log.WithContext(ctx).Debug("send signed articket to primary node")
		if err := task.connectedTo.RegisterArtwork.SendArtTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("failed to send signature to primary node %s at address %s %w", task.connectedTo.ID, task.connectedTo.Address, err)
		}
	}
	return nil
}

func (task *Task) verifyPeersSingature(ctx context.Context) error {
	log.WithContext(ctx).Debugf("all signature received so start validation")

	data, err := pastel.EncodeArtTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("failed to encoded art ticket %w", err)
	}
	for nodeID, signature := range task.peersArtTicketSignature {
		if ok, err := task.pastelClient.Verify(ctx, data, string(signature), nodeID, "ed448"); err != nil {
			return errors.Errorf("failed to verify signature %s of node %s", signature, nodeID)
		} else if !ok {
			return errors.Errorf("signature of node %s mistmatch", nodeID)
		}
	}
	return nil
}

func (task *Task) registerArt(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debugf("all signature received so start validation")

	req := pastel.RegisterArtRequest{
		Ticket: &pastel.ArtTicket{
			Version:       task.Ticket.Version,
			Author:        task.Ticket.Author,
			BlockNum:      task.Ticket.BlockNum,
			BlockHash:     task.Ticket.BlockHash,
			Copies:        task.Ticket.Copies,
			Royalty:       task.Ticket.Royalty,
			Green:         task.Ticket.Green,
			AppTicketData: task.Ticket.AppTicketData,
		},
		Signatures: &pastel.TicketSignatures{
			Artist: map[string]string{
				task.Ticket.AppTicketData.AuthorPastelID: string(task.artistSignature),
			},
			Mn1: map[string]string{
				task.config.PastelID: string(task.ownSignature),
			},
			Mn2: map[string]string{
				task.accepted[0].ID: string(task.peersArtTicketSignature[task.accepted[0].ID]),
			},
			Mn3: map[string]string{
				task.accepted[1].ID: string(task.peersArtTicketSignature[task.accepted[1].ID]),
			},
		},
		Mn1PastelID: task.config.PastelID,
		Pasphase:    task.config.PassPhrase,
		// TODO: fix this when how to get key1 and key2 are finalized
		Key1: task.key1,
		Key2: task.key2,
		Fee:  task.registrationFee,
	}

	artRegTxid, err := task.pastelClient.RegisterArtTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("failed to register art work %s", err)
	}
	return artRegTxid, nil
}

func (task *Task) storeRaptorQSymbols(ctx context.Context) error {
	img, err := task.Artwork.Bytes()
	if err != nil {
		return errors.Errorf("failed to read image data %w", err)
	}

	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("failed to connect to raptorq service %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	})

	encodeResp, err := rqService.Encode(ctx, img)
	if err != nil {
		return errors.Errorf("failed to create raptorq symbol from image %s %w", task.Artwork.Name(), err)
	}

	for name, data := range task.RQIDS {
		if _, err := task.p2pClient.Store(ctx, data); err != nil {
			return errors.Errorf("failed to store raptorq symbols id files %s %w", name, err)
		}
	}

	for id, symbol := range encodeResp.Symbols {
		if _, err := task.p2pClient.Store(ctx, symbol); err != nil {
			return errors.Errorf("failed to store symbolid %s into kamedila %w", id, err)
		}
	}

	return nil
}

func (task *Task) storeThumbnails(ctx context.Context) error {
	storeFn := func(ctx context.Context, artwork *artwork.File) (string, error) {
		data, err := artwork.Bytes()
		if err != nil {
			return "", errors.Errorf("failed to get data from artwork %s", artwork.Name())
		}
		return task.p2pClient.Store(ctx, data)
	}

	if _, err := storeFn(ctx, task.PreviewThumbnail); err != nil {
		return errors.Errorf("failed to store preview thumbnail into kamedila %w", err)
	}
	if _, err := storeFn(ctx, task.MediumThumbnail); err != nil {
		return errors.Errorf("failed to store medium thumbnail into kamedila %w", err)
	}
	if _, err := storeFn(ctx, task.SmallThumbnail); err != nil {
		return errors.Errorf("failed to store small thumbnail into kamedila %w", err)
	}

	return nil
}

func (task *Task) storeFingerprints(ctx context.Context) error {
	data := task.fingerprints.Bytes()
	if _, err := task.p2pClient.Store(ctx, data); err != nil {
		return errors.Errorf("failed to store fingerprints into kamedila")
	}
	return nil
}

func (task *Task) genFingerprintsData(ctx context.Context, img image.Image) (*pastel.FingerAndScores, pastel.Fingerprint, error) {
	// TODO: Fix this when implementation to dd service completed
	// fingerprints, err := task.probeTensor.Fingerprints(ctx, img)
	// if err != nil {
	// 	return nil, probe.Fingerprints{}, err
	// }

	fingerprint := pastel.Fingerprint{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1}
	fingerprintAndScores := pastel.FingerAndScores{
		DupeDectectionSystemVersion: "1.0",
		HashOfCandidateImageFile:    "HashOfCandidateImageFile",
		IsLikelyDupe:                0,
		OverallAverageRarenessScore: 0.5,
		IsRareOnInternet:            0,
		MatchesFoundOnFirstPage:     5,
		OpenNSFWScore:               0.5,
		ZstdCompressedFingerprint:   nil,
		AlternativeNSFWScore: pastel.AlternativeNSFWScore{
			Drawing: 0.5,
			Hentai:  0.5,
			Neutral: 0.5,
			Porn:    0.5,
			Sexy:    0.5,
		},
	}
	compressedFg, err := zstd.CompressLevel(nil, fingerprint.Bytes(), 22)
	if err != nil {
		return nil, nil, errors.Errorf("failed to compress fingerprint data: %w", err)
	}

	fingerprintAndScores.ZstdCompressedFingerprint = compressedFg
	return &fingerprintAndScores, fingerprint, nil
}

func (task *Task) compareRQSymbolID(ctx context.Context) error {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("failed to connect to raptorQ service %w", err)
	}
	defer conn.Close()

	content, err := task.Artwork.Bytes()
	if err != nil {
		return errors.Errorf("failed to read image contents")
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	})

	if len(task.RQIDS) == 0 {
		return errors.Errorf("no symbols identifiers file")
	}

	encodeInfo, err := rqService.EncodeInfo(ctx, content, uint32(len(task.RQIDS)), hex.EncodeToString(task.Ticket.BlockHash), task.Ticket.AppTicketData.AuthorPastelID)

	if err != nil {
		return errors.Errorf("failed to generate RaptorQ symbols' identifiers %w", err)
	}

	// pick just one file from wallnode to compare rq symbols
	var rqSymbolIDFile rq.SymbolIdFile
	for _, v := range task.RQIDS {
		if err := json.Unmarshal(v, &rqSymbolIDFile); err != nil {
			return errors.Errorf("failed to unmarshal raptorq symbols identifiers file %w", err)
		}
		break
	}

	// pick just one file generated to compare
	var rqRawSymbolIDFile rqnode.RawSymbolIdFile
	for _, v := range encodeInfo.SymbolIdFiles {
		rqRawSymbolIDFile = v
		break
	}
	sort.Strings(rqSymbolIDFile.SymbolIdentifiers)
	sort.Strings(rqRawSymbolIDFile.SymbolIdentifiers)
	if len(rqRawSymbolIDFile.SymbolIdentifiers) != len(rqSymbolIDFile.SymbolIdentifiers) {
		return errors.Errorf("number of raptorq symbols doesn't matched")
	}

	for i := range rqSymbolIDFile.SymbolIdentifiers {
		if rqSymbolIDFile.SymbolIdentifiers[i] != rqRawSymbolIDFile.SymbolIdentifiers[i] {
			return errors.Errorf("raptor symbol mismatched, index: %d, wallet:%s, super: %s", i, rqSymbolIDFile.SymbolIdentifiers[i], rqRawSymbolIDFile.SymbolIdentifiers[i])
		}
	}

	return nil
}

// UploadImageWithThumbnail uploads the image that contained image with pqsignature
// generate the image thumbnail from the coordinate provided for user and return
// the hash for the genreated thumbnail
func (task *Task) UploadImageWithThumbnail(_ context.Context, file *artwork.File, coordinate artwork.ThumbnailCoordinate) ([]byte, []byte, []byte, error) {
	var err error
	if err = task.RequiredStatus(StatusImageProbed); err != nil {
		return nil, nil, nil, errors.Errorf("require status %s not satisfied", StatusImageProbed)
	}

	previewThumbnailHash := make([]byte, 0)
	mediumThumbnailHash := make([]byte, 0)
	smallThumbnailHash := make([]byte, 0)
	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusImageAndThumbnailCoordinateUploaded)

		task.Artwork = file

		// Determine file size
		// TODO: improve it by call stats on file
		fileBytes, err := file.Bytes()
		if err != nil {
			return errors.Errorf("failed to read image file %w", err)
		}
		task.imageSizeBytes = len(fileBytes)

		previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, err = task.createAndHashThumbnails(coordinate)

		return err
	})

	return previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, nil
}

// AddPeerArticketSignature is called by supernode peers to primary first-rank supernode to send it's signature
func (task *Task) AddPeerArticketSignature(nodeID string, signature []byte) error {
	task.peersArtTicketSignatureMtx.Lock()
	defer task.peersArtTicketSignatureMtx.Unlock()

	if err := task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive art ticket signature from node %s", nodeID)
		if node := task.accepted.ByID(nodeID); node == nil {
			return errors.Errorf("node %s not in accepted list", nodeID)
		}

		task.peersArtTicketSignature[nodeID] = signature
		if len(task.peersArtTicketSignature) == len(task.accepted) {
			log.WithContext(ctx).Debug("all signature received")
			go func() {
				close(task.allSignaturesReceived)
			}()
		}
		return nil
	})
	return nil
}

func (task *Task) pastelNodeByExtKey(ctx context.Context, nodeID string) (*Node, error) {
	masterNodes, err := task.pastelClient.MasterNodesTop(ctx)
	log.WithContext(ctx).Debugf("master node %v", masterNodes)

	if err != nil {
		return nil, err
	}

	for _, masterNode := range masterNodes {
		if masterNode.ExtKey != nodeID {
			continue
		}
		node := &Node{
			client:  task.Service.nodeClient,
			ID:      masterNode.ExtKey,
			Address: masterNode.ExtAddress,
		}
		return node, nil
	}

	return nil, errors.Errorf("node %q not found", nodeID)
}

func (task *Task) context(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))
}

// NewTask returns a new Task instance.
func NewTask(service *Service) *Task {
	return &Task{
		Task:                       task.New(StatusTaskStarted),
		Service:                    service,
		peersArtTicketSignatureMtx: &sync.Mutex{},
		peersArtTicketSignature:    make(map[string][]byte),
		allSignaturesReceived:      make(chan struct{}),
	}
}
