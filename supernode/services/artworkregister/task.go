package artworkregister

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rq "github.com/pastelnetwork/gonode/raptorq"
	"golang.org/x/crypto/sha3"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Ticket           *pastel.ArtTicket
	ResampledArtwork *artwork.File
	Artwork          *artwork.File

	FingerprintsData []byte

	PreviewThumbnail *artwork.File
	MediumThumbnail  *artwork.File
	SmallThumbnail   *artwork.File

	acceptedMu sync.Mutex
	accpeted   Nodes

	RQIDS map[string][]byte
	Oti   []byte

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
	return task.accpeted, nil
}

// SessionNode accepts secondary node
func (task *Task) SessionNode(_ context.Context, nodeID string) error {
	task.acceptedMu.Lock()
	defer task.acceptedMu.Unlock()

	if err := task.RequiredStatus(StatusPrimaryMode); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if node := task.accpeted.ByID(nodeID); node != nil {
			return errors.Errorf("node %q is already registered", nodeID)
		}

		node, err := task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			return err
		}
		task.accpeted.Add(node)

		log.WithContext(ctx).WithField("nodeID", nodeID).Debugf("Accept secondary node")

		if len(task.accpeted) >= task.config.NumberConnectedNodes {
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
func (task *Task) ProbeImage(_ context.Context, file *artwork.File) ([]byte, error) {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		return nil, err
	}

	task.NewAction(func(ctx context.Context) error {
		task.ResampledArtwork = file
		defer task.ResampledArtwork.Remove()

		<-ctx.Done()
		return nil
	})

	var fingerprintData []byte

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusImageProbed)

		img, err := file.LoadImage()
		if err != nil {
			return errors.Errorf("failed to load image %s %w", file.Name(), err)
		}

		fingerprintData, err = task.genFingerprintsData(ctx, img)
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

	task.FingerprintsData = fingerprintData
	return fingerprintData, nil
}

// GetRegistrationFee get the fee to register artwork to bockchain
func (task *Task) GetRegistrationFee(ctx context.Context, ticket []byte, rqids map[string][]byte, oti []byte) (int64, error) {
	var err error
	if err = task.RequiredStatus(StatusImageAndThumbnailCoordinateUploaded); err != nil {
		return 0, errors.Errorf("require status %s not satisfied", StatusImageAndThumbnailCoordinateUploaded)
	}

	task.Oti = oti
	task.RQIDS = rqids

	var registrationFee int64
	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusRegistrationFeeCalculated)

		// TODO: fix this like how can we get the signature before calling cNode
		task.Ticket, err = pastel.DecodeArtTicket(ticket)
		if err != nil {
			return errors.Errorf("failed to decode art ticket %w", err)
		}

		// req := pastel.GetRegisterArtFeeRequest{
		// 	Ticket:     &ticket.RegTicketData.ArtTicketData,
		// 	Signatures: &ticket.RegTicketData.Signatures,
		// 	Mn1PastelId: func() {
		// 		if task.Status() == StatusPrimaryMode {
		// 			return
		// 		}
		// 	}() ,
		// }
		// task.pastelClient.GetRegisterArtFee(ctx, req)

		registrationFee = 15
		return nil
	})

	return registrationFee, nil
}

func (task *Task) ValidatePreBurnTransaction(ctx context.Context, txid string) error {
	var err error
	if err = task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return errors.Errorf("require status %s not satisfied", StatusRegistrationFeeCalculated)
	}

	// Only primary node start this action
	if task.connectedTo == nil {
		taskFn := func(ctx context.Context) {
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

					}
				}
			})
		}
		go taskFn(ctx)
	}

	<-task.NewAction(func(ctx context.Context) error {
		var confirmationChn chan error
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		validationFn := func(ctx context.Context) {
			var err error
			var confirmation *int64
			defer func() {
				if err != nil && err != context.Canceled {
					confirmationChn <- err
				} else {
					// Either error is nil or the context is cancelled
					// so no one read from confirmation
					close(confirmationChn)
				}
			}()

			const maxRetry = 3
			var txResult *pastel.GetTransactionResult
			for i := 0; i < maxRetry; i++ {
				txResult, err = task.pastelClient.GetTransaction(ctx, pastel.TxIDType(txid))
				if err != nil {
					err = errors.Errorf("failed to get transaction ID %w", err)
					return
				}
				if *confirmation = txResult.Confirmations; *confirmation >= 2 {
					err = nil
					return
				}

				// Sleep but responsive to context
				select {
				case <-ctx.Done():
					err = ctx.Err()
					return
				case <-time.After(150 * time.Second):
				}
			}

			err = errors.Errorf("timeout when wating for confirmation of transaction %s, confirmations %d ", txid, txResult.Confirmations)
		}
		go validationFn(ctx)

		// resample image and validate the finger prints data
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

		fingerprintsData, err := task.genFingerprintsData(ctx, resizeImg)
		if err != nil {
			return errors.Errorf("failed to generate fingerprints data %w", err)
		}

		if !bytes.Equal(task.FingerprintsData, fingerprintsData) {
			return errors.Errorf("fingerprints not matched")
		}

		// compare rqsymbols
		err = task.compareRQSymbolId(ctx)
		if err != nil {
			return errors.Errorf("failed to generate rqids %w", err)
		}

		// sign the ticket if not primary node
		if task.connectedTo != nil {
			ticket, err := pastel.EncodeArtTicket(task.Ticket)
			if err != nil {
				return errors.Errorf("failed to serialize art ticket %w", err)
			}
			signature, err := task.pastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase)
			if err != nil {
				return errors.Errorf("failed to sign ticket %w", err)
			}
			if err := task.connectedTo.RegisterArtwork.SendArtTicketSignature(ctx, signature); err != nil {
				return errors.Errorf("failed to send signature to primary node %s at address %s %w", task.connectedTo.ID, task.connectedTo.Address, err)
			}
		}

		err = <-confirmationChn
		if err != nil {
			return errors.Errorf("failed to validate preburn transaction validation %w", err)
		}
		return nil
	})

	return nil
}

func (task *Task) genFingerprintsData(ctx context.Context, img image.Image) ([]byte, error) {
	fingerprints, err := task.probeTensor.Fingerprints(ctx, img)
	if err != nil {
		return nil, err
	}

	fingerprintsData, err := zstd.CompressLevel(nil, fingerprints.Single().LSBTruncatedBytes(), 22)
	if err != nil {
		return nil, errors.Errorf("failed to compress fingerprint data: %w", err)
	}

	return fingerprintsData, nil
}

func (task *Task) compareRQSymbolId(ctx context.Context) error {
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

	rqService := conn.RaptorQ()
	encodeInfo, err := rqService.EncodeInfo(ctx, content)
	if err != nil {
		return errors.Errorf("failed to generate RaptorQ symbols' identifiers %w", err)
	}

	// pick just one file to compare rq symbols
	if len(task.RQIDS) == 0 {
		return errors.Errorf("no symbols identifiers file")
	}

	var rqSymbolIdFile rq.SymbolIdFile
	for _, v := range task.RQIDS {
		if err := json.Unmarshal(v, &rqSymbolIdFile); err != nil {
			return errors.Errorf("failed to unmarshal raptorq symbols identifiers file %w", err)
		}
		break
	}

	// sorting to make sure the comparision is correct
	sort.Strings(encodeInfo.SymbolIds)
	sort.Strings(rqSymbolIdFile.SymbolIdentifiers)
	if len(encodeInfo.SymbolIds) != len(rqSymbolIdFile.SymbolIdentifiers) {
		return errors.Errorf("number of raptorq symbols doesn't matched")
	}

	for i := range encodeInfo.SymbolIds {
		if encodeInfo.SymbolIds[i] != rqSymbolIdFile.SymbolIdentifiers[i] {
			return errors.Errorf("raptor symbol mismatched, index: %d, wallet:%s, super: %s", i, rqSymbolIdFile.SymbolIdentifiers[i], encodeInfo.SymbolIds[i])
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

		previewThumbnailHash, err = task.hashThumbnail(previewThumbnail, nil)
		if err != nil {
			return errors.Errorf("failed to generate thumbnail %w", err)
		}

		rect := image.Rect(int(coordinate.TopLeftX), int(coordinate.TopLeftY), int(coordinate.BottomRightX), int(coordinate.BottomRightY))
		mediumThumbnailHash, err = task.hashThumbnail(mediumThumbnail, &rect)
		if err != nil {
			return errors.Errorf("hash medium thumbnail failed %w", err)
		}

		smallThumbnailHash, err = task.hashThumbnail(smallThumbnail, &rect)
		if err != nil {
			return errors.Errorf("hash small thumbnail failed %w", err)
		}
		return nil
	})

	return previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, nil
}

func (task *Task) hashThumbnail(thumbnail thumbnailType, rect *image.Rectangle) ([]byte, error) {
	f := task.Storage.NewFile()
	if f == nil {
		return nil, errors.Errorf("failed to create thumbnail file")
	}

	previewFile, err := f.Create()
	if err != nil {
		return nil, errors.Errorf("failed to create file %s %w", f.Name(), err)
	}
	defer previewFile.Close()

	srcImg, err := task.Artwork.LoadImage()
	if err != nil {
		return nil, errors.Errorf("failed to load image from artwork %s %w", task.Artwork.Name(), err)
	}

	var thumbnailImg image.Image
	if rect != nil {
		thumbnailImg = imaging.Crop(srcImg, *rect)
	} else {
		thumbnailImg = srcImg
	}

	log.Debugf("Encode with target size %d and quality %f", thumbnails[thumbnail].targetSize, thumbnails[thumbnail].quality)
	encoderOptions, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, thumbnails[thumbnail].quality)
	encoderOptions.TargetSize = thumbnails[thumbnail].targetSize
	if err != nil {
		return nil, errors.Errorf("failed to create lossless encoder option %w", err)
	}
	if err := webp.Encode(previewFile, thumbnailImg, encoderOptions); err != nil {
		return nil, errors.Errorf("failed to encode to webp format %w", err)
	}
	log.Debugf("preview thumbnail %s", f.Name())

	previewFile.Seek(0, io.SeekStart)
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, previewFile); err != nil {
		return nil, errors.Errorf("hash failed %w", err)
	}

	return hasher.Sum(nil), nil
}

func (task *Task) AddPeerArticketSignature(nodeID string, signature []byte) error {
	task.peersArtTicketSignatureMtx.Lock()
	defer task.peersArtTicketSignatureMtx.Unlock()

	if err := task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive art ticket signature from node %s", nodeID)
		if node := task.accpeted.ByID(nodeID); node == nil {
			return errors.Errorf("node %s not in accepted list", nodeID)
		}

		task.peersArtTicketSignature[nodeID] = signature
		if len(task.peersArtTicketSignature) == len(task.accpeted) {
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
