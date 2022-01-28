package common

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/node"
	"sync"
)

type NetworkHandler struct {
	task          *SuperNodeTask
	pastelHandler *mixins.PastelHandler

	nodeMaker  node.NodeMaker
	nodeClient node.ClientInterface

	acceptedMu sync.Mutex
	Accepted   SuperNodePeerList

	meshedNodes []types.MeshedSuperNode
	// valid only for secondary node
	ConnectedTo *SuperNodePeer

	pastelID                string
	minNumberConnectedNodes int
}

func NewNetworkHandler(task *SuperNodeTask,
	nodeClient node.ClientInterface, nodeMaker node.NodeMaker,
	pastelClient pastel.Client,
	pastelID string,
	minNumberConnectedNodes int,
) *NetworkHandler {
	return &NetworkHandler{
		task:                    task,
		nodeMaker:               nodeMaker,
		pastelHandler:           mixins.NewPastelHandler(pastelClient),
		nodeClient:              nodeClient,
		pastelID:                pastelID,
		minNumberConnectedNodes: minNumberConnectedNodes,
	}
}

func (h *NetworkHandler) MeshedNodesPastelID() []string {
	var ids []string
	for _, peer := range h.meshedNodes {
		ids = append(ids, peer.NodeID)
	}
	return ids
}

// Session is handshake wallet to supernode
func (h *NetworkHandler) Session(_ context.Context, isPrimary bool) error {
	if err := h.task.RequiredStatus(StatusTaskStarted); err != nil {
		return err
	}

	<-h.task.NewAction(func(ctx context.Context) error {
		if isPrimary {
			log.WithContext(ctx).Debug("Acts as primary node")
			h.task.UpdateStatus(StatusPrimaryMode)
			return nil
		}

		log.WithContext(ctx).Debug("Acts as secondary node")
		h.task.UpdateStatus(StatusSecondaryMode)

		return nil
	})
	return nil
}

// AcceptedNodes waits for connection supernodes, as soon as there is the required amount returns them.
func (h *NetworkHandler) AcceptedNodes(serverCtx context.Context) (SuperNodePeerList, error) {
	if err := h.task.RequiredStatus(StatusPrimaryMode); err != nil {
		return nil, err
	}

	<-h.task.NewAction(func(ctx context.Context) error {
		log.WithContext(ctx).Debug("Waiting for supernodes to connect")

		sub := h.task.SubscribeStatus()
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
	return h.Accepted, nil
}

// SessionNode accepts secondary node
func (h *NetworkHandler) SessionNode(_ context.Context, nodeID string) error {
	h.acceptedMu.Lock()
	defer h.acceptedMu.Unlock()

	if err := h.task.RequiredStatus(StatusPrimaryMode); err != nil {
		return err
	}

	var err error

	<-h.task.NewAction(func(ctx context.Context) error {
		if node := h.Accepted.ByID(nodeID); node != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).Errorf("node is already registered")
			err = errors.Errorf("node %q is already registered", nodeID)
			return nil
		}

		var someNode *SuperNodePeer
		someNode, err = h.PastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("get node by extID")
			err = errors.Errorf("get node by extID %s: %w", nodeID, err)
			return nil
		}
		h.Accepted.Add(someNode)

		log.WithContext(ctx).WithField("nodeID", nodeID).Debug("Accept secondary node")

		if len(h.Accepted) >= h.minNumberConnectedNodes {
			h.task.UpdateStatus(StatusConnected)
		}
		return nil
	})
	return err
}

// ConnectTo connects to primary node
func (h *NetworkHandler) ConnectTo(_ context.Context, nodeID, sessID string) error {
	if err := h.task.RequiredStatus(StatusSecondaryMode); err != nil {
		return err
	}

	var err error

	<-h.task.NewAction(func(ctx context.Context) error {
		var someNode *SuperNodePeer
		someNode, err = h.PastelNodeByExtKey(ctx, nodeID)
		if err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("get node by extID")
			return nil
		}

		if err = someNode.Connect(ctx); err != nil {
			log.WithContext(ctx).WithField("nodeID", nodeID).WithError(err).Errorf("connect to node")
			return nil
		}

		if err = someNode.Session(ctx, h.pastelID, sessID); err != nil {
			log.WithContext(ctx).WithField("sessID", sessID).WithField("pastelID", h.pastelID).WithError(err).Errorf("handsake with peer")
			return nil
		}

		h.ConnectedTo = someNode
		h.task.UpdateStatus(StatusConnected)
		return nil
	})
	return err
}

// MeshNodes to set info of all meshed supernodes - that will be to send
func (h *NetworkHandler) MeshNodes(_ context.Context, meshedNodes []types.MeshedSuperNode) error {
	if err := h.task.RequiredStatus(StatusConnected); err != nil {
		return err
	}
	h.meshedNodes = meshedNodes

	return nil
}

func (h *NetworkHandler) CheckNodeInMeshedNodes(nodeID string) error {
	if h.meshedNodes == nil {
		return errors.New("nil meshedNodes")
	}

	for _, node := range h.meshedNodes {
		if node.NodeID == nodeID {
			return nil
		}
	}

	return errors.New("nodeID not found")
}

func (h *NetworkHandler) PastelNodeByExtKey(ctx context.Context, nodeID string) (*SuperNodePeer, error) {
	masterNodes, err := h.pastelHandler.PastelClient.MasterNodesTop(ctx)
	log.WithContext(ctx).Debugf("master node %v", masterNodes)

	if err != nil {
		return nil, err
	}

	for _, masterNode := range masterNodes {
		if masterNode.ExtKey != nodeID {
			continue
		}
		someNode := NewSuperNode(h.nodeClient, masterNode.ExtKey, masterNode.ExtAddress, h.nodeMaker)
		return someNode, nil
	}

	return nil, errors.Errorf("node %q not found", nodeID)
}
