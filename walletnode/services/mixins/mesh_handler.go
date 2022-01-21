package mixins

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"sync"
	"time"
)

type MeshHandler struct {
	task          *common.WalletNodeTask
	pastelHandler *PastelHandler

	nodeMaker  node.NodeMaker
	nodeClient node.ClientInterface

	// TODO: make config interface and use it here instead of individual items
	minNumberSuperNodes    int
	connectToNodeTimeout   time.Duration
	acceptNodesTimeout     time.Duration
	connectToNextNodeDelay time.Duration

	callersPastelID string
	passphrase      string

	Nodes common.SuperNodeList
}

func NewMeshHandler(task *common.WalletNodeTask,
	nodeClient node.ClientInterface, nodeMaker node.NodeMaker,
	pastelHandler *PastelHandler,
	pastelID string, passphrase string,
	minSNs int, connectToNodeTimeout time.Duration, acceptNodesTimeout time.Duration, connectToNextNodeDelay time.Duration,
) *MeshHandler {
	return &MeshHandler{
		task:                   task,
		nodeMaker:              nodeMaker,
		pastelHandler:          pastelHandler,
		nodeClient:             nodeClient,
		callersPastelID:        pastelID,
		passphrase:             passphrase,
		minNumberSuperNodes:    minSNs,
		connectToNodeTimeout:   connectToNodeTimeout,
		acceptNodesTimeout:     acceptNodesTimeout,
		connectToNextNodeDelay: connectToNextNodeDelay,
	}
}

// ConnectToTopRankNodes - find 3 supernodes and create mesh of 3 nodes
func (m *MeshHandler) ConnectToTopRankNodes(ctx context.Context) (int, string, error) {

	// Retrieve supernodes with the highest ranks.
	var tmpNodes common.SuperNodeList
	candidatesNodes, err := m.getTopNodes(ctx, 0)
	if err != nil {
		return -1, "", errors.Errorf("call masternode top: %w", err)
	}

	if len(candidatesNodes) < m.minNumberSuperNodes {
		m.task.UpdateStatus(common.StatusErrorNotEnoughSuperNode)
		return -1, "", errors.New("unable to find enough Supernodes with acceptable storage fee")
	}

	// Get current block height & hash
	blockNum, blockHash, err := m.pastelHandler.GetBlock(ctx)
	if err != nil {
		return -1, "", errors.Errorf("get block: %v", err)
	}

	// Connect to top nodes to find 3SN and validate their infor
	err = m.validateMNsInfo(ctx)
	if err != nil {
		m.task.UpdateStatus(common.StatusErrorFindRespondingSNs)
		return -1, "", errors.Errorf("validate MNs info: %v", err)
	}

	/* Establish mesh with 3 SNs */
	var meshedNodes common.SuperNodeList
	var errs error

	for primaryRank := range tmpNodes {
		meshedNodes, err = m.findNodesForMesh(ctx, tmpNodes, primaryRank)
		if err != nil {
			// close connected connections
			tmpNodes.DisconnectAll()

			if errors.IsContextCanceled(err) {
				return -1, "", err
			}
			errs = errors.Append(errs, err)
			log.WithContext(ctx).WithError(err).Warn("Could not create a mesh of the nodes")
			continue
		}
		break
	}
	if len(meshedNodes) < m.minNumberSuperNodes {
		// close connected connections
		meshedNodes.DisconnectAll()
		return -1, "", errors.Errorf("Could not create a mesh of %d nodes: %w", m.minNumberSuperNodes, errs)
	}

	// Activate supernodes that are in the mesh.
	meshedNodes.Activate()

	// Cancel context when any connection is broken.
	m.task.UpdateStatus(common.StatusConnected)
	m.Nodes = meshedNodes

	// Send all meshed supernode info to nodes - that will be used to node send info to other nodes
	var meshedSNInfo []types.MeshedSuperNode
	for _, someNode := range meshedNodes {
		meshedSNInfo = append(meshedSNInfo, types.MeshedSuperNode{
			NodeID: someNode.PastelID(),
			SessID: someNode.SessID(),
		})
	}

	for _, someNode := range meshedNodes {
		err = someNode.MeshNodes(ctx, meshedSNInfo)
		if err != nil {
			meshedNodes.DisconnectAll()
			return -1, "", errors.Errorf("could not send info of meshed nodes: %w", err)
		}
	}
	return blockNum, blockHash, nil
}

// GetTopNodes get list of current top masternodes and validate them (in the order received from pasteld)
func (m *MeshHandler) getTopNodes(ctx context.Context, maximumFee float64) (common.SuperNodeList, error) {

	var nodes common.SuperNodeList

	mns, err := m.pastelHandler.PastelClient.MasterNodesTop(ctx)
	if err != nil {
		return nil, err
	}
	for _, mn := range mns {
		if mn.Fee > maximumFee {
			continue
		}
		if mn.ExtKey == "" || mn.ExtAddress == "" {
			continue
		}

		// Ensures that the PastelId(mn.ExtKey) of MN node is registered
		_, err = m.pastelHandler.PastelClient.FindTicketByID(ctx, mn.ExtKey)
		if err != nil {
			log.WithContext(ctx).WithField("mn", mn).Warn("FindTicketByID() failed")
			continue
		}
		nodes.AddNewNode(m.nodeClient, mn.ExtAddress, mn.ExtKey, m.nodeMaker)
	}
	return nodes, nil
}

// validateMNsInfo - validate MNs info, until found at least 3 valid MNs
func (m *MeshHandler) validateMNsInfo(ctx context.Context) error {
	count := 0

	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  "ed448",
	}

	var nodes common.SuperNodeList
	for _, someNode := range m.Nodes {
		if err := someNode.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
			continue
		}

		count++
		nodes = append(nodes, someNode)
		if count == m.minNumberSuperNodes {
			break
		}
	}

	// Close all connected connections
	nodes.DisconnectAll()

	if count < m.minNumberSuperNodes {
		return errors.Errorf("validate %d Supernodes from pastel network", m.minNumberSuperNodes)
	}

	m.Nodes = nodes
	return nil
}

// meshNodes establishes communication between supernodes.
func (m *MeshHandler) findNodesForMesh(ctx context.Context, candidatesNodes common.SuperNodeList, primaryIndex int) (common.SuperNodeList, error) {

	secInfo := &alts.SecInfo{
		PastelID:   m.callersPastelID,
		PassPhrase: m.passphrase,
		Algorithm:  "ed448",
	}

	primary := candidatesNodes[primaryIndex]
	log.WithContext(ctx).Debugf("Trying to connect to primary node %q", primary)
	if err := primary.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
		return nil, err
	}
	if err := primary.Session(ctx, true); err != nil {
		return nil, err
	}

	nextConnCtx, nextConnCancel := context.WithCancel(ctx)
	defer nextConnCancel()

	// FIXME: ugly hack here. Need to make the Node and List to be safer
	secondariesMtx := &sync.Mutex{}
	var secondaries common.SuperNodeList
	go func() {
		for i, someNode := range candidatesNodes {
			sameNode := someNode

			if i == primaryIndex {
				continue
			}

			select {
			case <-nextConnCtx.Done():
				return
			case <-time.After(m.connectToNextNodeDelay):
				go func() {
					defer errors.Recover(log.Fatal)

					if err := sameNode.Connect(ctx, m.connectToNodeTimeout, secInfo); err != nil {
						return
					}
					if err := sameNode.Session(ctx, false); err != nil {
						return
					}
					// Should not run this code in go routine
					func() {
						secondariesMtx.Lock()
						defer secondariesMtx.Unlock()
						secondaries.Add(sameNode)
					}()

					if err := sameNode.ConnectTo(ctx, types.MeshedSuperNode{
						NodeID: primary.PastelID(),
						SessID: primary.SessID(),
					}); err != nil {
						return
					}
					log.WithContext(ctx).Debugf("Secondary %q connected to primary", sameNode)
				}()
			}
		}
	}()

	acceptCtx, acceptCancel := context.WithTimeout(ctx, m.acceptNodesTimeout)
	defer acceptCancel()

	accepted, err := primary.AcceptedNodes(acceptCtx)
	if err != nil {
		return nil, err
	}

	var meshedNodes common.SuperNodeList
	primary.SetPrimary(true)
	meshedNodes.Add(primary)

	secondariesMtx.Lock()
	defer secondariesMtx.Unlock()

	for _, pastelID := range accepted {
		log.WithContext(ctx).Debugf("Primary accepted %q secondary node", pastelID)

		someNode := secondaries.FindByPastelID(pastelID)
		if someNode == nil {
			return nil, errors.New("not found accepted node")
		}
		meshedNodes.Add(someNode)
	}

	return meshedNodes, nil
}

// CloseSNsConnections closes connections to all nodes
func (m *MeshHandler) CloseSNsConnections(ctx context.Context, nodesDone chan struct{}) error {
	var err error

	close(nodesDone)

	log.WithContext(ctx).Debug("close connections to supernodes")

	for i := range m.Nodes {
		if err := m.Nodes[i].ConnectionInterface.Close(); err != nil {
			log.WithContext(ctx).WithFields(log.Fields{
				"pastelId": m.Nodes[i].PastelID(),
				"addr":     m.Nodes[i].String(),
			}).WithError(err).Errorf("close supernode connection failed")
		}
	}

	return err
}

// ValidBurnTxID returns whether the burn txid is valid at ALL SNs
func (m *MeshHandler) CheckSNReportedState() bool {
	for _, someNode := range m.Nodes {
		if !someNode.IsRemoteState() {
			return false
		}
	}
	return true
}
