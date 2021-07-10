package artworkregister

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/pastel"
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
	rqcoti  int64
	rqssoti int64

	// TODO: call cNodeAPI to get the following info
	blockTxID string
	blockNum  int
	burnTxId  pastel.TxIDType

	// TODO: the following info is supposed return by supernodes
	rarenessScore int
	nsfwScore     int
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
	if err := nodes.MatchFingerprints(); err != nil {
		task.UpdateStatus(StatusErrorFingerprintsNotMatch)
		return err
	}
	task.UpdateStatus(StatusImageProbed)
	task.fingerprints = nodes.Fingerprint()
	if task.fingerprintsHash, err = sha3256hash(task.fingerprints); err != nil {
		return errors.Errorf("failed to hash fingerprints %w", err)
	}

	// task.Ticket.Fingerprint = nodes.Fingerprint()
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
	rqidsList, err := task.genRQIDSList(ctx)
	if err != nil {
		task.UpdateStatus(StatusErrorGenRaptorQSymbolsFailed)
		return errors.Errorf("gen RaptorQ symbols' identifiers failed %w", err)
	}
	task.rqids = rqidsList.Identifiers()

	ticket, err := task.createTicket()
	if err != nil {
		return errors.Errorf("failed to create ticket %w", err)
	}

	buf, err := json.MarshalIndent(ticket, "", "  ")
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
	if err := nodes.UploadSignedTicket(ctx, buf, artistSignature); err != nil {
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

func (task *Task) genRQIDSList(ctx context.Context) (RQIDSList, error) {
	log.Debugf("Connect to %s", task.config.RaptorQServiceAddress)
	conn, err := task.rqClient.Connect(ctx, task.config.RaptorQServiceAddress)
	if err != nil {
		return nil, errors.Errorf("failed to connect to raptorQ service %w", err)
	}
	defer conn.Close()

	content, err := task.Request.Image.Bytes()
	if err != nil {
		return nil, errors.Errorf("failed to read image contents")
	}

	rq := conn.RaptorQ()
	encodeInfo, err := rq.EncodeInfo(ctx, content)
	if err != nil {
		return nil, errors.Errorf("failed to generate RaptorQ symbols' identifiers %w", err)
	}

	var list RQIDSList
	for i := 0; i < task.config.NumberRQIDSFiles; i++ {
		item, err := func() (*RQIDS, error) {
			// FIXME: fill the real block_hash after it is confirmed
			rqidsFiles, err := createRQIDSFile(encodeInfo.SymbolIds, []byte("BLOCK_HASH"))
			if err != nil {
				return nil, errors.Errorf("failed to create rqids file %w", err)
			}

			file, err := os.OpenFile(rqidsFiles, os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				return nil, errors.Errorf("failed to open rqids file %s %w", rqidsFiles, err)
			}
			defer file.Close()
			defer os.Remove(file.Name())

			buf, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, errors.Errorf("failed to read rqids file %s %w", rqidsFiles, err)
			}

			// sign the file
			signature, err := task.pastelClient.Sign(ctx, buf, task.Request.ArtistPastelID, task.Request.ArtistPastelIDPassphrase)
			if err != nil {
				return nil, errors.Errorf("failed to sign rqids file %s %w", rqidsFiles, err)
			}

			// FIXME: remove this after confimartion that file identifier = hash(file + signature)
			// zstd.BestCompression is 20
			// Maximum level compression is 22 and it was use to compress the fingerprints so here we do the same
			// compressedSignature, err := zstd.CompressLevel(nil, signature, 22)
			// if err != nil {
			// 	return nil, errors.Errorf("failed to compress signature of rqids file %s %w", rqidsFiles, err)
			// }

			// append the signature to end of the file
			if _, err := file.Write(signature); err != nil {
				return nil, errors.Errorf("failed to")
			}

			// hash the file and sum
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				return nil, errors.Errorf("failed to send to begin of file %s %w", file.Name(), err)
			}
			hasher := sha3.New256()
			if _, err := io.Copy(hasher, file); err != nil {
				return nil, errors.Errorf("failed to hash rqids file %s %w", rqidsFiles, err)
			}

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				return nil, errors.Errorf("failed to send to begin of file %s %w", file.Name(), err)
			}
			rqids := &RQIDS{
				Id:      base64.StdEncoding.EncodeToString(hasher.Sum(nil)),
				Content: []byte{},
			}

			if rqids.Content, err = ioutil.ReadAll(file); err != nil {
				return nil, errors.Errorf("failed to read content of rqids file %s %w", file.Name(), err)
			}

			return rqids, nil
		}()

		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (task *Task) createTicket() (*pastel.ArtTicket, error) {
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

	// TODO: fill all 0 and "TBD" value with real values when other API ready
	ticket := &pastel.ArtTicket{
		Version:  1,
		Blocknum: 0,
		Author:   pastelID,
		DataHash: task.datahash,
		Copies:   task.Request.IssuedCopies,
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
			NSFWScore:                      task.nsfwScore,
			SeenScore:                      task.seenScore,
			RQIDs:                          task.rqids,
			RQCoti:                         0,
			RQSsoti:                        0,
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
