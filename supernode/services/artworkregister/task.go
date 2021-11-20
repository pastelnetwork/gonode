package artworkregister

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	Ticket           *pastel.NFTTicket
	ResampledArtwork *artwork.File
	Artwork          *artwork.File
	imageSizeBytes   int

	fingerAndScores *pastel.FingerAndScores
	fingerprints    pastel.Fingerprint

	PreviewThumbnail *artwork.File
	MediumThumbnail  *artwork.File
	SmallThumbnail   *artwork.File

	acceptedMu sync.Mutex
	accepted   Nodes

	RQIDS map[string][]byte
	Oti   []byte

	// signature of ticket data signed by node's pateslID
	ownSignature []byte

	creatorSignature []byte
	key1             string
	key2             string
	registrationFee  int64

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
		log.WithContext(ctx).WithField("status", status.String()).Debug("States updated")
	})

	defer task.removeArtifacts()
	return task.RunAction(ctx)
}

// Session is handshake wallet to supernode
func (task *Task) Session(_ context.Context, isPrimary bool) error {
	if err := task.RequiredStatus(StatusTaskStarted); err != nil {
		return err
	}

	<-task.NewAction(func(ctx context.Context) error {
		if isPrimary {
			log.WithContext(ctx).Debug("Acts as primary node")
			task.UpdateStatus(StatusPrimaryMode)
			return nil
		}

		log.WithContext(ctx).Debug("Acts as secondary node")
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
		log.WithContext(ctx).Debug("Waiting for supernodes to connect")

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

	var err error

	<-task.NewAction(func(ctx context.Context) error {
		if node := task.accepted.ByID(nodeID); node != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).Errorf("node is already registered")
			err = errors.Errorf("node %q is already registered", nodeID)
			return nil
		}

		var node *Node
		node, err = task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("get node by extID")
			err = errors.Errorf("get node by extID %s: %w", nodeID, err)
			return nil
		}
		task.accepted.Add(node)

		log.WithContext(ctx).WithField("nodeID", nodeID).Debug("Accept secondary node")

		if len(task.accepted) >= task.config.NumberConnectedNodes {
			task.UpdateStatus(StatusConnected)
		}
		return nil
	})
	return err
}

// ConnectTo connects to primary node
func (task *Task) ConnectTo(_ context.Context, nodeID, sessID string) error {
	if err := task.RequiredStatus(StatusSecondaryMode); err != nil {
		return err
	}

	var err error

	<-task.NewAction(func(ctx context.Context) error {
		var node *Node
		node, err = task.pastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("get node by extID")
			return nil
		}

		if err = node.connect(ctx); err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("connect to node")
			return nil
		}

		if err = node.Session(ctx, task.config.PastelID, sessID); err != nil {
			log.WithContext(ctx).WithField("sessID", sessID).WithField("pastelID", task.config.PastelID).WithError(err).Errorf("handsake with peer")
			return nil
		}

		task.connectedTo = node
		task.UpdateStatus(StatusConnected)
		return nil
	})
	return err
}

// ProbeImage uploads the resampled image compute and return a fingerpirnt.
func (task *Task) ProbeImage(_ context.Context, file *artwork.File) (*pastel.FingerAndScores, error) {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		return nil, err
	}

	var err error

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusImageProbed)

		task.ResampledArtwork = file
		task.fingerAndScores, task.fingerprints, err = task.genFingerprintsData(ctx, task.ResampledArtwork)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("generate fingerprints data")
			err = errors.Errorf("generate fingerprints data: %w", err)
			return nil
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

	return task.fingerAndScores, err
}

// GetRegistrationFee get the fee to register artwork to bockchain
func (task *Task) GetRegistrationFee(_ context.Context, ticket []byte, creatorSignature []byte, key1 string, key2 string, rqids map[string][]byte, oti []byte) (int64, error) {
	var err error
	if err = task.RequiredStatus(StatusImageAndThumbnailCoordinateUploaded); err != nil {
		return 0, errors.Errorf("require status %s not satisfied", StatusImageAndThumbnailCoordinateUploaded)
	}

	task.Oti = oti
	task.RQIDS = rqids
	task.creatorSignature = creatorSignature
	task.key1 = key1
	task.key2 = key2

	<-task.NewAction(func(ctx context.Context) error {
		task.UpdateStatus(StatusRegistrationFeeCalculated)

		// TODO: fix this like how can we get the signature before calling cNode
		task.Ticket, err = pastel.DecodeNFTTicket(ticket)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("decode NFT ticket")
			err = errors.Errorf("decode NFT ticket %w", err)
			return nil
		}

		// Assume passphase is 16-bytes length

		getFeeRequest := pastel.GetRegisterNFTFeeRequest{
			Ticket: task.Ticket,
			Signatures: &pastel.TicketSignatures{
				Creator: map[string]string{
					task.Ticket.AppTicketData.AuthorPastelID: string(creatorSignature),
				},
				Mn2: map[string]string{
					task.Service.config.PastelID: string(creatorSignature),
				},
				Mn3: map[string]string{
					task.Service.config.PastelID: string(creatorSignature),
				},
			},
			Mn1PastelID: task.Service.config.PastelID, // all ID has same lenght, so can use any id here
			Passphrase:  task.config.PassPhrase,
			Key1:        key1,
			Key2:        key2,
			Fee:         0, // fake data
			ImgSizeInMb: int64(task.imageSizeBytes) / (1024 * 1024),
		}

		task.registrationFee, err = task.pastelClient.GetRegisterNFTFee(ctx, getFeeRequest)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("get register NFT fee")
			err = errors.Errorf("get register NFT fee %w", err)
		}
		return nil
	})

	return task.registrationFee, err
}

// ValidatePreBurnTransaction will get pre-burnt transaction fee txid, wait until it's confirmations meet expectation.
func (task *Task) ValidatePreBurnTransaction(ctx context.Context, txid string) (string, error) {
	var err error
	if err = task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return "", errors.Errorf("require status %s not satisfied", StatusRegistrationFeeCalculated)
	}

	log.WithContext(ctx).Debugf("preburn-txid: %s", txid)
	<-task.NewAction(func(ctx context.Context) error {
		confirmationChn := task.waitConfirmation(ctx, txid, int64(task.config.PreburntTxMinConfirmations), task.config.PreburntTxConfirmationTimeout, 15*time.Second)

		// compare rqsymbols
		if err = task.compareRQSymbolID(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("generate rqids")
			err = errors.Errorf("generate rqids: %w", err)
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.connectedTo == nil)
		if err = task.signAndSendArtTicket(ctx, task.connectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("signed and send NFT ticket")
			err = errors.Errorf("signed and send NFT ticket")
			return nil
		}

		log.WithContext(ctx).Debug("waiting for confimation")
		if err = <-confirmationChn; err != nil {
			log.WithContext(ctx).WithError(err).Errorf("validate preburn transaction validation")
			err = errors.Errorf("validate preburn transaction validation :%w", err)
			return nil
		}
		log.WithContext(ctx).Debug("confirmation done")

		return nil
	})

	if err != nil {
		return "", err
	}

	// only primary node start this action
	var nftRegTxid string
	if task.connectedTo == nil {
		<-task.NewAction(func(ctx context.Context) error {
			log.WithContext(ctx).Debug("waiting for signature from peers")
			for {
				select {
				case <-ctx.Done():
					err = ctx.Err()
					if err != nil {
						log.WithContext(ctx).Debug("waiting for signature from peers cancelled or timeout")
					}
					return nil
				case <-task.allSignaturesReceived:
					log.WithContext(ctx).Debug("all signature received so start validation")

					if err = task.verifyPeersSingature(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("peers' singature mismatched: %w", err)
						return nil
					}

					nftRegTxid, err = task.registerArt(ctx)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("register NFT: %w", err)
						return nil
					}

					// confirmations := task.waitConfirmation(ctx, nftRegTxid, 10, 30*time.Second, 55)
					// err = <-confirmations
					// if err != nil {
					// 	return errors.Errorf("wait for confirmation of reg-art ticket %w", err)
					// }

					if err = task.storeRaptorQSymbols(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("store raptor symbols")
						err = errors.Errorf("store raptor symbols: %w", err)
						return nil
					}

					if err = task.storeThumbnails(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("store thumbnails")
						err = errors.Errorf("store thumbnails: %w", err)
						return nil
					}

					if err = task.storeFingerprints(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("store fingerprints")
						err = errors.Errorf("store fingerprints: %w", err)
						return nil
					}

					return nil
				}
			}
		})
	}

	return nftRegTxid, err
}

func (task *Task) waitConfirmation(ctx context.Context, txid string, minConfirmation int64, timeout time.Duration, interval time.Duration) <-chan error {
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
			case <-time.After(interval):
				log.WithContext(ctx).Debugf("retry: %d", retry)
				retry++
				txResult, _ := task.pastelClient.GetRawTransactionVerbose1(ctx, txid)
				if txResult.Confirmations >= minConfirmation {
					log.WithContext(ctx).Debug("transaction confirmed")
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

// sign and send NFT ticket if not primary
func (task *Task) signAndSendArtTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeNFTTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("serialize NFT ticket: %w", err)
	}
	task.ownSignature, err = task.pastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket: %w", err)
	}
	if !isPrimary {
		log.WithContext(ctx).Debug("send signed articket to primary node")
		if err := task.connectedTo.RegisterArtwork.SendArtTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s: %w", task.connectedTo.ID, task.connectedTo.Address, err)
		}
	}
	return nil
}

func (task *Task) verifyPeersSingature(ctx context.Context) error {
	log.WithContext(ctx).Debug("all signature received so start validation")

	data, err := pastel.EncodeNFTTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("encoded NFT ticket: %w", err)
	}
	for nodeID, signature := range task.peersArtTicketSignature {
		if ok, err := task.pastelClient.Verify(ctx, data, string(signature), nodeID, pastel.SignAlgorithmED448); err != nil {
			return errors.Errorf("verify signature %s of node %s", signature, nodeID)
		} else if !ok {
			return errors.Errorf("signature of node %s mistmatch", nodeID)
		}
	}
	return nil
}

func (task *Task) registerArt(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debug("all signature received so start validation")

	req := pastel.RegisterNFTRequest{
		Ticket: &pastel.NFTTicket{
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
			Creator: map[string]string{
				task.Ticket.AppTicketData.AuthorPastelID: string(task.creatorSignature),
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

	nftRegTxid, err := task.pastelClient.RegisterNFTTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("register NFT work: %w", err)
	}
	return nftRegTxid, nil
}

func (task *Task) storeRaptorQSymbols(ctx context.Context) error {
	img, err := task.Artwork.Bytes()
	if err != nil {
		return errors.Errorf("read image data: %w", err)
	}

	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("connect to raptorq service: %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	})

	encodeResp, err := rqService.Encode(ctx, img)
	if err != nil {
		return errors.Errorf("create raptorq symbol from image %s: %w", task.Artwork.Name(), err)
	}

	for name, data := range task.RQIDS {
		if _, err := task.p2pClient.Store(ctx, data); err != nil {
			return errors.Errorf("store raptorq symbols id files %s: %w", name, err)
		}
	}

	for id, symbol := range encodeResp.Symbols {
		if _, err := task.p2pClient.Store(ctx, symbol); err != nil {
			return errors.Errorf("store symbolid %s into kamedila: %w", id, err)
		}
	}

	return nil
}

func (task *Task) storeThumbnails(ctx context.Context) error {
	storeFn := func(ctx context.Context, artwork *artwork.File) (string, error) {
		data, err := artwork.Bytes()
		if err != nil {
			return "", errors.Errorf("get data from artwork %s", artwork.Name())
		}
		return task.p2pClient.Store(ctx, data)
	}

	if _, err := storeFn(ctx, task.PreviewThumbnail); err != nil {
		return errors.Errorf("store preview thumbnail into kamedila: %w", err)
	}
	if _, err := storeFn(ctx, task.MediumThumbnail); err != nil {
		return errors.Errorf("store medium thumbnail into kamedila: %w", err)
	}
	if _, err := storeFn(ctx, task.SmallThumbnail); err != nil {
		return errors.Errorf("store small thumbnail into kamedila: %w", err)
	}

	return nil
}

func (task *Task) storeFingerprints(ctx context.Context) error {
	compressFringerprints, err := zstd.CompressLevel(nil, pastel.Fingerprint(task.fingerprints).Bytes(), 22)
	if err != nil {
		return errors.Errorf("compress fingerprint")
	}

	if _, err := task.p2pClient.Store(ctx, compressFringerprints); err != nil {
		return errors.Errorf("store fingerprints into kamedila")
	}
	return nil
}

func (task *Task) genFingerprintsData(ctx context.Context, file *artwork.File) (*pastel.FingerAndScores, pastel.Fingerprint, error) {
	img, err := file.Bytes()
	if err != nil {
		return nil, nil, errors.Errorf("get content of image %s: %w", file.Name(), err)
	}

	ddResult, err := task.ddClient.ImageRarenessScore(ctx, img, file.Format().String())

	if err != nil {
		return nil, nil, errors.Errorf("get dupe detection result from dd-server: %w", err)
	}

	fingerprint := ddResult.Fingerprints

	fingerprintAndScores := pastel.FingerAndScores{
		DupeDectectionSystemVersion: ddResult.DupeDetectionSystemVer,
		HashOfCandidateImageFile:    ddResult.ImageHash,
		OverallAverageRarenessScore: ddResult.PastelRarenessScore,
		IsLikelyDupe:                ddResult.IsLikelyDupe,
		IsRareOnInternet:            ddResult.IsRareOnInternet,
		NumberOfPagesOfResults:      ddResult.NumberOfResultPages,
		MatchesFoundOnFirstPage:     ddResult.MatchesFoundOnFirstPage,
		URLOfFirstMatchInPage:       ddResult.FirstMatchURL,
		OpenNSFWScore:               ddResult.OpenNSFWScore,
		ZstdCompressedFingerprint:   nil,
		AlternativeNSFWScore: pastel.AlternativeNSFWScore{
			Drawing: ddResult.AlternateNSFWScores.Drawings,
			Hentai:  ddResult.AlternateNSFWScores.Hentai,
			Neutral: ddResult.AlternateNSFWScores.Neutral,
			Porn:    ddResult.AlternateNSFWScores.Porn,
			Sexy:    ddResult.AlternateNSFWScores.Sexy,
		},
		PerceptualImageHashes: pastel.PerceptualImageHashes{
			PerceptualHash: ddResult.PerceptualImageHashes.PerceptualHash,
			AverageHash:    ddResult.PerceptualImageHashes.AverageHash,
			DifferenceHash: ddResult.PerceptualImageHashes.DifferenceHash,
			PDQHash:        ddResult.PerceptualImageHashes.PDQHash,
			NeuralHash:     ddResult.PerceptualImageHashes.NeuralHash,
		},
		PerceptualHashOverlapCount:                   ddResult.PerceptualHashOverlapCount,
		NumberOfFingerprintsRequiringFurtherTesting1: ddResult.NumberOfFingerprintsRequiringFurtherTesting1,
		NumberOfFingerprintsRequiringFurtherTesting2: ddResult.NumberOfFingerprintsRequiringFurtherTesting2,
		NumberOfFingerprintsRequiringFurtherTesting3: ddResult.NumberOfFingerprintsRequiringFurtherTesting3,
		NumberOfFingerprintsRequiringFurtherTesting4: ddResult.NumberOfFingerprintsRequiringFurtherTesting4,
		NumberOfFingerprintsRequiringFurtherTesting5: ddResult.NumberOfFingerprintsRequiringFurtherTesting5,
		NumberOfFingerprintsRequiringFurtherTesting6: ddResult.NumberOfFingerprintsRequiringFurtherTesting6,
		NumberOfFingerprintsOfSuspectedDupes:         ddResult.NumberOfFingerprintsOfSuspectedDupes,
		PearsonMax:                                   ddResult.PearsonMax,
		SpearmanMax:                                  ddResult.SpearmanMax,
		KendallMax:                                   ddResult.KendallMax,
		HoeffdingMax:                                 ddResult.HoeffdingMax,
		MutualInformationMax:                         ddResult.MutualInformationMax,
		HsicMax:                                      ddResult.HsicMax,
		XgbimportanceMax:                             ddResult.XgbimportanceMax,
		PearsonTop1BpsPercentile:                     ddResult.PearsonTop1BpsPercentile,
		SpearmanTop1BpsPercentile:                    ddResult.SpearmanTop1BpsPercentile,
		KendallTop1BpsPercentile:                     ddResult.KendallTop1BpsPercentile,
		HoeffdingTop10BpsPercentile:                  ddResult.HoeffdingTop10BpsPercentile,
		MutualInformationTop100BpsPercentile:         ddResult.MutualInformationTop100BpsPercentile,
		HsicTop100BpsPercentile:                      ddResult.HsicTop100BpsPercentile,
		XgbimportanceTop100BpsPercentile:             ddResult.XgbimportanceTop100BpsPercentile,
		CombinedRarenessScore:                        ddResult.CombinedRarenessScore,
		XgboostPredictedRarenessScore:                ddResult.XgboostPredictedRarenessScore,
		NnPredictedRarenessScore:                     ddResult.NnPredictedRarenessScore,
	}
	compressedFg, err := zstd.CompressLevel(nil, pastel.Fingerprint(fingerprint).Bytes(), 22)
	if err != nil {
		return nil, nil, errors.Errorf("compress fingerprint data: %w", err)
	}

	fingerprintAndScores.ZstdCompressedFingerprint = compressedFg
	return &fingerprintAndScores, fingerprint, nil
}

func (task *Task) compareRQSymbolID(ctx context.Context) error {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return errors.Errorf("connect to raptorQ service: %w", err)
	}
	defer conn.Close()

	content, err := task.Artwork.Bytes()
	if err != nil {
		return errors.Errorf("read image contents: %w", err)
	}

	rqService := conn.RaptorQ(&rqnode.Config{
		RqFilesDir: task.Service.config.RqFilesDir,
	})

	if len(task.RQIDS) == 0 {
		return errors.Errorf("no symbols identifiers file")
	}

	encodeInfo, err := rqService.EncodeInfo(ctx, content, uint32(len(task.RQIDS)), hex.EncodeToString([]byte(task.Ticket.BlockHash)), task.Ticket.AppTicketData.AuthorPastelID)

	if err != nil {
		return errors.Errorf("generate RaptorQ symbols' identifiers: %w", err)
	}

	// pick just one file from wallnode to compare rq symbols
	// the one send from walletnode is compressed so need to decompress first
	var fromWalletNode []byte
	for _, v := range task.RQIDS {
		fromWalletNode = v
		break
	}
	rawData, err := zstd.Decompress(nil, fromWalletNode)
	if err != nil {
		return errors.Errorf("decompress symbols id file: %w", err)
	}
	scaner := bufio.NewScanner(bytes.NewReader(rawData))
	scaner.Split(bufio.ScanLines)
	for i := 0; i < 3; i++ {
		scaner.Scan()
		if scaner.Err() != nil {
			return errors.Errorf("failed when bypass headers: %w", err)
		}
	}
	scaner.Text()

	// pick just one file generated to compare
	var rawEncodeFile rqnode.RawSymbolIDFile
	for _, v := range encodeInfo.SymbolIDFiles {
		rawEncodeFile = v
		break
	}

	for i := range rawEncodeFile.SymbolIdentifiers {
		scaner.Scan()
		if scaner.Err() != nil {
			return errors.Errorf("scan next symbol identifiers: %w", err)
		}
		symbolFromWalletNode := scaner.Text()
		if rawEncodeFile.SymbolIdentifiers[i] != symbolFromWalletNode {
			return errors.Errorf("raptor symbol mismatched, index: %d, wallet:%s, super: %s", i, symbolFromWalletNode, rawEncodeFile.SymbolIdentifiers[i])
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
		var fileBytes []byte
		fileBytes, err = file.Bytes()
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("read image file")
			err = errors.Errorf("read image file: %w", err)
			return nil
		}
		task.imageSizeBytes = len(fileBytes)

		previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, err = task.createAndHashThumbnails(coordinate)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("create and hash thumbnail")
			return nil
		}

		return nil
	})

	return previewThumbnailHash, mediumThumbnailHash, smallThumbnailHash, err
}

// AddPeerArticketSignature is called by supernode peers to primary first-rank supernode to send it's signature
func (task *Task) AddPeerArticketSignature(nodeID string, signature []byte) error {
	task.peersArtTicketSignatureMtx.Lock()
	defer task.peersArtTicketSignatureMtx.Unlock()

	if err := task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return err
	}

	var err error

	<-task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debugf("receive NFT ticket signature from node %s", nodeID)
		if node := task.accepted.ByID(nodeID); node == nil {
			log.WithContext(ctx).WithField("node", nodeID).Errorf("node is not in accepted list")
			err = errors.Errorf("node %s not in accepted list", nodeID)
			return nil
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
	return err
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

func (task *Task) removeArtifacts() {
	removeFn := func(file *artwork.File) {
		if file != nil {
			log.Debugf("remove file: %s", file.Name())
			if err := file.Remove(); err != nil {
				log.Errorf("failed to remove file: %s", err.Error())
			}
		}
	}

	removeFn(task.ResampledArtwork)
	removeFn(task.Artwork)
	removeFn(task.PreviewThumbnail)
	removeFn(task.MediumThumbnail)
	removeFn(task.SmallThumbnail)
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
