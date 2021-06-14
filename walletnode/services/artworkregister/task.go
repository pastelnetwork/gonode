package artworkregister

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/image/qrsignature"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister/node"
)

// Task is the task of registering new artwork.
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

	if ok, err := task.isSuitableStorageFee(ctx); err != nil {
		return err
	} else if !ok {
		task.UpdateStatus(ErrorInsufficientFee)
		return errors.Errorf("network storage fee is higher than specified in the ticket: %v", task.Ticket.MaximumFee)
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

	// Create a thumbnail copy of the image.
	log.WithContext(ctx).WithField("filename", task.Ticket.Image.Name()).Debugf("Copy image")
	thumbnail, err := task.Ticket.Image.Copy()
	if err != nil {
		return err
	}

	log.WithContext(ctx).WithField("filename", thumbnail.Name()).Debugf("Resize image to %dx%d", task.config.thumbnailWidth, task.config.thumbnailHeight)
	if err := thumbnail.ResizeImage(thumbnailWidth, thumbnailHeight); err != nil {
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

	// Sign fingerprint
	fingerprint := nodes.Fingerprint()
	ed448Signature, err := task.pastelClient.Sign(ctx, fingerprint, task.Ticket.ArtistPastelID, task.Ticket.ArtistPastelIDPassphrase)
	if err != nil {
		return err
	}

	qrSig := qrsignature.New(
		qrsignature.PostQuantumPubKey(pqPubKey),
		qrsignature.Fingerprint(fingerprint),
		qrsignature.Fingerprint(finger),
		qrsignature.Ed448Signature(ed448Signature),
		qrsignature.Ed448PubKey([]byte(task.Ticket.ArtistPastelID)),
	)

	// var qrSig qrsignature.Signature
	// for i := 0; i < 15; i++ {
	// 	length := rand.Intn(10000)
	// 	if length == 0 {
	// 		continue
	// 	}

	// 	str, _ := random.String(length, random.Base62Chars)
	// 	qrSig = append(qrSig, qrsignature.Fingerprint([]byte(str)))
	// }

	if err := task.Ticket.Image.Encode(qrSig); err != nil {
		return err
	}

	fmt.Println(string(ed448Signature))

	// Wait for all connections to disconnect.
	return groupConnClose.Wait()
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
	return fee <= task.Ticket.MaximumFee, nil
}

func (task *Task) pastelTopNodes(ctx context.Context) (node.List, error) {
	var nodes node.List

	mns, err := task.pastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > task.Ticket.MaximumFee {
			continue
		}
		nodes = append(nodes, node.NewNode(task.Service.nodeClient, mn.ExtAddress, mn.ExtKey))
	}

	return nodes, nil
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Ticket) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
		Ticket:  Ticket,
	}
}
