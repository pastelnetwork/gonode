package artworkdownload

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload/node"
)

// Task is the task of downloading artwork.
type Task struct {
	task.Task
	*Service

	Ticket *Ticket
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

	ttxid, err := task.validateTicketOwnership(ctx, task.Ticket.TXID, task.Ticket.PastelID, task.Ticket.PastelIDPassphrase)
	if err != nil {
		return err
	}

	// Sign current-timestamp with PsstelID passed in request
	timestamp := time.Now().Format(time.RFC3339)
	signature, err := task.signTimestamp(ctx, timestamp, task.Ticket.PastelID, task.Ticket.PastelIDPassphrase)
	if err != nil {
		return err
	}

	// Retrieve supernodes with highest ranks.
	topNodes, err := task.pastelTopNodes(ctx)
	if err != nil {
		return err
	}
	if len(topNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		return errors.New("Unable to find enough Supernodes")
	}

	var connectedNodes node.List
	var errs error
	// Connect to top supernodes
	for _, node := range topNodes {
		if err := node.Connect(ctx, task.config.connectTimeout); err == nil {
			connectedNodes.Add(node)
		} else {
			errs = errors.Append(errs, err)
		}
	}

	if len(connectedNodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughSuperNode)
		return errors.Errorf("Could not connect enough %d supernodes: %w", task.config.NumberSuperNodes, errs)
	}
	task.UpdateStatus(StatusConnected)

	// Send download request to all top supernodes.
	var nodes node.List
	var downloadErrs error
	errsChan := make(chan error, len(connectedNodes))
	nodesChan := make(chan *node.Node, len(connectedNodes))
	err = connectedNodes.Download(ctx, task.Ticket.TXID, string(timestamp), string(signature), ttxid, nodesChan, errsChan)
	if err != nil {
		task.UpdateStatus(StatusErrorDownloadFailed)
		return err
	}
	close(errsChan)
	close(nodesChan)
	for node := range nodesChan {
		nodes = append(nodes, node)
	}
	for err := range errsChan {
		downloadErrs = errors.Append(errs, err)
	}
	if len(nodes) < task.config.NumberSuperNodes {
		task.UpdateStatus(StatusErrorNotEnoughFiles)
		return errors.Errorf("Not downloading enough files from %d supernodes: %w", task.config.NumberSuperNodes, downloadErrs)
	}
	task.UpdateStatus(StatusDownloaded)

	// Activate supernodes that returned file successfully.
	nodes.Activate()
	// Disconnect supernodes that did not return file.
	topNodes.DisconnectInactive()

	// Cancel context when any connection is broken.
	groupConnClose, _ := errgroup.WithContext(ctx)
	groupConnClose.Go(func() error {
		defer cancel()
		return nodes.WaitConnClose(ctx)
	})

	// Check files are the same
	err = nodes.MatchFiles()
	if err != nil {
		task.UpdateStatus(StatusErrorFilesNotMatch)
		return err
	}

	// TODO: Send file to the caller

	// Wait for all connections to disconnect.
	return groupConnClose.Wait()
}

func (task *Task) signTimestamp(ctx context.Context, timestamp, pastelID, passphrase string) ([]byte, error) {
	return task.pastelClient.Sign(ctx, []byte(timestamp), pastelID, passphrase)
}

func (task *Task) validateTicketOwnership(ctx context.Context, txid, pastelID, passphrase string) (ttxid string, err error) {
	ttxid, err = task.pastelClient.TicketOwnership(ctx, txid, pastelID, passphrase)
	return
}

// func (task *Task) encodeFingerprint(ctx context.Context, fingerprint []byte, image *artwork.File) error {
// 	// Sign fingerprint
// 	ed448PubKey := []byte(task.Ticket.ArtistPastelID)
// 	ed448Signature, err := task.pastelClient.Sign(ctx, fingerprint, task.Ticket.ArtistPastelID, task.Ticket.ArtistPastelIDPassphrase)
// 	if err != nil {
// 		return err
// 	}

// 	// TODO: Should be replaced with real data from the Pastel API.
// 	pqPubKey := []byte("pqPubKey")
// 	pqSignature := []byte("pqSignature")

// 	// Encode data to the image.
// 	encSig := qrsignature.New(
// 		qrsignature.Fingerprint(fingerprint),
// 		qrsignature.PostQuantumSignature(pqSignature),
// 		qrsignature.PostQuantumPubKey(pqPubKey),
// 		qrsignature.Ed448Signature(ed448Signature),
// 		qrsignature.Ed448PubKey(ed448PubKey),
// 	)
// 	if err := image.Encode(encSig); err != nil {
// 		return err
// 	}

// 	// Decode data from the image, to make sure their integrity.
// 	decSig := qrsignature.New()
// 	if err := image.Decode(decSig); err != nil {
// 		return err
// 	}

// 	if !bytes.Equal(fingerprint, decSig.Fingerprint()) {
// 		return errors.Errorf("fingerprints do not match, original len:%d, decoded len:%d\n", len(fingerprint), len(decSig.Fingerprint()))
// 	}
// 	if !bytes.Equal(pqSignature, decSig.PostQuantumSignature()) {
// 		return errors.Errorf("post quantum signatures do not match, original len:%d, decoded len:%d\n", len(pqSignature), len(decSig.PostQuantumSignature()))
// 	}
// 	if !bytes.Equal(pqPubKey, decSig.PostQuantumPubKey()) {
// 		return errors.Errorf("post quantum public keys do not match, original len:%d, decoded len:%d\n", len(pqPubKey), len(decSig.PostQuantumPubKey()))
// 	}
// 	if !bytes.Equal(ed448Signature, decSig.Ed448Signature()) {
// 		return errors.Errorf("ed448 signatures do not match, original len:%d, decoded len:%d\n", len(ed448Signature), len(decSig.Ed448Signature()))
// 	}
// 	if !bytes.Equal(ed448PubKey, decSig.Ed448PubKey()) {
// 		return errors.Errorf("ed448 public keys do not match, original len:%d, decoded len:%d\n", len(ed448PubKey), len(decSig.Ed448PubKey()))
// 	}
// 	return nil
// }

// // meshNodes establishes communication between supernodes.
// func (task *Task) meshNodes(ctx context.Context, nodes node.List, primaryIndex int) (node.List, error) {
// 	var meshNodes node.List

// 	primary := nodes[primaryIndex]
// 	if err := primary.Connect(ctx, task.config.connectTimeout); err != nil {
// 		return nil, err
// 	}
// 	if err := primary.Session(ctx, true); err != nil {
// 		return nil, err
// 	}

// 	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
// 	defer nextConnCancel()

// 	var secondaries node.List
// 	go func() {
// 		for i, node := range nodes {
// 			node := node

// 			if i == primaryIndex {
// 				continue
// 			}

// 			select {
// 			case <-nextConnCtx.Done():
// 				return
// 			case <-time.After(task.config.connectToNextNodeDelay):
// 				go func() {
// 					defer errors.Recover(log.Fatal)

// 					if err := node.Connect(ctx, task.config.connectTimeout); err != nil {
// 						return
// 					}
// 					if err := node.Session(ctx, false); err != nil {
// 						return
// 					}
// 					secondaries.Add(node)

// 					if err := node.ConnectTo(ctx, primary.PastelID(), primary.SessID()); err != nil {
// 						return
// 					}
// 					log.WithContext(ctx).Debugf("Seconary %q connected to primary", node)
// 				}()
// 			}
// 		}
// 	}()

// 	acceptCtx, acceptCancel := context.WithTimeout(ctx, task.config.acceptNodesTimeout)
// 	defer acceptCancel()

// 	accepted, err := primary.AcceptedNodes(acceptCtx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meshNodes.Add(primary)
// 	for _, pastelID := range accepted {
// 		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

// 		node := secondaries.FindByPastelID(pastelID)
// 		if node == nil {
// 			return nil, errors.New("not found accepted node")
// 		}
// 		meshNodes.Add(node)
// 	}
// 	return meshNodes, nil
// }

// func (task *Task) isSuitableStorageFee(ctx context.Context) (bool, error) {
// 	fee, err := task.pastelClient.StorageNetworkFee(ctx)
// 	if err != nil {
// 		return false, err
// 	}
// 	return fee <= task.Ticket.MaximumFee, nil
// }

func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		nodes.Add(node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, ticket *Ticket) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Ticket:  ticket,
	}
}
