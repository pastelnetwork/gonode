package cascaderegister

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	nftRegMetadata *types.ActionRegMetadata
	Ticket         *pastel.ActionTicket
	Artwork        *artwork.File

	meshedNodes []types.MeshedSuperNode

	// DDAndFingerprints created by node its self
	myDDAndFingerprints                   *pastel.DDAndFingerprints
	calculatedDDAndFingerprints           *pastel.DDAndFingerprints
	allDDAndFingerprints                  map[string]*pastel.DDAndFingerprints
	allSignedDDAndFingerprintsReceivedChn chan struct{}
	ddMtx                                 sync.Mutex

	acceptedMu sync.Mutex
	accepted   Nodes

	// signature of ticket data signed by node's pateslID
	ownSignature []byte

	creatorSignature []byte
	registrationFee  int64

	// valid only for a task run as primary
	peersArtTicketSignatureMtx *sync.Mutex
	peersArtTicketSignature    map[string][]byte
	allSignaturesReceivedChn   chan struct{}

	// valid only for secondary node
	connectedTo *Node

	ddFpFiles [][]byte
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

// MeshNodes to set info of all meshed supernodes - that will be to send
func (task *Task) MeshNodes(_ context.Context, meshedNodes []types.MeshedSuperNode) error {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		return err
	}
	task.meshedNodes = meshedNodes

	return nil
}

// SendRegMetadata sends reg metadata
func (task *Task) SendRegMetadata(_ context.Context, regMetadata *types.ActionRegMetadata) error {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		return err
	}
	task.nftRegMetadata = regMetadata

	return nil
}

func (task *Task) compressSignedDDAndFingerprints(ctx context.Context, ddData *pastel.DDAndFingerprints) ([]byte, error) {
	// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
	ddDataBytes, err := json.Marshal(ddData)
	if err != nil {
		return nil, errors.Errorf("marshal DDAndFingerprints: %w", err)
	}

	// sign it
	signature, err := task.pastelClient.Sign(ctx, ddDataBytes, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, errors.Errorf("sign DDAndFingerprints: %w", err)
	}

	compressed, err := pastel.ToCompressSignedDDAndFingerprints(ddData, signature)
	if err != nil {
		return nil, errors.Errorf("compress SignedDDAndFingerprints: %w", err)
	}

	return compressed, nil
}

// ProbeImage uploads the resampled image compute and return a compression of pastel.DDAndFingerprints
func (task *Task) ProbeImage(ctx context.Context, file *artwork.File) ([]byte, error) {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		log.WithContext(ctx).WithField("CurrentStatus", task.Status()).Error("ProbeImage() failed with non-satisfied task status")
		return nil, err
	}

	var err error

	<-task.NewAction(func(ctx context.Context) error {
		defer errors.Recover(func(recErr error) {
			log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).WithError(recErr).Error("PanicWhenProbeImage")
		})
		task.UpdateStatus(StatusImageProbed)
		task.Artwork = file
		task.myDDAndFingerprints, err = task.genFingerprintsData(ctx, task.Artwork)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("generate fingerprints data")
			err = errors.Errorf("generate fingerprints data: %w", err)
			return nil
		}

		// send signed DDAndFingerprints to other SNs
		if task.meshedNodes == nil || len(task.meshedNodes) != 3 {
			log.WithContext(ctx).Error("Not enough meshed SuperNodes")
			err = errors.New("not enough meshed SuperNodes")
			return nil
		}

		for _, nodeInfo := range task.meshedNodes {
			// Don't send to itself
			if nodeInfo.NodeID == task.config.PastelID {
				continue
			}

			var node *Node
			node, err = task.pastelNodeByExtKey(ctx, nodeInfo.NodeID)
			if err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("get node by extID")
				return nil
			}

			if err = node.connect(ctx); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("connect to node")
				return nil
			}

			if err = node.SendSignedDDAndFingerprints(ctx, nodeInfo.SessID, task.config.PastelID, task.myDDAndFingerprints.ZstdCompressedFingerprint); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"nodeID":  nodeInfo.NodeID,
					"sessID":  nodeInfo.SessID,
					"context": "SendDDAndFingerprints",
				}).WithError(err).Errorf("send signed DDAndFingerprints failed")
				return nil
			}
		}

		// wait for other SNs shared their signed DDAndFingerprints
		select {
		case <-ctx.Done():
			err = ctx.Err()
			log.WithContext(ctx).WithError(err).Error("ctx.Done()")
			if err != nil {
				log.WithContext(ctx).Error("waiting for DDAndFingerprints from peers context cancelled")
			}
			return nil
		case <-task.allSignedDDAndFingerprintsReceivedChn:
			log.WithContext(ctx).Debug("all DDAndFingerprints received so start calculate final DDAndFingerprints")
			task.allDDAndFingerprints[task.config.PastelID] = task.myDDAndFingerprints

			// get list of DDAndFingerprints in order of node rank
			dDAndFingerprintsList := []*pastel.DDAndFingerprints{}
			for _, node := range task.meshedNodes {
				v, ok := task.allDDAndFingerprints[node.NodeID]
				if !ok {
					err = errors.Errorf("not found DDAndFingerprints of node : %s", node.NodeID)
					log.WithContext(ctx).WithFields(log.Fields{
						"nodeID": node.NodeID,
					}).Errorf("DDAndFingerprints of node not found")
					return nil
				}
				dDAndFingerprintsList = append(dDAndFingerprintsList, v)
			}

			// calculate final result from DdDAndFingerprints from all SNs node
			if len(dDAndFingerprintsList) != 3 {
				err = errors.Errorf("not enough DDAndFingerprints, len: %d", len(dDAndFingerprintsList))
				log.WithContext(ctx).WithFields(log.Fields{
					"list": dDAndFingerprintsList,
				}).Errorf("not enough DDAndFingerprints")
				return nil
			}

			task.calculatedDDAndFingerprints, err = pastel.CombineFingerPrintAndScores(
				dDAndFingerprintsList[0],
				dDAndFingerprintsList[1],
				dDAndFingerprintsList[2],
			)

			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("call CombineFingerPrintAndScores() failed")
				return nil
			}

			// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
			var compressed []byte
			compressed, err = task.compressSignedDDAndFingerprints(ctx, task.calculatedDDAndFingerprints)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("compress combine DDAndFingerPrintAndScore failed")

				err = errors.Errorf("compress combine DDAndFingerPrintAndScore failed: %w", err)
				return nil
			}
			task.calculatedDDAndFingerprints.ZstdCompressedFingerprint = compressed
			return nil
		case <-time.After(30 * time.Second):
			log.WithContext(ctx).Error("waiting for DDAndFingerprints from peers timeout")
			err = errors.New("waiting for DDAndFingerprints timeout")
			return nil
		}
	})

	if err != nil {
		return nil, err
	}

	return task.calculatedDDAndFingerprints.ZstdCompressedFingerprint, nil
}

func (task *Task) validateDdFpIds(ctx context.Context, dd []byte) error {
	// FIXME: add Cascade sequences
	// validate := func(ctx context.Context, data []byte, ic uint32, max uint32, ids []string) error {
	// 	length := 4

	// 	dec, err := utils.B64Decode(data)
	// 	if err != nil {
	// 		return errors.Errorf("decode data: %w", err)
	// 	}

	// 	decData, err := zstd.Decompress(nil, dec)
	// 	if err != nil {
	// 		return errors.Errorf("decompress: %w", err)
	// 	}

	// 	splits := bytes.Split(decData, []byte{pastel.SeparatorByte})
	// 	if len(splits) != length {
	// 		return errors.New("invalid data")
	// 	}

	// 	file, err := utils.B64Decode(splits[0])
	// 	if err != nil {
	// 		errors.Errorf("decode file: %w", err)
	// 	}

	// 	verifications := 0
	// 	verifiedNodes := make(map[int]bool)
	// 	for i := 1; i < length; i++ {
	// 		for j := 0; j < len(task.meshedNodes); j++ {
	// 			if _, ok := verifiedNodes[j]; ok {
	// 				continue
	// 			}

	// 			verified, err := task.pastelClient.Verify(ctx, file, string(splits[i]), task.meshedNodes[j].NodeID, pastel.SignAlgorithmED448)
	// 			if err != nil {
	// 				return errors.Errorf("verify file signature %w", err)
	// 			}

	// 			if verified {
	// 				verifiedNodes[j] = true
	// 				verifications++
	// 				break
	// 			}
	// 		}
	// 	}

	// 	if verifications != 3 {
	// 		return errors.Errorf("file verification failed: need %d verifications, got %d", 3, verifications)
	// 	}

	// 	gotIDs, idFiles, err := pastel.GetIDFiles(decData, ic, max)
	// 	if err != nil {
	// 		return errors.Errorf("get ids: %w", err)
	// 	}

	// 	task.ddFpFiles = idFiles

	// 	if err := utils.EqualStrList(gotIDs, ids); err != nil {
	// 		return errors.Errorf("IDs don't match: %w", err)
	// 	}

	// 	return nil
	// }

	// apiCascadeTicket, err := task.Ticket.APICascadeTicket()
	// if err != nil {
	// 	return errors.Errorf("invalid Cascade ticket: %w", err)
	// }

	// if err := validate(ctx, dd, apiCascadeTicket.DDAndFingerprintsIc, apiCascadeTicket.DDAndFingerprintsMax,
	// 	apiCascadeTicket.DDAndFingerprintsIDs); err != nil {

	// 	return errors.Errorf("validate dd_and_fingerprints: %w", err)
	// }

	return nil
}

// GetRegistrationFee get the fee to register artwork to bockchain
func (task *Task) validateSignedTicketFromWN(ctx context.Context, ticket []byte, creatorSignature []byte, ddFpFile []byte) error {
	var err error
	task.creatorSignature = creatorSignature

	// TODO: fix this like how can we get the signature before calling cNode
	task.Ticket, err = pastel.DecodeActionTicket(ticket)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("decode action ticket")
		return errors.Errorf("decode action ticket: %w", err)
	}

	// Verify APICascadeTicket
	_, err = task.Ticket.APICascadeTicket()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("invalid api Cascade ticket")
		return errors.Errorf("invalid api Cascade ticket: %w", err)
	}

	verified, err := task.pastelClient.Verify(ctx, ticket, string(creatorSignature), task.Ticket.Caller, pastel.SignAlgorithmED448)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("verify ticket signature")
		return errors.Errorf("verify ticket signature %w", err)
	}

	if !verified {
		err = errors.New("ticket verification failed")
		log.WithContext(ctx).WithError(err).Errorf("verification failure")
		return err
	}

	if err := task.validateDdFpIds(ctx, ddFpFile); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("validate rq & dd id files")

		return errors.Errorf("validate rq & dd id files %w", err)
	}
	return nil
}

// ValidateBurnTxID - will validate the pre-burnt transaction ID created by 3rd party
func (task *Task) ValidateBurnTxID(_ context.Context) error {
	var err error

	<-task.NewAction(func(ctx context.Context) error {
		confirmationChn := task.waitConfirmation(ctx, task.nftRegMetadata.BurnTxID, int64(task.config.PreburntTxMinConfirmations), 15*time.Second)

		log.WithContext(ctx).Debug("waiting for confimation")
		if err = <-confirmationChn; err != nil {
			task.UpdateStatus(StatusErrorInvalidBurnTxID)
			log.WithContext(ctx).WithError(err).Errorf("validate preburn transaction validation")
			err = errors.Errorf("validate preburn transaction validation :%w", err)
			return err
		}
		log.WithContext(ctx).Debug("confirmation done")
		return nil
	})

	return err
}

// ValidateAndRegister will get signed ticket from fee txid, wait until it's confirmations meet expectation.
func (task *Task) ValidateAndRegister(_ context.Context, ticket []byte, creatorSignature []byte, ddFpFile []byte) (string, error) {
	var err error
	if err = task.RequiredStatus(StatusImageProbed); err != nil {
		return "", errors.Errorf("require status %s not satisfied", StatusImageProbed)
	}

	task.creatorSignature = creatorSignature

	<-task.NewAction(func(ctx context.Context) error {
		if err = task.validateSignedTicketFromWN(ctx, ticket, creatorSignature, ddFpFile); err != nil {
			return nil
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %t", task.connectedTo == nil)
		if err = task.signAndSendArtTicket(ctx, task.connectedTo == nil); err != nil {
			log.WithContext(ctx).WithError(err).Errorf("signed and send NFT ticket")
			err = errors.Errorf("signed and send NFT ticket")
			return nil
		}

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
				case <-task.allSignaturesReceivedChn:
					log.WithContext(ctx).Debug("all signature received so start validation")

					if err = task.verifyPeersSingature(ctx); err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("peers' singature mismatched: %w", err)
						return nil
					}

					nftRegTxid, err = task.registerAction(ctx)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("peers' singature mismatched")
						err = errors.Errorf("register NFT: %w", err)
						return nil
					}

					return nil
				}
			}
		})
	}

	return nftRegTxid, err
}

// ValidateActionActAndStore informs actionRegTxID to trigger store IDfiles in case of actionRegTxID was
func (task *Task) ValidateActionActAndStore(ctx context.Context, actionRegTxID string) error {
	var err error

	// Wait for action ticket to be activated by walletnode
	confirmations := task.waitActionActivation(ctx, actionRegTxID, 2, 30*time.Second)
	err = <-confirmations
	if err != nil {
		return errors.Errorf("wait for confirmation of reg-art ticket %w", err)
	}

	// Store dd_and_fingerprints into Kademlia
	if err = task.storeIDFiles(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Errorf("store id files")
		err = errors.Errorf("store id files: %w", err)
		return nil
	}

	return nil
}

func (task *Task) waitActionActivation(ctx context.Context, txid string, timeoutInBlock int64, interval time.Duration) <-chan error {
	ch := make(chan error)

	go func(ctx context.Context, txid string) {
		defer close(ch)
		blockTracker := blocktracker.New(task.pastelClient)
		baseBlkCnt, err := blockTracker.GetBlockCount()
		if err != nil {
			log.WithContext(ctx).WithError(err).Warn("failed to get block count")
			ch <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				// context cancelled or abort by caller so no need to return anything
				log.WithContext(ctx).Debugf("context done: %s", ctx.Err())
				ch <- ctx.Err()
				return
			case <-time.After(interval):
				txResult, err := task.pastelClient.FindActionActByActionRegTxid(ctx, txid)
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("FindActionActByActionRegTxid err")
				} else {
					if txResult != nil {
						log.WithContext(ctx).Debug("action reg is activated")
						ch <- nil
						return
					}
				}

				currentBlkCnt, err := blockTracker.GetBlockCount()
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("failed to get block count")
					continue
				}

				if currentBlkCnt-baseBlkCnt >= int32(timeoutInBlock)+2 {
					ch <- errors.Errorf("timeout when wating for confirmation of transaction %s", txid)
					return
				}
			}

		}
	}(ctx, txid)
	return ch
}

func (task *Task) waitConfirmation(ctx context.Context, txid string, minConfirmation int64, interval time.Duration) <-chan error {
	ch := make(chan error)

	go func(ctx context.Context, txid string) {
		defer close(ch)
		blockTracker := blocktracker.New(task.pastelClient)
		baseBlkCnt, err := blockTracker.GetBlockCount()
		if err != nil {
			log.WithContext(ctx).WithError(err).Warn("failed to get block count")
			ch <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				// context cancelled or abort by caller so no need to return anything
				log.WithContext(ctx).Debugf("context done: %s", ctx.Err())
				ch <- ctx.Err()
				return
			case <-time.After(interval):
				txResult, err := task.pastelClient.GetRawTransactionVerbose1(ctx, txid)
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("GetRawTransactionVerbose1 err")
				} else {
					if txResult.Confirmations >= minConfirmation {
						log.WithContext(ctx).Debug("transaction confirmed")
						ch <- nil
						return
					}
				}

				currentBlkCnt, err := blockTracker.GetBlockCount()
				if err != nil {
					log.WithContext(ctx).WithError(err).Warn("failed to get block count")
					continue
				}

				if currentBlkCnt-baseBlkCnt >= int32(minConfirmation)+2 {
					ch <- errors.Errorf("timeout when wating for confirmation of transaction %s", txid)
					return
				}
			}

		}
	}(ctx, txid)
	return ch
}

// sign and send NFT ticket if not primary
func (task *Task) signAndSendArtTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeActionTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("serialize NFT ticket: %w", err)
	}
	task.ownSignature, err = task.pastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return errors.Errorf("sign ticket: %w", err)
	}
	if !isPrimary {
		log.WithContext(ctx).Debug("send signed articket to primary node")
		if err := task.connectedTo.RegisterCascade.SendArtTicketSignature(ctx, task.config.PastelID, task.ownSignature); err != nil {
			return errors.Errorf("send signature to primary node %s at address %s: %w", task.connectedTo.ID, task.connectedTo.Address, err)
		}
	}
	return nil
}

func (task *Task) verifyPeersSingature(ctx context.Context) error {
	log.WithContext(ctx).Debug("all signature received so start validation")

	data, err := pastel.EncodeActionTicket(task.Ticket)
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

func (task *Task) registerAction(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debug("all signature received so start validation")

	req := pastel.RegisterActionRequest{
		Ticket: &pastel.ActionTicket{
			Version:       task.Ticket.Version,
			Caller:        task.Ticket.Caller,
			BlockNum:      task.Ticket.BlockNum,
			BlockHash:     task.Ticket.BlockHash,
			ActionType:    task.Ticket.ActionType,
			APITicketData: task.Ticket.APITicketData,
		},
		Signatures: &pastel.ActionTicketSignatures{
			Caller: map[string]string{
				task.Ticket.Caller: string(task.creatorSignature),
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
		Passphrase:  task.config.PassPhrase,
		Fee:         task.registrationFee,
	}

	nftRegTxid, err := task.pastelClient.RegisterActionTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("register action ticket: %w", err)
	}
	return nftRegTxid, nil
}

func (task *Task) storeIDFiles(ctx context.Context) error {
	store := func(ctx context.Context, files [][]byte) error {
		for _, idFile := range files {
			if _, err := task.p2pClient.Store(ctx, idFile); err != nil {
				return err
			}
		}

		return nil
	}

	if err := store(ctx, task.ddFpFiles); err != nil {
		return fmt.Errorf("store ddAndFp files: %w", err)
	}

	return nil
}

func (task *Task) genFingerprintsData(ctx context.Context, file *artwork.File) (*pastel.DDAndFingerprints, error) {
	img, err := file.Bytes()
	if err != nil {
		return nil, errors.Errorf("get content of image %s: %w", file.Name(), err)
	}

	if task.nftRegMetadata == nil || task.nftRegMetadata.BlockHash == "" || task.nftRegMetadata.CreatorPastelID == "" {
		return nil, errors.Errorf("invalid nftRegMetadata")
	}

	// Get DDAndFingerprints
	ddAndFingerprints, err := task.ddClient.ImageRarenessScore(
		ctx,
		img,
		file.Format().String(),
		task.nftRegMetadata.BlockHash,
		task.nftRegMetadata.CreatorPastelID,
	)

	if err != nil {
		return nil, errors.Errorf("call ImageRarenessScore(): %w", err)
	}

	// Creates compress(Base64(dd_and_fingerprints).Base64(signature))
	compressed, err := task.compressSignedDDAndFingerprints(ctx, ddAndFingerprints)
	if err != nil {
		return nil, errors.Errorf("call compressSignedDDAndFingerprints failed: %w", err)
	}

	ddAndFingerprints.ZstdCompressedFingerprint = compressed
	return ddAndFingerprints, nil
}

func (task *Task) checkNodeInMeshedNodes(nodeID string) error {
	if task.meshedNodes == nil {
		return errors.New("nil meshedNodes")
	}

	for _, node := range task.meshedNodes {
		if node.NodeID == nodeID {
			return nil
		}
	}

	return errors.New("nodeID not found")
}

// AddSignedDDAndFingerprints adds signed dd and fp
func (task *Task) AddSignedDDAndFingerprints(nodeID string, compressedSignedDDAndFingerprints []byte) error {
	task.ddMtx.Lock()
	defer task.ddMtx.Unlock()

	var err error

	<-task.NewAction(func(ctx context.Context) error {
		defer errors.Recover(func(recErr error) {
			log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).WithError(recErr).Error("PanicWhenProbeImage")
		})
		log.WithContext(ctx).Debugf("receive compressedSignedDDAndFingerprints from node %s", nodeID)
		// Check nodeID should in task.meshedNodes
		err = task.checkNodeInMeshedNodes(nodeID)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("check if node existed in meshedNodes failed")
			return nil
		}

		// extract data
		var ddAndFingerprints *pastel.DDAndFingerprints
		var ddAndFingerprintsBytes []byte
		var signature []byte

		ddAndFingerprints, ddAndFingerprintsBytes, signature, err = pastel.ExtractCompressSignedDDAndFingerprints(compressedSignedDDAndFingerprints)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("extract compressedSignedDDAndFingerprints failed")
			return nil
		}

		var ok bool
		ok, err = task.pastelClient.Verify(ctx, ddAndFingerprintsBytes, string(signature), nodeID, pastel.SignAlgorithmED448)
		if err != nil || !ok {
			err = errors.New("signature verification failed")
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("verify signature of ddAndFingerprintsBytes failed")
			return nil
		}

		task.allDDAndFingerprints[nodeID] = ddAndFingerprints

		// if enough ddAndFingerprints received (not include from itself)
		if len(task.allDDAndFingerprints) == 2 {
			log.WithContext(ctx).Debug("all ddAndFingerprints received")
			go func() {
				close(task.allSignedDDAndFingerprintsReceivedChn)
			}()
		}
		return nil
	})

	return err
}

// AddPeerArticketSignature is called by supernode peers to primary first-rank supernode to send it's signature
func (task *Task) AddPeerArticketSignature(nodeID string, signature []byte) error {
	task.peersArtTicketSignatureMtx.Lock()
	defer task.peersArtTicketSignatureMtx.Unlock()

	if err := task.RequiredStatus(StatusImageProbed); err != nil {
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
				close(task.allSignaturesReceivedChn)
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

	removeFn(task.Artwork)
}

// NewTask returns a new Task instance.
func NewTask(service *Service) *Task {
	return &Task{
		Task:                                  task.New(StatusTaskStarted),
		Service:                               service,
		peersArtTicketSignatureMtx:            &sync.Mutex{},
		peersArtTicketSignature:               make(map[string][]byte),
		allSignaturesReceivedChn:              make(chan struct{}),
		allDDAndFingerprints:                  map[string]*pastel.DDAndFingerprints{},
		allSignedDDAndFingerprintsReceivedChn: make(chan struct{}),
	}
}
