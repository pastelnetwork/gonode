package userdataprocess

import (
	"bytes"
	"context"
	"encoding/hex"
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
	"github.com/pastelnetwork/gonode/probe"
	rq "github.com/pastelnetwork/gonode/raptorq"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"golang.org/x/crypto/sha3"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	acceptedMu sync.Mutex
	accepted   Nodes


	// signature of ticket data signed by node's pateslID
	ownSignature []byte

	artistSignature []byte

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

func (task *Task) ValidatePreBurnTransaction(ctx context.Context, txid string) (string, error) {
	var err error
	if err = task.RequiredStatus(StatusRegistrationFeeCalculated); err != nil {
		return "", errors.Errorf("require status %s not satisfied", StatusRegistrationFeeCalculated)
	}

	log.WithContext(ctx).Debugf("preburn-txid: %s", txid)
	<-task.NewAction(func(ctx context.Context) error {
		confirmationChn := task.waitConfirmation(ctx, txid, 3, 150*time.Second, 4)

		if err := task.matchFingersPrintAndScores(ctx); err != nil {
			return errors.Errorf("fingerprints or scores don't matched")
		}

		// compare rqsymbols
		if err := task.compareRQSymbolId(ctx); err != nil {
			return errors.Errorf("failed to generate rqids %w", err)
		}

		// sign the ticket if not primary node
		log.WithContext(ctx).Debugf("isPrimary: %d", task.connectedTo == nil)
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

					confirmations := task.waitConfirmation(ctx, artRegTxid, 10, 150*time.Second, 11)
					err = <-confirmations
					if err != nil {
						return errors.Errorf("failed to wait for confirmation of reg-art ticket %w", err)
					}

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

// sign and send art ticket if not primary
func (task *Task) signAndSendArtTicket(ctx context.Context, isPrimary bool) error {
	ticket, err := pastel.EncodeArtTicket(task.Ticket)
	if err != nil {
		return errors.Errorf("failed to serialize art ticket %w", err)
	}
	task.ownSignature, err = task.pastelClient.Sign(ctx, ticket, task.config.PastelID, task.config.PassPhrase)
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
		if ok, err := task.pastelClient.Verify(ctx, data, string(signature), nodeID); err != nil {
			return errors.Errorf("failed to verify signature %s of node %s", signature, nodeID)
		} else {
			if !ok {
				return errors.Errorf("signature of node %s mistmatch", nodeID)
			}
		}
	}
	return nil
}

func (task *Task) registerArt(ctx context.Context) (string, error) {
	log.WithContext(ctx).Debugf("all signature received so start validation")
	// verify signature from secondary nodes
	data, err := pastel.EncodeArtTicket(task.Ticket)
	if err != nil {
		return "", errors.Errorf("failed to encoded art ticket %w", err)
	}

	req := pastel.RegisterArtRequest{
		Ticket: &pastel.ArtTicket{
			Version:   task.Ticket.Version,
			Author:    task.Ticket.Author,
			BlockNum:  task.Ticket.BlockNum,
			BlockHash: task.Ticket.BlockHash,
			Copies:    task.Ticket.Copies,
			Royalty:   task.Ticket.Royalty,
			Green:     task.Ticket.Green,
			AppTicket: data,
		},
		Signatures: &pastel.TicketSignatures{
			Artist: map[string][]byte{
				task.Ticket.AppTicketData.AuthorPastelID: task.artistSignature,
			},
			Mn1: map[string][]byte{
				task.config.PastelID: task.ownSignature,
			},
			Mn2: map[string][]byte{
				task.accepted[0].ID: task.peersArtTicketSignature[task.accepted[0].ID],
			},
			Mn3: map[string][]byte{
				task.accepted[1].ID: task.peersArtTicketSignature[task.accepted[1].ID],
			},
		},
		Mn1PastelId: task.config.PastelID,
		Pasphase:    task.config.PassPhrase,
		Key1:        task.key1,
		Key2:        task.key2,
		Fee:         task.registrationFee,
	}

	artRegTxid, err := task.pastelClient.RegisterArtTicket(ctx, req)
	if err != nil {
		return "", errors.Errorf("failed to register art work %s", err)
	}
	return artRegTxid, nil
}
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
	log.WithContext(ctx).Debugf("master node %s", masterNodes)

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
